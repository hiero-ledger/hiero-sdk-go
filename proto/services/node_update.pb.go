// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: node_update.proto

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
// Transaction body to modify address book node attributes.
//
//   - This transaction SHALL enable the node operator, as identified by the
//     `admin_key`, to modify operational attributes of the node.
//   - This transaction MUST be signed by the active `admin_key` for the node.
//   - If this transaction sets a new value for the `admin_key`, then both the
//     current `admin_key`, and the new `admin_key` MUST sign this transaction.
//   - This transaction SHALL NOT change any field that is not set (is null) in
//     this transaction body.
//   - This SHALL create a pending update to the node, but the change SHALL NOT
//     be immediately applied to the active configuration.
//   - All pending node updates SHALL be applied to the active network
//     configuration during the next `freeze` transaction with the field
//     `freeze_type` set to `PREPARE_UPGRADE`.
//
// ### Block Stream Effects
// None.
type NodeUpdateTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A consensus node identifier in the network state.
	// <p>
	// The node identified MUST exist in the network address book.<br/>
	// The node identified MUST NOT be deleted.<br/>
	// This value is REQUIRED.
	NodeId uint64 `protobuf:"varint,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	// *
	// An account identifier.
	// <p>
	// If set, this SHALL replace the node account identifier.<br/>
	// If set, this transaction MUST be signed by the active `key` for _both_
	// the current node account _and_ the identified new node account.
	AccountId *AccountID `protobuf:"bytes,2,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// *
	// A short description of the node.
	// <p>
	// This value, if set, MUST NOT exceed `transaction.maxMemoUtf8Bytes`
	// (default 100) bytes when encoded as UTF-8.<br/>
	// If set, this value SHALL replace the previous value.
	Description *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// *
	// A list of service endpoints for gossip.
	// <p>
	// If set, this list MUST meet the following requirements.
	// <hr/>
	// These endpoints SHALL represent the published endpoints to which other
	// consensus nodes may _gossip_ transactions.<br/>
	// These endpoints SHOULD NOT specify both address and DNS name.<br/>
	// This list MUST NOT be empty.<br/>
	// This list MUST NOT contain more than `10` entries.<br/>
	// The first two entries in this list SHALL be the endpoints published to
	// all consensus nodes.<br/>
	// All other entries SHALL be reserved for future use.
	// <p>
	// Each network may have additional requirements for these endpoints.
	// A client MUST check network-specific documentation for those
	// details.<br/>
	// <blockquote>Example<blockquote>
	// Hedera Mainnet _requires_ that address be specified, and does not
	// permit DNS name (FQDN) to be specified.
	// </blockquote>
	// <blockquote>
	// Solo, however, _requires_ DNS name (FQDN) but also permits
	// address.
	// </blockquote></blockquote>
	// <p>
	// If set, the new list SHALL replace the existing list.
	GossipEndpoint []*ServiceEndpoint `protobuf:"bytes,4,rep,name=gossip_endpoint,json=gossipEndpoint,proto3" json:"gossip_endpoint,omitempty"`
	// *
	// A list of service endpoints for gRPC calls.
	// <p>
	// If set, this list MUST meet the following requirements.
	// <hr/>
	// These endpoints SHALL represent the published endpoints to which clients
	// may submit transactions.<br/>
	// These endpoints SHOULD specify address and port.<br/>
	// These endpoints MAY specify a DNS name.<br/>
	// These endpoints SHOULD NOT specify both address and DNS name.<br/>
	// This list MUST NOT be empty.<br/>
	// This list MUST NOT contain more than `8` entries.
	// <p>
	// Each network may have additional requirements for these endpoints.
	// A client MUST check network-specific documentation for those
	// details.
	// <p>
	// If set, the new list SHALL replace the existing list.
	ServiceEndpoint []*ServiceEndpoint `protobuf:"bytes,5,rep,name=service_endpoint,json=serviceEndpoint,proto3" json:"service_endpoint,omitempty"`
	// *
	// A certificate used to sign gossip events.
	// <p>
	// This value MUST be a certificate of a type permitted for gossip
	// signatures.<br/>
	// This value MUST be the DER encoding of the certificate presented.
	// <p>
	// If set, the new value SHALL replace the existing bytes value.
	GossipCaCertificate *wrapperspb.BytesValue `protobuf:"bytes,6,opt,name=gossip_ca_certificate,json=gossipCaCertificate,proto3" json:"gossip_ca_certificate,omitempty"`
	// *
	// A hash of the node gRPC TLS certificate.
	// <p>
	// This value MAY be used to verify the certificate presented by the node
	// during TLS negotiation for gRPC.<br/>
	// This value MUST be a SHA-384 hash.<br/>
	// The TLS certificate to be hashed MUST first be in PEM format and MUST be
	// encoded with UTF-8 NFKD encoding to a stream of bytes provided to
	// the hash algorithm.<br/>
	// <p>
	// If set, the new value SHALL replace the existing hash value.
	GrpcCertificateHash *wrapperspb.BytesValue `protobuf:"bytes,7,opt,name=grpc_certificate_hash,json=grpcCertificateHash,proto3" json:"grpc_certificate_hash,omitempty"`
	// *
	// An administrative key controlled by the node operator.
	// <p>
	// This field is OPTIONAL.<br/>
	// If set, this key MUST sign this transaction.<br/>
	// If set, this key MUST sign each subsequent transaction to
	// update this node.<br/>
	// If set, this field MUST contain a valid `Key` value.<br/>
	// If set, this field MUST NOT be set to an empty `KeyList`.
	AdminKey *Key `protobuf:"bytes,8,opt,name=admin_key,json=adminKey,proto3" json:"admin_key,omitempty"`
	// *
	// A boolean indicating that this node has chosen to decline node rewards
	// distributed at the end of staking period.
	// <p>
	// This node SHALL NOT receive reward if this value is set, and `true`.
	DeclineReward *wrapperspb.BoolValue `protobuf:"bytes,9,opt,name=decline_reward,json=declineReward,proto3" json:"decline_reward,omitempty"`
	// *
	// A web proxy for gRPC from non-gRPC clients.
	// <p>
	// This endpoint SHALL be a Fully Qualified Domain Name (FQDN) using the HTTPS
	// protocol, and SHALL support gRPC-Web for use by browser-based clients.<br/>
	// This endpoint MUST be signed by a trusted certificate authority.<br/>
	// This endpoint MUST use a valid port and SHALL be reachable over TLS.<br/>
	// This field MAY be omitted if the node does not support gRPC-Web access.<br/>
	// This field MUST be updated if the gRPC-Web endpoint changes.<br/>
	// This field SHALL enable frontend clients to avoid hard-coded proxy endpoints.
	GrpcProxyEndpoint *ServiceEndpoint `protobuf:"bytes,10,opt,name=grpc_proxy_endpoint,json=grpcProxyEndpoint,proto3" json:"grpc_proxy_endpoint,omitempty"`
}

