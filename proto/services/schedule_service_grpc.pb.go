//*
// # Schedule Service
// gRPC service definitions for the Schedule Service.
//
// ### Keywords
// The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
// "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
// document are to be interpreted as described in
// [RFC2119](https://www.ietf.org/rfc/rfc2119) and clarified in
// [RFC8174](https://www.ietf.org/rfc/rfc8174).

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: schedule_service.proto

package services

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ScheduleService_CreateSchedule_FullMethodName  = "/proto.ScheduleService/createSchedule"
	ScheduleService_SignSchedule_FullMethodName    = "/proto.ScheduleService/signSchedule"
	ScheduleService_DeleteSchedule_FullMethodName  = "/proto.ScheduleService/deleteSchedule"
	ScheduleService_GetScheduleInfo_FullMethodName = "/proto.ScheduleService/getScheduleInfo"
)

// ScheduleServiceClient is the client API for ScheduleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// *
// Transactions and queries for the Schedule Service.<br/>
// The Schedule Service enables transactions to be submitted without all
// required signatures and offers a `scheduleSign` transaction to provide
// additional signatures independently after the schedule is created. The
// scheduled transaction may be executed immediately when all required
// signatures are present, or at expiration if "long term" schedules
// are enabled in network configuration.
//
// ### Execution
// Scheduled transactions SHALL be executed under the following conditions.
//  1. If "long term" schedules are enabled and `wait_for_expiry` is set for
//     that schedule then the transaction SHALL NOT be executed before the
//     network consensus time matches or exceeds the `expiration_time` field
//     for that schedule.
//  1. If "long term" schedules are enabled and `wait_for_expiry` is _not_ set
//     for that schedule, then the transaction SHALL be executed when all
//     signatures required by the scheduled transaction are active for that
//     schedule. This MAY be immediately after the `scheduleCreate` or a
//     subsequent `scheduleSign` transaction, or MAY be at expiration if
//     the signature requirements are met at that time.
//  1. If "long term" schedules are _disabled_, then the scheduled transaction
//     SHALL be executed immediately after all signature requirements for the
//     scheduled transaction are met during the `scheduleCreate` or a subsequent
//     `scheduleSign` transaction. The scheduled transaction SHALL NOT be
//     on expiration when "long term" schedules are disabled.
//
// A schedule SHALL remain in state and MAY be queried with a `getScheduleInfo`
// transaction after execution, until the schedule expires.<br/>
// When network consensus time matches or exceeds the `expiration_time` for
// a schedule, that schedule SHALL be removed from state, whether it has
// executed or not.<br/>
// If "long term" schedules are _disabled_, the maximum expiration time SHALL
// be the consensus time of the `scheduleCreate` transaction extended by
// the network configuration value `ledger.scheduleTxExpiryTimeSecs`.
//
// ### Block Stream Effects
// When a scheduled transaction is executed, the timestamp in the transaction
// identifier for that transaction SHALL be 1 nanosecond after the consensus
// timestamp for the transaction that resulted in its execution. If execution
// occurred at expiration, that transaction may be almost any transaction,
// including another scheduled transaction that executed at expiration.<br/>
// The transaction identifier for a scheduled transaction that is executed
// SHALL have the `scheduled` flag set and SHALL inherit the `accountID` and
// `transactionValidStart` values from the `scheduleCreate` that created the
// schedule.<br/>
// The `scheduleRef` property of the record for a scheduled transaction SHALL
// be populated with the schedule identifier of the schedule that executed.
type ScheduleServiceClient interface {
	// *
	// Create a new Schedule.
	// <p>
	// If all signature requirements are met with this transaction, the
	// scheduled transaction MAY execute immediately.
	CreateSchedule(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*TransactionResponse, error)
	// *
	// Add signatures to an existing schedule.
	// <p>
	// Signatures on this transaction SHALL be added to the set of active
	// signatures on the schedule, and MAY result in execution of the
	// scheduled transaction if all signature requirements are met.
	SignSchedule(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*TransactionResponse, error)
	// *
	// Mark an existing schedule deleted.
	// <p>
	// Once deleted a schedule SHALL NOT be executed and any subsequent
	// `scheduleSign` transaction SHALL fail.
	DeleteSchedule(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*TransactionResponse, error)
	// *
	// Retrieve the metadata for a schedule.
	GetScheduleInfo(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error)
}

type scheduleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewScheduleServiceClient(cc grpc.ClientConnInterface) ScheduleServiceClient {
	return &scheduleServiceClient{cc}
}

