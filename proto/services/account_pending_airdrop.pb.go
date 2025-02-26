//*
// # Account Pending Airdrop.
// A single pending airdrop awaiting claim by the recipient.
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
// source: account_pending_airdrop.proto

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
// A node within a doubly linked list of pending airdrop references.<br/>
// This internal state message forms the entries in a doubly-linked list
// of references to pending airdrop entries that are "owed" by a particular
// account as "sender".
//
// Each entry in this list MUST refer to an existing pending airdrop.<br/>
// The pending airdrop MUST NOT be claimed.<br/>
// The pending airdrop MUST NOT be canceled.<br/>
// The pending airdrop `sender` account's `head_pending_airdrop_id` field
// MUST match the `pending_airdrop_id` field in this message.
type AccountPendingAirdrop struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// An amount of fungible tokens to be sent for this pending airdrop.
	// <p>
	// This field SHALL NOT be set for non-fungible/unique tokens.
	PendingAirdropValue *PendingAirdropValue `protobuf:"bytes,1,opt,name=pending_airdrop_value,json=pendingAirdropValue,proto3" json:"pending_airdrop_value,omitempty"`
	// *
	// A pending airdrop identifier.
	// <p>
	// This field SHALL identify the specific pending airdrop that
	// precedes this position within the doubly linked list of pending
	// airdrops "owed" by the sending account associated with this
	// account airdrop "list".<br/>
	// This SHALL match `pending_airdrop_id` if this is the only entry
	// in the "list".
	PreviousAirdrop *PendingAirdropId `protobuf:"bytes,2,opt,name=previous_airdrop,json=previousAirdrop,proto3" json:"previous_airdrop,omitempty"`
	// *
	// A pending airdrop identifier.<br/>
	// <p>
	// This field SHALL identify the specific pending airdrop that
	// follows this position within the doubly linked list of pending
	// airdrops "owed" by the sending account associated with this
	// account airdrop "list".<br/>
	// This SHALL match `pending_airdrop_id` if this is the only entry
	// in the "list".
	NextAirdrop *PendingAirdropId `protobuf:"bytes,3,opt,name=next_airdrop,json=nextAirdrop,proto3" json:"next_airdrop,omitempty"`
}

func (x *AccountPendingAirdrop) Reset() {
	*x = AccountPendingAirdrop{}
	mi := &file_account_pending_airdrop_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AccountPendingAirdrop) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountPendingAirdrop) ProtoMessage() {}

func (x *AccountPendingAirdrop) ProtoReflect() protoreflect.Message {
	mi := &file_account_pending_airdrop_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountPendingAirdrop.ProtoReflect.Descriptor instead.
func (*AccountPendingAirdrop) Descriptor() ([]byte, []int) {
	return file_account_pending_airdrop_proto_rawDescGZIP(), []int{0}
}

func (x *AccountPendingAirdrop) GetPendingAirdropValue() *PendingAirdropValue {
	if x != nil {
		return x.PendingAirdropValue
	}
	return nil
}

func (x *AccountPendingAirdrop) GetPreviousAirdrop() *PendingAirdropId {
	if x != nil {
		return x.PreviousAirdrop
	}
	return nil
}

func (x *AccountPendingAirdrop) GetNextAirdrop() *PendingAirdropId {
	if x != nil {
		return x.NextAirdrop
	}
	return nil
}

var File_account_pending_airdrop_proto protoreflect.FileDescriptor

var file_account_pending_airdrop_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x70, 0x65, 0x6e, 0x64, 0x69, 0x6e,
	0x67, 0x5f, 0x61, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe7, 0x01, 0x0a, 0x15, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x41, 0x69, 0x72, 0x64,
	0x72, 0x6f, 0x70, 0x12, 0x4e, 0x0a, 0x15, 0x70, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x61,
	0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x65, 0x6e, 0x64, 0x69,
	0x6e, 0x67, 0x41, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x13,
	0x70, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x41, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x42, 0x0a, 0x10, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f,
	0x61, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x41, 0x69, 0x72,
	0x64, 0x72, 0x6f, 0x70, 0x49, 0x64, 0x52, 0x0f, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73,
	0x41, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x12, 0x3a, 0x0a, 0x0c, 0x6e, 0x65, 0x78, 0x74, 0x5f,
	0x61, 0x69, 0x72, 0x64, 0x72, 0x6f, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x41, 0x69, 0x72,
	0x64, 0x72, 0x6f, 0x70, 0x49, 0x64, 0x52, 0x0b, 0x6e, 0x65, 0x78, 0x74, 0x41, 0x69, 0x72, 0x64,
	0x72, 0x6f, 0x70, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72,
	0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_account_pending_airdrop_proto_rawDescOnce sync.Once
	file_account_pending_airdrop_proto_rawDescData = file_account_pending_airdrop_proto_rawDesc
)

func file_account_pending_airdrop_proto_rawDescGZIP() []byte {
	file_account_pending_airdrop_proto_rawDescOnce.Do(func() {
		file_account_pending_airdrop_proto_rawDescData = protoimpl.X.CompressGZIP(file_account_pending_airdrop_proto_rawDescData)
	})
	return file_account_pending_airdrop_proto_rawDescData
}

var file_account_pending_airdrop_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_account_pending_airdrop_proto_goTypes = []any{
	(*AccountPendingAirdrop)(nil), // 0: proto.AccountPendingAirdrop
	(*PendingAirdropValue)(nil),   // 1: proto.PendingAirdropValue
	(*PendingAirdropId)(nil),      // 2: proto.PendingAirdropId
}
var file_account_pending_airdrop_proto_depIdxs = []int32{
	1, // 0: proto.AccountPendingAirdrop.pending_airdrop_value:type_name -> proto.PendingAirdropValue
	2, // 1: proto.AccountPendingAirdrop.previous_airdrop:type_name -> proto.PendingAirdropId
	2, // 2: proto.AccountPendingAirdrop.next_airdrop:type_name -> proto.PendingAirdropId
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_account_pending_airdrop_proto_init() }
func file_account_pending_airdrop_proto_init() {
	if File_account_pending_airdrop_proto != nil {
		return
	}
	file_basic_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_account_pending_airdrop_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_account_pending_airdrop_proto_goTypes,
		DependencyIndexes: file_account_pending_airdrop_proto_depIdxs,
		MessageInfos:      file_account_pending_airdrop_proto_msgTypes,
	}.Build()
	File_account_pending_airdrop_proto = out.File
	file_account_pending_airdrop_proto_rawDesc = nil
	file_account_pending_airdrop_proto_goTypes = nil
	file_account_pending_airdrop_proto_depIdxs = nil
}