func (x *NodeUpdateTransactionBody) Reset() {
	*x = NodeUpdateTransactionBody{}
	mi := &file_node_update_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NodeUpdateTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeUpdateTransactionBody) ProtoMessage() {}

func (x *NodeUpdateTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_node_update_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeUpdateTransactionBody.ProtoReflect.Descriptor instead.
func (*NodeUpdateTransactionBody) Descriptor() ([]byte, []int) {
	return file_node_update_proto_rawDescGZIP(), []int{0}
}

func (x *NodeUpdateTransactionBody) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *NodeUpdateTransactionBody) GetAccountId() *AccountID {
	if x != nil {
		return x.AccountId
	}
	return nil
}

func (x *NodeUpdateTransactionBody) GetDescription() *wrapperspb.StringValue {
	if x != nil {
		return x.Description
	}
	return nil
}

func (x *NodeUpdateTransactionBody) GetGossipEndpoint() []*ServiceEndpoint {
	if x != nil {
		return x.GossipEndpoint
	}
	return nil
}

func (x *NodeUpdateTransactionBody) GetServiceEndpoint() []*ServiceEndpoint {
	if x != nil {
		return x.ServiceEndpoint
	}
	return nil
}

func (x *NodeUpdateTransactionBody) GetGossipCaCertificate() *wrapperspb.BytesValue {
	if x != nil {
		return x.GossipCaCertificate
	}
	return nil
}

func (x *NodeUpdateTransactionBody) GetGrpcCertificateHash() *wrapperspb.BytesValue {
	if x != nil {
		return x.GrpcCertificateHash
	}
	return nil
}

