//*
// # Metadata Proof Vote Transaction
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
// source: history_proof_vote.proto

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
// A transaction body to publish a node's vote for a
// proof of history associated to a construction id.
type HistoryProofVoteTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The id of the proof construction this vote is for.
	ConstructionId uint64 `protobuf:"varint,1,opt,name=construction_id,json=constructionId,proto3" json:"construction_id,omitempty"`
	// *
	// The submitting node's vote on the history proof.
	Vote *HistoryProofVote `protobuf:"bytes,2,opt,name=vote,proto3" json:"vote,omitempty"`
}

func (x *HistoryProofVoteTransactionBody) Reset() {
	*x = HistoryProofVoteTransactionBody{}
	mi := &file_history_proof_vote_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HistoryProofVoteTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoryProofVoteTransactionBody) ProtoMessage() {}

func (x *HistoryProofVoteTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_history_proof_vote_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoryProofVoteTransactionBody.ProtoReflect.Descriptor instead.
func (*HistoryProofVoteTransactionBody) Descriptor() ([]byte, []int) {
	return file_history_proof_vote_proto_rawDescGZIP(), []int{0}
}

func (x *HistoryProofVoteTransactionBody) GetConstructionId() uint64 {
	if x != nil {
		return x.ConstructionId
	}
	return 0
}

func (x *HistoryProofVoteTransactionBody) GetVote() *HistoryProofVote {
	if x != nil {
		return x.Vote
	}
	return nil
}

var File_history_proof_vote_proto protoreflect.FileDescriptor

var file_history_proof_vote_proto_rawDesc = []byte{
	0x0a, 0x18, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x5f,
	0x76, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x2a, 0x63, 0x6f, 0x6d, 0x2e,
	0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x75, 0x78, 0x69, 0x6c, 0x69, 0x61, 0x72, 0x79, 0x2e, 0x68,
	0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x1a, 0x13, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x94, 0x01, 0x0a, 0x1f,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x56, 0x6f, 0x74, 0x65,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12,
	0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x48, 0x0a, 0x04, 0x76, 0x6f, 0x74, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64,
	0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x2e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x48, 0x69, 0x73, 0x74,
	0x6f, 0x72, 0x79, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x04, 0x76, 0x6f,
	0x74, 0x65, 0x42, 0x35, 0x0a, 0x31, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61,
	0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x61,
	0x75, 0x78, 0x69, 0x6c, 0x69, 0x61, 0x72, 0x79, 0x2e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_history_proof_vote_proto_rawDescOnce sync.Once
	file_history_proof_vote_proto_rawDescData = file_history_proof_vote_proto_rawDesc
)

func file_history_proof_vote_proto_rawDescGZIP() []byte {
	file_history_proof_vote_proto_rawDescOnce.Do(func() {
		file_history_proof_vote_proto_rawDescData = protoimpl.X.CompressGZIP(file_history_proof_vote_proto_rawDescData)
	})
	return file_history_proof_vote_proto_rawDescData
}

var file_history_proof_vote_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_history_proof_vote_proto_goTypes = []any{
	(*HistoryProofVoteTransactionBody)(nil), // 0: com.hedera.hapi.services.auxiliary.history.HistoryProofVoteTransactionBody
	(*HistoryProofVote)(nil),                // 1: com.hedera.hapi.node.state.history.HistoryProofVote
}
var file_history_proof_vote_proto_depIdxs = []int32{
	1, // 0: com.hedera.hapi.services.auxiliary.history.HistoryProofVoteTransactionBody.vote:type_name -> com.hedera.hapi.node.state.history.HistoryProofVote
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_history_proof_vote_proto_init() }
func file_history_proof_vote_proto_init() {
	if File_history_proof_vote_proto != nil {
		return
	}
	file_history_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_history_proof_vote_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_history_proof_vote_proto_goTypes,
		DependencyIndexes: file_history_proof_vote_proto_depIdxs,
		MessageInfos:      file_history_proof_vote_proto_msgTypes,
	}.Build()
	File_history_proof_vote_proto = out.File
	file_history_proof_vote_proto_rawDesc = nil
	file_history_proof_vote_proto_goTypes = nil
	file_history_proof_vote_proto_depIdxs = nil
}
