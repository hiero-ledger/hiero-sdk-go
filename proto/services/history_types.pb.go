// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: history_types.proto

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
// A set of proof keys for a node; that is, the key the node is
// currently using and the key it wants to use in assembling the
// next address book in the ledger id's chain of trust.
type ProofKeySet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The consensus time when the network adopted the active
	// proof key in this set. An adoption time that is sufficiently
	// tardy relative to the latest assembly start time may result
	// in the node's key being omitted from the address book.
	AdoptionTime *Timestamp `protobuf:"bytes,2,opt,name=adoption_time,json=adoptionTime,proto3" json:"adoption_time,omitempty"`
	// *
	// The proof key the node is using.
	Key []byte `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	// *
	// If set, the proof key the node wants to start using in the
	// address book.
	NextKey []byte `protobuf:"bytes,4,opt,name=next_key,json=nextKey,proto3" json:"next_key,omitempty"`
}

func (x *ProofKeySet) Reset() {
	*x = ProofKeySet{}
	mi := &file_history_types_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProofKeySet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProofKeySet) ProtoMessage() {}

func (x *ProofKeySet) ProtoReflect() protoreflect.Message {
	mi := &file_history_types_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProofKeySet.ProtoReflect.Descriptor instead.
func (*ProofKeySet) Descriptor() ([]byte, []int) {
	return file_history_types_proto_rawDescGZIP(), []int{0}
}

func (x *ProofKeySet) GetAdoptionTime() *Timestamp {
	if x != nil {
		return x.AdoptionTime
	}
	return nil
}

func (x *ProofKeySet) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *ProofKeySet) GetNextKey() []byte {
	if x != nil {
		return x.NextKey
	}
	return nil
}

// *
// A record of the proof key a node had in a particular address
// book. Necessary to keep at each point history so that nodes
// can verify the correct key was used to sign in transitions
// starting from the current address book; no matter how keys
// have been rotated from the time the address book was created.
type ProofKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The node id.
	NodeId uint64 `protobuf:"varint,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	// *
	// The key.
	Key []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *ProofKey) Reset() {
	*x = ProofKey{}
	mi := &file_history_types_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProofKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProofKey) ProtoMessage() {}

func (x *ProofKey) ProtoReflect() protoreflect.Message {
	mi := &file_history_types_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProofKey.ProtoReflect.Descriptor instead.
func (*ProofKey) Descriptor() ([]byte, []int) {
	return file_history_types_proto_rawDescGZIP(), []int{1}
}

func (x *ProofKey) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *ProofKey) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

