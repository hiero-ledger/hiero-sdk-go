//*
// # Get Stakers
// Query all of the accounts proxy staking _to_ a specified account.
//
// > Important
// >> This query is obsolete and not supported.<br/>
// >> Any query of this type that is submitted SHALL fail with a `PRE_CHECK`
// >> result of `NOT_SUPPORTED`.
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
// source: crypto_get_stakers.proto

package services

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// *
// Get all the accounts that are proxy staking to this account. For each of
// them, give the amount currently staked. This was never implemented.
//
// Deprecated: Marked as deprecated in crypto_get_stakers.proto.
type CryptoGetStakersQuery struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// *
	// Standard information sent with every query operation.<br/>
	// This includes the signed payment and what kind of response is requested
	// (cost, state proof, both, or neither).
	Header *QueryHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// The Account ID for which the records should be retrieved
	AccountID     *AccountID `protobuf:"bytes,2,opt,name=accountID,proto3" json:"accountID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CryptoGetStakersQuery) Reset() {
	*x = CryptoGetStakersQuery{}
	mi := &file_crypto_get_stakers_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CryptoGetStakersQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoGetStakersQuery) ProtoMessage() {}

func (x *CryptoGetStakersQuery) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_get_stakers_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoGetStakersQuery.ProtoReflect.Descriptor instead.
func (*CryptoGetStakersQuery) Descriptor() ([]byte, []int) {
	return file_crypto_get_stakers_proto_rawDescGZIP(), []int{0}
}

func (x *CryptoGetStakersQuery) GetHeader() *QueryHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *CryptoGetStakersQuery) GetAccountID() *AccountID {
	if x != nil {
		return x.AccountID
	}
	return nil
}

// *
// information about a single account that is proxy staking
//
// Deprecated: Marked as deprecated in crypto_get_stakers.proto.
type ProxyStaker struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// *
	// The Account ID that is proxy staking
	AccountID *AccountID `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	// *
	// The number of hbars that are currently proxy staked
	Amount        int64 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProxyStaker) Reset() {
	*x = ProxyStaker{}
	mi := &file_crypto_get_stakers_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProxyStaker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProxyStaker) ProtoMessage() {}

func (x *ProxyStaker) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_get_stakers_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProxyStaker.ProtoReflect.Descriptor instead.
func (*ProxyStaker) Descriptor() ([]byte, []int) {
	return file_crypto_get_stakers_proto_rawDescGZIP(), []int{1}
}

func (x *ProxyStaker) GetAccountID() *AccountID {
	if x != nil {
		return x.AccountID
	}
	return nil
}

func (x *ProxyStaker) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

// *
// All of the accounts proxy staking to a given account, and the amounts proxy
// staked
//
// Deprecated: Marked as deprecated in crypto_get_stakers.proto.
type AllProxyStakers struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// *
	// The Account ID that is being proxy staked to
	AccountID *AccountID `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	// *
	// Each of the proxy staking accounts, and the amount they are proxy staking
	ProxyStaker   []*ProxyStaker `protobuf:"bytes,2,rep,name=proxyStaker,proto3" json:"proxyStaker,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AllProxyStakers) Reset() {
	*x = AllProxyStakers{}
	mi := &file_crypto_get_stakers_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AllProxyStakers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllProxyStakers) ProtoMessage() {}

func (x *AllProxyStakers) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_get_stakers_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllProxyStakers.ProtoReflect.Descriptor instead.
func (*AllProxyStakers) Descriptor() ([]byte, []int) {
	return file_crypto_get_stakers_proto_rawDescGZIP(), []int{2}
}

func (x *AllProxyStakers) GetAccountID() *AccountID {
	if x != nil {
		return x.AccountID
	}
	return nil
}

func (x *AllProxyStakers) GetProxyStaker() []*ProxyStaker {
	if x != nil {
		return x.ProxyStaker
	}
	return nil
}

