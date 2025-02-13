//*
// # Get Contract Bytecode
// A standard query to read the current bytecode for a smart contract.
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
// source: contract_get_bytecode.proto

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
// A transaction body to request the current bytecode for a smart contract.
type ContractGetBytecodeQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// Standard information sent with every query operation.<br/>
	// This includes the signed payment and what kind of response is requested
	// (cost, state proof, both, or neither).
	Header *QueryHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// A smart contract ID.
	// <p>
	// The network SHALL return bytecode for this smart contract, if successful.
	ContractID *ContractID `protobuf:"bytes,2,opt,name=contractID,proto3" json:"contractID,omitempty"`
}

func (x *ContractGetBytecodeQuery) Reset() {
	*x = ContractGetBytecodeQuery{}
	mi := &file_contract_get_bytecode_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContractGetBytecodeQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractGetBytecodeQuery) ProtoMessage() {}

func (x *ContractGetBytecodeQuery) ProtoReflect() protoreflect.Message {
	mi := &file_contract_get_bytecode_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractGetBytecodeQuery.ProtoReflect.Descriptor instead.
func (*ContractGetBytecodeQuery) Descriptor() ([]byte, []int) {
	return file_contract_get_bytecode_proto_rawDescGZIP(), []int{0}
}

func (x *ContractGetBytecodeQuery) GetHeader() *QueryHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *ContractGetBytecodeQuery) GetContractID() *ContractID {
	if x != nil {
		return x.ContractID
	}
	return nil
}

// *
// Information returned in response to a "get bytecode" query for a
// smart contract.
type ContractGetBytecodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The standard response information for queries.<br/>
	// This includes the values requested in the `QueryHeader`
	// (cost, state proof, both, or neither).
	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// The current bytecode of the requested smart contract.
	Bytecode []byte `protobuf:"bytes,6,opt,name=bytecode,proto3" json:"bytecode,omitempty"`
}

func (x *ContractGetBytecodeResponse) Reset() {
	*x = ContractGetBytecodeResponse{}
	mi := &file_contract_get_bytecode_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContractGetBytecodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractGetBytecodeResponse) ProtoMessage() {}

func (x *ContractGetBytecodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contract_get_bytecode_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractGetBytecodeResponse.ProtoReflect.Descriptor instead.
func (*ContractGetBytecodeResponse) Descriptor() ([]byte, []int) {
	return file_contract_get_bytecode_proto_rawDescGZIP(), []int{1}
}

func (x *ContractGetBytecodeResponse) GetHeader() *ResponseHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *ContractGetBytecodeResponse) GetBytecode() []byte {
	if x != nil {
		return x.Bytecode
	}
	return nil
}

var File_contract_get_bytecode_proto protoreflect.FileDescriptor

var file_contract_get_bytecode_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x67, 0x65, 0x74, 0x5f, 0x62,
	0x79, 0x74, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x79, 0x0a, 0x18, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x47, 0x65,
	0x74, 0x42, 0x79, 0x74, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x2a,
	0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x0a, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49,
	0x44, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x44, 0x22, 0x68, 0x0a,
	0x1b, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x47, 0x65, 0x74, 0x42, 0x79, 0x74, 0x65,
	0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x62,
	0x79, 0x74, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x62,
	0x79, 0x74, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68,
	0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contract_get_bytecode_proto_rawDescOnce sync.Once
	file_contract_get_bytecode_proto_rawDescData = file_contract_get_bytecode_proto_rawDesc
)

func file_contract_get_bytecode_proto_rawDescGZIP() []byte {
	file_contract_get_bytecode_proto_rawDescOnce.Do(func() {
		file_contract_get_bytecode_proto_rawDescData = protoimpl.X.CompressGZIP(file_contract_get_bytecode_proto_rawDescData)
	})
	return file_contract_get_bytecode_proto_rawDescData
}

var file_contract_get_bytecode_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_contract_get_bytecode_proto_goTypes = []any{
	(*ContractGetBytecodeQuery)(nil),    // 0: proto.ContractGetBytecodeQuery
	(*ContractGetBytecodeResponse)(nil), // 1: proto.ContractGetBytecodeResponse
	(*QueryHeader)(nil),                 // 2: proto.QueryHeader
	(*ContractID)(nil),                  // 3: proto.ContractID
	(*ResponseHeader)(nil),              // 4: proto.ResponseHeader
}
var file_contract_get_bytecode_proto_depIdxs = []int32{
	2, // 0: proto.ContractGetBytecodeQuery.header:type_name -> proto.QueryHeader
	3, // 1: proto.ContractGetBytecodeQuery.contractID:type_name -> proto.ContractID
	4, // 2: proto.ContractGetBytecodeResponse.header:type_name -> proto.ResponseHeader
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_contract_get_bytecode_proto_init() }
func file_contract_get_bytecode_proto_init() {
	if File_contract_get_bytecode_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_query_header_proto_init()
	file_response_header_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contract_get_bytecode_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_contract_get_bytecode_proto_goTypes,
		DependencyIndexes: file_contract_get_bytecode_proto_depIdxs,
		MessageInfos:      file_contract_get_bytecode_proto_msgTypes,
	}.Build()
	File_contract_get_bytecode_proto = out.File
	file_contract_get_bytecode_proto_rawDesc = nil
	file_contract_get_bytecode_proto_goTypes = nil
	file_contract_get_bytecode_proto_depIdxs = nil
}
