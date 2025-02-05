//*
// # Contract Call
// Transaction body for calls to a Smart Contract.
//
// ### Keywords
// The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
// "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
// document are to be interpreted as described in
// [RFC2119](https://www.ietf.org/rfc/rfc2119) and clarified in
// [RFC8174](https://www.ietf.org/rfc/rfc8174).

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: contract_call.proto

package services

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// *
// Call a function of a given smart contract, providing function parameter
// inputs as needed.
//
// Resource ("gas") charges SHALL include all relevant fees incurred by the
// contract execution, including any storage required.<br/>
// The total transaction fee SHALL incorporate all of the "gas" actually
// consumed as well as the standard fees for transaction handling, data
// transfers, signature verification, etc...<br/>
// The response SHALL contain the output returned by the function call.
//
// ### Block Stream Effects
// A `CallContractOutput` message SHALL be emitted for each transaction.
type ContractCallTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The ID of a smart contract to call.
	ContractID *ContractID `protobuf:"bytes,1,opt,name=contractID,proto3" json:"contractID,omitempty"`
	// *
	// A maximum limit to the amount of gas to use for this call.
	// <p>
	// The network SHALL charge the greater of the following, but
	// SHALL NOT charge more than the value of this field.
	// <ol>
	//
	//	<li>The actual gas consumed by the smart contract call.</li>
	//	<li>`80%` of this value.</li>
	//
	// </ol>
	// The `80%` factor encourages reasonable estimation, while allowing for
	// some overage to ensure successful execution.
	Gas int64 `protobuf:"varint,2,opt,name=gas,proto3" json:"gas,omitempty"`
	// *
	// An amount of tinybar sent via this contract call.
	// <p>
	// If this is non-zero, the function MUST be `payable`.
	Amount int64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	// *
	// The smart contract function to call.
	// <p>
	// This MUST contain The application binary interface (ABI) encoding of the
	// function call per the Ethereum contract ABI standard, giving the
	// function signature and arguments being passed to the function.
	FunctionParameters []byte `protobuf:"bytes,4,opt,name=functionParameters,proto3" json:"functionParameters,omitempty"`
}

func (x *ContractCallTransactionBody) Reset() {
	*x = ContractCallTransactionBody{}
	mi := &file_contract_call_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContractCallTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractCallTransactionBody) ProtoMessage() {}

func (x *ContractCallTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_contract_call_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractCallTransactionBody.ProtoReflect.Descriptor instead.
func (*ContractCallTransactionBody) Descriptor() ([]byte, []int) {
	return file_contract_call_proto_rawDescGZIP(), []int{0}
}

func (x *ContractCallTransactionBody) GetContractID() *ContractID {
	if x != nil {
		return x.ContractID
	}
	return nil
}

func (x *ContractCallTransactionBody) GetGas() int64 {
	if x != nil {
		return x.Gas
	}
	return 0
}

func (x *ContractCallTransactionBody) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *ContractCallTransactionBody) GetFunctionParameters() []byte {
	if x != nil {
		return x.FunctionParameters
	}
	return nil
}

var File_contract_call_proto protoreflect.FileDescriptor

var file_contract_call_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x63, 0x61, 0x6c, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61,
	0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xaa, 0x01, 0x0a, 0x1b, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x43, 0x61, 0x6c, 0x6c,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12,
	0x31, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x49, 0x44, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x49, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x61, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x03, 0x67, 0x61, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x12,
	0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x12, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x42, 0x26, 0x0a, 0x22,
	0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61,
	0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contract_call_proto_rawDescOnce sync.Once
	file_contract_call_proto_rawDescData = file_contract_call_proto_rawDesc
)

func file_contract_call_proto_rawDescGZIP() []byte {
	file_contract_call_proto_rawDescOnce.Do(func() {
		file_contract_call_proto_rawDescData = protoimpl.X.CompressGZIP(file_contract_call_proto_rawDescData)
	})
	return file_contract_call_proto_rawDescData
}

var file_contract_call_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_contract_call_proto_goTypes = []any{
	(*ContractCallTransactionBody)(nil), // 0: proto.ContractCallTransactionBody
	(*ContractID)(nil),                  // 1: proto.ContractID
}
var file_contract_call_proto_depIdxs = []int32{
	1, // 0: proto.ContractCallTransactionBody.contractID:type_name -> proto.ContractID
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_contract_call_proto_init() }
func file_contract_call_proto_init() {
	if File_contract_call_proto != nil {
		return
	}
	file_basic_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contract_call_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_contract_call_proto_goTypes,
		DependencyIndexes: file_contract_call_proto_depIdxs,
		MessageInfos:      file_contract_call_proto_msgTypes,
	}.Build()
	File_contract_call_proto = out.File
	file_contract_call_proto_rawDesc = nil
	file_contract_call_proto_goTypes = nil
	file_contract_call_proto_depIdxs = nil
}
