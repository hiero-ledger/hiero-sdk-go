//*
// # Token Reject
// Messages used to implement a transaction to reject a token type from an
// account.
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
// source: token_reject.proto

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
// Reject undesired token(s).<br/>
// Transfer one or more token balances held by the requesting account to the
// treasury for each token type.
//
// Each transfer SHALL be one of the following
// - A single non-fungible/unique token.
// - The full balance held for a fungible/common token.
// A single `tokenReject` transaction SHALL support a maximum
// of 10 transfers.<br/>
// A token that is `pause`d MUST NOT be rejected.<br/>
// If the `owner` account is `frozen` with respect to the identified token(s)
// the token(s) MUST NOT be rejected.<br/>
// The `payer` for this transaction, and `owner` if set, SHALL NOT be charged
// any custom fees or other fees beyond the `tokenReject` transaction fee.
//
// ### Block Stream Effects
//   - Each successful transfer from `payer` to `treasury` SHALL be recorded in
//     the `token_transfer_list` for the transaction record.
type TokenRejectTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// An account identifier.<br/>
	// This OPTIONAL field identifies the account holding the
	// tokens to be rejected.
	// <p>
	// If set, this account MUST sign this transaction.
	// If not set, the `payer` for this transaction SHALL be the effective
	// `owner` for this transaction.
	Owner *AccountID `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	// *
	// A list of one or more token rejections.
	// <p>
	// On success each rejected token serial number or balance SHALL be
	// transferred from the requesting account to the treasury account for
	// that token type.<br/>
	// After rejection the requesting account SHALL continue to be associated
	// with the token.<br/>
	// If dissociation is desired then a separate `TokenDissociate` transaction
	// MUST be submitted to remove the association.<br/>
	// This list MUST contain at least one (1) entry and MUST NOT contain more
	// than ten (10) entries.
	Rejections []*TokenReference `protobuf:"bytes,2,rep,name=rejections,proto3" json:"rejections,omitempty"`
}

func (x *TokenRejectTransactionBody) Reset() {
	*x = TokenRejectTransactionBody{}
	mi := &file_token_reject_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenRejectTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenRejectTransactionBody) ProtoMessage() {}

func (x *TokenRejectTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_token_reject_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenRejectTransactionBody.ProtoReflect.Descriptor instead.
func (*TokenRejectTransactionBody) Descriptor() ([]byte, []int) {
	return file_token_reject_proto_rawDescGZIP(), []int{0}
}

func (x *TokenRejectTransactionBody) GetOwner() *AccountID {
	if x != nil {
		return x.Owner
	}
	return nil
}

func (x *TokenRejectTransactionBody) GetRejections() []*TokenReference {
	if x != nil {
		return x.Rejections
	}
	return nil
}

// *
// A union token identifier.
//
// Identify a fungible/common token type, or a single
// non-fungible/unique token serial.
type TokenReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to TokenIdentifier:
	//
	//	*TokenReference_FungibleToken
	//	*TokenReference_Nft
	TokenIdentifier isTokenReference_TokenIdentifier `protobuf_oneof:"token_identifier"`
}

func (x *TokenReference) Reset() {
	*x = TokenReference{}
	mi := &file_token_reject_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenReference) ProtoMessage() {}

func (x *TokenReference) ProtoReflect() protoreflect.Message {
	mi := &file_token_reject_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenReference.ProtoReflect.Descriptor instead.
func (*TokenReference) Descriptor() ([]byte, []int) {
	return file_token_reject_proto_rawDescGZIP(), []int{1}
}

func (m *TokenReference) GetTokenIdentifier() isTokenReference_TokenIdentifier {
	if m != nil {
		return m.TokenIdentifier
	}
	return nil
}

func (x *TokenReference) GetFungibleToken() *TokenID {
	if x, ok := x.GetTokenIdentifier().(*TokenReference_FungibleToken); ok {
		return x.FungibleToken
	}
	return nil
}

func (x *TokenReference) GetNft() *NftID {
	if x, ok := x.GetTokenIdentifier().(*TokenReference_Nft); ok {
		return x.Nft
	}
	return nil
}

type isTokenReference_TokenIdentifier interface {
	isTokenReference_TokenIdentifier()
}

type TokenReference_FungibleToken struct {
	// *
	// A fungible/common token type.
	FungibleToken *TokenID `protobuf:"bytes,1,opt,name=fungible_token,json=fungibleToken,proto3,oneof"`
}

type TokenReference_Nft struct {
	// *
	// A single specific serialized non-fungible/unique token.
	Nft *NftID `protobuf:"bytes,2,opt,name=nft,proto3,oneof"`
}

func (*TokenReference_FungibleToken) isTokenReference_TokenIdentifier() {}

func (*TokenReference_Nft) isTokenReference_TokenIdentifier() {}

var File_token_reject_proto protoreflect.FileDescriptor

var file_token_reject_proto_rawDesc = []byte{
	0x0a, 0x12, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73,
	0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7b,
	0x0a, 0x1a, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x26, 0x0a, 0x05,
	0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x05, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x0a, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52,
	0x0a, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x7f, 0x0a, 0x0e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x37, 0x0a,
	0x0e, 0x66, 0x75, 0x6e, 0x67, 0x69, 0x62, 0x6c, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x49, 0x44, 0x48, 0x00, 0x52, 0x0d, 0x66, 0x75, 0x6e, 0x67, 0x69, 0x62, 0x6c,
	0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x20, 0x0a, 0x03, 0x6e, 0x66, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x66, 0x74, 0x49,
	0x44, 0x48, 0x00, 0x52, 0x03, 0x6e, 0x66, 0x74, 0x42, 0x12, 0x0a, 0x10, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x42, 0x26, 0x0a, 0x22,
	0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61,
	0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_token_reject_proto_rawDescOnce sync.Once
	file_token_reject_proto_rawDescData = file_token_reject_proto_rawDesc
)

func file_token_reject_proto_rawDescGZIP() []byte {
	file_token_reject_proto_rawDescOnce.Do(func() {
		file_token_reject_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_reject_proto_rawDescData)
	})
	return file_token_reject_proto_rawDescData
}

var file_token_reject_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_token_reject_proto_goTypes = []any{
	(*TokenRejectTransactionBody)(nil), // 0: proto.TokenRejectTransactionBody
	(*TokenReference)(nil),             // 1: proto.TokenReference
	(*AccountID)(nil),                  // 2: proto.AccountID
	(*TokenID)(nil),                    // 3: proto.TokenID
	(*NftID)(nil),                      // 4: proto.NftID
}
var file_token_reject_proto_depIdxs = []int32{
	2, // 0: proto.TokenRejectTransactionBody.owner:type_name -> proto.AccountID
	1, // 1: proto.TokenRejectTransactionBody.rejections:type_name -> proto.TokenReference
	3, // 2: proto.TokenReference.fungible_token:type_name -> proto.TokenID
	4, // 3: proto.TokenReference.nft:type_name -> proto.NftID
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_token_reject_proto_init() }
func file_token_reject_proto_init() {
	if File_token_reject_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_token_reject_proto_msgTypes[1].OneofWrappers = []any{
		(*TokenReference_FungibleToken)(nil),
		(*TokenReference_Nft)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_token_reject_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_token_reject_proto_goTypes,
		DependencyIndexes: file_token_reject_proto_depIdxs,
		MessageInfos:      file_token_reject_proto_msgTypes,
	}.Build()
	File_token_reject_proto = out.File
	file_token_reject_proto_rawDesc = nil
	file_token_reject_proto_goTypes = nil
	file_token_reject_proto_depIdxs = nil
}
