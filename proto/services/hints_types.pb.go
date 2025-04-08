// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: hints_types.proto

package services

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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
// The stage of a CRS construction.
type CRSStage int32

const (
	// *
	// The network is gathering contributions to the CRS from all nodes.
	CRSStage_GATHERING_CONTRIBUTIONS CRSStage = 0
	// *
	// The network is waiting for some grace period to allow the verification future
	// to be completed after the last node has contributed to the CRS.
	CRSStage_WAITING_FOR_ADOPTING_FINAL_CRS CRSStage = 1
	// *
	// The network has completed the CRS construction and is set in the CrsState.
	CRSStage_COMPLETED CRSStage = 2
)

// Enum value maps for CRSStage.
var (
	CRSStage_name = map[int32]string{
		0: "GATHERING_CONTRIBUTIONS",
		1: "WAITING_FOR_ADOPTING_FINAL_CRS",
		2: "COMPLETED",
	}
	CRSStage_value = map[string]int32{
		"GATHERING_CONTRIBUTIONS":        0,
		"WAITING_FOR_ADOPTING_FINAL_CRS": 1,
		"COMPLETED":                      2,
	}
)

func (x CRSStage) Enum() *CRSStage {
	p := new(CRSStage)
	*p = x
	return p
}

func (x CRSStage) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CRSStage) Descriptor() protoreflect.EnumDescriptor {
	return file_hints_types_proto_enumTypes[0].Descriptor()
}

func (CRSStage) Type() protoreflect.EnumType {
	return &file_hints_types_proto_enumTypes[0]
}

func (x CRSStage) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CRSStage.Descriptor instead.
func (CRSStage) EnumDescriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{0}
}

// *
// The id of a party in a hinTS scheme with a certain
// number of parties.
type HintsPartyId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The party id, in the range [0, num_parties).
	PartyId uint32 `protobuf:"varint,1,opt,name=party_id,json=partyId,proto3" json:"party_id,omitempty"`
	// *
	// The number of parties in the hinTS scheme.
	NumParties uint32 `protobuf:"varint,2,opt,name=num_parties,json=numParties,proto3" json:"num_parties,omitempty"`
}

func (x *HintsPartyId) Reset() {
	*x = HintsPartyId{}
	mi := &file_hints_types_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HintsPartyId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HintsPartyId) ProtoMessage() {}

func (x *HintsPartyId) ProtoReflect() protoreflect.Message {
	mi := &file_hints_types_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HintsPartyId.ProtoReflect.Descriptor instead.
func (*HintsPartyId) Descriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{0}
}

func (x *HintsPartyId) GetPartyId() uint32 {
	if x != nil {
		return x.PartyId
	}
	return 0
}

func (x *HintsPartyId) GetNumParties() uint32 {
	if x != nil {
		return x.NumParties
	}
	return 0
}

// *
// A set of hinTS keys submitted by a node.
type HintsKeySet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The id of the node submitting these keys.
	NodeId uint64 `protobuf:"varint,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	// *
	// The consensus time at which the network adopted the active
	// hinTS key in this set.
	AdoptionTime *Timestamp `protobuf:"bytes,2,opt,name=adoption_time,json=adoptionTime,proto3" json:"adoption_time,omitempty"`
	// *
	// The party's active hinTS key.
	Key []byte `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	// *
	// If set, the new hinTS key the node wants to use when
	// the next construction begins.
	NextKey []byte `protobuf:"bytes,4,opt,name=next_key,json=nextKey,proto3" json:"next_key,omitempty"`
}

func (x *HintsKeySet) Reset() {
	*x = HintsKeySet{}
	mi := &file_hints_types_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HintsKeySet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HintsKeySet) ProtoMessage() {}

func (x *HintsKeySet) ProtoReflect() protoreflect.Message {
	mi := &file_hints_types_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HintsKeySet.ProtoReflect.Descriptor instead.
func (*HintsKeySet) Descriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{1}
}

func (x *HintsKeySet) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *HintsKeySet) GetAdoptionTime() *Timestamp {
	if x != nil {
		return x.AdoptionTime
	}
	return nil
}

