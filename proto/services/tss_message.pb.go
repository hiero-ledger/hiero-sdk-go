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
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: tss_message.proto

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
	state protoimpl.MessageState `protogen:"open.v1"`
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
	TssMessage    []byte `protobuf:"bytes,4,opt,name=tss_message,json=tssMessage,proto3" json:"tss_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
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

const file_tss_message_proto_rawDesc = "" +
	"\n" +
	"\x11tss_message.proto\x12&com.hedera.hapi.services.auxiliary.tss\"\xb9\x01\n" +
	"\x19TssMessageTransactionBody\x12,\n" +
	"\x12source_roster_hash\x18\x01 \x01(\fR\x10sourceRosterHash\x12,\n" +
	"\x12target_roster_hash\x18\x02 \x01(\fR\x10targetRosterHash\x12\x1f\n" +
	"\vshare_index\x18\x03 \x01(\x04R\n" +
	"shareIndex\x12\x1f\n" +
	"\vtss_message\x18\x04 \x01(\fR\n" +
	"tssMessageB1\n" +
	"-com.hedera.hapi.services.auxiliary.tss.legacyP\x01b\x06proto3"

var (
	file_tss_message_proto_rawDescOnce sync.Once
	file_tss_message_proto_rawDescData []byte
)

func file_tss_message_proto_rawDescGZIP() []byte {
	file_tss_message_proto_rawDescOnce.Do(func() {
		file_tss_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_tss_message_proto_rawDesc), len(file_tss_message_proto_rawDesc)))
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
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_tss_message_proto_rawDesc), len(file_tss_message_proto_rawDesc)),
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
	file_tss_message_proto_goTypes = nil
	file_tss_message_proto_depIdxs = nil
}