func (x *NodeUpdateTransactionBody) GetAdminKey() *Key {
	if x != nil {
		return x.AdminKey
	}
	return nil
}

func (x *NodeUpdateTransactionBody) GetDeclineReward() *wrapperspb.BoolValue {
	if x != nil {
		return x.DeclineReward
	}
	return nil
}

func (x *NodeUpdateTransactionBody) GetGrpcWebProxyEndpoint() *ServiceEndpoint {
	if x != nil {
		return x.GrpcProxyEndpoint
	}
	return nil
}

var File_node_update_proto protoreflect.FileDescriptor

var file_node_update_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x20, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e,
	0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x62, 0x6f, 0x6f, 0x6b, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xff, 0x04, 0x0a, 0x19, 0x4e, 0x6f, 0x64,
	0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12,
	0x2f, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x3e, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x3f, 0x0a, 0x0f, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x52, 0x0e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x12, 0x41, 0x0a, 0x10, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x52, 0x0f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x64, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x12, 0x4f, 0x0a, 0x15, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x5f, 0x63,
	0x61, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x13, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x43, 0x61, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x4f, 0x0a, 0x15, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x13, 0x67, 0x72, 0x70, 0x63, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x27, 0x0a, 0x09, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f,
	0x6b, 0x65, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x4b, 0x65, 0x79, 0x12,
	0x41, 0x0a, 0x0e, 0x64, 0x65, 0x63, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x72, 0x65, 0x77, 0x61, 0x72,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x0d, 0x64, 0x65, 0x63, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x77, 0x61,
	0x72, 0x64, 0x12, 0x46, 0x0a, 0x13, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x70, 0x72, 0x6f, 0x78, 0x79,
	0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x45,
	0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x11, 0x67, 0x72, 0x70, 0x63, 0x50, 0x72, 0x6f,
	0x78, 0x79, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f,
	0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70,
	0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61,
	0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_update_proto_rawDescOnce sync.Once
	file_node_update_proto_rawDescData = file_node_update_proto_rawDesc
)

func file_node_update_proto_rawDescGZIP() []byte {
	file_node_update_proto_rawDescOnce.Do(func() {
		file_node_update_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_update_proto_rawDescData)
	})
	return file_node_update_proto_rawDescData
}

var file_node_update_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_node_update_proto_goTypes = []any{
	(*NodeUpdateTransactionBody)(nil), // 0: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody
	(*AccountID)(nil),                 // 1: proto.AccountID
	(*wrapperspb.StringValue)(nil),    // 2: google.protobuf.StringValue
	(*ServiceEndpoint)(nil),           // 3: proto.ServiceEndpoint
	(*wrapperspb.BytesValue)(nil),     // 4: google.protobuf.BytesValue
	(*Key)(nil),                       // 5: proto.Key
	(*wrapperspb.BoolValue)(nil),      // 6: google.protobuf.BoolValue
}
var file_node_update_proto_depIdxs = []int32{
	1, // 0: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody.account_id:type_name -> proto.AccountID
	2, // 1: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody.description:type_name -> google.protobuf.StringValue
	3, // 2: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody.gossip_endpoint:type_name -> proto.ServiceEndpoint
	3, // 3: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody.service_endpoint:type_name -> proto.ServiceEndpoint
	4, // 4: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody.gossip_ca_certificate:type_name -> google.protobuf.BytesValue
	4, // 5: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody.grpc_certificate_hash:type_name -> google.protobuf.BytesValue
	5, // 6: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody.admin_key:type_name -> proto.Key
	6, // 7: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody.decline_reward:type_name -> google.protobuf.BoolValue
	3, // 8: com.hedera.hapi.node.addressbook.NodeUpdateTransactionBody.grpc_proxy_endpoint:type_name -> proto.ServiceEndpoint
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_node_update_proto_init() }
func file_node_update_proto_init() {
	if File_node_update_proto != nil {
		return
	}
	file_basic_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_node_update_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_node_update_proto_goTypes,
		DependencyIndexes: file_node_update_proto_depIdxs,
		MessageInfos:      file_node_update_proto_msgTypes,
	}.Build()
	File_node_update_proto = out.File
	file_node_update_proto_rawDesc = nil
	file_node_update_proto_goTypes = nil
	file_node_update_proto_depIdxs = nil
}