func (x *HintsKeySet) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *HintsKeySet) GetNextKey() []byte {
	if x != nil {
		return x.NextKey
	}
	return nil
}

// *
// The output of the hinTS preprocessing algorithm; that is, a
// linear-size aggregation key and a succinct verification key.
type PreprocessedKeys struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The aggregation key for the hinTS scheme
	AggregationKey []byte `protobuf:"bytes,1,opt,name=aggregation_key,json=aggregationKey,proto3" json:"aggregation_key,omitempty"`
	// *
	// The succinct verification key for the hinTS scheme.
	VerificationKey []byte `protobuf:"bytes,2,opt,name=verification_key,json=verificationKey,proto3" json:"verification_key,omitempty"`
}

func (x *PreprocessedKeys) Reset() {
	*x = PreprocessedKeys{}
	mi := &file_hints_types_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PreprocessedKeys) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreprocessedKeys) ProtoMessage() {}

func (x *PreprocessedKeys) ProtoReflect() protoreflect.Message {
	mi := &file_hints_types_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreprocessedKeys.ProtoReflect.Descriptor instead.
func (*PreprocessedKeys) Descriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{2}
}

func (x *PreprocessedKeys) GetAggregationKey() []byte {
	if x != nil {
		return x.AggregationKey
	}
	return nil
}

func (x *PreprocessedKeys) GetVerificationKey() []byte {
	if x != nil {
		return x.VerificationKey
	}
	return nil
}

// *
// The id for a node's vote for the output of the
// preprocessing output of a hinTS construction.
type PreprocessingVoteId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The construction this vote is connected to.
	ConstructionId uint64 `protobuf:"varint,1,opt,name=construction_id,json=constructionId,proto3" json:"construction_id,omitempty"`
	// *
	// The id of the node submitting the vote.
	NodeId uint64 `protobuf:"varint,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *PreprocessingVoteId) Reset() {
	*x = PreprocessingVoteId{}
	mi := &file_hints_types_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PreprocessingVoteId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreprocessingVoteId) ProtoMessage() {}

func (x *PreprocessingVoteId) ProtoReflect() protoreflect.Message {
	mi := &file_hints_types_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreprocessingVoteId.ProtoReflect.Descriptor instead.
func (*PreprocessingVoteId) Descriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{3}
}

func (x *PreprocessingVoteId) GetConstructionId() uint64 {
	if x != nil {
		return x.ConstructionId
	}
	return 0
}

func (x *PreprocessingVoteId) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

// *
// A node's vote for the consensus output of a hinTS preprocessing
// algorithm.
type PreprocessingVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Vote:
	//
	//	*PreprocessingVote_PreprocessedKeys
	//	*PreprocessingVote_CongruentNodeId
	Vote isPreprocessingVote_Vote `protobuf_oneof:"vote"`
}

func (x *PreprocessingVote) Reset() {
	*x = PreprocessingVote{}
	mi := &file_hints_types_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PreprocessingVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreprocessingVote) ProtoMessage() {}

func (x *PreprocessingVote) ProtoReflect() protoreflect.Message {
	mi := &file_hints_types_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreprocessingVote.ProtoReflect.Descriptor instead.
func (*PreprocessingVote) Descriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{4}
}

func (m *PreprocessingVote) GetVote() isPreprocessingVote_Vote {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (x *PreprocessingVote) GetPreprocessedKeys() *PreprocessedKeys {
	if x, ok := x.GetVote().(*PreprocessingVote_PreprocessedKeys); ok {
		return x.PreprocessedKeys
	}
	return nil
}

func (x *PreprocessingVote) GetCongruentNodeId() uint64 {
	if x, ok := x.GetVote().(*PreprocessingVote_CongruentNodeId); ok {
		return x.CongruentNodeId
	}
	return 0
}

type isPreprocessingVote_Vote interface {
	isPreprocessingVote_Vote()
}

type PreprocessingVote_PreprocessedKeys struct {
	// *
	// The preprocessed keys this node is voting for.
	PreprocessedKeys *PreprocessedKeys `protobuf:"bytes,1,opt,name=preprocessed_keys,json=preprocessedKeys,proto3,oneof"`
}

type PreprocessingVote_CongruentNodeId struct {
	// *
	// The id of any node that already voted for the exact keys
	// that this node wanted to vote for.
	CongruentNodeId uint64 `protobuf:"varint,2,opt,name=congruent_node_id,json=congruentNodeId,proto3,oneof"`
}

func (*PreprocessingVote_PreprocessedKeys) isPreprocessingVote_Vote() {}

func (*PreprocessingVote_CongruentNodeId) isPreprocessingVote_Vote() {}

// *
// A node's hinTS party id.
type NodePartyId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The node id.
	NodeId uint64 `protobuf:"varint,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	// *
	// The party id.
	PartyId uint32 `protobuf:"varint,2,opt,name=party_id,json=partyId,proto3" json:"party_id,omitempty"`
}

