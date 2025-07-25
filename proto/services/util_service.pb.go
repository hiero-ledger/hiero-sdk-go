//*
// # Utility Service
// This service provides a transaction to generate a deterministic
// pseudo-random value, either a 32-bit integer within a requested range
// or a 384-bit byte array.
//
// ### Keywords
// The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
// "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
// document are to be interpreted as described in
// [RFC2119](https://www.ietf.org/rfc/rfc2119) and clarified in
// [RFC8174](https://www.ietf.org/rfc/rfc8174).

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: util_service.proto

package services

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_util_service_proto protoreflect.FileDescriptor

const file_util_service_proto_rawDesc = "" +
	"\n" +
	"\x12util_service.proto\x12\x05proto\x1a\x1atransaction_response.proto\x1a\x11transaction.proto2\x84\x01\n" +
	"\vUtilService\x126\n" +
	"\x04prng\x12\x12.proto.Transaction\x1a\x1a.proto.TransactionResponse\x12=\n" +
	"\vatomicBatch\x12\x12.proto.Transaction\x1a\x1a.proto.TransactionResponseB(\n" +
	"&com.hederahashgraph.service.proto.javab\x06proto3"

var file_util_service_proto_goTypes = []any{
	(*Transaction)(nil),         // 0: proto.Transaction
	(*TransactionResponse)(nil), // 1: proto.TransactionResponse
}
var file_util_service_proto_depIdxs = []int32{
	0, // 0: proto.UtilService.prng:input_type -> proto.Transaction
	0, // 1: proto.UtilService.atomicBatch:input_type -> proto.Transaction
	1, // 2: proto.UtilService.prng:output_type -> proto.TransactionResponse
	1, // 3: proto.UtilService.atomicBatch:output_type -> proto.TransactionResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_util_service_proto_init() }
func file_util_service_proto_init() {
	if File_util_service_proto != nil {
		return
	}
	file_transaction_response_proto_init()
	file_transaction_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_util_service_proto_rawDesc), len(file_util_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_util_service_proto_goTypes,
		DependencyIndexes: file_util_service_proto_depIdxs,
	}.Build()
	File_util_service_proto = out.File
	file_util_service_proto_goTypes = nil
	file_util_service_proto_depIdxs = nil
}
