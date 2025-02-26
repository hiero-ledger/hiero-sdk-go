//*
// # Tss Message Transaction
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
// source: tss_message.proto

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

// * A transaction body to to send a Threshold Signature Scheme (TSS)
// Message.<br/>
// This is a wrapper around several different TSS message types that a node
// might communicate with other nodes in the network.
//
//   - A `TssMessageTransactionBody` MUST identify the hash of the roster
//     containing the node generating this TssMessage
//   - A `TssMessageTransactionBody` MUST identify the hash of the roster that
//     the TSS messages is for
//   - A `TssMessageTransactionBody` SHALL contain the specificc TssMessage data
//     that has been generated by the node for the share_index.
type TssMessageTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A hash of the roster containing the node generating the TssMessage.<br/>
	// This hash uniquely identifies the source roster, which will include
	// an entry for the node generating this TssMessage.
	// <p>
	// This value MUST be set.<br/>
	// This value MUST NOT be empty.<br/>
	// This value MUST contain a valid hash.
	SourceRosterHash []byte `protobuf:"bytes,1,opt,name=source_roster_hash,json=sourceRosterHash,proto3" json:"source_roster_hash,omitempty"`
	// *
	// A hash of the roster that the TssMessage is for.
	// <p>
	// This value MUST be set.<br/>
	// This value MUST NOT be empty.<br/>
	// This value MUST contain a valid hash.
	TargetRosterHash []byte `protobuf:"bytes,2,opt,name=target_roster_hash,json=targetRosterHash,proto3" json:"target_roster_hash,omitempty"`
	// *
	// An index to order shares.
	// <p>
	// A share index SHALL establish a global ordering of shares across all
	// shares in the network.<br/>
	// A share index MUST correspond to the index of the public share in the list
	// returned from the TSS library when the share was created for the source
	// roster.
	ShareIndex uint64 `protobuf:"varint,3,opt,name=share_index,json=shareIndex,proto3" json:"share_index,omitempty"`
	// *
	// A byte array.
	// <p>
	// This field SHALL contain the TssMessage data generated by the node
	// for the specified `share_index`.
	TssMessage []byte `protobuf:"bytes,4,opt,name=tss_message,json=tssMessage,proto3" json:"tss_message,omitempty"`
}

func (x *TssMessageTransactionBody) Reset() {
	*x = TssMessageTransactionBody{}
	mi := &file_tss_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TssMessageTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TssMessageTransactionBody) ProtoMessage() {}

func (x *TssMessageTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_tss_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TssMessageTransactionBody.ProtoReflect.Descriptor instead.
func (*TssMessageTransactionBody) Descriptor() ([]byte, []int) {
	return file_tss_message_proto_rawDescGZIP(), []int{0}
}

func (x *TssMessageTransactionBody) GetSourceRosterHash() []byte {
	if x != nil {
		return x.SourceRosterHash
	}
	return nil
}

func (x *TssMessageTransactionBody) GetTargetRosterHash() []byte {
	if x != nil {
		return x.TargetRosterHash
	}
	return nil
}

func (x *TssMessageTransactionBody) GetShareIndex() uint64 {
	if x != nil {
		return x.ShareIndex
	}
	return 0
}

func (x *TssMessageTransactionBody) GetTssMessage() []byte {
	if x != nil {
		return x.TssMessage
	}
	return nil
}

var File_tss_message_proto protoreflect.FileDescriptor

var file_tss_message_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x73, 0x73, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x26, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e,
	0x68, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x75,
	0x78, 0x69, 0x6c, 0x69, 0x61, 0x72, 0x79, 0x2e, 0x74, 0x73, 0x73, 0x22, 0xb9, 0x01, 0x0a, 0x19,
	0x54, 0x73, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x2c, 0x0a, 0x12, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x6f, 0x73,
	0x74, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x12, 0x2c, 0x0a, 0x12, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x5f, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x10, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x6f, 0x73, 0x74, 0x65,
	0x72, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x73, 0x73, 0x5f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x74, 0x73, 0x73,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x31, 0x0a, 0x2d, 0x63, 0x6f, 0x6d, 0x2e, 0x68,
	0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x61, 0x75, 0x78, 0x69, 0x6c, 0x69, 0x61, 0x72, 0x79, 0x2e, 0x74, 0x73,
	0x73, 0x2e, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_tss_message_proto_rawDescOnce sync.Once
	file_tss_message_proto_rawDescData = file_tss_message_proto_rawDesc
)

func file_tss_message_proto_rawDescGZIP() []byte {
	file_tss_message_proto_rawDescOnce.Do(func() {
		file_tss_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_tss_message_proto_rawDescData)
	})
	return file_tss_message_proto_rawDescData
}

var file_tss_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_tss_message_proto_goTypes = []any{
	(*TssMessageTransactionBody)(nil), // 0: com.hedera.hapi.services.auxiliary.tss.TssMessageTransactionBody
}
var file_tss_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tss_message_proto_init() }
func file_tss_message_proto_init() {
	if File_tss_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tss_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tss_message_proto_goTypes,
		DependencyIndexes: file_tss_message_proto_depIdxs,
		MessageInfos:      file_tss_message_proto_msgTypes,
	}.Build()
	File_tss_message_proto = out.File
	file_tss_message_proto_rawDesc = nil
	file_tss_message_proto_goTypes = nil
	file_tss_message_proto_depIdxs = nil
}