func (x *NodePartyId) Reset() {
	*x = NodePartyId{}
	mi := &file_hints_types_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NodePartyId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodePartyId) ProtoMessage() {}

func (x *NodePartyId) ProtoReflect() protoreflect.Message {
	mi := &file_hints_types_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodePartyId.ProtoReflect.Descriptor instead.
func (*NodePartyId) Descriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{5}
}

func (x *NodePartyId) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *NodePartyId) GetPartyId() uint32 {
	if x != nil {
		return x.PartyId
	}
	return 0
}

// *
// The information constituting the hinTS scheme Hiero TSS.
type HintsScheme struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The aggregation and verification keys for the scheme.
	PreprocessedKeys *PreprocessedKeys `protobuf:"bytes,1,opt,name=preprocessed_keys,json=preprocessedKeys,proto3" json:"preprocessed_keys,omitempty"`
	// *
	// The final party ids assigned to each node in the target roster.
	NodePartyIds []*NodePartyId `protobuf:"bytes,2,rep,name=node_party_ids,json=nodePartyIds,proto3" json:"node_party_ids,omitempty"`
}

func (x *HintsScheme) Reset() {
	*x = HintsScheme{}
	mi := &file_hints_types_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HintsScheme) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HintsScheme) ProtoMessage() {}

func (x *HintsScheme) ProtoReflect() protoreflect.Message {
	mi := &file_hints_types_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HintsScheme.ProtoReflect.Descriptor instead.
func (*HintsScheme) Descriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{6}
}

func (x *HintsScheme) GetPreprocessedKeys() *PreprocessedKeys {
	if x != nil {
		return x.PreprocessedKeys
	}
	return nil
}

func (x *HintsScheme) GetNodePartyIds() []*NodePartyId {
	if x != nil {
		return x.NodePartyIds
	}
	return nil
}

// *
// A summary of progress in constructing a hinTS scheme.
type HintsConstruction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The id of the construction.
	ConstructionId uint64 `protobuf:"varint,1,opt,name=construction_id,json=constructionId,proto3" json:"construction_id,omitempty"`
	// *
	// The hash of the roster whose weights are used to determine when
	// the >=1/3 weight signing threshold is reached.
	SourceRosterHash []byte `protobuf:"bytes,2,opt,name=source_roster_hash,json=sourceRosterHash,proto3" json:"source_roster_hash,omitempty"`
	// *
	// The hash of the roster whose weights are used to determine when
	// the >2/3 weight availability threshold is reached.
	TargetRosterHash []byte `protobuf:"bytes,3,opt,name=target_roster_hash,json=targetRosterHash,proto3" json:"target_roster_hash,omitempty"`
	// Types that are assignable to PreprocessingState:
	//
	//	*HintsConstruction_GracePeriodEndTime
	//	*HintsConstruction_PreprocessingStartTime
	//	*HintsConstruction_HintsScheme
	PreprocessingState isHintsConstruction_PreprocessingState `protobuf_oneof:"preprocessing_state"`
}

func (x *HintsConstruction) Reset() {
	*x = HintsConstruction{}
	mi := &file_hints_types_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HintsConstruction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HintsConstruction) ProtoMessage() {}

