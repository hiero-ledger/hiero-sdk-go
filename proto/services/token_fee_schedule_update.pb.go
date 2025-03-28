//*
// # Fee Schedule Update
// Transaction to update the fee schedule for a token. A token creator may
// wish to charge custom transaction fees for a token type, and if a
// `fee_schedule_key` is assigned, this transaction enables adding, removing,
// or updating those custom transaction fees.
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
// source: token_fee_schedule_update.proto

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
// Update the custom fee schedule for a token type.
//
// The token MUST have a `fee_schedule_key` set and that key MUST NOT
// be an empty `KeyList`.<br/>
// The token `fee_schedule_key` MUST sign this transaction.<br/>
// The token MUST exist, MUST NOT be deleted, and MUST NOT be expired.<br/>
//
// If the custom_fees list is empty, clears the fee schedule or resolves to
// CUSTOM_SCHEDULE_ALREADY_HAS_NO_FEES if the fee schedule was already empty.
//
// ### Block Stream Effects
// None
type TokenFeeScheduleUpdateTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A token identifier.
	// <p>
	// This SHALL identify the token type to modify with an updated
	// custom fee schedule.<br/>
	// The identified token MUST exist, and MUST NOT be deleted.
	TokenId *TokenID `protobuf:"bytes,1,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
	// *
	// A list of custom fees representing a fee schedule.
	// <p>
	// This list MAY be empty to remove custom fees from a token.<br/>
	// If the identified token is a non-fungible/unique type, the entries
	// in this list MUST NOT declare a `fractional_fee`.<br/>
	// If the identified token is a fungible/common type, the entries in this
	// list MUST NOT declare a `royalty_fee`.<br/>
	// Any token type MAY include entries that declare a `fixed_fee`.
	CustomFees []*CustomFee `protobuf:"bytes,2,rep,name=custom_fees,json=customFees,proto3" json:"custom_fees,omitempty"`
}

func (x *TokenFeeScheduleUpdateTransactionBody) Reset() {
	*x = TokenFeeScheduleUpdateTransactionBody{}
	mi := &file_token_fee_schedule_update_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenFeeScheduleUpdateTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenFeeScheduleUpdateTransactionBody) ProtoMessage() {}

func (x *TokenFeeScheduleUpdateTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_token_fee_schedule_update_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenFeeScheduleUpdateTransactionBody.ProtoReflect.Descriptor instead.
func (*TokenFeeScheduleUpdateTransactionBody) Descriptor() ([]byte, []int) {
	return file_token_fee_schedule_update_proto_rawDescGZIP(), []int{0}
}

func (x *TokenFeeScheduleUpdateTransactionBody) GetTokenId() *TokenID {
	if x != nil {
		return x.TokenId
	}
	return nil
}

func (x *TokenFeeScheduleUpdateTransactionBody) GetCustomFees() []*CustomFee {
	if x != nil {
		return x.CustomFees
	}
	return nil
}

var File_token_fee_schedule_update_proto protoreflect.FileDescriptor

var file_token_fee_schedule_update_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x66, 0x65, 0x65, 0x5f, 0x73, 0x63, 0x68, 0x65,
	0x64, 0x75, 0x6c, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x63, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x5f, 0x66, 0x65, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85,
	0x01, 0x0a, 0x25, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x46, 0x65, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x29, 0x0a, 0x08, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x44, 0x52, 0x07, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x66, 0x65,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x46, 0x65, 0x65, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74,
	0x6f, 0x6d, 0x46, 0x65, 0x65, 0x73, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65,
	0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_token_fee_schedule_update_proto_rawDescOnce sync.Once
	file_token_fee_schedule_update_proto_rawDescData = file_token_fee_schedule_update_proto_rawDesc
)

func file_token_fee_schedule_update_proto_rawDescGZIP() []byte {
	file_token_fee_schedule_update_proto_rawDescOnce.Do(func() {
		file_token_fee_schedule_update_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_fee_schedule_update_proto_rawDescData)
	})
	return file_token_fee_schedule_update_proto_rawDescData
}

var file_token_fee_schedule_update_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_token_fee_schedule_update_proto_goTypes = []any{
	(*TokenFeeScheduleUpdateTransactionBody)(nil), // 0: proto.TokenFeeScheduleUpdateTransactionBody
	(*TokenID)(nil),   // 1: proto.TokenID
	(*CustomFee)(nil), // 2: proto.CustomFee
}
var file_token_fee_schedule_update_proto_depIdxs = []int32{
	1, // 0: proto.TokenFeeScheduleUpdateTransactionBody.token_id:type_name -> proto.TokenID
	2, // 1: proto.TokenFeeScheduleUpdateTransactionBody.custom_fees:type_name -> proto.CustomFee
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_token_fee_schedule_update_proto_init() }
func file_token_fee_schedule_update_proto_init() {
	if File_token_fee_schedule_update_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_custom_fees_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_token_fee_schedule_update_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_token_fee_schedule_update_proto_goTypes,
		DependencyIndexes: file_token_fee_schedule_update_proto_depIdxs,
		MessageInfos:      file_token_fee_schedule_update_proto_msgTypes,
	}.Build()
	File_token_fee_schedule_update_proto = out.File
	file_token_fee_schedule_update_proto_rawDesc = nil
	file_token_fee_schedule_update_proto_goTypes = nil
	file_token_fee_schedule_update_proto_depIdxs = nil
}