// *
// A piece of new history in the form of an address book hash and
// associated metadata.
type History struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The address book hash of the new history.
	AddressBookHash []byte `protobuf:"bytes,1,opt,name=address_book_hash,json=addressBookHash,proto3" json:"address_book_hash,omitempty"`
	// *
	// The metadata associated to the address book.
	Metadata []byte `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *History) Reset() {
	*x = History{}
	mi := &file_history_types_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *History) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*History) ProtoMessage() {}

func (x *History) ProtoReflect() protoreflect.Message {
	mi := &file_history_types_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use History.ProtoReflect.Descriptor instead.
func (*History) Descriptor() ([]byte, []int) {
	return file_history_types_proto_rawDescGZIP(), []int{2}
}

func (x *History) GetAddressBookHash() []byte {
	if x != nil {
		return x.AddressBookHash
	}
	return nil
}

func (x *History) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// *
// A proof that some address book history belongs to the ledger id's
// chain of trust.
type HistoryProof struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The hash of the source address book.
	SourceAddressBookHash []byte `protobuf:"bytes,1,opt,name=source_address_book_hash,json=sourceAddressBookHash,proto3" json:"source_address_book_hash,omitempty"`
	// *
	// The proof keys for the target address book, needed to keep
	// constructing proofs after adopting the target address book's
	// roster at a handoff.
	TargetProofKeys []*ProofKey `protobuf:"bytes,2,rep,name=target_proof_keys,json=targetProofKeys,proto3" json:"target_proof_keys,omitempty"`
	// *
	// The target history of the proof.
	TargetHistory *History `protobuf:"bytes,3,opt,name=target_history,json=targetHistory,proto3" json:"target_history,omitempty"`
	// *
	// The proof of chain of trust from the ledger id.
	Proof []byte `protobuf:"bytes,4,opt,name=proof,proto3" json:"proof,omitempty"`
}

func (x *HistoryProof) Reset() {
	*x = HistoryProof{}
	mi := &file_history_types_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HistoryProof) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoryProof) ProtoMessage() {}

func (x *HistoryProof) ProtoReflect() protoreflect.Message {
	mi := &file_history_types_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoryProof.ProtoReflect.Descriptor instead.
func (*HistoryProof) Descriptor() ([]byte, []int) {
	return file_history_types_proto_rawDescGZIP(), []int{3}
}

func (x *HistoryProof) GetSourceAddressBookHash() []byte {
	if x != nil {
		return x.SourceAddressBookHash
	}
	return nil
}

func (x *HistoryProof) GetTargetProofKeys() []*ProofKey {
	if x != nil {
		return x.TargetProofKeys
	}
	return nil
}

func (x *HistoryProof) GetTargetHistory() *History {
	if x != nil {
		return x.TargetHistory
	}
	return nil
}

func (x *HistoryProof) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

// *
// Summary of the status of constructing a metadata proof, necessary to
// ensure deterministic construction ending in a roster with sufficient
// weight to enact its own constructions.
type HistoryProofConstruction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The construction id.
	ConstructionId uint64 `protobuf:"varint,1,opt,name=construction_id,json=constructionId,proto3" json:"construction_id,omitempty"`
	// *
	// The hash of the roster whose weights are used to determine when
	// certain thresholds are during construction.
	SourceRosterHash []byte `protobuf:"bytes,2,opt,name=source_roster_hash,json=sourceRosterHash,proto3" json:"source_roster_hash,omitempty"`
	// *
	// If set, the proof that the address book of the source roster belongs
	// to the the ledger id's chain of trust; if not set, the source roster's
	// address book must *be* the ledger id.
	SourceProof *HistoryProof `protobuf:"bytes,3,opt,name=source_proof,json=sourceProof,proto3" json:"source_proof,omitempty"`
	// *
	// The hash of the roster whose weights are used to assess progress
	// toward obtaining proof keys for parties that hold at least a
	// strong minority of the stake in that roster.
	TargetRosterHash []byte `protobuf:"bytes,4,opt,name=target_roster_hash,json=targetRosterHash,proto3" json:"target_roster_hash,omitempty"`
	// Types that are assignable to ProofState:
	//
	//	*HistoryProofConstruction_GracePeriodEndTime
	//	*HistoryProofConstruction_AssemblyStartTime
	//	*HistoryProofConstruction_TargetProof
	//	*HistoryProofConstruction_FailureReason
	ProofState isHistoryProofConstruction_ProofState `protobuf_oneof:"proof_state"`
}

func (x *HistoryProofConstruction) Reset() {
	*x = HistoryProofConstruction{}
	mi := &file_history_types_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HistoryProofConstruction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoryProofConstruction) ProtoMessage() {}

func (x *HistoryProofConstruction) ProtoReflect() protoreflect.Message {
	mi := &file_history_types_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoryProofConstruction.ProtoReflect.Descriptor instead.
func (*HistoryProofConstruction) Descriptor() ([]byte, []int) {
	return file_history_types_proto_rawDescGZIP(), []int{4}
}

func (x *HistoryProofConstruction) GetConstructionId() uint64 {
	if x != nil {
		return x.ConstructionId
	}
	return 0
}

func (x *HistoryProofConstruction) GetSourceRosterHash() []byte {
	if x != nil {
		return x.SourceRosterHash
	}
	return nil
}

func (x *HistoryProofConstruction) GetSourceProof() *HistoryProof {
	if x != nil {
		return x.SourceProof
	}
	return nil
}

func (x *HistoryProofConstruction) GetTargetRosterHash() []byte {
	if x != nil {
		return x.TargetRosterHash
	}
	return nil
}

func (m *HistoryProofConstruction) GetProofState() isHistoryProofConstruction_ProofState {
	if m != nil {
		return m.ProofState
	}
	return nil
}

func (x *HistoryProofConstruction) GetGracePeriodEndTime() *Timestamp {
	if x, ok := x.GetProofState().(*HistoryProofConstruction_GracePeriodEndTime); ok {
		return x.GracePeriodEndTime
	}
	return nil
}

func (x *HistoryProofConstruction) GetAssemblyStartTime() *Timestamp {
	if x, ok := x.GetProofState().(*HistoryProofConstruction_AssemblyStartTime); ok {
		return x.AssemblyStartTime
	}
	return nil
}

func (x *HistoryProofConstruction) GetTargetProof() *HistoryProof {
	if x, ok := x.GetProofState().(*HistoryProofConstruction_TargetProof); ok {
		return x.TargetProof
	}
	return nil
}

func (x *HistoryProofConstruction) GetFailureReason() string {
	if x, ok := x.GetProofState().(*HistoryProofConstruction_FailureReason); ok {
		return x.FailureReason
	}
	return ""
}

type isHistoryProofConstruction_ProofState interface {
	isHistoryProofConstruction_ProofState()
}

type HistoryProofConstruction_GracePeriodEndTime struct {
	// *
	// If the network is still gathering proof keys for this
	// construction, the next time at which nodes should stop waiting
	// for tardy proof keys and assembly the history to be proven as
	// soon as it has the associated metadata and proof keys for nodes
	// with >2/3 weight in the target roster.
	GracePeriodEndTime *Timestamp `protobuf:"bytes,5,opt,name=grace_period_end_time,json=gracePeriodEndTime,proto3,oneof"`
}

type HistoryProofConstruction_AssemblyStartTime struct {
	// *
	// If the network has gathered enough proof keys to assemble the
	// history for this construction, the cutoff time at which those
	// keys must have been adopted to be included in the final history.
	AssemblyStartTime *Timestamp `protobuf:"bytes,6,opt,name=assembly_start_time,json=assemblyStartTime,proto3,oneof"`
}

type HistoryProofConstruction_TargetProof struct {
	// *
	// When this construction is complete, the recursive proof that
	// the target roster's address book and associated metadata belong
	// to the ledger id's chain of trust.
	TargetProof *HistoryProof `protobuf:"bytes,7,opt,name=target_proof,json=targetProof,proto3,oneof"`
}

type HistoryProofConstruction_FailureReason struct {
	// *
	// If set, the reason the construction failed.
	FailureReason string `protobuf:"bytes,8,opt,name=failure_reason,json=failureReason,proto3,oneof"`
}

func (*HistoryProofConstruction_GracePeriodEndTime) isHistoryProofConstruction_ProofState() {}

func (*HistoryProofConstruction_AssemblyStartTime) isHistoryProofConstruction_ProofState() {}

func (*HistoryProofConstruction_TargetProof) isHistoryProofConstruction_ProofState() {}

func (*HistoryProofConstruction_FailureReason) isHistoryProofConstruction_ProofState() {}

// *
// A construction-scoped node id.
type ConstructionNodeId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The unique id of a history proof construction.
	ConstructionId uint64 `protobuf:"varint,1,opt,name=construction_id,json=constructionId,proto3" json:"construction_id,omitempty"`
	// *
	// The unique id of a node.
	NodeId uint64 `protobuf:"varint,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *ConstructionNodeId) Reset() {
	*x = ConstructionNodeId{}
	mi := &file_history_types_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConstructionNodeId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConstructionNodeId) ProtoMessage() {}

