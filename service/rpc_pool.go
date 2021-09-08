package service

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
	"time"
)

type RpcPoolConf struct {
	Tcp     string
	MaxCons int
	MinCons int
}

type RpcPool struct {
	rpcConf      *RpcPoolConf
	connQueue    chan *grpc.ClientConn
	curCons      int                //当前连接数
	maxCons      int                //最大连接数
	minCons      int                //最小连接数
	freeCons     int                //空闲连接数
	debitCardMap map[int]*DebitCard //借记卡
	poolLock     sync.RWMutex       //线程锁
}

type DebitCard struct {
	Conn      *grpc.ClientConn
	Key       int
	IsExpired bool
}

func NewDebitCard(conn *grpc.ClientConn, key int) *DebitCard {
	return &DebitCard{
		Conn:      conn,
		Key:       key,
		IsExpired: false,
	}
}

var once sync.Once

func NewRpcPool(rpcConf *RpcPoolConf) *RpcPool {
	var rpcPool *RpcPool

	once.Do(func() {
		rpcPool = initRpcPool(rpcConf)
		fmt.Println("初始化rpc线程完成")
	})

	return rpcPool
}

func initRpcPool(rpcConf *RpcPoolConf) *RpcPool {
	var rpcPool = RpcPool{
		rpcConf:      rpcConf,
		connQueue:    make(chan *grpc.ClientConn, rpcConf.MaxCons),
		curCons:      rpcConf.MinCons,
		freeCons:     rpcConf.MinCons,
		minCons:      rpcConf.MinCons,
		maxCons:      rpcConf.MaxCons,
		debitCardMap: make(map[int]*DebitCard),
	}

	fmt.Println("初始化线程池...")
	for i := 1; i <= rpcConf.MinCons; i++ {
		fmt.Println(fmt.Sprintf("初始化线程 %v ", i))
		conn, err := grpc.Dial(rpcConf.Tcp, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		rpcPool.connQueue <- conn
	}

	go rpcPool.metrics()
	go rpcPool.dynamicConf()

	return &rpcPool
}

func (r *RpcPool) GetConn() (debitCard *DebitCard, err error) {
	var depNum int
	var conn *grpc.ClientConn

	conn = <-r.connQueue

	r.poolLock.Lock()
	defer r.poolLock.Unlock()

	fmt.Println(fmt.Sprintf("获取到conn %v", conn.Target()))

	for {
		depNum = rand.Intn(1000000)
		fmt.Println(fmt.Sprintf("depNum = %v", depNum))
		if _, ok := r.debitCardMap[depNum]; ok {
			continue
		}
		break
	}

	r.debitCardMap[depNum] = NewDebitCard(conn, depNum)
	r.freeCons--

	return r.debitCardMap[depNum], nil
}

func (r *RpcPool) ReturnBack(debitCard *DebitCard) error {
	r.poolLock.Lock()
	defer r.poolLock.Unlock()

	if _, ok := r.debitCardMap[debitCard.Key]; !ok {
		fmt.Println("借记卡不对")
		return errors.New("借记卡不对")
	}

	if r.debitCardMap[debitCard.Key] != debitCard {
		fmt.Println("链接不对")
		return errors.New("链接不对")
	}

	delete(r.debitCardMap, debitCard.Key)

	if debitCard.IsExpired {
		fmt.Println("该链接已经过期")
		return nil
	}

	r.freeCons++
	r.connQueue <- debitCard.Conn

	return nil
}

func (r *RpcPool) dynamicConf() {
	for {
		time.Sleep(1 * time.Second)

		if r.freeCons <= 5 {
			increaseNum := r.curCons / 2
			for {
				if r.curCons+increaseNum <= r.maxCons {
					break
				}
				increaseNum = increaseNum / 2
			}

			r.increaseRpcPoolConn(increaseNum)
			continue
		}

		if r.freeCons > r.curCons/2 {
			reduceNum := r.curCons / 2
			if r.curCons-reduceNum <= r.minCons {
				continue
			}
			r.reduceRpcPoolConn(reduceNum)
			continue
		}
	}

}

func (r *RpcPool) increaseRpcPoolConn(increaseNum int) {
	r.poolLock.Lock()
	defer r.poolLock.Unlock()

	if increaseNum <= 0 {
		return
	}

	fmt.Println(fmt.Sprintf("开始扩容，扩容数量为 %v", increaseNum))

	for i := 1; i <= increaseNum; i++ {
		fmt.Println(fmt.Sprintf("新增线程 %v ", i))
		conn, err := grpc.Dial(r.rpcConf.Tcp, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		r.connQueue <- conn
	}

	r.curCons = r.curCons + increaseNum
	r.freeCons = r.freeCons + increaseNum
	fmt.Println(fmt.Sprintf("统计当前数据 %v", r))
}

func (r *RpcPool) reduceRpcPoolConn(reduceNum int) {
	r.poolLock.Lock()
	defer r.poolLock.Unlock()

	if r.curCons-reduceNum <= r.minCons {
		return
	}

	if r.curCons <= reduceNum {
		return
	}

	for i := 1; i <= reduceNum; i++ {
		conn := <-r.connQueue
		_ = conn.Close()
	}

	r.curCons = r.curCons - reduceNum
	r.freeCons = r.freeCons - reduceNum
	fmt.Println(fmt.Sprintf("统计当前数据 %v", r))
}

func (r *RpcPool) metrics() {
	for {
		time.Sleep(time.Second * 2)
		r.poolLock.RLock()

		fmt.Println("统计信息")
		fmt.Println(r)

		r.poolLock.RUnlock()
	}

}

func (r *RpcPool) String() string {
	return fmt.Sprintf("当前连接数 %v, 当前空闲连接数 %v, 当前最大连接数 %v, 当前最小连接数 %v", r.curCons, r.freeCons, r.maxCons, r.minCons)
}