func (x *HintsConstruction) ProtoReflect() protoreflect.Message {
	mi := &file_hints_types_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HintsConstruction.ProtoReflect.Descriptor instead.
func (*HintsConstruction) Descriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{7}
}

func (x *HintsConstruction) GetConstructionId() uint64 {
	if x != nil {
		return x.ConstructionId
	}
	return 0
}

func (x *HintsConstruction) GetSourceRosterHash() []byte {
	if x != nil {
		return x.SourceRosterHash
	}
	return nil
}

func (x *HintsConstruction) GetTargetRosterHash() []byte {
	if x != nil {
		return x.TargetRosterHash
	}
	return nil
}

func (m *HintsConstruction) GetPreprocessingState() isHintsConstruction_PreprocessingState {
	if m != nil {
		return m.PreprocessingState
	}
	return nil
}

func (x *HintsConstruction) GetGracePeriodEndTime() *Timestamp {
	if x, ok := x.GetPreprocessingState().(*HintsConstruction_GracePeriodEndTime); ok {
		return x.GracePeriodEndTime
	}
	return nil
}

func (x *HintsConstruction) GetPreprocessingStartTime() *Timestamp {
	if x, ok := x.GetPreprocessingState().(*HintsConstruction_PreprocessingStartTime); ok {
		return x.PreprocessingStartTime
	}
	return nil
}

func (x *HintsConstruction) GetHintsScheme() *HintsScheme {
	if x, ok := x.GetPreprocessingState().(*HintsConstruction_HintsScheme); ok {
		return x.HintsScheme
	}
	return nil
}

type isHintsConstruction_PreprocessingState interface {
	isHintsConstruction_PreprocessingState()
}

type HintsConstruction_GracePeriodEndTime struct {
	// *
	// If the network is still gathering hinTS keys for this construction,
	// the time at which honest nodes should stop waiting for tardy
	// publications and begin preprocessing as soon as there are valid
	// hinTS keys for nodes with >2/3 weight in the target roster.
	GracePeriodEndTime *Timestamp `protobuf:"bytes,4,opt,name=grace_period_end_time,json=gracePeriodEndTime,proto3,oneof"`
}

type HintsConstruction_PreprocessingStartTime struct {
	// *
	// If the network has gathered enough hinTS keys for this construction
	// to begin preprocessing, the cutoff time by which keys must have been
	// adopted to be included as input to the preprocessing algorithm.
	PreprocessingStartTime *Timestamp `protobuf:"bytes,5,opt,name=preprocessing_start_time,json=preprocessingStartTime,proto3,oneof"`
}

type HintsConstruction_HintsScheme struct {
	// *
	// If set, the completed hinTS scheme.
	HintsScheme *HintsScheme `protobuf:"bytes,6,opt,name=hints_scheme,json=hintsScheme,proto3,oneof"`
}

func (*HintsConstruction_GracePeriodEndTime) isHintsConstruction_PreprocessingState() {}

func (*HintsConstruction_PreprocessingStartTime) isHintsConstruction_PreprocessingState() {}

func (*HintsConstruction_HintsScheme) isHintsConstruction_PreprocessingState() {}

// *
// The state of a CRS construction.
type CRSState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The bytes of the CRS. Based on the CRSStage, this may be the initial CRS
	// or the final CRS.
	Crs []byte `protobuf:"bytes,1,opt,name=crs,proto3" json:"crs,omitempty"`
	// *
	// The stage of the CRS construction.
	Stage CRSStage `protobuf:"varint,2,opt,name=stage,proto3,enum=com.hedera.hapi.node.state.hints.CRSStage" json:"stage,omitempty"`
	// *
	// The id of the next node that should contribute to the CRS. This is used
	// to ensure that all nodes contribute to the CRS in a round-robin fashion.
	// If this is null, then all nodes in the network have contributed to the CRS.
	NextContributingNodeId *wrapperspb.UInt64Value `protobuf:"bytes,3,opt,name=next_contributing_node_id,json=nextContributingNodeId,proto3" json:"next_contributing_node_id,omitempty"`
	// *
	// The time at which the network should stop waiting for the node's contributions
	// and move on to the next node in the round-robin fashion.
	ContributionEndTime *Timestamp `protobuf:"bytes,4,opt,name=contribution_end_time,json=contributionEndTime,proto3" json:"contribution_end_time,omitempty"`
}

