package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testproject/api/v1"
	"testproject/service"
	"time"
)

func main() {
	rpcConf := &service.RpcPoolConf{
		MinCons: 10,
		MaxCons: 200,
		Tcp:     "localhost:9000",
	}

	demo(rpcConf)
}

func demo(rpcConf *service.RpcPoolConf) {
	var wg sync.WaitGroup
	waitChan := make(chan bool, 1)
	rpcPool := service.NewRpcPool(rpcConf)

	defer func() {
		if err:=recover();err!=nil{
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			os.Exit(1)
		}
	}()

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			debitCard, err := rpcPool.GetConn()
			defer func() {
				err = rpcPool.ReturnBack(debitCard)
				fmt.Println(fmt.Sprintf("归还err %v", err))
			}()

			if err != nil {
				fmt.Println(fmt.Sprintf("获取链接失败 %v", err))
				return
			}
			dialRpc(debitCard)

		}()
	}

	wg.Wait()
	fmt.Println("All goroutines finished!")

	<-waitChan
}

func dialRpc(debitCard *service.DebitCard) {
	start := time.Now().Unix()

	defer func() {
		if err:=recover();err!=nil{
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
			os.Exit(1)
		}
	}()

	client := v1.NewRightsClient(debitCard.Conn)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	reply, err := client.CheckRights(ctx, &v1.CheckRightsRequest{
		UserId:      "1111111",
		ActionValue: "22222222",
	})

	if err != nil && reply == nil {
		fmt.Println("该连接已经失效,请注意回收")
		debitCard.IsExpired = true
	}

	fmt.Println(fmt.Sprintf("总共经历时间 %v reply== %v rpc err %v", time.Now().Unix()-start, reply, err))
}

//func dialRpc1() {
//
//	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
//
//	start := time.Now().Unix()
//
//	client := v1.NewRightsClient(conn)
//	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
//	defer cancel()
//	reply, err := client.CheckRights(ctx, &v1.CheckRightsRequest{
//		UserId:      "1111111",
//		ActionValue: "22222222",
//	})
//
//	fmt.Println(fmt.Sprintf("rpc err %v", err))
//	fmt.Println(reply)
//	fmt.Println(fmt.Sprintf("总共经历时间 %v", time.Now().Unix()-start))
//}
