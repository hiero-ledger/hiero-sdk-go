//*
// # Token Update NFTs
// Given a token identifier and a metadata block, change the metadata for
// one or more non-fungible/unique token instances.
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
// source: token_update_nfts.proto

package services

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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
// Modify the metadata field for an individual non-fungible/unique token (NFT).
//
// Updating the metadata of an NFT SHALL NOT affect ownership or
// the ability to transfer that NFT.<br/>
// This transaction SHALL affect only the specific serial numbered tokens
// identified.
// This transaction SHALL modify individual token metadata.<br/>
// This transaction MUST be signed by the token `metadata_key`.<br/>
// The token `metadata_key` MUST be a valid `Key`.<br/>
// The token `metadata_key` MUST NOT be an empty `KeyList`.
//
// ### Block Stream Effects
// None
type TokenUpdateNftsTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A token identifier.<br/>
	// This is the token type (i.e. collection) for which to update NFTs.
	// <p>
	// This field is REQUIRED.<br/>
	// The identified token MUST exist, MUST NOT be paused, MUST have the type
	// non-fungible/unique, and MUST have a valid `metadata_key`.
	Token *TokenID `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	// *
	// A list of serial numbers to be updated.
	// <p>
	// This field is REQUIRED.<br/>
	// This list MUST have at least one(1) entry.<br/>
	// This list MUST NOT have more than ten(10) entries.
	SerialNumbers []int64 `protobuf:"varint,2,rep,packed,name=serial_numbers,json=serialNumbers,proto3" json:"serial_numbers,omitempty"`
	// *
	// A new value for the metadata.
	// <p>
	// If this field is not set, the metadata SHALL NOT change.<br/>
	// This value, if set, MUST NOT exceed 100 bytes.
	Metadata *wrapperspb.BytesValue `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *TokenUpdateNftsTransactionBody) Reset() {
	*x = TokenUpdateNftsTransactionBody{}
	mi := &file_token_update_nfts_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenUpdateNftsTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenUpdateNftsTransactionBody) ProtoMessage() {}

func (x *TokenUpdateNftsTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_token_update_nfts_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenUpdateNftsTransactionBody.ProtoReflect.Descriptor instead.
func (*TokenUpdateNftsTransactionBody) Descriptor() ([]byte, []int) {
	return file_token_update_nfts_proto_rawDescGZIP(), []int{0}
}

func (x *TokenUpdateNftsTransactionBody) GetToken() *TokenID {
	if x != nil {
		return x.Token
	}
	return nil
}

func (x *TokenUpdateNftsTransactionBody) GetSerialNumbers() []int64 {
	if x != nil {
		return x.SerialNumbers
	}
	return nil
}

func (x *TokenUpdateNftsTransactionBody) GetMetadata() *wrapperspb.BytesValue {
	if x != nil {
		return x.Metadata
	}
	return nil
}

var File_token_update_nfts_proto protoreflect.FileDescriptor

var file_token_update_nfts_proto_rawDesc = []byte{
	0x0a, 0x17, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x6e,
	0x66, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xa6, 0x01, 0x0a, 0x1e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4e, 0x66, 0x74, 0x73, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x24, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x49, 0x44, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x25, 0x0a, 0x0e,
	0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x03, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x73, 0x12, 0x37, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x42, 0x26, 0x0a, 0x22,
	0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61,
	0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_token_update_nfts_proto_rawDescOnce sync.Once
	file_token_update_nfts_proto_rawDescData = file_token_update_nfts_proto_rawDesc
)

func file_token_update_nfts_proto_rawDescGZIP() []byte {
	file_token_update_nfts_proto_rawDescOnce.Do(func() {
		file_token_update_nfts_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_update_nfts_proto_rawDescData)
	})
	return file_token_update_nfts_proto_rawDescData
}

var file_token_update_nfts_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_token_update_nfts_proto_goTypes = []any{
	(*TokenUpdateNftsTransactionBody)(nil), // 0: proto.TokenUpdateNftsTransactionBody
	(*TokenID)(nil),                        // 1: proto.TokenID
	(*wrapperspb.BytesValue)(nil),          // 2: google.protobuf.BytesValue
}
var file_token_update_nfts_proto_depIdxs = []int32{
	1, // 0: proto.TokenUpdateNftsTransactionBody.token:type_name -> proto.TokenID
	2, // 1: proto.TokenUpdateNftsTransactionBody.metadata:type_name -> google.protobuf.BytesValue
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_token_update_nfts_proto_init() }
func file_token_update_nfts_proto_init() {
	if File_token_update_nfts_proto != nil {
		return
	}
	file_basic_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_token_update_nfts_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_token_update_nfts_proto_goTypes,
		DependencyIndexes: file_token_update_nfts_proto_depIdxs,
		MessageInfos:      file_token_update_nfts_proto_msgTypes,
	}.Build()
	File_token_update_nfts_proto = out.File
	file_token_update_nfts_proto_rawDesc = nil
	file_token_update_nfts_proto_goTypes = nil
	file_token_update_nfts_proto_depIdxs = nil
}
