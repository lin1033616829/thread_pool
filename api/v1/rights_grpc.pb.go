// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RightsClient is the client API for Rights service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RightsClient interface {
	CheckRights(ctx context.Context, in *CheckRightsRequest, opts ...grpc.CallOption) (*CheckRightsReply, error)
}

type rightsClient struct {
	cc grpc.ClientConnInterface
}

func NewRightsClient(cc grpc.ClientConnInterface) RightsClient {
	return &rightsClient{cc}
}

func (c *rightsClient) CheckRights(ctx context.Context, in *CheckRightsRequest, opts ...grpc.CallOption) (*CheckRightsReply, error) {
	out := new(CheckRightsReply)
	err := c.cc.Invoke(ctx, "/api.mind.v1.Rights/CheckRights", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RightsServer is the server API for Rights service.
// All implementations must embed UnimplementedRightsServer
// for forward compatibility
type RightsServer interface {
	CheckRights(context.Context, *CheckRightsRequest) (*CheckRightsReply, error)
	mustEmbedUnimplementedRightsServer()
}

// UnimplementedRightsServer must be embedded to have forward compatible implementations.
type UnimplementedRightsServer struct {
}

func (UnimplementedRightsServer) CheckRights(context.Context, *CheckRightsRequest) (*CheckRightsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckRights not implemented")
}
func (UnimplementedRightsServer) mustEmbedUnimplementedRightsServer() {}

// UnsafeRightsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RightsServer will
// result in compilation errors.
type UnsafeRightsServer interface {
	mustEmbedUnimplementedRightsServer()
}

func RegisterRightsServer(s grpc.ServiceRegistrar, srv RightsServer) {
	s.RegisterService(&Rights_ServiceDesc, srv)
}

func _Rights_CheckRights_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRightsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RightsServer).CheckRights(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.mind.v1.Rights/CheckRights",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RightsServer).CheckRights(ctx, req.(*CheckRightsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Rights_ServiceDesc is the grpc.ServiceDesc for Rights service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rights_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.mind.v1.Rights",
	HandlerType: (*RightsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckRights",
			Handler:    _Rights_CheckRights_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/mind/v1/rights.proto",
}