// *
// Response when the client sends the node CryptoGetStakersQuery
//
// Deprecated: Marked as deprecated in crypto_get_stakers.proto.
type CryptoGetStakersResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// *
	// The standard response information for queries.<br/>
	// This includes the values requested in the `QueryHeader`
	// (cost, state proof, both, or neither).
	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// List of accounts proxy staking to this account, and the amount each is
	// currently proxy staking
	Stakers       *AllProxyStakers `protobuf:"bytes,3,opt,name=stakers,proto3" json:"stakers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CryptoGetStakersResponse) Reset() {
	*x = CryptoGetStakersResponse{}
	mi := &file_crypto_get_stakers_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CryptoGetStakersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoGetStakersResponse) ProtoMessage() {}

func (x *CryptoGetStakersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_get_stakers_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoGetStakersResponse.ProtoReflect.Descriptor instead.
func (*CryptoGetStakersResponse) Descriptor() ([]byte, []int) {
	return file_crypto_get_stakers_proto_rawDescGZIP(), []int{3}
}

func (x *CryptoGetStakersResponse) GetHeader() *ResponseHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *CryptoGetStakersResponse) GetStakers() *AllProxyStakers {
	if x != nil {
		return x.Stakers
	}
	return nil
}

var File_crypto_get_stakers_proto protoreflect.FileDescriptor

const file_crypto_get_stakers_proto_rawDesc = "" +
	"\n" +
	"\x18crypto_get_stakers.proto\x12\x05proto\x1a\x11basic_types.proto\x1a\x12query_header.proto\x1a\x15response_header.proto\"w\n" +
	"\x15CryptoGetStakersQuery\x12*\n" +
	"\x06header\x18\x01 \x01(\v2\x12.proto.QueryHeaderR\x06header\x12.\n" +
	"\taccountID\x18\x02 \x01(\v2\x10.proto.AccountIDR\taccountID:\x02\x18\x01\"Y\n" +
	"\vProxyStaker\x12.\n" +
	"\taccountID\x18\x01 \x01(\v2\x10.proto.AccountIDR\taccountID\x12\x16\n" +
	"\x06amount\x18\x02 \x01(\x03R\x06amount:\x02\x18\x01\"{\n" +
	"\x0fAllProxyStakers\x12.\n" +
	"\taccountID\x18\x01 \x01(\v2\x10.proto.AccountIDR\taccountID\x124\n" +
	"\vproxyStaker\x18\x02 \x03(\v2\x12.proto.ProxyStakerR\vproxyStaker:\x02\x18\x01\"\x7f\n" +
	"\x18CryptoGetStakersResponse\x12-\n" +
	"\x06header\x18\x01 \x01(\v2\x15.proto.ResponseHeaderR\x06header\x120\n" +
	"\astakers\x18\x03 \x01(\v2\x16.proto.AllProxyStakersR\astakers:\x02\x18\x01B&\n" +
	"\"com.hederahashgraph.api.proto.javaP\x01b\x06proto3"

var (
	file_crypto_get_stakers_proto_rawDescOnce sync.Once
	file_crypto_get_stakers_proto_rawDescData []byte
)

func file_crypto_get_stakers_proto_rawDescGZIP() []byte {
	file_crypto_get_stakers_proto_rawDescOnce.Do(func() {
		file_crypto_get_stakers_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_crypto_get_stakers_proto_rawDesc), len(file_crypto_get_stakers_proto_rawDesc)))
	})
	return file_crypto_get_stakers_proto_rawDescData
}

var file_crypto_get_stakers_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_crypto_get_stakers_proto_goTypes = []any{
	(*CryptoGetStakersQuery)(nil),    // 0: proto.CryptoGetStakersQuery
	(*ProxyStaker)(nil),              // 1: proto.ProxyStaker
	(*AllProxyStakers)(nil),          // 2: proto.AllProxyStakers
	(*CryptoGetStakersResponse)(nil), // 3: proto.CryptoGetStakersResponse
	(*QueryHeader)(nil),              // 4: proto.QueryHeader
	(*AccountID)(nil),                // 5: proto.AccountID
	(*ResponseHeader)(nil),           // 6: proto.ResponseHeader
}
var file_crypto_get_stakers_proto_depIdxs = []int32{
	4, // 0: proto.CryptoGetStakersQuery.header:type_name -> proto.QueryHeader
	5, // 1: proto.CryptoGetStakersQuery.accountID:type_name -> proto.AccountID
	5, // 2: proto.ProxyStaker.accountID:type_name -> proto.AccountID
	5, // 3: proto.AllProxyStakers.accountID:type_name -> proto.AccountID
	1, // 4: proto.AllProxyStakers.proxyStaker:type_name -> proto.ProxyStaker
	6, // 5: proto.CryptoGetStakersResponse.header:type_name -> proto.ResponseHeader
	2, // 6: proto.CryptoGetStakersResponse.stakers:type_name -> proto.AllProxyStakers
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_crypto_get_stakers_proto_init() }
func file_crypto_get_stakers_proto_init() {
	if File_crypto_get_stakers_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_query_header_proto_init()
	file_response_header_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_crypto_get_stakers_proto_rawDesc), len(file_crypto_get_stakers_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_crypto_get_stakers_proto_goTypes,
		DependencyIndexes: file_crypto_get_stakers_proto_depIdxs,
		MessageInfos:      file_crypto_get_stakers_proto_msgTypes,
	}.Build()
	File_crypto_get_stakers_proto = out.File
	file_crypto_get_stakers_proto_goTypes = nil
	file_crypto_get_stakers_proto_depIdxs = nil
}