func (x *CRSState) Reset() {
	*x = CRSState{}
	mi := &file_hints_types_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CRSState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CRSState) ProtoMessage() {}

func (x *CRSState) ProtoReflect() protoreflect.Message {
	mi := &file_hints_types_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CRSState.ProtoReflect.Descriptor instead.
func (*CRSState) Descriptor() ([]byte, []int) {
	return file_hints_types_proto_rawDescGZIP(), []int{8}
}

func (x *CRSState) GetCrs() []byte {
	if x != nil {
		return x.Crs
	}
	return nil
}

func (x *CRSState) GetStage() CRSStage {
	if x != nil {
		return x.Stage
	}
	return CRSStage_GATHERING_CONTRIBUTIONS
}

func (x *CRSState) GetNextContributingNodeId() *wrapperspb.UInt64Value {
	if x != nil {
		return x.NextContributingNodeId
	}
	return nil
}

func (x *CRSState) GetContributionEndTime() *Timestamp {
	if x != nil {
		return x.ContributionEndTime
	}
	return nil
}

var File_hints_types_proto protoreflect.FileDescriptor

var file_hints_types_proto_rawDesc = []byte{
	0x0a, 0x11, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x20, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e,
	0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e,
	0x68, 0x69, 0x6e, 0x74, 0x73, 0x1a, 0x0f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x0c, 0x48, 0x69, 0x6e, 0x74, 0x73, 0x50,
	0x61, 0x72, 0x74, 0x79, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x61, 0x72, 0x74, 0x79, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x70, 0x61, 0x72, 0x74, 0x79, 0x49,
	0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x75, 0x6d, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x6e, 0x75, 0x6d, 0x50, 0x61, 0x72, 0x74, 0x69,
	0x65, 0x73, 0x22, 0x8a, 0x01, 0x0a, 0x0b, 0x48, 0x69, 0x6e, 0x74, 0x73, 0x4b, 0x65, 0x79, 0x53,
	0x65, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x0d, 0x61,
	0x64, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x61, 0x64, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x6b, 0x65, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6e, 0x65, 0x78, 0x74, 0x4b, 0x65, 0x79, 0x22,
	0x66, 0x0a, 0x10, 0x50, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x4b,
	0x65, 0x79, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x61, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x10,
	0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x22, 0x57, 0x0a, 0x13, 0x50, 0x72, 0x65, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x56, 0x6f, 0x74, 0x65, 0x49, 0x64, 0x12, 0x27,
	0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x22, 0xac, 0x01, 0x0a, 0x11, 0x50, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69,
	0x6e, 0x67, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x61, 0x0a, 0x11, 0x70, 0x72, 0x65, 0x70, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x32, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68,
	0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68,
	0x69, 0x6e, 0x74, 0x73, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65,
	0x64, 0x4b, 0x65, 0x79, 0x73, 0x48, 0x00, 0x52, 0x10, 0x70, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x63, 0x6f, 0x6e,
	0x67, 0x72, 0x75, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x65, 0x6e,
	0x74, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x76, 0x6f, 0x74, 0x65, 0x22,
	0x41, 0x0a, 0x0b, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x61, 0x72, 0x74, 0x79, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x61, 0x72, 0x74, 0x79,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x70, 0x61, 0x72, 0x74, 0x79,
	0x49, 0x64, 0x22, 0xc3, 0x01, 0x0a, 0x0b, 0x48, 0x69, 0x6e, 0x74, 0x73, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x65, 0x12, 0x5f, 0x0a, 0x11, 0x70, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x65, 0x64, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68, 0x69, 0x6e, 0x74, 0x73,
	0x2e, 0x50, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x4b, 0x65, 0x79,
	0x73, 0x52, 0x10, 0x70, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x4b,
	0x65, 0x79, 0x73, 0x12, 0x53, 0x0a, 0x0e, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x70, 0x61, 0x72, 0x74,
	0x79, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x2e, 0x4e,
	0x6f, 0x64, 0x65, 0x50, 0x61, 0x72, 0x74, 0x79, 0x49, 0x64, 0x52, 0x0c, 0x6e, 0x6f, 0x64, 0x65,
	0x50, 0x61, 0x72, 0x74, 0x79, 0x49, 0x64, 0x73, 0x22, 0x98, 0x03, 0x0a, 0x11, 0x48, 0x69, 0x6e,
	0x74, 0x73, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27,
	0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x12, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x10, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x6f, 0x73, 0x74, 0x65,
	0x72, 0x48, 0x61, 0x73, 0x68, 0x12, 0x2c, 0x0a, 0x12, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f,
	0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x10, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x48,
	0x61, 0x73, 0x68, 0x12, 0x45, 0x0a, 0x15, 0x67, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x70, 0x65, 0x72,
	0x69, 0x6f, 0x64, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x12, 0x67, 0x72, 0x61, 0x63, 0x65, 0x50, 0x65, 0x72,
	0x69, 0x6f, 0x64, 0x45, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x4c, 0x0a, 0x18, 0x70, 0x72,
	0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00,
	0x52, 0x16, 0x70, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x52, 0x0a, 0x0c, 0x68, 0x69, 0x6e, 0x74,
	0x73, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d,
	0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69,
	0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68, 0x69, 0x6e, 0x74,
	0x73, 0x2e, 0x48, 0x69, 0x6e, 0x74, 0x73, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x48, 0x00, 0x52,
	0x0b, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x42, 0x15, 0x0a, 0x13,
	0x70, 0x72, 0x65, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x22, 0xfd, 0x01, 0x0a, 0x08, 0x43, 0x52, 0x53, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x63, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x63,
	0x72, 0x73, 0x12, 0x40, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x2a, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68,
	0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68,
	0x69, 0x6e, 0x74, 0x73, 0x2e, 0x43, 0x52, 0x53, 0x53, 0x74, 0x61, 0x67, 0x65, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x67, 0x65, 0x12, 0x57, 0x0a, 0x19, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x16, 0x6e, 0x65, 0x78, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x44, 0x0a,
	0x15, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x6e,
	0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x13,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x54,
	0x69, 0x6d, 0x65, 0x2a, 0x5a, 0x0a, 0x08, 0x43, 0x52, 0x53, 0x53, 0x74, 0x61, 0x67, 0x65, 0x12,
	0x1b, 0x0a, 0x17, 0x47, 0x41, 0x54, 0x48, 0x45, 0x52, 0x49, 0x4e, 0x47, 0x5f, 0x43, 0x4f, 0x4e,
	0x54, 0x52, 0x49, 0x42, 0x55, 0x54, 0x49, 0x4f, 0x4e, 0x53, 0x10, 0x00, 0x12, 0x22, 0x0a, 0x1e,
	0x57, 0x41, 0x49, 0x54, 0x49, 0x4e, 0x47, 0x5f, 0x46, 0x4f, 0x52, 0x5f, 0x41, 0x44, 0x4f, 0x50,
	0x54, 0x49, 0x4e, 0x47, 0x5f, 0x46, 0x49, 0x4e, 0x41, 0x4c, 0x5f, 0x43, 0x52, 0x53, 0x10, 0x01,
	0x12, 0x0d, 0x0a, 0x09, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x02, 0x42,
	0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73,
	0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hints_types_proto_rawDescOnce sync.Once
	file_hints_types_proto_rawDescData = file_hints_types_proto_rawDesc
)

func file_hints_types_proto_rawDescGZIP() []byte {
	file_hints_types_proto_rawDescOnce.Do(func() {
		file_hints_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_hints_types_proto_rawDescData)
	})
	return file_hints_types_proto_rawDescData
}

