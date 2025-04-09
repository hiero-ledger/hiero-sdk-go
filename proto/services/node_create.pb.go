// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: node_create.proto

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
// A transaction body to add a new consensus node to the network address book.
//
// This transaction body SHALL be considered a "privileged transaction".
//
// This message supports a transaction to create a new node in the network
// address book. The transaction, once complete, enables a new consensus node
// to join the network, and requires governing council authorization.
//
//   - A `NodeCreateTransactionBody` MUST be signed by the `Key` assigned to the
//     `admin_key` field and one of those keys: treasure account (2) key,
//     systemAdmin(50) key, or addressBookAdmin(55) key.
//   - The newly created node information SHALL be added to the network address
//     book information in the network state.
//   - The new entry SHALL be created in "state" but SHALL NOT participate in
//     network consensus and SHALL NOT be present in network "configuration"
//     until the next "upgrade" transaction (as noted below).
//   - All new address book entries SHALL be added to the active network
//     configuration during the next `freeze` transaction with the field
//     `freeze_type` set to `PREPARE_UPGRADE`.
//
// ### Block Stream Effects
// Upon completion the newly assigned `node_id` SHALL be recorded in
// the transaction receipt.<br/>
// This value SHALL be the next available node identifier.<br/>
// Node identifiers SHALL NOT be reused.
type NodeCreateTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A Node account identifier.
	// <p>
	// This account identifier MUST be in the "account number" form.<br/>
	// This account identifier MUST NOT use the alias field.<br/>
	// If the identified account does not exist, this transaction SHALL fail.<br/>
	// Multiple nodes MAY share the same node account.<br/>
	// This field is REQUIRED.
	AccountId *AccountID `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// *
	// A short description of the node.
	// <p>
	// This value, if set, MUST NOT exceed `transaction.maxMemoUtf8Bytes`
	// (default 100) bytes when encoded as UTF-8.<br/>
	// This field is OPTIONAL.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// *
	// A list of service endpoints for gossip.
	// <p>
	// These endpoints SHALL represent the published endpoints to which other
	// consensus nodes may _gossip_ transactions.<br/>
	// These endpoints MUST specify a port.<br/>
	// This list MUST NOT be empty.<br/>
	// This list MUST NOT contain more than `10` entries.<br/>
	// The first two entries in this list SHALL be the endpoints published to
	// all consensus nodes.<br/>
	// All other entries SHALL be reserved for future use.
	// <p>
	// Each network may have additional requirements for these endpoints.
	// A client MUST check network-specific documentation for those
	// details.<br/>
	// If the network configuration value `gossipFqdnRestricted` is set, then
	// all endpoints in this list MUST supply only IP address.<br/>
	// If the network configuration value `gossipFqdnRestricted` is _not_ set,
	// then endpoints in this list MAY supply either IP address or FQDN, but
	// MUST NOT supply both values for the same endpoint.
	GossipEndpoint []*ServiceEndpoint `protobuf:"bytes,3,rep,name=gossip_endpoint,json=gossipEndpoint,proto3" json:"gossip_endpoint,omitempty"`
	// *
	// A list of service endpoints for gRPC calls.
	// <p>
	// These endpoints SHALL represent the published gRPC endpoints to which
	// clients may submit transactions.<br/>
	// These endpoints MUST specify a port.<br/>
	// Endpoints in this list MAY supply either IP address or FQDN, but MUST
	// NOT supply both values for the same endpoint.<br/>
	// This list MUST NOT be empty.<br/>
	// This list MUST NOT contain more than `8` entries.
	ServiceEndpoint []*ServiceEndpoint `protobuf:"bytes,4,rep,name=service_endpoint,json=serviceEndpoint,proto3" json:"service_endpoint,omitempty"`
	// *
	// A certificate used to sign gossip events.
	// <p>
	// This value MUST be a certificate of a type permitted for gossip
	// signatures.<br/>
	// This value MUST be the DER encoding of the certificate presented.<br/>
	// This field is REQUIRED and MUST NOT be empty.
	GossipCaCertificate []byte `protobuf:"bytes,5,opt,name=gossip_ca_certificate,json=gossipCaCertificate,proto3" json:"gossip_ca_certificate,omitempty"`
	// *
	// A hash of the node gRPC TLS certificate.
	// <p>
	// This value MAY be used to verify the certificate presented by the node
	// during TLS negotiation for gRPC.<br/>
	// This value MUST be a SHA-384 hash.<br/>
	// The TLS certificate to be hashed MUST first be in PEM format and MUST be
	// encoded with UTF-8 NFKD encoding to a stream of bytes provided to
	// the hash algorithm.<br/>
	// This field is OPTIONAL.
	GrpcCertificateHash []byte `protobuf:"bytes,6,opt,name=grpc_certificate_hash,json=grpcCertificateHash,proto3" json:"grpc_certificate_hash,omitempty"`
	// *
	// An administrative key controlled by the node operator.
	// <p>
	// This key MUST sign this transaction.<br/>
	// This key MUST sign each transaction to update this node.<br/>
	// This field MUST contain a valid `Key` value.<br/>
	// This field is REQUIRED and MUST NOT be set to an empty `KeyList`.
	AdminKey *Key `protobuf:"bytes,7,opt,name=admin_key,json=adminKey,proto3" json:"admin_key,omitempty"`
	// *
	// A boolean flag indicating whether the node operator declines to receive
	// node rewards.
	// <p>
	// If this flag is set to `true`, the node operator declines to receive
	// node rewards.<br/>
	DeclineReward bool `protobuf:"varint,8,opt,name=decline_reward,json=declineReward,proto3" json:"decline_reward,omitempty"`
}

func (x *NodeCreateTransactionBody) Reset() {
	*x = NodeCreateTransactionBody{}
	mi := &file_node_create_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NodeCreateTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeCreateTransactionBody) ProtoMessage() {}

func (x *NodeCreateTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_node_create_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeCreateTransactionBody.ProtoReflect.Descriptor instead.
func (*NodeCreateTransactionBody) Descriptor() ([]byte, []int) {
	return file_node_create_proto_rawDescGZIP(), []int{0}
}

func (x *NodeCreateTransactionBody) GetAccountId() *AccountID {
	if x != nil {
		return x.AccountId
	}
	return nil
}

func (x *NodeCreateTransactionBody) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *NodeCreateTransactionBody) GetGossipEndpoint() []*ServiceEndpoint {
	if x != nil {
		return x.GossipEndpoint
	}
	return nil
}

func (x *NodeCreateTransactionBody) GetServiceEndpoint() []*ServiceEndpoint {
	if x != nil {
		return x.ServiceEndpoint
	}
	return nil
}

func (x *NodeCreateTransactionBody) GetGossipCaCertificate() []byte {
	if x != nil {
		return x.GossipCaCertificate
	}
	return nil
}

func (x *NodeCreateTransactionBody) GetGrpcCertificateHash() []byte {
	if x != nil {
		return x.GrpcCertificateHash
	}
	return nil
}

func (x *NodeCreateTransactionBody) GetAdminKey() *Key {
	if x != nil {
		return x.AdminKey
	}
	return nil
}

func (x *NodeCreateTransactionBody) GetDeclineReward() bool {
	if x != nil {
		return x.DeclineReward
	}
	return false
}

var File_node_create_proto protoreflect.FileDescriptor

var file_node_create_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x20, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e,
	0x68, 0x61, 0x70, 0x69, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x62, 0x6f, 0x6f, 0x6b, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaa, 0x03, 0x0a, 0x19, 0x4e, 0x6f, 0x64,
	0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x2f, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x09, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3f, 0x0a, 0x0f, 0x67, 0x6f, 0x73,
	0x73, 0x69, 0x70, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x0e, 0x67, 0x6f, 0x73, 0x73,
	0x69, 0x70, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x41, 0x0a, 0x10, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x0f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x32, 0x0a,
	0x15, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x5f, 0x63, 0x61, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x13, 0x67, 0x6f,
	0x73, 0x73, 0x69, 0x70, 0x43, 0x61, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x12, 0x32, 0x0a, 0x15, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x13, 0x67, 0x72, 0x70, 0x63, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x27, 0x0a, 0x09, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4b, 0x65, 0x79, 0x52, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x25,
	0x0a, 0x0e, 0x64, 0x65, 0x63, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x72, 0x65, 0x77, 0x61, 0x72, 0x64,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x64, 0x65, 0x63, 0x6c, 0x69, 0x6e, 0x65, 0x52,
	0x65, 0x77, 0x61, 0x72, 0x64, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64,
	0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_node_create_proto_rawDescOnce sync.Once
	file_node_create_proto_rawDescData = file_node_create_proto_rawDesc
)

func file_node_create_proto_rawDescGZIP() []byte {
	file_node_create_proto_rawDescOnce.Do(func() {
		file_node_create_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_create_proto_rawDescData)
	})
	return file_node_create_proto_rawDescData
}

var file_node_create_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_node_create_proto_goTypes = []any{
	(*NodeCreateTransactionBody)(nil), // 0: com.hedera.hapi.node.addressbook.NodeCreateTransactionBody
	(*AccountID)(nil),                 // 1: proto.AccountID
	(*ServiceEndpoint)(nil),           // 2: proto.ServiceEndpoint
	(*Key)(nil),                       // 3: proto.Key
}
var file_node_create_proto_depIdxs = []int32{
	1, // 0: com.hedera.hapi.node.addressbook.NodeCreateTransactionBody.account_id:type_name -> proto.AccountID
	2, // 1: com.hedera.hapi.node.addressbook.NodeCreateTransactionBody.gossip_endpoint:type_name -> proto.ServiceEndpoint
	2, // 2: com.hedera.hapi.node.addressbook.NodeCreateTransactionBody.service_endpoint:type_name -> proto.ServiceEndpoint
	3, // 3: com.hedera.hapi.node.addressbook.NodeCreateTransactionBody.admin_key:type_name -> proto.Key
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_node_create_proto_init() }
func file_node_create_proto_init() {
	if File_node_create_proto != nil {
		return
	}
	file_basic_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_node_create_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_node_create_proto_goTypes,
		DependencyIndexes: file_node_create_proto_depIdxs,
		MessageInfos:      file_node_create_proto_msgTypes,
	}.Build()
	File_node_create_proto = out.File
	file_node_create_proto_rawDesc = nil
	file_node_create_proto_goTypes = nil
	file_node_create_proto_depIdxs = nil
}
