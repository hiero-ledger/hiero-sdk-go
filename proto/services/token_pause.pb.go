//*
// # Token Pause
// A transaction to "pause" all activity for a token. While a token is paused
// it cannot be transferred between accounts by any transaction other than
// `rejectToken`.
//
// ### Keywords
// The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
// "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
// document are to be interpreted as described in
// [RFC2119](https://www.ietf.org/rfc/rfc2119) and clarified in
// [RFC8174](https://www.ietf.org/rfc/rfc8174).

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: token_pause.proto

package services

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// *
// Pause transaction activity for a token.
//
// This transaction MUST be signed by the Token `pause_key`.<br/>
// The `token` identified MUST exist, and MUST NOT be deleted.<br/>
// The `token` identified MAY be paused; if the token is already paused,
// this transaction SHALL have no effect.
// The `token` identified MUST have a `pause_key` set, the `pause_key` MUST be
// a valid `Key`, and the `pause_key` MUST NOT be an empty `KeyList`.<br/>
// A `paused` token SHALL NOT be transferred or otherwise modified except to
// "up-pause" the token with `unpauseToken` or in a `rejectToken` transaction.
//
// ### Block Stream Effects
// None
type TokenPauseTransactionBody struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// *
	// A token identifier.
	// <p>
	// The identified token SHALL be paused. Subsequent transactions
	// involving that token SHALL fail until the token is "unpaused".
	Token         *TokenID `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TokenPauseTransactionBody) Reset() {
	*x = TokenPauseTransactionBody{}
	mi := &file_token_pause_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenPauseTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenPauseTransactionBody) ProtoMessage() {}

func (x *TokenPauseTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_token_pause_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenPauseTransactionBody.ProtoReflect.Descriptor instead.
func (*TokenPauseTransactionBody) Descriptor() ([]byte, []int) {
	return file_token_pause_proto_rawDescGZIP(), []int{0}
}

func (x *TokenPauseTransactionBody) GetToken() *TokenID {
	if x != nil {
		return x.Token
	}
	return nil
}

var File_token_pause_proto protoreflect.FileDescriptor

const file_token_pause_proto_rawDesc = "" +
	"\n" +
	"\x11token_pause.proto\x12\x05proto\x1a\x11basic_types.proto\"A\n" +
	"\x19TokenPauseTransactionBody\x12$\n" +
	"\x05token\x18\x01 \x01(\v2\x0e.proto.TokenIDR\x05tokenB&\n" +
	"\"com.hederahashgraph.api.proto.javaP\x01b\x06proto3"

var (
	file_token_pause_proto_rawDescOnce sync.Once
	file_token_pause_proto_rawDescData []byte
)

func file_token_pause_proto_rawDescGZIP() []byte {
	file_token_pause_proto_rawDescOnce.Do(func() {
		file_token_pause_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_token_pause_proto_rawDesc), len(file_token_pause_proto_rawDesc)))
	})
	return file_token_pause_proto_rawDescData
}

var file_token_pause_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_token_pause_proto_goTypes = []any{
	(*TokenPauseTransactionBody)(nil), // 0: proto.TokenPauseTransactionBody
	(*TokenID)(nil),                   // 1: proto.TokenID
}
var file_token_pause_proto_depIdxs = []int32{
	1, // 0: proto.TokenPauseTransactionBody.token:type_name -> proto.TokenID
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_token_pause_proto_init() }
func file_token_pause_proto_init() {
	if File_token_pause_proto != nil {
		return
	}
	file_basic_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_token_pause_proto_rawDesc), len(file_token_pause_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_token_pause_proto_goTypes,
		DependencyIndexes: file_token_pause_proto_depIdxs,
		MessageInfos:      file_token_pause_proto_msgTypes,
	}.Build()
	File_token_pause_proto = out.File
	file_token_pause_proto_goTypes = nil
	file_token_pause_proto_depIdxs = nil
}