func (x *ConstructionNodeId) ProtoReflect() protoreflect.Message {
	mi := &file_history_types_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConstructionNodeId.ProtoReflect.Descriptor instead.
func (*ConstructionNodeId) Descriptor() ([]byte, []int) {
	return file_history_types_proto_rawDescGZIP(), []int{5}
}

func (x *ConstructionNodeId) GetConstructionId() uint64 {
	if x != nil {
		return x.ConstructionId
	}
	return 0
}

func (x *ConstructionNodeId) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

// *
// A node's vote for a particular history proof; either by explicitly
// giving the proof, or by identifying a node that already voted for it.
type HistoryProofVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Vote:
	//
	//	*HistoryProofVote_Proof
	//	*HistoryProofVote_CongruentNodeId
	Vote isHistoryProofVote_Vote `protobuf_oneof:"vote"`
}

func (x *HistoryProofVote) Reset() {
	*x = HistoryProofVote{}
	mi := &file_history_types_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HistoryProofVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoryProofVote) ProtoMessage() {}

func (x *HistoryProofVote) ProtoReflect() protoreflect.Message {
	mi := &file_history_types_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoryProofVote.ProtoReflect.Descriptor instead.
func (*HistoryProofVote) Descriptor() ([]byte, []int) {
	return file_history_types_proto_rawDescGZIP(), []int{6}
}

func (m *HistoryProofVote) GetVote() isHistoryProofVote_Vote {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (x *HistoryProofVote) GetProof() *HistoryProof {
	if x, ok := x.GetVote().(*HistoryProofVote_Proof); ok {
		return x.Proof
	}
	return nil
}

func (x *HistoryProofVote) GetCongruentNodeId() uint64 {
	if x, ok := x.GetVote().(*HistoryProofVote_CongruentNodeId); ok {
		return x.CongruentNodeId
	}
	return 0
}

type isHistoryProofVote_Vote interface {
	isHistoryProofVote_Vote()
}

type HistoryProofVote_Proof struct {
	// *
	// The history proof the submitting node is voting for.
	Proof *HistoryProof `protobuf:"bytes,1,opt,name=proof,proto3,oneof"`
}

type HistoryProofVote_CongruentNodeId struct {
	// *
	// The id of another node that already voted for the exact proof
	// the submitting node is voting for.
	CongruentNodeId uint64 `protobuf:"varint,2,opt,name=congruent_node_id,json=congruentNodeId,proto3,oneof"`
}

func (*HistoryProofVote_Proof) isHistoryProofVote_Vote() {}

func (*HistoryProofVote_CongruentNodeId) isHistoryProofVote_Vote() {}

// *
// A node's signature blessing some new history.
type HistorySignature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The new history the node is signing.
	History *History `protobuf:"bytes,1,opt,name=history,proto3" json:"history,omitempty"`
	// *
	// The node's signature on the canonical serialization of
	// the new history.
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *HistorySignature) Reset() {
	*x = HistorySignature{}
	mi := &file_history_types_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HistorySignature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistorySignature) ProtoMessage() {}

func (x *HistorySignature) ProtoReflect() protoreflect.Message {
	mi := &file_history_types_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistorySignature.ProtoReflect.Descriptor instead.
func (*HistorySignature) Descriptor() ([]byte, []int) {
	return file_history_types_proto_rawDescGZIP(), []int{7}
}

func (x *HistorySignature) GetHistory() *History {
	if x != nil {
		return x.History
	}
	return nil
}

func (x *HistorySignature) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

// *
// A signature on some new history recorded at a certain time.
type RecordedHistorySignature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The time at which the signature was recorded.
	SigningTime *Timestamp `protobuf:"bytes,1,opt,name=signing_time,json=signingTime,proto3" json:"signing_time,omitempty"`
	// *
	// The signature on some new history.
	HistorySignature *HistorySignature `protobuf:"bytes,2,opt,name=history_signature,json=historySignature,proto3" json:"history_signature,omitempty"`
}

func (x *RecordedHistorySignature) Reset() {
	*x = RecordedHistorySignature{}
	mi := &file_history_types_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecordedHistorySignature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordedHistorySignature) ProtoMessage() {}

func (x *RecordedHistorySignature) ProtoReflect() protoreflect.Message {
	mi := &file_history_types_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordedHistorySignature.ProtoReflect.Descriptor instead.
func (*RecordedHistorySignature) Descriptor() ([]byte, []int) {
	return file_history_types_proto_rawDescGZIP(), []int{8}
}

func (x *RecordedHistorySignature) GetSigningTime() *Timestamp {
	if x != nil {
		return x.SigningTime
	}
	return nil
}

func (x *RecordedHistorySignature) GetHistorySignature() *HistorySignature {
	if x != nil {
		return x.HistorySignature
	}
	return nil
}

var File_history_types_proto protoreflect.FileDescriptor

var file_history_types_proto_rawDesc = []byte{
	0x0a, 0x13, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72,
	0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x2e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x1a, 0x0f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x71, 0x0a, 0x0b, 0x50, 0x72,
	0x6f, 0x6f, 0x66, 0x4b, 0x65, 0x79, 0x53, 0x65, 0x74, 0x12, 0x35, 0x0a, 0x0d, 0x61, 0x64, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0c, 0x61, 0x64, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6e, 0x65, 0x78, 0x74, 0x4b, 0x65, 0x79, 0x22, 0x35, 0x0a,
	0x08, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x4b, 0x65, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65,
	0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x22, 0x51, 0x0a, 0x07, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x12,
	0x2a, 0x0a, 0x11, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x62, 0x6f, 0x6f, 0x6b, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x42, 0x6f, 0x6f, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x8b, 0x02, 0x0a, 0x0c, 0x48, 0x69, 0x73, 0x74,
	0x6f, 0x72, 0x79, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x37, 0x0a, 0x18, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x62, 0x6f, 0x6f, 0x6b, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x15, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x42, 0x6f, 0x6f, 0x6b, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x58, 0x0a, 0x11, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6f,
	0x66, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72,
	0x79, 0x2e, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x4b, 0x65, 0x79, 0x52, 0x0f, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x52, 0x0a, 0x0e, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61,
	0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x2e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79,
	0x52, 0x0d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x70, 0x72, 0x6f, 0x6f, 0x66, 0x22, 0x8e, 0x04, 0x0a, 0x18, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72,
	0x79, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x63, 0x6f, 0x6e,
	0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x12, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52,
	0x6f, 0x73, 0x74, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x12, 0x53, 0x0a, 0x0c, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x30, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70,
	0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68, 0x69, 0x73,
	0x74, 0x6f, 0x72, 0x79, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x50, 0x72, 0x6f, 0x6f,
	0x66, 0x52, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x2c,
	0x0a, 0x12, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x72, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x52, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x12, 0x45, 0x0a, 0x15,
	0x67, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x5f, 0x65, 0x6e, 0x64,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52,
	0x12, 0x67, 0x72, 0x61, 0x63, 0x65, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x45, 0x6e, 0x64, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x42, 0x0a, 0x13, 0x61, 0x73, 0x73, 0x65, 0x6d, 0x62, 0x6c, 0x79, 0x5f,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x48, 0x00, 0x52, 0x11, 0x61, 0x73, 0x73, 0x65, 0x6d, 0x62, 0x6c, 0x79, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x55, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68, 0x69, 0x73, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x48,
	0x00, 0x52, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x27,
	0x0a, 0x0e, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0d, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72,
	0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x42, 0x0d, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x6f, 0x66,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x56, 0x0a, 0x12, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f,
	0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x92,
	0x01, 0x0a, 0x10, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x56,
	0x6f, 0x74, 0x65, 0x12, 0x48, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x30, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e,
	0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e,
	0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x50,
	0x72, 0x6f, 0x6f, 0x66, 0x48, 0x00, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x2c, 0x0a,
	0x11, 0x63, 0x6f, 0x6e, 0x67, 0x72, 0x75, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x67,
	0x72, 0x75, 0x65, 0x6e, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x76,
	0x6f, 0x74, 0x65, 0x22, 0x77, 0x0a, 0x10, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x45, 0x0a, 0x07, 0x68, 0x69, 0x73, 0x74, 0x6f,
	0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68,
	0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x48, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x07, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0xb2, 0x01, 0x0a,
	0x18, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x65, 0x64, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79,
	0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x33, 0x0a, 0x0c, 0x73, 0x69, 0x67,
	0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x0b, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x61,
	0x0a, 0x11, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x48,
	0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52,
	0x10, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68,
	0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_history_types_proto_rawDescOnce sync.Once
	file_history_types_proto_rawDescData = file_history_types_proto_rawDesc
)

func file_history_types_proto_rawDescGZIP() []byte {
	file_history_types_proto_rawDescOnce.Do(func() {
		file_history_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_history_types_proto_rawDescData)
	})
	return file_history_types_proto_rawDescData
}

var file_history_types_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_history_types_proto_goTypes = []any{
	(*ProofKeySet)(nil),              // 0: com.hedera.hapi.node.state.history.ProofKeySet
	(*ProofKey)(nil),                 // 1: com.hedera.hapi.node.state.history.ProofKey
	(*History)(nil),                  // 2: com.hedera.hapi.node.state.history.History
	(*HistoryProof)(nil),             // 3: com.hedera.hapi.node.state.history.HistoryProof
	(*HistoryProofConstruction)(nil), // 4: com.hedera.hapi.node.state.history.HistoryProofConstruction
	(*ConstructionNodeId)(nil),       // 5: com.hedera.hapi.node.state.history.ConstructionNodeId
	(*HistoryProofVote)(nil),         // 6: com.hedera.hapi.node.state.history.HistoryProofVote
	(*HistorySignature)(nil),         // 7: com.hedera.hapi.node.state.history.HistorySignature
	(*RecordedHistorySignature)(nil), // 8: com.hedera.hapi.node.state.history.RecordedHistorySignature
	(*Timestamp)(nil),                // 9: proto.Timestamp
}
var file_history_types_proto_depIdxs = []int32{
	9,  // 0: com.hedera.hapi.node.state.history.ProofKeySet.adoption_time:type_name -> proto.Timestamp
	1,  // 1: com.hedera.hapi.node.state.history.HistoryProof.target_proof_keys:type_name -> com.hedera.hapi.node.state.history.ProofKey
	2,  // 2: com.hedera.hapi.node.state.history.HistoryProof.target_history:type_name -> com.hedera.hapi.node.state.history.History
	3,  // 3: com.hedera.hapi.node.state.history.HistoryProofConstruction.source_proof:type_name -> com.hedera.hapi.node.state.history.HistoryProof
	9,  // 4: com.hedera.hapi.node.state.history.HistoryProofConstruction.grace_period_end_time:type_name -> proto.Timestamp
	9,  // 5: com.hedera.hapi.node.state.history.HistoryProofConstruction.assembly_start_time:type_name -> proto.Timestamp
	3,  // 6: com.hedera.hapi.node.state.history.HistoryProofConstruction.target_proof:type_name -> com.hedera.hapi.node.state.history.HistoryProof
	3,  // 7: com.hedera.hapi.node.state.history.HistoryProofVote.proof:type_name -> com.hedera.hapi.node.state.history.HistoryProof
	2,  // 8: com.hedera.hapi.node.state.history.HistorySignature.history:type_name -> com.hedera.hapi.node.state.history.History
	9,  // 9: com.hedera.hapi.node.state.history.RecordedHistorySignature.signing_time:type_name -> proto.Timestamp
	7,  // 10: com.hedera.hapi.node.state.history.RecordedHistorySignature.history_signature:type_name -> com.hedera.hapi.node.state.history.HistorySignature
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_history_types_proto_init() }
func file_history_types_proto_init() {
	if File_history_types_proto != nil {
		return
	}
	file_timestamp_proto_init()
	file_history_types_proto_msgTypes[4].OneofWrappers = []any{
		(*HistoryProofConstruction_GracePeriodEndTime)(nil),
		(*HistoryProofConstruction_AssemblyStartTime)(nil),
		(*HistoryProofConstruction_TargetProof)(nil),
		(*HistoryProofConstruction_FailureReason)(nil),
	}
	file_history_types_proto_msgTypes[6].OneofWrappers = []any{
		(*HistoryProofVote_Proof)(nil),
		(*HistoryProofVote_CongruentNodeId)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_history_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_history_types_proto_goTypes,
		DependencyIndexes: file_history_types_proto_depIdxs,
		MessageInfos:      file_history_types_proto_msgTypes,
	}.Build()
	File_history_types_proto = out.File
	file_history_types_proto_rawDesc = nil
	file_history_types_proto_goTypes = nil
	file_history_types_proto_depIdxs = nil
}