func (c *scheduleServiceClient) CreateSchedule(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*TransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, ScheduleService_CreateSchedule_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scheduleServiceClient) SignSchedule(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*TransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, ScheduleService_SignSchedule_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scheduleServiceClient) DeleteSchedule(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*TransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, ScheduleService_DeleteSchedule_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scheduleServiceClient) GetScheduleInfo(ctx context.Context, in *Query, opts ...grpc.CallOption) (*Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Response)
	err := c.cc.Invoke(ctx, ScheduleService_GetScheduleInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScheduleServiceServer is the server API for ScheduleService service.
// All implementations must embed UnimplementedScheduleServiceServer
// for forward compatibility.
//
// *
// Transactions and queries for the Schedule Service.<br/>
// The Schedule Service enables transactions to be submitted without all
// required signatures and offers a `scheduleSign` transaction to provide
// additional signatures independently after the schedule is created. The
// scheduled transaction may be executed immediately when all required
// signatures are present, or at expiration if "long term" schedules
// are enabled in network configuration.
//
// ### Execution
// Scheduled transactions SHALL be executed under the following conditions.
//  1. If "long term" schedules are enabled and `wait_for_expiry` is set for
//     that schedule then the transaction SHALL NOT be executed before the
//     network consensus time matches or exceeds the `expiration_time` field
//     for that schedule.
//  1. If "long term" schedules are enabled and `wait_for_expiry` is _not_ set
//     for that schedule, then the transaction SHALL be executed when all
//     signatures required by the scheduled transaction are active for that
//     schedule. This MAY be immediately after the `scheduleCreate` or a
//     subsequent `scheduleSign` transaction, or MAY be at expiration if
//     the signature requirements are met at that time.
//  1. If "long term" schedules are _disabled_, then the scheduled transaction
//     SHALL be executed immediately after all signature requirements for the
//     scheduled transaction are met during the `scheduleCreate` or a subsequent
//     `scheduleSign` transaction. The scheduled transaction SHALL NOT be
//     on expiration when "long term" schedules are disabled.
//
// A schedule SHALL remain in state and MAY be queried with a `getScheduleInfo`
// transaction after execution, until the schedule expires.<br/>
// When network consensus time matches or exceeds the `expiration_time` for
// a schedule, that schedule SHALL be removed from state, whether it has
// executed or not.<br/>
// If "long term" schedules are _disabled_, the maximum expiration time SHALL
// be the consensus time of the `scheduleCreate` transaction extended by
// the network configuration value `ledger.scheduleTxExpiryTimeSecs`.
//
// ### Block Stream Effects
// When a scheduled transaction is executed, the timestamp in the transaction
// identifier for that transaction SHALL be 1 nanosecond after the consensus
// timestamp for the transaction that resulted in its execution. If execution
// occurred at expiration, that transaction may be almost any transaction,
// including another scheduled transaction that executed at expiration.<br/>
// The transaction identifier for a scheduled transaction that is executed
// SHALL have the `scheduled` flag set and SHALL inherit the `accountID` and
// `transactionValidStart` values from the `scheduleCreate` that created the
// schedule.<br/>
// The `scheduleRef` property of the record for a scheduled transaction SHALL
// be populated with the schedule identifier of the schedule that executed.
type ScheduleServiceServer interface {
	// *
	// Create a new Schedule.
	// <p>
	// If all signature requirements are met with this transaction, the
	// scheduled transaction MAY execute immediately.
	CreateSchedule(context.Context, *Transaction) (*TransactionResponse, error)
	// *
	// Add signatures to an existing schedule.
	// <p>
	// Signatures on this transaction SHALL be added to the set of active
	// signatures on the schedule, and MAY result in execution of the
	// scheduled transaction if all signature requirements are met.
	SignSchedule(context.Context, *Transaction) (*TransactionResponse, error)
	// *
	// Mark an existing schedule deleted.
	// <p>
	// Once deleted a schedule SHALL NOT be executed and any subsequent
	// `scheduleSign` transaction SHALL fail.
	DeleteSchedule(context.Context, *Transaction) (*TransactionResponse, error)
	// *
	// Retrieve the metadata for a schedule.
	GetScheduleInfo(context.Context, *Query) (*Response, error)
	mustEmbedUnimplementedScheduleServiceServer()
}

// UnimplementedScheduleServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedScheduleServiceServer struct{}

func (UnimplementedScheduleServiceServer) CreateSchedule(context.Context, *Transaction) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSchedule not implemented")
}
func (UnimplementedScheduleServiceServer) SignSchedule(context.Context, *Transaction) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignSchedule not implemented")
}
func (UnimplementedScheduleServiceServer) DeleteSchedule(context.Context, *Transaction) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSchedule not implemented")
}
func (UnimplementedScheduleServiceServer) GetScheduleInfo(context.Context, *Query) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetScheduleInfo not implemented")
}
func (UnimplementedScheduleServiceServer) mustEmbedUnimplementedScheduleServiceServer() {}
func (UnimplementedScheduleServiceServer) testEmbeddedByValue()                         {}

// UnsafeScheduleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScheduleServiceServer will
// result in compilation errors.
type UnsafeScheduleServiceServer interface {
	mustEmbedUnimplementedScheduleServiceServer()
}

func RegisterScheduleServiceServer(s grpc.ServiceRegistrar, srv ScheduleServiceServer) {
	// If the following call pancis, it indicates UnimplementedScheduleServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ScheduleService_ServiceDesc, srv)
}

func _ScheduleService_CreateSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScheduleServiceServer).CreateSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScheduleService_CreateSchedule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScheduleServiceServer).CreateSchedule(ctx, req.(*Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScheduleService_SignSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScheduleServiceServer).SignSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScheduleService_SignSchedule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScheduleServiceServer).SignSchedule(ctx, req.(*Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScheduleService_DeleteSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScheduleServiceServer).DeleteSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScheduleService_DeleteSchedule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScheduleServiceServer).DeleteSchedule(ctx, req.(*Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScheduleService_GetScheduleInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScheduleServiceServer).GetScheduleInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScheduleService_GetScheduleInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScheduleServiceServer).GetScheduleInfo(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

// ScheduleService_ServiceDesc is the grpc.ServiceDesc for ScheduleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ScheduleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ScheduleService",
	HandlerType: (*ScheduleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "createSchedule",
			Handler:    _ScheduleService_CreateSchedule_Handler,
		},
		{
			MethodName: "signSchedule",
			Handler:    _ScheduleService_SignSchedule_Handler,
		},
		{
			MethodName: "deleteSchedule",
			Handler:    _ScheduleService_DeleteSchedule_Handler,
		},
		{
			MethodName: "getScheduleInfo",
			Handler:    _ScheduleService_GetScheduleInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "schedule_service.proto",
}
