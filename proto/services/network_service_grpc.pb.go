// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: network_service.proto

package services

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

// NetworkServiceClient is the client API for NetworkService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetworkServiceClient interface {
	// *
	// Retrieve the active versions of Hedera Services and API messages.
	GetVersionInfo(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error)
	// *
	// Request detail information about an account.
	// <p>
	// The returned information SHALL include balance and allowances.<br/>
	// The returned information SHALL NOT include a list of account records.
	GetAccountDetails(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error)
	// Deprecated: Do not use.
	// *
	// Retrieve the time, in nanoseconds, spent in direct processing for one or
	// more recent transactions.
	// <p>
	// For each transaction identifier provided, if that transaction is
	// sufficiently recent (that is, it is within the range of the
	// configuration value `stats.executionTimesToTrack`), the node SHALL
	// return the time, in nanoseconds, spent to directly process that
	// transaction (that is, excluding time to reach consensus).<br/>
	// Note that because each node processes every transaction for the Hedera
	// network, this query MAY be sent to any node.
	// <p>
	// <blockquote>Important<blockquote>
	// This query is obsolete, not supported, and SHALL fail with a pre-check
	// result of `NOT_SUPPORTED`.</blockquote></blockquote>
	GetExecutionTime(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error)
	// Deprecated: Do not use.
	// *
	// Submit a transaction that wraps another transaction which will
	// skip most validation.
	// <p>
	// <blockquote>Important<blockquote>
	// This query is obsolete, not supported, and SHALL fail with a pre-check
	// result of `NOT_SUPPORTED`.
	// </blockquote></blockquote>
	UncheckedSubmit(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*TransactionResponse, error)
}

type networkServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkServiceClient(cc grpc.ClientConnInterface) NetworkServiceClient {
	return &networkServiceClient{cc}
}

func (c *networkServiceClient) GetVersionInfo(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.NetworkService/getVersionInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkServiceClient) GetAccountDetails(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.NetworkService/getAccountDetails", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Deprecated: Do not use.
func (c *networkServiceClient) GetExecutionTime(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.NetworkService/getExecutionTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Deprecated: Do not use.
func (c *networkServiceClient) UncheckedSubmit(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*TransactionResponse, error) {
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, "/proto.NetworkService/uncheckedSubmit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkServiceServer is the server API for NetworkService service.
// All implementations must embed UnimplementedNetworkServiceServer
// for forward compatibility
type NetworkServiceServer interface {
	// *
	// Retrieve the active versions of Hedera Services and API messages.
	GetVersionInfo(context.Context, *Query) (*Response, error)
	// *
	// Request detail information about an account.
	// <p>
	// The returned information SHALL include balance and allowances.<br/>
	// The returned information SHALL NOT include a list of account records.
	GetAccountDetails(context.Context, *Query) (*Response, error)
	// Deprecated: Do not use.
	// *
	// Retrieve the time, in nanoseconds, spent in direct processing for one or
	// more recent transactions.
	// <p>
	// For each transaction identifier provided, if that transaction is
	// sufficiently recent (that is, it is within the range of the
	// configuration value `stats.executionTimesToTrack`), the node SHALL
	// return the time, in nanoseconds, spent to directly process that
	// transaction (that is, excluding time to reach consensus).<br/>
	// Note that because each node processes every transaction for the Hedera
	// network, this query MAY be sent to any node.
	// <p>
	// <blockquote>Important<blockquote>
	// This query is obsolete, not supported, and SHALL fail with a pre-check
	// result of `NOT_SUPPORTED`.</blockquote></blockquote>
	GetExecutionTime(context.Context, *Query) (*Response, error)
	// Deprecated: Do not use.
	// *
	// Submit a transaction that wraps another transaction which will
	// skip most validation.
	// <p>
	// <blockquote>Important<blockquote>
	// This query is obsolete, not supported, and SHALL fail with a pre-check
	// result of `NOT_SUPPORTED`.
	// </blockquote></blockquote>
	UncheckedSubmit(context.Context, *Transaction) (*TransactionResponse, error)
	mustEmbedUnimplementedNetworkServiceServer()
}

// UnimplementedNetworkServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNetworkServiceServer struct {
}

func (UnimplementedNetworkServiceServer) GetVersionInfo(context.Context, *Query) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersionInfo not implemented")
}
func (UnimplementedNetworkServiceServer) GetAccountDetails(context.Context, *Query) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccountDetails not implemented")
}
func (UnimplementedNetworkServiceServer) GetExecutionTime(context.Context, *Query) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExecutionTime not implemented")
}
func (UnimplementedNetworkServiceServer) UncheckedSubmit(context.Context, *Transaction) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UncheckedSubmit not implemented")
}
func (UnimplementedNetworkServiceServer) mustEmbedUnimplementedNetworkServiceServer() {}

// UnsafeNetworkServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetworkServiceServer will
// result in compilation errors.
type UnsafeNetworkServiceServer interface {
	mustEmbedUnimplementedNetworkServiceServer()
}

func RegisterNetworkServiceServer(s grpc.ServiceRegistrar, srv NetworkServiceServer) {
	s.RegisterService(&NetworkService_ServiceDesc, srv)
}

func _NetworkService_GetVersionInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceServer).GetVersionInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NetworkService/getVersionInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceServer).GetVersionInfo(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkService_GetAccountDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceServer).GetAccountDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NetworkService/getAccountDetails",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceServer).GetAccountDetails(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkService_GetExecutionTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceServer).GetExecutionTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NetworkService/getExecutionTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceServer).GetExecutionTime(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkService_UncheckedSubmit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceServer).UncheckedSubmit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NetworkService/uncheckedSubmit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceServer).UncheckedSubmit(ctx, req.(*Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

// NetworkService_ServiceDesc is the grpc.ServiceDesc for NetworkService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetworkService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.NetworkService",
	HandlerType: (*NetworkServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getVersionInfo",
			Handler:    _NetworkService_GetVersionInfo_Handler,
		},
		{
			MethodName: "getAccountDetails",
			Handler:    _NetworkService_GetAccountDetails_Handler,
		},
		{
			MethodName: "getExecutionTime",
			Handler:    _NetworkService_GetExecutionTime_Handler,
		},
		{
			MethodName: "uncheckedSubmit",
			Handler:    _NetworkService_UncheckedSubmit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "network_service.proto",
}
