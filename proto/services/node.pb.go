// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: node.proto

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
// A single address book node in the network state.
//
// Each node in the network address book SHALL represent a single actual
// consensus node that is eligible to participate in network consensus.
//
// Address book nodes SHALL NOT be _globally_ uniquely identified. A given node
// is only valid within a single realm and shard combination, so the identifier
// for a network node SHALL only be unique within a single realm and shard
// combination.
type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A consensus node identifier.
	// <p>
	// Node identifiers SHALL be globally unique for a given ledger.
	NodeId uint64 `protobuf:"varint,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	// *
	// An account identifier.
	// <p>
	// This account SHALL be owned by the entity responsible for the node.<br/>
	// This account SHALL be charged transaction fees for any transactions
	// that are submitted to the network by this node and
	// fail due diligence checks.<br/>
	// This account SHALL be paid the node portion of transaction fees
	// for transactions submitted by this node.
	AccountId *AccountID `protobuf:"bytes,2,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// *
	// A short description of the node.
	// <p>
	// This value, if set, MUST NOT exceed `transaction.maxMemoUtf8Bytes`
	// (default 100) bytes when encoded as UTF-8.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// *
	// A list of service endpoints for gossip.
	// <p>
	// These endpoints SHALL represent the published endpoints to which other
	// consensus nodes may _gossip_ transactions.<br/>
	// If the network configuration value `gossipFqdnRestricted` is set, then
	// all endpoints in this list SHALL supply only IP address.<br/>
	// If the network configuration value `gossipFqdnRestricted` is _not_ set,
	// then endpoints in this list MAY supply either IP address or FQDN, but
	// SHALL NOT supply both values for the same endpoint.<br/>
	// This list SHALL NOT be empty.<br/>
	// This list SHALL NOT contain more than `10` entries.<br/>
	// The first two entries in this list SHALL be the endpoints published to
	// all consensus nodes.<br/>
	// All other entries SHALL be reserved for future use.
	GossipEndpoint []*ServiceEndpoint `protobuf:"bytes,4,rep,name=gossip_endpoint,json=gossipEndpoint,proto3" json:"gossip_endpoint,omitempty"`
	// *
	// A list of service endpoints for client calls.
	// <p>
	// These endpoints SHALL represent the published endpoints to which clients
	// may submit transactions.<br/>
	// These endpoints SHALL specify a port.<br/>
	// Endpoints in this list MAY supply either IP address or FQDN, but SHALL
	// NOT supply both values for the same endpoint.<br/>
	// This list SHALL NOT be empty.<br/>
	// This list SHALL NOT contain more than `8` entries.
	ServiceEndpoint []*ServiceEndpoint `protobuf:"bytes,5,rep,name=service_endpoint,json=serviceEndpoint,proto3" json:"service_endpoint,omitempty"`
	// *
	// A certificate used to sign gossip events.
	// <p>
	// This value SHALL be a certificate of a type permitted for gossip
	// signatures.<br/>
	// This value SHALL be the DER encoding of the certificate presented.<br/>
	// This field is REQUIRED and MUST NOT be empty.
	GossipCaCertificate []byte `protobuf:"bytes,6,opt,name=gossip_ca_certificate,json=gossipCaCertificate,proto3" json:"gossip_ca_certificate,omitempty"`
	// *
	// A hash of the node gRPC certificate.
	// <p>
	// This value MAY be used to verify the certificate presented by the node
	// during TLS negotiation for gRPC.<br/>
	// This value SHALL be a SHA-384 hash.<br/>
	// The TLS certificate to be hashed SHALL first be in PEM format and SHALL
	// be encoded with UTF-8 NFKD encoding to a stream of bytes provided to
	// the hash algorithm.<br/>
	// This field is OPTIONAL.
	GrpcCertificateHash []byte `protobuf:"bytes,7,opt,name=grpc_certificate_hash,json=grpcCertificateHash,proto3" json:"grpc_certificate_hash,omitempty"`
	// *
	// A consensus weight.
	// <p>
	// Each node SHALL have a weight in consensus calculations.<br/>
	// The consensus weight of a node SHALL be calculated based on the amount
	// of HBAR staked to that node.<br/>
	// Consensus SHALL be calculated based on agreement of greater than `2/3`
	// of the total `weight` value of all nodes on the network.
	// <p>
	// This field is deprecated and SHALL NOT be used when RosterLifecycle
	// is enabled.
	//
	// Deprecated: Marked as deprecated in node.proto.
	Weight uint64 `protobuf:"varint,8,opt,name=weight,proto3" json:"weight,omitempty"`
	// *
	// A flag indicating this node is deleted.
	// <p>
	// If this field is set, then this node SHALL NOT be included in the next
	// update of the network address book.<br/>
	// If this field is set, then this node SHALL be immutable and SHALL NOT
	// be modified.<br/>
	// If this field is set, then any `nodeUpdate` transaction to modify this
	// node SHALL fail.
	Deleted bool `protobuf:"varint,9,opt,name=deleted,proto3" json:"deleted,omitempty"`
	// *
	// An administrative key controlled by the node operator.
	// <p>
	// This key MUST sign each transaction to update this node.<br/>
	// This field MUST contain a valid `Key` value.<br/>
	// This field is REQUIRED and MUST NOT be set to an empty `KeyList`.
	AdminKey *Key `protobuf:"bytes,10,opt,name=admin_key,json=adminKey,proto3" json:"admin_key,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	mi := &file_node_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_node_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_node_proto_rawDescGZIP(), []int{0}
}

func (x *Node) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *Node) GetAccountId() *AccountID {
	if x != nil {
		return x.AccountId
	}
	return nil
}

func (x *Node) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Node) GetGossipEndpoint() []*ServiceEndpoint {
	if x != nil {
		return x.GossipEndpoint
	}
	return nil
}

func (x *Node) GetServiceEndpoint() []*ServiceEndpoint {
	if x != nil {
		return x.ServiceEndpoint
	}
	return nil
}

func (x *Node) GetGossipCaCertificate() []byte {
	if x != nil {
		return x.GossipCaCertificate
	}
	return nil
}

func (x *Node) GetGrpcCertificateHash() []byte {
	if x != nil {
		return x.GrpcCertificateHash
	}
	return nil
}

// Deprecated: Marked as deprecated in node.proto.
func (x *Node) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *Node) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

func (x *Node) GetAdminKey() *Key {
	if x != nil {
		return x.AdminKey
	}
	return nil
}

var File_node_proto protoreflect.FileDescriptor

var file_node_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x26, 0x63, 0x6f,
	0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x62, 0x6f, 0x6f, 0x6b, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbd, 0x03, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65,
	0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x0a, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52,
	0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3f, 0x0a, 0x0f,
	0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x0e, 0x67,
	0x6f, 0x73, 0x73, 0x69, 0x70, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x41, 0x0a,
	0x10, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52,
	0x0f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x12, 0x32, 0x0a, 0x15, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x5f, 0x63, 0x61, 0x5f, 0x63, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x13, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x43, 0x61, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x12, 0x32, 0x0a, 0x15, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x13, 0x67, 0x72, 0x70, 0x63, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1a, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x18, 0x01, 0x52, 0x06, 0x77, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x27,
	0x0a, 0x09, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x08, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x4b, 0x65, 0x79, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68,
	0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_proto_rawDescOnce sync.Once
	file_node_proto_rawDescData = file_node_proto_rawDesc
)

func file_node_proto_rawDescGZIP() []byte {
	file_node_proto_rawDescOnce.Do(func() {
		file_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_proto_rawDescData)
	})
	return file_node_proto_rawDescData
}

var file_node_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_node_proto_goTypes = []any{
	(*Node)(nil),            // 0: com.hedera.hapi.node.state.addressbook.Node
	(*AccountID)(nil),       // 1: proto.AccountID
	(*ServiceEndpoint)(nil), // 2: proto.ServiceEndpoint
	(*Key)(nil),             // 3: proto.Key
}
var file_node_proto_depIdxs = []int32{
	1, // 0: com.hedera.hapi.node.state.addressbook.Node.account_id:type_name -> proto.AccountID
	2, // 1: com.hedera.hapi.node.state.addressbook.Node.gossip_endpoint:type_name -> proto.ServiceEndpoint
	2, // 2: com.hedera.hapi.node.state.addressbook.Node.service_endpoint:type_name -> proto.ServiceEndpoint
	3, // 3: com.hedera.hapi.node.state.addressbook.Node.admin_key:type_name -> proto.Key
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_node_proto_init() }
func file_node_proto_init() {
	if File_node_proto != nil {
		return
	}
	file_basic_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_node_proto_goTypes,
		DependencyIndexes: file_node_proto_depIdxs,
		MessageInfos:      file_node_proto_msgTypes,
	}.Build()
	File_node_proto = out.File
	file_node_proto_rawDesc = nil
	file_node_proto_goTypes = nil
	file_node_proto_depIdxs = nil
}
