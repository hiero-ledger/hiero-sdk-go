//*
// # Tss Vote Transaction
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
// source: tss_vote.proto

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
// A transaction body to vote on the validity of Threshold Signature Scheme
// (TSS) Messages for a candidate roster.
//
//   - A `TssVoteTransactionBody` MUST identify the hash of the roster containing
//     the node generating this TssVote
//   - A `TssVoteTransactionBody` MUST identify the hash of the roster that the
//     TSS messages is for
//   - If the candidate roster has received enough yes votes, the candidate
//     roster SHALL be adopted.
//   - Switching to the candidate roster MUST not happen until enough nodes have
//     voted that they have verified a threshold number of TSS messages from the
//     active roster.
//   - A vote consists of a bit vector of message statuses where each bit
//     corresponds to the order of TssMessages as they have come through
//     consensus.
//   - The threshold for votes to adopt a candidate roster SHALL be at least 1/3
//     of the consensus weight of the active roster to ensure that at least 1
//     honest node has validated the TSS key material.
type TssVoteTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A hash of the roster containing the node generating this TssVote.
	SourceRosterHash []byte `protobuf:"bytes,1,opt,name=source_roster_hash,json=sourceRosterHash,proto3" json:"source_roster_hash,omitempty"`
	// *
	// A hash of the roster that this TssVote is for.
	TargetRosterHash []byte `protobuf:"bytes,2,opt,name=target_roster_hash,json=targetRosterHash,proto3" json:"target_roster_hash,omitempty"`
	// *
	// An identifier (and public key) computed from the TssMessages for the target
	// roster.
	LedgerId []byte `protobuf:"bytes,3,opt,name=ledger_id,json=ledgerId,proto3" json:"ledger_id,omitempty"`
	// *
	// A signature produced by the node.
	// <p>
	// This signature SHALL be produced using the node RSA signing key to sign
	// the ledger_id.<br/>
	// This signature SHALL be used to establish a chain of trust in the ledger id.
	NodeSignature []byte `protobuf:"bytes,4,opt,name=node_signature,json=nodeSignature,proto3" json:"node_signature,omitempty"`
	// *
	// A bit vector of message statuses.
	// <p>
	// #### Example
	// <ul><li>The least significant bit of byte[0] SHALL be the 0th item in the sequence.</li>
	//
	//	<li>The most significant bit of byte[0] SHALL be the 7th item in the sequence.</li>
	//	<li>The least significant bit of byte[1] SHALL be the 8th item in the sequence.</li>
	//	<li>The most significant bit of byte[1] SHALL be the 15th item in the sequence.</li>
	//
	// </ul>
	// A bit SHALL be set if the `TssMessage` for the `TssMessageTransaction`
	// with a sequence number matching that bit index has been
	// received, and is valid.<br/>
	// A bit SHALL NOT be set if the `TssMessage` has not been received or was
	// received but not valid.
	TssVote []byte `protobuf:"bytes,5,opt,name=tss_vote,json=tssVote,proto3" json:"tss_vote,omitempty"`
}

func (x *TssVoteTransactionBody) Reset() {
	*x = TssVoteTransactionBody{}
	mi := &file_tss_vote_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TssVoteTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TssVoteTransactionBody) ProtoMessage() {}

func (x *TssVoteTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_tss_vote_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TssVoteTransactionBody.ProtoReflect.Descriptor instead.
func (*TssVoteTransactionBody) Descriptor() ([]byte, []int) {
	return file_tss_vote_proto_rawDescGZIP(), []int{0}
}

func (x *TssVoteTransactionBody) GetSourceRosterHash() []byte {
	if x != nil {
		return x.SourceRosterHash
	}
	return nil
}

func (x *TssVoteTransactionBody) GetTargetRosterHash() []byte {
	if x != nil {
		return x.TargetRosterHash
	}
	return nil
}

func (x *TssVoteTransactionBody) GetLedgerId() []byte {
	if x != nil {
		return x.LedgerId
	}
	return nil
}

func (x *TssVoteTransactionBody) GetNodeSignature() []byte {
	if x != nil {
		return x.NodeSignature
	}
	return nil
}

func (x *TssVoteTransactionBody) GetTssVote() []byte {
	if x != nil {
		return x.TssVote
	}
	return nil
}

var File_tss_vote_proto protoreflect.FileDescriptor

var file_tss_vote_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x74, 0x73, 0x73, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x26, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70,
	0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x75, 0x78, 0x69, 0x6c,
	0x69, 0x61, 0x72, 0x79, 0x2e, 0x74, 0x73, 0x73, 0x22, 0xd3, 0x01, 0x0a, 0x16, 0x54, 0x73, 0x73,
	0x56, 0x6f, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x6f, 0x64, 0x79, 0x12, 0x2c, 0x0a, 0x12, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x72, 0x6f,
	0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x10, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x2c, 0x0a, 0x12, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x72, 0x6f, 0x73, 0x74,
	0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x1b, 0x0a, 0x09, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x08, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e,
	0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x6e, 0x6f, 0x64, 0x65, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x73, 0x73, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x74, 0x73, 0x73, 0x56, 0x6f, 0x74, 0x65, 0x42, 0x31,
	0x0a, 0x2d, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70,
	0x69, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x75, 0x78, 0x69, 0x6c,
	0x69, 0x61, 0x72, 0x79, 0x2e, 0x74, 0x73, 0x73, 0x2e, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x50,
	0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tss_vote_proto_rawDescOnce sync.Once
	file_tss_vote_proto_rawDescData = file_tss_vote_proto_rawDesc
)

func file_tss_vote_proto_rawDescGZIP() []byte {
	file_tss_vote_proto_rawDescOnce.Do(func() {
		file_tss_vote_proto_rawDescData = protoimpl.X.CompressGZIP(file_tss_vote_proto_rawDescData)
	})
	return file_tss_vote_proto_rawDescData
}

var file_tss_vote_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_tss_vote_proto_goTypes = []any{
	(*TssVoteTransactionBody)(nil), // 0: com.hedera.hapi.services.auxiliary.tss.TssVoteTransactionBody
}
var file_tss_vote_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tss_vote_proto_init() }
func file_tss_vote_proto_init() {
	if File_tss_vote_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tss_vote_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tss_vote_proto_goTypes,
		DependencyIndexes: file_tss_vote_proto_depIdxs,
		MessageInfos:      file_tss_vote_proto_msgTypes,
	}.Build()
	File_tss_vote_proto = out.File
	file_tss_vote_proto_rawDesc = nil
	file_tss_vote_proto_goTypes = nil
	file_tss_vote_proto_depIdxs = nil
}
