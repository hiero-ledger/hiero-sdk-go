// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: crypto_delete_live_hash.proto

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
// At consensus, deletes a livehash associated to the given account. The transaction must be signed
// by either the key of the owning account, or at least one of the keys associated to the livehash.
type CryptoDeleteLiveHashTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The account owning the livehash
	AccountOfLiveHash *AccountID `protobuf:"bytes,1,opt,name=accountOfLiveHash,proto3" json:"accountOfLiveHash,omitempty"`
	// *
	// The SHA-384 livehash to delete from the account
	LiveHashToDelete []byte `protobuf:"bytes,2,opt,name=liveHashToDelete,proto3" json:"liveHashToDelete,omitempty"`
}

func (x *CryptoDeleteLiveHashTransactionBody) Reset() {
	*x = CryptoDeleteLiveHashTransactionBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crypto_delete_live_hash_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CryptoDeleteLiveHashTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoDeleteLiveHashTransactionBody) ProtoMessage() {}

func (x *CryptoDeleteLiveHashTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_delete_live_hash_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoDeleteLiveHashTransactionBody.ProtoReflect.Descriptor instead.
func (*CryptoDeleteLiveHashTransactionBody) Descriptor() ([]byte, []int) {
	return file_crypto_delete_live_hash_proto_rawDescGZIP(), []int{0}
}

func (x *CryptoDeleteLiveHashTransactionBody) GetAccountOfLiveHash() *AccountID {
	if x != nil {
		return x.AccountOfLiveHash
	}
	return nil
}

func (x *CryptoDeleteLiveHashTransactionBody) GetLiveHashToDelete() []byte {
	if x != nil {
		return x.LiveHashToDelete
	}
	return nil
}

var File_crypto_delete_live_hash_proto protoreflect.FileDescriptor

var file_crypto_delete_live_hash_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f,
	0x6c, 0x69, 0x76, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x91, 0x01, 0x0a, 0x23, 0x43, 0x72,
	0x79, 0x70, 0x74, 0x6f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x69, 0x76, 0x65, 0x48, 0x61,
	0x73, 0x68, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64,
	0x79, 0x12, 0x3e, 0x0a, 0x11, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4f, 0x66, 0x4c, 0x69,
	0x76, 0x65, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x11,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4f, 0x66, 0x4c, 0x69, 0x76, 0x65, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x2a, 0x0a, 0x10, 0x6c, 0x69, 0x76, 0x65, 0x48, 0x61, 0x73, 0x68, 0x54, 0x6f, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x6c, 0x69, 0x76,
	0x65, 0x48, 0x61, 0x73, 0x68, 0x54, 0x6f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x26, 0x0a,
	0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67,
	0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a,
	0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_crypto_delete_live_hash_proto_rawDescOnce sync.Once
	file_crypto_delete_live_hash_proto_rawDescData = file_crypto_delete_live_hash_proto_rawDesc
)

func file_crypto_delete_live_hash_proto_rawDescGZIP() []byte {
	file_crypto_delete_live_hash_proto_rawDescOnce.Do(func() {
		file_crypto_delete_live_hash_proto_rawDescData = protoimpl.X.CompressGZIP(file_crypto_delete_live_hash_proto_rawDescData)
	})
	return file_crypto_delete_live_hash_proto_rawDescData
}

var file_crypto_delete_live_hash_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_crypto_delete_live_hash_proto_goTypes = []interface{}{
	(*CryptoDeleteLiveHashTransactionBody)(nil), // 0: proto.CryptoDeleteLiveHashTransactionBody
	(*AccountID)(nil), // 1: proto.AccountID
}
var file_crypto_delete_live_hash_proto_depIdxs = []int32{
	1, // 0: proto.CryptoDeleteLiveHashTransactionBody.accountOfLiveHash:type_name -> proto.AccountID
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_crypto_delete_live_hash_proto_init() }
func file_crypto_delete_live_hash_proto_init() {
	if File_crypto_delete_live_hash_proto != nil {
		return
	}
	file_basic_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_crypto_delete_live_hash_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CryptoDeleteLiveHashTransactionBody); i {
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
			RawDescriptor: file_crypto_delete_live_hash_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_crypto_delete_live_hash_proto_goTypes,
		DependencyIndexes: file_crypto_delete_live_hash_proto_depIdxs,
		MessageInfos:      file_crypto_delete_live_hash_proto_msgTypes,
	}.Build()
	File_crypto_delete_live_hash_proto = out.File
	file_crypto_delete_live_hash_proto_rawDesc = nil
	file_crypto_delete_live_hash_proto_goTypes = nil
	file_crypto_delete_live_hash_proto_depIdxs = nil
}
