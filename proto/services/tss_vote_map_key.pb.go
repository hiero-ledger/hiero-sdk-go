//*
// # Tss Vote Map Key
//
// ### Keywords
// The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
// "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
// document are to be interpreted as described in
// [RFC2119](https://www.ietf.org/rfc/rfc2119) and clarified in
// [RFC8174](https://www.ietf.org/rfc/rfc8174).

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: tss_vote_map_key.proto

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
// A key for use in the Threshold Signature Scheme (TSS) TssVoteMaps.
//
// This key SHALL be used to uniquely identify entries in the Vote Maps.
type TssVoteMapKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A hash of the target roster for the associated value in the map. <br/>
	// This hash uniquely identifies the target roster.
	// <p>
	// This value MUST be set.<br/>
	// This value MUST contain a valid hash.
	RosterHash []byte `protobuf:"bytes,1,opt,name=roster_hash,json=rosterHash,proto3" json:"roster_hash,omitempty"`
	// * The node id of the node that created the TssVote.</br>
	// This id uniquely identifies the node.
	// <p>
	// This value MUST be set.<br/>
	// This value MUST be a valid node id.
	NodeId uint64 `protobuf:"varint,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *TssVoteMapKey) Reset() {
	*x = TssVoteMapKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tss_vote_map_key_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TssVoteMapKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TssVoteMapKey) ProtoMessage() {}

func (x *TssVoteMapKey) ProtoReflect() protoreflect.Message {
	mi := &file_tss_vote_map_key_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TssVoteMapKey.ProtoReflect.Descriptor instead.
func (*TssVoteMapKey) Descriptor() ([]byte, []int) {
	return file_tss_vote_map_key_proto_rawDescGZIP(), []int{0}
}

func (x *TssVoteMapKey) GetRosterHash() []byte {
	if x != nil {
		return x.RosterHash
	}
	return nil
}

func (x *TssVoteMapKey) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

var File_tss_vote_map_key_proto protoreflect.FileDescriptor

var file_tss_vote_map_key_proto_rawDesc = []byte{
	0x0a, 0x16, 0x74, 0x73, 0x73, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x5f, 0x6d, 0x61, 0x70, 0x5f, 0x6b,
	0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65,
	0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x2e, 0x74, 0x73, 0x73, 0x22, 0x49, 0x0a, 0x0d, 0x54, 0x73, 0x73, 0x56,
	0x6f, 0x74, 0x65, 0x4d, 0x61, 0x70, 0x4b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x6f, 0x73,
	0x74, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a,
	0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f,
	0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64,
	0x65, 0x49, 0x64, 0x42, 0x29, 0x0a, 0x25, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72,
	0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x2e, 0x74, 0x73, 0x73, 0x2e, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x50, 0x01, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tss_vote_map_key_proto_rawDescOnce sync.Once
	file_tss_vote_map_key_proto_rawDescData = file_tss_vote_map_key_proto_rawDesc
)

func file_tss_vote_map_key_proto_rawDescGZIP() []byte {
	file_tss_vote_map_key_proto_rawDescOnce.Do(func() {
		file_tss_vote_map_key_proto_rawDescData = protoimpl.X.CompressGZIP(file_tss_vote_map_key_proto_rawDescData)
	})
	return file_tss_vote_map_key_proto_rawDescData
}

var file_tss_vote_map_key_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_tss_vote_map_key_proto_goTypes = []interface{}{
	(*TssVoteMapKey)(nil), // 0: com.hedera.hapi.node.state.tss.TssVoteMapKey
}
var file_tss_vote_map_key_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tss_vote_map_key_proto_init() }
func file_tss_vote_map_key_proto_init() {
	if File_tss_vote_map_key_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tss_vote_map_key_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TssVoteMapKey); i {
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
			RawDescriptor: file_tss_vote_map_key_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tss_vote_map_key_proto_goTypes,
		DependencyIndexes: file_tss_vote_map_key_proto_depIdxs,
		MessageInfos:      file_tss_vote_map_key_proto_msgTypes,
	}.Build()
	File_tss_vote_map_key_proto = out.File
	file_tss_vote_map_key_proto_rawDesc = nil
	file_tss_vote_map_key_proto_goTypes = nil
	file_tss_vote_map_key_proto_depIdxs = nil
}
