//*
// # Token Claim Airdrop
// Messages used to implement a transaction to claim a pending airdrop.
//
// ### Keywords
// The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
// "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
// document are to be interpreted as described in [RFC2119](https://www.ietf.org/rfc/rfc2119).

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: token_claim_airdrop.proto

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
// Token claim airdrop<br/>
// Complete one or more pending transfers on behalf of the
// recipient(s) for an airdrop.
//
// The sender MUST have sufficient balance to fulfill the airdrop at the
// time of claim. If the sender does not have sufficient balance, the
// claim SHALL fail.<br/>
// Each pending airdrop successfully claimed SHALL be removed from state and
// SHALL NOT be available to claim again.<br/>
// Each claim SHALL be represented in the transaction body and
// SHALL NOT be restated in the record file.<br/>
// All claims MUST succeed for this transaction to succeed.
//
// ### Record Stream Effects
// The completed transfers SHALL be present in the transfer list.
type TokenClaimAirdropTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A list of one or more pending airdrop identifiers.
	// <p>
	// This transaction MUST be signed by the account identified by
	// the `receiver_id` for each entry in this list.<br/>
	// This list MUST contain between 1 and 10 entries, inclusive.<br/>
	// This list MUST NOT have any duplicate entries.
	PendingAirdrops []*PendingAirdropId `protobuf:"bytes,1,rep,name=pending_airdrops,json=pendingAirdrops,proto3" json:"pending_airdrops,omitempty"`
}

func (x *TokenClaimAirdropTransactionBody) Reset() {
	*x = TokenClaimAirdropTransactionBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_claim_airdrop_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenClaimAirdropTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenClaimAirdropTransactionBody) ProtoMessage() {}

func (x *TokenClaimAirdropTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_token_claim_airdrop_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenClaimAirdropTransactionBody.ProtoReflect.Descriptor instead.
func (*TokenClaimAirdropTransactionBody) Descriptor() ([]byte, []int) {
	return file_token_claim_airdrop_proto_rawDescGZIP(), []int{0}
}

func (x *TokenClaimAirdropTransactionBody) GetPendingAirdrops() []*PendingAirdropId {
	if x != nil {
		return x.PendingAirdrops
	}
	return nil
}

var File_token_claim_airdrop_proto protoreflect.FileDescriptor

var file_token_claim_airdrop_proto_rawDesc = []byte{
	0x0a, 0x19, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x5f, 0x61, 0x69,
	0x72, 0x64, 0x72, 0x6f, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x66, 0x0a, 0x20, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x43, 0x6c,
	0x61, 0x69, 0x6d, 0x41, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x42, 0x0a, 0x10, 0x70, 0x65, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x5f, 0x61, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x65, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x41, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x49, 0x64, 0x52, 0x0f, 0x70, 0x65,
	0x6e, 0x64, 0x69, 0x6e, 0x67, 0x41, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x73, 0x42, 0x26, 0x0a,
	0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67,
	0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a,
	0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_token_claim_airdrop_proto_rawDescOnce sync.Once
	file_token_claim_airdrop_proto_rawDescData = file_token_claim_airdrop_proto_rawDesc
)

func file_token_claim_airdrop_proto_rawDescGZIP() []byte {
	file_token_claim_airdrop_proto_rawDescOnce.Do(func() {
		file_token_claim_airdrop_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_claim_airdrop_proto_rawDescData)
	})
	return file_token_claim_airdrop_proto_rawDescData
}

var file_token_claim_airdrop_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_token_claim_airdrop_proto_goTypes = []interface{}{
	(*TokenClaimAirdropTransactionBody)(nil), // 0: proto.TokenClaimAirdropTransactionBody
	(*PendingAirdropId)(nil),                 // 1: proto.PendingAirdropId
}
var file_token_claim_airdrop_proto_depIdxs = []int32{
	1, // 0: proto.TokenClaimAirdropTransactionBody.pending_airdrops:type_name -> proto.PendingAirdropId
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_token_claim_airdrop_proto_init() }
func file_token_claim_airdrop_proto_init() {
	if File_token_claim_airdrop_proto != nil {
		return
	}
	file_basic_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_token_claim_airdrop_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenClaimAirdropTransactionBody); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_token_claim_airdrop_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_token_claim_airdrop_proto_goTypes,
		DependencyIndexes: file_token_claim_airdrop_proto_depIdxs,
		MessageInfos:      file_token_claim_airdrop_proto_msgTypes,
	}.Build()
	File_token_claim_airdrop_proto = out.File
	file_token_claim_airdrop_proto_rawDesc = nil
	file_token_claim_airdrop_proto_goTypes = nil
	file_token_claim_airdrop_proto_depIdxs = nil
}
