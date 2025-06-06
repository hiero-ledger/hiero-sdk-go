//*
// # File Append
// A transaction body message to append data to a "file" in state.
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
// source: file_append.proto

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
// A transaction body for an `appendContent` transaction.<br/>
// This transaction body provides a mechanism to append content to a "file" in
// network state. Hedera transactions are limited in size, but there are many
// uses for in-state byte arrays (e.g. smart contract bytecode) which require
// more than may fit within a single transaction. The `appendFile` transaction
// exists to support these requirements. The typical pattern is to create a
// file, append more data until the full content is stored, verify the file is
// correct, then update the file entry with any final metadata changes (e.g.
// adding threshold keys and removing the initial upload key).
//
// Each append transaction MUST remain within the total transaction size limit
// for the network (typically 6144 bytes).<br/>
// The total size of a file MUST remain within the maximum file size limit for
// the network (typically 1048576 bytes).
//
// #### Signature Requirements
// Append transactions MUST have signatures from _all_ keys in the `KeyList`
// assigned to the `keys` field of the file.<br/>
// See the [File Service](#FileService) specification for a detailed
// explanation of the signature requirements for all file transactions.
//
// ### Block Stream Effects
// None
type FileAppendTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A file identifier.<br/>
	// This identifies the file to which the `contents` will be appended.
	// <p>
	// This field is REQUIRED.<br/>
	// The identified file MUST exist.<br/>
	// The identified file MUST NOT be larger than the current maximum file
	// size limit.<br/>
	// The identified file MUST NOT be deleted.<br/>
	// The identified file MUST NOT be immutable.
	FileID *FileID `protobuf:"bytes,2,opt,name=fileID,proto3" json:"fileID,omitempty"`
	// *
	// An array of bytes to append.<br/>
	// <p>
	// This content SHALL be appended to the identified file if this
	// transaction succeeds.<br/>
	// This field is REQUIRED.<br/>
	// This field MUST NOT be empty.
	Contents []byte `protobuf:"bytes,4,opt,name=contents,proto3" json:"contents,omitempty"`
}

func (x *FileAppendTransactionBody) Reset() {
	*x = FileAppendTransactionBody{}
	mi := &file_file_append_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileAppendTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileAppendTransactionBody) ProtoMessage() {}

func (x *FileAppendTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_file_append_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileAppendTransactionBody.ProtoReflect.Descriptor instead.
func (*FileAppendTransactionBody) Descriptor() ([]byte, []int) {
	return file_file_append_proto_rawDescGZIP(), []int{0}
}

func (x *FileAppendTransactionBody) GetFileID() *FileID {
	if x != nil {
		return x.FileID
	}
	return nil
}

func (x *FileAppendTransactionBody) GetContents() []byte {
	if x != nil {
		return x.Contents
	}
	return nil
}

var File_file_append_proto protoreflect.FileDescriptor

var file_file_append_proto_rawDesc = []byte{
	0x0a, 0x11, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69,
	0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5e, 0x0a,
	0x19, 0x46, 0x69, 0x6c, 0x65, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x25, 0x0a, 0x06, 0x66, 0x69,
	0x6c, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49,
	0x44, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x42, 0x26, 0x0a,
	0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67,
	0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a,
	0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_file_append_proto_rawDescOnce sync.Once
	file_file_append_proto_rawDescData = file_file_append_proto_rawDesc
)

func file_file_append_proto_rawDescGZIP() []byte {
	file_file_append_proto_rawDescOnce.Do(func() {
		file_file_append_proto_rawDescData = protoimpl.X.CompressGZIP(file_file_append_proto_rawDescData)
	})
	return file_file_append_proto_rawDescData
}

var file_file_append_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_file_append_proto_goTypes = []any{
	(*FileAppendTransactionBody)(nil), // 0: proto.FileAppendTransactionBody
	(*FileID)(nil),                    // 1: proto.FileID
}
var file_file_append_proto_depIdxs = []int32{
	1, // 0: proto.FileAppendTransactionBody.fileID:type_name -> proto.FileID
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_file_append_proto_init() }
func file_file_append_proto_init() {
	if File_file_append_proto != nil {
		return
	}
	file_basic_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_file_append_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_file_append_proto_goTypes,
		DependencyIndexes: file_file_append_proto_depIdxs,
		MessageInfos:      file_file_append_proto_msgTypes,
	}.Build()
	File_file_append_proto = out.File
	file_file_append_proto_rawDesc = nil
	file_file_append_proto_goTypes = nil
	file_file_append_proto_depIdxs = nil
}
