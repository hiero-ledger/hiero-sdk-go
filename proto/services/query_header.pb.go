//*
// # Query Header
// Messages that comprise a header sent with each query request.
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
// source: query_header.proto

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
// The type of query response.<br/>
//
// This SHALL be answer-only as a default.<br/>
// This value SHALL support an "estimated cost" type.<br/>
// This value SHOULD support a "state proof" type, when available.
type ResponseType int32

const (
	// *
	// A response with the query answer.
	ResponseType_ANSWER_ONLY ResponseType = 0
	// *
	// A response with both the query answer and a state proof.
	ResponseType_ANSWER_STATE_PROOF ResponseType = 1
	// *
	// A response with the estimated cost to answer the query.
	ResponseType_COST_ANSWER ResponseType = 2
	// *
	// A response with the estimated cost to answer and a state proof.
	ResponseType_COST_ANSWER_STATE_PROOF ResponseType = 3
)

// Enum value maps for ResponseType.
var (
	ResponseType_name = map[int32]string{
		0: "ANSWER_ONLY",
		1: "ANSWER_STATE_PROOF",
		2: "COST_ANSWER",
		3: "COST_ANSWER_STATE_PROOF",
	}
	ResponseType_value = map[string]int32{
		"ANSWER_ONLY":             0,
		"ANSWER_STATE_PROOF":      1,
		"COST_ANSWER":             2,
		"COST_ANSWER_STATE_PROOF": 3,
	}
)

func (x ResponseType) Enum() *ResponseType {
	p := new(ResponseType)
	*p = x
	return p
}

func (x ResponseType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseType) Descriptor() protoreflect.EnumDescriptor {
	return file_query_header_proto_enumTypes[0].Descriptor()
}

func (ResponseType) Type() protoreflect.EnumType {
	return &file_query_header_proto_enumTypes[0]
}

func (x ResponseType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseType.Descriptor instead.
func (ResponseType) EnumDescriptor() ([]byte, []int) {
	return file_query_header_proto_rawDescGZIP(), []int{0}
}

// *
// A standard query header.<br/>
// Each query from the client to the node must contain a QueryHeader, which
// specifies the desired response type, and includes a payment transaction
// that will compensate the network for responding to the query.
// The payment may be blank if the query is free.
//
// The payment transaction MUST be a `cryptoTransfer` from the payer account
// to the account of the node where the query is submitted.<br/>
// If the payment is sufficient, the network SHALL respond with the response
// type requested.<br/>
// If the response type is `COST_ANSWER` the payment MUST be unset.
// A state proof SHALL be available for some types of information.<br/>
// A state proof SHALL be available for a Record, but not a receipt, and the
// response entry for each supported "get info" query.
type QueryHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A signed `CryptoTransferTransaction` to pay query fees.
	// <p>
	// This MUST transfer HBAR from the "payer" to the responding node account
	// sufficient to pay the query fees.
	Payment *Transaction `protobuf:"bytes,1,opt,name=payment,proto3" json:"payment,omitempty"`
	// *
	// A type of query response requested.
	ResponseType ResponseType `protobuf:"varint,2,opt,name=responseType,proto3,enum=proto.ResponseType" json:"responseType,omitempty"`
}

func (x *QueryHeader) Reset() {
	*x = QueryHeader{}
	mi := &file_query_header_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryHeader) ProtoMessage() {}

func (x *QueryHeader) ProtoReflect() protoreflect.Message {
	mi := &file_query_header_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryHeader.ProtoReflect.Descriptor instead.
func (*QueryHeader) Descriptor() ([]byte, []int) {
	return file_query_header_proto_rawDescGZIP(), []int{0}
}

func (x *QueryHeader) GetPayment() *Transaction {
	if x != nil {
		return x.Payment
	}
	return nil
}

func (x *QueryHeader) GetResponseType() ResponseType {
	if x != nil {
		return x.ResponseType
	}
	return ResponseType_ANSWER_ONLY
}

var File_query_header_proto protoreflect.FileDescriptor

var file_query_header_proto_rawDesc = []byte{
	0x0a, 0x12, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x74,
	0x0a, 0x0b, 0x51, 0x75, 0x65, 0x72, 0x79, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x2c, 0x0a,
	0x07, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x37, 0x0a, 0x0c, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x2a, 0x65, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x41, 0x4e, 0x53, 0x57, 0x45, 0x52, 0x5f, 0x4f,
	0x4e, 0x4c, 0x59, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x41, 0x4e, 0x53, 0x57, 0x45, 0x52, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x4f, 0x46, 0x10, 0x01, 0x12, 0x0f, 0x0a,
	0x0b, 0x43, 0x4f, 0x53, 0x54, 0x5f, 0x41, 0x4e, 0x53, 0x57, 0x45, 0x52, 0x10, 0x02, 0x12, 0x1b,
	0x0a, 0x17, 0x43, 0x4f, 0x53, 0x54, 0x5f, 0x41, 0x4e, 0x53, 0x57, 0x45, 0x52, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x4f, 0x46, 0x10, 0x03, 0x42, 0x26, 0x0a, 0x22, 0x63,
	0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61,
	0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76,
	0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_query_header_proto_rawDescOnce sync.Once
	file_query_header_proto_rawDescData = file_query_header_proto_rawDesc
)

func file_query_header_proto_rawDescGZIP() []byte {
	file_query_header_proto_rawDescOnce.Do(func() {
		file_query_header_proto_rawDescData = protoimpl.X.CompressGZIP(file_query_header_proto_rawDescData)
	})
	return file_query_header_proto_rawDescData
}

var file_query_header_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_query_header_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_query_header_proto_goTypes = []any{
	(ResponseType)(0),   // 0: proto.ResponseType
	(*QueryHeader)(nil), // 1: proto.QueryHeader
	(*Transaction)(nil), // 2: proto.Transaction
}
var file_query_header_proto_depIdxs = []int32{
	2, // 0: proto.QueryHeader.payment:type_name -> proto.Transaction
	0, // 1: proto.QueryHeader.responseType:type_name -> proto.ResponseType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_query_header_proto_init() }
func file_query_header_proto_init() {
	if File_query_header_proto != nil {
		return
	}
	file_transaction_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_query_header_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_query_header_proto_goTypes,
		DependencyIndexes: file_query_header_proto_depIdxs,
		EnumInfos:         file_query_header_proto_enumTypes,
		MessageInfos:      file_query_header_proto_msgTypes,
	}.Build()
	File_query_header_proto = out.File
	file_query_header_proto_rawDesc = nil
	file_query_header_proto_goTypes = nil
	file_query_header_proto_depIdxs = nil
}
