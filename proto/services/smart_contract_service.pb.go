//*
// # Smart Contract Service
// gRPC service definitions for calls to the Hedera EVM-compatible
// Smart Contract service.
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
// source: smart_contract_service.proto

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

var File_smart_contract_service_proto protoreflect.FileDescriptor

const file_smart_contract_service_proto_rawDesc = "" +
	"\n" +
	"\x1csmart_contract_service.proto\x12\x05proto\x1a\x1atransaction_response.proto\x1a\vquery.proto\x1a\x0eresponse.proto\x1a\x11transaction.proto2\x86\x06\n" +
	"\x14SmartContractService\x12@\n" +
	"\x0ecreateContract\x12\x12.proto.Transaction\x1a\x1a.proto.TransactionResponse\x12@\n" +
	"\x0eupdateContract\x12\x12.proto.Transaction\x1a\x1a.proto.TransactionResponse\x12D\n" +
	"\x12contractCallMethod\x12\x12.proto.Transaction\x1a\x1a.proto.TransactionResponse\x128\n" +
	"\x17contractCallLocalMethod\x12\f.proto.Query\x1a\x0f.proto.Response\x120\n" +
	"\x0fgetContractInfo\x12\f.proto.Query\x1a\x0f.proto.Response\x124\n" +
	"\x13ContractGetBytecode\x12\f.proto.Query\x1a\x0f.proto.Response\x125\n" +
	"\x0fgetBySolidityID\x12\f.proto.Query\x1a\x0f.proto.Response\"\x03\x88\x02\x01\x12=\n" +
	"\x17getTxRecordByContractID\x12\f.proto.Query\x1a\x0f.proto.Response\"\x03\x88\x02\x01\x12@\n" +
	"\x0edeleteContract\x12\x12.proto.Transaction\x1a\x1a.proto.TransactionResponse\x12C\n" +
	"\fsystemDelete\x12\x12.proto.Transaction\x1a\x1a.proto.TransactionResponse\"\x03\x88\x02\x01\x12E\n" +
	"\x0esystemUndelete\x12\x12.proto.Transaction\x1a\x1a.proto.TransactionResponse\"\x03\x88\x02\x01\x12>\n" +
	"\fcallEthereum\x12\x12.proto.Transaction\x1a\x1a.proto.TransactionResponseB(\n" +
	"&com.hederahashgraph.service.proto.javab\x06proto3"

var file_smart_contract_service_proto_goTypes = []any{
	(*Transaction)(nil),         // 0: proto.Transaction
	(*Query)(nil),               // 1: proto.Query
	(*TransactionResponse)(nil), // 2: proto.TransactionResponse
	(*Response)(nil),            // 3: proto.Response
}
var file_smart_contract_service_proto_depIdxs = []int32{
	0,  // 0: proto.SmartContractService.createContract:input_type -> proto.Transaction
	0,  // 1: proto.SmartContractService.updateContract:input_type -> proto.Transaction
	0,  // 2: proto.SmartContractService.contractCallMethod:input_type -> proto.Transaction
	1,  // 3: proto.SmartContractService.contractCallLocalMethod:input_type -> proto.Query
	1,  // 4: proto.SmartContractService.getContractInfo:input_type -> proto.Query
	1,  // 5: proto.SmartContractService.ContractGetBytecode:input_type -> proto.Query
	1,  // 6: proto.SmartContractService.getBySolidityID:input_type -> proto.Query
	1,  // 7: proto.SmartContractService.getTxRecordByContractID:input_type -> proto.Query
	0,  // 8: proto.SmartContractService.deleteContract:input_type -> proto.Transaction
	0,  // 9: proto.SmartContractService.systemDelete:input_type -> proto.Transaction
	0,  // 10: proto.SmartContractService.systemUndelete:input_type -> proto.Transaction
	0,  // 11: proto.SmartContractService.callEthereum:input_type -> proto.Transaction
	2,  // 12: proto.SmartContractService.createContract:output_type -> proto.TransactionResponse
	2,  // 13: proto.SmartContractService.updateContract:output_type -> proto.TransactionResponse
	2,  // 14: proto.SmartContractService.contractCallMethod:output_type -> proto.TransactionResponse
	3,  // 15: proto.SmartContractService.contractCallLocalMethod:output_type -> proto.Response
	3,  // 16: proto.SmartContractService.getContractInfo:output_type -> proto.Response
	3,  // 17: proto.SmartContractService.ContractGetBytecode:output_type -> proto.Response
	3,  // 18: proto.SmartContractService.getBySolidityID:output_type -> proto.Response
	3,  // 19: proto.SmartContractService.getTxRecordByContractID:output_type -> proto.Response
	2,  // 20: proto.SmartContractService.deleteContract:output_type -> proto.TransactionResponse
	2,  // 21: proto.SmartContractService.systemDelete:output_type -> proto.TransactionResponse
	2,  // 22: proto.SmartContractService.systemUndelete:output_type -> proto.TransactionResponse
	2,  // 23: proto.SmartContractService.callEthereum:output_type -> proto.TransactionResponse
	12, // [12:24] is the sub-list for method output_type
	0,  // [0:12] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_smart_contract_service_proto_init() }
func file_smart_contract_service_proto_init() {
	if File_smart_contract_service_proto != nil {
		return
	}
	file_transaction_response_proto_init()
	file_query_proto_init()
	file_response_proto_init()
	file_transaction_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_smart_contract_service_proto_rawDesc), len(file_smart_contract_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_smart_contract_service_proto_goTypes,
		DependencyIndexes: file_smart_contract_service_proto_depIdxs,
	}.Build()
	File_smart_contract_service_proto = out.File
	file_smart_contract_service_proto_goTypes = nil
	file_smart_contract_service_proto_depIdxs = nil
}