var file_hints_types_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_hints_types_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_hints_types_proto_goTypes = []any{
	(CRSStage)(0),                  // 0: com.hedera.hapi.node.state.hints.CRSStage
	(*HintsPartyId)(nil),           // 1: com.hedera.hapi.node.state.hints.HintsPartyId
	(*HintsKeySet)(nil),            // 2: com.hedera.hapi.node.state.hints.HintsKeySet
	(*PreprocessedKeys)(nil),       // 3: com.hedera.hapi.node.state.hints.PreprocessedKeys
	(*PreprocessingVoteId)(nil),    // 4: com.hedera.hapi.node.state.hints.PreprocessingVoteId
	(*PreprocessingVote)(nil),      // 5: com.hedera.hapi.node.state.hints.PreprocessingVote
	(*NodePartyId)(nil),            // 6: com.hedera.hapi.node.state.hints.NodePartyId
	(*HintsScheme)(nil),            // 7: com.hedera.hapi.node.state.hints.HintsScheme
	(*HintsConstruction)(nil),      // 8: com.hedera.hapi.node.state.hints.HintsConstruction
	(*CRSState)(nil),               // 9: com.hedera.hapi.node.state.hints.CRSState
	(*Timestamp)(nil),              // 10: proto.Timestamp
	(*wrapperspb.UInt64Value)(nil), // 11: google.protobuf.UInt64Value
}
var file_hints_types_proto_depIdxs = []int32{
	10, // 0: com.hedera.hapi.node.state.hints.HintsKeySet.adoption_time:type_name -> proto.Timestamp
	3,  // 1: com.hedera.hapi.node.state.hints.PreprocessingVote.preprocessed_keys:type_name -> com.hedera.hapi.node.state.hints.PreprocessedKeys
	3,  // 2: com.hedera.hapi.node.state.hints.HintsScheme.preprocessed_keys:type_name -> com.hedera.hapi.node.state.hints.PreprocessedKeys
	6,  // 3: com.hedera.hapi.node.state.hints.HintsScheme.node_party_ids:type_name -> com.hedera.hapi.node.state.hints.NodePartyId
	10, // 4: com.hedera.hapi.node.state.hints.HintsConstruction.grace_period_end_time:type_name -> proto.Timestamp
	10, // 5: com.hedera.hapi.node.state.hints.HintsConstruction.preprocessing_start_time:type_name -> proto.Timestamp
	7,  // 6: com.hedera.hapi.node.state.hints.HintsConstruction.hints_scheme:type_name -> com.hedera.hapi.node.state.hints.HintsScheme
	0,  // 7: com.hedera.hapi.node.state.hints.CRSState.stage:type_name -> com.hedera.hapi.node.state.hints.CRSStage
	11, // 8: com.hedera.hapi.node.state.hints.CRSState.next_contributing_node_id:type_name -> google.protobuf.UInt64Value
	10, // 9: com.hedera.hapi.node.state.hints.CRSState.contribution_end_time:type_name -> proto.Timestamp
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_hints_types_proto_init() }
func file_hints_types_proto_init() {
	if File_hints_types_proto != nil {
		return
	}
	file_timestamp_proto_init()
	file_hints_types_proto_msgTypes[4].OneofWrappers = []any{
		(*PreprocessingVote_PreprocessedKeys)(nil),
		(*PreprocessingVote_CongruentNodeId)(nil),
	}
	file_hints_types_proto_msgTypes[7].OneofWrappers = []any{
		(*HintsConstruction_GracePeriodEndTime)(nil),
		(*HintsConstruction_PreprocessingStartTime)(nil),
		(*HintsConstruction_HintsScheme)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_hints_types_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_hints_types_proto_goTypes,
		DependencyIndexes: file_hints_types_proto_depIdxs,
		EnumInfos:         file_hints_types_proto_enumTypes,
		MessageInfos:      file_hints_types_proto_msgTypes,
	}.Build()
	File_hints_types_proto = out.File
	file_hints_types_proto_rawDesc = nil
	file_hints_types_proto_goTypes = nil
	file_hints_types_proto_depIdxs = nil
}
