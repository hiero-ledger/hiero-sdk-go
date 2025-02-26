// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: roster_state.proto

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
// The current state of platform rosters.<br/>
// This message stores a roster data for the platform in network state.
//
// The roster state SHALL encapsulate the incoming candidate roster's hash,
// and a list of pairs of round number and active roster hash.<br/>
// This data SHALL be used to track round numbers and the rosters used in determining the consensus.<br/>
type RosterState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The SHA-384 hash of a candidate roster.
	// <p>
	// This is the hash of the roster that is currently being considered
	// for adoption.<br/>
	// A Node SHALL NOT, ever, have more than one candidate roster
	// at the same time.
	CandidateRosterHash []byte `protobuf:"bytes,1,opt,name=candidate_roster_hash,json=candidateRosterHash,proto3" json:"candidate_roster_hash,omitempty"`
	// *
	// A list of round numbers and roster hashes.<br/>
	// The round number indicates the round in which the corresponding roster became active
	// <p>
	// This list SHALL be ordered by round numbers in descending order.
	RoundRosterPairs []*RoundRosterPair `protobuf:"bytes,2,rep,name=round_roster_pairs,json=roundRosterPairs,proto3" json:"round_roster_pairs,omitempty"`
}

func (x *RosterState) Reset() {
	*x = RosterState{}
	mi := &file_roster_state_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RosterState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RosterState) ProtoMessage() {}

func (x *RosterState) ProtoReflect() protoreflect.Message {
	mi := &file_roster_state_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RosterState.ProtoReflect.Descriptor instead.
func (*RosterState) Descriptor() ([]byte, []int) {
	return file_roster_state_proto_rawDescGZIP(), []int{0}
}

func (x *RosterState) GetCandidateRosterHash() []byte {
	if x != nil {
		return x.CandidateRosterHash
	}
	return nil
}

func (x *RosterState) GetRoundRosterPairs() []*RoundRosterPair {
	if x != nil {
		return x.RoundRosterPairs
	}
	return nil
}

// *
// A pair of round number and active roster hash.
// <p>
// This message SHALL encapsulate the round number and the hash of the
// active roster used for that round.
type RoundRosterPair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The round number.
	// <p>
	// This value SHALL be the round number of the consensus round in which this roster became active.
	RoundNumber uint64 `protobuf:"varint,1,opt,name=round_number,json=roundNumber,proto3" json:"round_number,omitempty"`
	// *
	// The SHA-384 hash of the active roster for the given round number.
	// <p>
	// This value SHALL be the hash of the active roster used for the round.
	ActiveRosterHash []byte `protobuf:"bytes,2,opt,name=active_roster_hash,json=activeRosterHash,proto3" json:"active_roster_hash,omitempty"`
}

func (x *RoundRosterPair) Reset() {
	*x = RoundRosterPair{}
	mi := &file_roster_state_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RoundRosterPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoundRosterPair) ProtoMessage() {}

func (x *RoundRosterPair) ProtoReflect() protoreflect.Message {
	mi := &file_roster_state_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoundRosterPair.ProtoReflect.Descriptor instead.
func (*RoundRosterPair) Descriptor() ([]byte, []int) {
	return file_roster_state_proto_rawDescGZIP(), []int{1}
}

func (x *RoundRosterPair) GetRoundNumber() uint64 {
	if x != nil {
		return x.RoundNumber
	}
	return 0
}

func (x *RoundRosterPair) GetActiveRosterHash() []byte {
	if x != nil {
		return x.ActiveRosterHash
	}
	return nil
}

var File_roster_state_proto protoreflect.FileDescriptor

var file_roster_state_proto_rawDesc = []byte{
	0x0a, 0x12, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x21, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61,
	0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x2e, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x22, 0xa3, 0x01, 0x0a, 0x0b, 0x52, 0x6f, 0x73, 0x74,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x32, 0x0a, 0x15, 0x63, 0x61, 0x6e, 0x64, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x13, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x12, 0x60, 0x0a, 0x12, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x69, 0x72,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65,
	0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x2e, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x52, 0x6f, 0x75, 0x6e,
	0x64, 0x52, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x69, 0x72, 0x52, 0x10, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x52, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x69, 0x72, 0x73, 0x22, 0x62, 0x0a,
	0x0f, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x69, 0x72,
	0x12, 0x21, 0x0a, 0x0c, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x12, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x72, 0x6f,
	0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x10, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x48, 0x61, 0x73,
	0x68, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68,
	0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_roster_state_proto_rawDescOnce sync.Once
	file_roster_state_proto_rawDescData = file_roster_state_proto_rawDesc
)

func file_roster_state_proto_rawDescGZIP() []byte {
	file_roster_state_proto_rawDescOnce.Do(func() {
		file_roster_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_roster_state_proto_rawDescData)
	})
	return file_roster_state_proto_rawDescData
}

var file_roster_state_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_roster_state_proto_goTypes = []any{
	(*RosterState)(nil),     // 0: com.hedera.hapi.node.state.roster.RosterState
	(*RoundRosterPair)(nil), // 1: com.hedera.hapi.node.state.roster.RoundRosterPair
}
var file_roster_state_proto_depIdxs = []int32{
	1, // 0: com.hedera.hapi.node.state.roster.RosterState.round_roster_pairs:type_name -> com.hedera.hapi.node.state.roster.RoundRosterPair
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_roster_state_proto_init() }
func file_roster_state_proto_init() {
	if File_roster_state_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_roster_state_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_roster_state_proto_goTypes,
		DependencyIndexes: file_roster_state_proto_depIdxs,
		MessageInfos:      file_roster_state_proto_msgTypes,
	}.Build()
	File_roster_state_proto = out.File
	file_roster_state_proto_rawDesc = nil
	file_roster_state_proto_goTypes = nil
	file_roster_state_proto_depIdxs = nil
}
