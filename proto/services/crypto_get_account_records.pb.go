//*
// # Crypto Get Account Records
// Messages for a query to retrieve recent transaction records involving a
// specified account as effective `payer`.<br/>
// A "recent" transaction is typically one that reached consensus within
// the previous three(`3`) minutes of _consensus_ time. Additionally, the
// network only stores records in state when
// `ledger.keepRecordsInState=true` was true during transaction handling.
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
// source: crypto_get_account_records.proto

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
// Request records of all "recent" transactions for which the specified
// account is the effective payer.
type CryptoGetAccountRecordsQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// Standard information sent with every query operation.<br/>
	// This includes the signed payment and what kind of response is requested
	// (cost, state proof, both, or neither).
	Header *QueryHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// An account identifier.<br/>
	// This identifies the account to use when filtering the
	// transaction record lists.
	// <p>
	// This field is REQUIRED.
	AccountID *AccountID `protobuf:"bytes,2,opt,name=accountID,proto3" json:"accountID,omitempty"`
}

func (x *CryptoGetAccountRecordsQuery) Reset() {
	*x = CryptoGetAccountRecordsQuery{}
	mi := &file_crypto_get_account_records_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CryptoGetAccountRecordsQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoGetAccountRecordsQuery) ProtoMessage() {}

func (x *CryptoGetAccountRecordsQuery) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_get_account_records_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoGetAccountRecordsQuery.ProtoReflect.Descriptor instead.
func (*CryptoGetAccountRecordsQuery) Descriptor() ([]byte, []int) {
	return file_crypto_get_account_records_proto_rawDescGZIP(), []int{0}
}

func (x *CryptoGetAccountRecordsQuery) GetHeader() *QueryHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *CryptoGetAccountRecordsQuery) GetAccountID() *AccountID {
	if x != nil {
		return x.AccountID
	}
	return nil
}

// *
// Return records of all "recent" transactions for which the specified
// account is the effective payer.
type CryptoGetAccountRecordsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The standard response information for queries.<br/>
	// This includes the values requested in the `QueryHeader`
	// (cost, state proof, both, or neither).
	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// An account identifier.<br/>
	// This identifies the account used when filtering the
	// transaction record lists.
	// <p>
	// This field SHALL match the requested account identifier.
	AccountID *AccountID `protobuf:"bytes,2,opt,name=accountID,proto3" json:"accountID,omitempty"`
	// *
	// A list of records.
	// <p>
	// This list SHALL contain all available and "recent" records in which
	// the account identified in the `accountID` field acted as effective payer.
	Records []*TransactionRecord `protobuf:"bytes,3,rep,name=records,proto3" json:"records,omitempty"`
}

func (x *CryptoGetAccountRecordsResponse) Reset() {
	*x = CryptoGetAccountRecordsResponse{}
	mi := &file_crypto_get_account_records_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CryptoGetAccountRecordsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoGetAccountRecordsResponse) ProtoMessage() {}

func (x *CryptoGetAccountRecordsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_get_account_records_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoGetAccountRecordsResponse.ProtoReflect.Descriptor instead.
func (*CryptoGetAccountRecordsResponse) Descriptor() ([]byte, []int) {
	return file_crypto_get_account_records_proto_rawDescGZIP(), []int{1}
}

func (x *CryptoGetAccountRecordsResponse) GetHeader() *ResponseHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *CryptoGetAccountRecordsResponse) GetAccountID() *AccountID {
	if x != nil {
		return x.AccountID
	}
	return nil
}

func (x *CryptoGetAccountRecordsResponse) GetRecords() []*TransactionRecord {
	if x != nil {
		return x.Records
	}
	return nil
}

var File_crypto_get_account_records_proto protoreflect.FileDescriptor

var file_crypto_get_account_records_proto_rawDesc = []byte{
	0x0a, 0x20, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x5f, 0x67, 0x65, 0x74, 0x5f, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x7a, 0x0a, 0x1c, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x47, 0x65, 0x74, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x12, 0x2a, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x2e, 0x0a,
	0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x49, 0x44, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x22, 0xb4, 0x01,
	0x0a, 0x1f, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2d, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x2e, 0x0a, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44,
	0x12, 0x32, 0x0a, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x73, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65,
	0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_crypto_get_account_records_proto_rawDescOnce sync.Once
	file_crypto_get_account_records_proto_rawDescData = file_crypto_get_account_records_proto_rawDesc
)

func file_crypto_get_account_records_proto_rawDescGZIP() []byte {
	file_crypto_get_account_records_proto_rawDescOnce.Do(func() {
		file_crypto_get_account_records_proto_rawDescData = protoimpl.X.CompressGZIP(file_crypto_get_account_records_proto_rawDescData)
	})
	return file_crypto_get_account_records_proto_rawDescData
}

var file_crypto_get_account_records_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_crypto_get_account_records_proto_goTypes = []any{
	(*CryptoGetAccountRecordsQuery)(nil),    // 0: proto.CryptoGetAccountRecordsQuery
	(*CryptoGetAccountRecordsResponse)(nil), // 1: proto.CryptoGetAccountRecordsResponse
	(*QueryHeader)(nil),                     // 2: proto.QueryHeader
	(*AccountID)(nil),                       // 3: proto.AccountID
	(*ResponseHeader)(nil),                  // 4: proto.ResponseHeader
	(*TransactionRecord)(nil),               // 5: proto.TransactionRecord
}
var file_crypto_get_account_records_proto_depIdxs = []int32{
	2, // 0: proto.CryptoGetAccountRecordsQuery.header:type_name -> proto.QueryHeader
	3, // 1: proto.CryptoGetAccountRecordsQuery.accountID:type_name -> proto.AccountID
	4, // 2: proto.CryptoGetAccountRecordsResponse.header:type_name -> proto.ResponseHeader
	3, // 3: proto.CryptoGetAccountRecordsResponse.accountID:type_name -> proto.AccountID
	5, // 4: proto.CryptoGetAccountRecordsResponse.records:type_name -> proto.TransactionRecord
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_crypto_get_account_records_proto_init() }
func file_crypto_get_account_records_proto_init() {
	if File_crypto_get_account_records_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_transaction_record_proto_init()
	file_query_header_proto_init()
	file_response_header_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_crypto_get_account_records_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_crypto_get_account_records_proto_goTypes,
		DependencyIndexes: file_crypto_get_account_records_proto_depIdxs,
		MessageInfos:      file_crypto_get_account_records_proto_msgTypes,
	}.Build()
	File_crypto_get_account_records_proto = out.File
	file_crypto_get_account_records_proto_rawDesc = nil
	file_crypto_get_account_records_proto_goTypes = nil
	file_crypto_get_account_records_proto_depIdxs = nil
}
