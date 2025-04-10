//*
// # hinTS Key Publication Transaction
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
// source: hints_key_publication.proto

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
// A transaction body to publish a node's hinTS key for a certain
// party id and number of parties. A hinTS key is an extended
// public key; that is, a BLS public key combined with "hints"
// derived from the matching private key that a signature
// aggregator can use to prove well-formedness of an aggregate
// public key by an efficiently verifiable SNARK.
type HintsKeyPublicationTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The party id for which the hinTS key is being published;
	// must be in the range [0, num_parties).
	// <p>
	// This value MUST be set to a non-negative integer.<br/>
	PartyId uint32 `protobuf:"varint,1,opt,name=party_id,json=partyId,proto3" json:"party_id,omitempty"`
	// *
	// The number of parties in the hinTS scheme.
	NumParties uint32 `protobuf:"varint,2,opt,name=num_parties,json=numParties,proto3" json:"num_parties,omitempty"`
	// *
	// The party's hinTS key.
	HintsKey []byte `protobuf:"bytes,3,opt,name=hints_key,json=hintsKey,proto3" json:"hints_key,omitempty"`
}

func (x *HintsKeyPublicationTransactionBody) Reset() {
	*x = HintsKeyPublicationTransactionBody{}
	mi := &file_hints_key_publication_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HintsKeyPublicationTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HintsKeyPublicationTransactionBody) ProtoMessage() {}

func (x *HintsKeyPublicationTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_hints_key_publication_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HintsKeyPublicationTransactionBody.ProtoReflect.Descriptor instead.
func (*HintsKeyPublicationTransactionBody) Descriptor() ([]byte, []int) {
	return file_hints_key_publication_proto_rawDescGZIP(), []int{0}
}

func (x *HintsKeyPublicationTransactionBody) GetPartyId() uint32 {
	if x != nil {
		return x.PartyId
	}
	return 0
}

func (x *HintsKeyPublicationTransactionBody) GetNumParties() uint32 {
	if x != nil {
		return x.NumParties
	}
	return 0
}

func (x *HintsKeyPublicationTransactionBody) GetHintsKey() []byte {
	if x != nil {
		return x.HintsKey
	}
	return nil
}

var File_hints_key_publication_proto protoreflect.FileDescriptor

var file_hints_key_publication_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x28, 0x63,
	0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x75, 0x78, 0x69, 0x6c, 0x69, 0x61, 0x72,
	0x79, 0x2e, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x1a, 0x11, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7d, 0x0a, 0x22, 0x48, 0x69,
	0x6e, 0x74, 0x73, 0x4b, 0x65, 0x79, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79,
	0x12, 0x19, 0x0a, 0x08, 0x70, 0x61, 0x72, 0x74, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x70, 0x61, 0x72, 0x74, 0x79, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x6e,
	0x75, 0x6d, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x6e, 0x75, 0x6d, 0x50, 0x61, 0x72, 0x74, 0x69, 0x65, 0x73, 0x12, 0x1b, 0x0a, 0x09,
	0x68, 0x69, 0x6e, 0x74, 0x73, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x08, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x4b, 0x65, 0x79, 0x42, 0x33, 0x0a, 0x2f, 0x63, 0x6f, 0x6d,
	0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x61, 0x75, 0x78, 0x69, 0x6c, 0x69, 0x61, 0x72, 0x79, 0x2e,
	0x68, 0x69, 0x6e, 0x74, 0x73, 0x2e, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x50, 0x01, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hints_key_publication_proto_rawDescOnce sync.Once
	file_hints_key_publication_proto_rawDescData = file_hints_key_publication_proto_rawDesc
)

func file_hints_key_publication_proto_rawDescGZIP() []byte {
	file_hints_key_publication_proto_rawDescOnce.Do(func() {
		file_hints_key_publication_proto_rawDescData = protoimpl.X.CompressGZIP(file_hints_key_publication_proto_rawDescData)
	})
	return file_hints_key_publication_proto_rawDescData
}

var file_hints_key_publication_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_hints_key_publication_proto_goTypes = []any{
	(*HintsKeyPublicationTransactionBody)(nil), // 0: com.hedera.hapi.services.auxiliary.hints.HintsKeyPublicationTransactionBody
}
var file_hints_key_publication_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_hints_key_publication_proto_init() }
func file_hints_key_publication_proto_init() {
	if File_hints_key_publication_proto != nil {
		return
	}
	file_hints_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_hints_key_publication_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_hints_key_publication_proto_goTypes,
		DependencyIndexes: file_hints_key_publication_proto_depIdxs,
		MessageInfos:      file_hints_key_publication_proto_msgTypes,
	}.Build()
	File_hints_key_publication_proto = out.File
	file_hints_key_publication_proto_rawDesc = nil
	file_hints_key_publication_proto_goTypes = nil
	file_hints_key_publication_proto_depIdxs = nil
}
