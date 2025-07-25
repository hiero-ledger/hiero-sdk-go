//*
// # System Undelete
// A system transaction to "undo" a `systemDelete` transaction.<br/>
// This transaction is a privileged operation restricted to "system"
// accounts.
//
// > Note
// >> System undelete is defined here for a smart contract (to delete
// >> the bytecode), but was never implemented.
// >
// >> Currently, system delete and system undelete specifying a smart
// >> contract identifier SHALL return `INVALID_FILE_ID`
// >> or `MISSING_ENTITY_ID`.
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
// source: system_undelete.proto

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
// Recover a file or contract bytecode deleted from the Hedera File
// System (HFS) by a `systemDelete` transaction.
//
// > Note
// >> A system delete/undelete for a `contractID` is not supported and
// >> SHALL return `INVALID_FILE_ID` or `MISSING_ENTITY_ID`.
//
// This transaction can _only_ recover a file removed with the `systemDelete`
// transaction. A file deleted via `fileDelete` SHALL be irrecoverable.<br/>
// This transaction MUST be signed by an Hedera administrative ("system")
// account.
//
// ### What is a "system" file
// A "system" file is any file with a file number less than or equal to the
// current configuration value for `ledger.numReservedSystemEntities`,
// typically `750`.
//
// ### Block Stream Effects
// None
type SystemUndeleteTransactionBody struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Id:
	//
	//	*SystemUndeleteTransactionBody_FileID
	//	*SystemUndeleteTransactionBody_ContractID
	Id            isSystemUndeleteTransactionBody_Id `protobuf_oneof:"id"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SystemUndeleteTransactionBody) Reset() {
	*x = SystemUndeleteTransactionBody{}
	mi := &file_system_undelete_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SystemUndeleteTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SystemUndeleteTransactionBody) ProtoMessage() {}

func (x *SystemUndeleteTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_system_undelete_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SystemUndeleteTransactionBody.ProtoReflect.Descriptor instead.
func (*SystemUndeleteTransactionBody) Descriptor() ([]byte, []int) {
	return file_system_undelete_proto_rawDescGZIP(), []int{0}
}

func (x *SystemUndeleteTransactionBody) GetId() isSystemUndeleteTransactionBody_Id {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *SystemUndeleteTransactionBody) GetFileID() *FileID {
	if x != nil {
		if x, ok := x.Id.(*SystemUndeleteTransactionBody_FileID); ok {
			return x.FileID
		}
	}
	return nil
}

func (x *SystemUndeleteTransactionBody) GetContractID() *ContractID {
	if x != nil {
		if x, ok := x.Id.(*SystemUndeleteTransactionBody_ContractID); ok {
			return x.ContractID
		}
	}
	return nil
}

type isSystemUndeleteTransactionBody_Id interface {
	isSystemUndeleteTransactionBody_Id()
}

type SystemUndeleteTransactionBody_FileID struct {
	// *
	// A file identifier.
	// <p>
	// The identified file MUST exist in the HFS.<br/>
	// The identified file MUST be deleted.<br/>
	// The identified file deletion MUST be a result of a
	// `systemDelete` transaction.<br/>
	// The identified file MUST NOT be a "system" file.<br/>
	// This field is REQUIRED.
	FileID *FileID `protobuf:"bytes,1,opt,name=fileID,proto3,oneof"`
}

type SystemUndeleteTransactionBody_ContractID struct {
	// *
	// A contract identifier.
	// <p>
	// The identified contract MUST exist in network state.<br/>
	// The identified contract bytecode MUST be deleted.<br/>
	// The identified contract deletion MUST be a result of a
	// `systemDelete` transaction.
	// <p>
	// This option is _unsupported_.
	ContractID *ContractID `protobuf:"bytes,2,opt,name=contractID,proto3,oneof"`
}

func (*SystemUndeleteTransactionBody_FileID) isSystemUndeleteTransactionBody_Id() {}

func (*SystemUndeleteTransactionBody_ContractID) isSystemUndeleteTransactionBody_Id() {}

var File_system_undelete_proto protoreflect.FileDescriptor

const file_system_undelete_proto_rawDesc = "" +
	"\n" +
	"\x15system_undelete.proto\x12\x05proto\x1a\x11basic_types.proto\"\x83\x01\n" +
	"\x1dSystemUndeleteTransactionBody\x12'\n" +
	"\x06fileID\x18\x01 \x01(\v2\r.proto.FileIDH\x00R\x06fileID\x123\n" +
	"\n" +
	"contractID\x18\x02 \x01(\v2\x11.proto.ContractIDH\x00R\n" +
	"contractIDB\x04\n" +
	"\x02idB&\n" +
	"\"com.hederahashgraph.api.proto.javaP\x01b\x06proto3"

var (
	file_system_undelete_proto_rawDescOnce sync.Once
	file_system_undelete_proto_rawDescData []byte
)

func file_system_undelete_proto_rawDescGZIP() []byte {
	file_system_undelete_proto_rawDescOnce.Do(func() {
		file_system_undelete_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_system_undelete_proto_rawDesc), len(file_system_undelete_proto_rawDesc)))
	})
	return file_system_undelete_proto_rawDescData
}

var file_system_undelete_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_system_undelete_proto_goTypes = []any{
	(*SystemUndeleteTransactionBody)(nil), // 0: proto.SystemUndeleteTransactionBody
	(*FileID)(nil),                        // 1: proto.FileID
	(*ContractID)(nil),                    // 2: proto.ContractID
}
var file_system_undelete_proto_depIdxs = []int32{
	1, // 0: proto.SystemUndeleteTransactionBody.fileID:type_name -> proto.FileID
	2, // 1: proto.SystemUndeleteTransactionBody.contractID:type_name -> proto.ContractID
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_system_undelete_proto_init() }
func file_system_undelete_proto_init() {
	if File_system_undelete_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_system_undelete_proto_msgTypes[0].OneofWrappers = []any{
		(*SystemUndeleteTransactionBody_FileID)(nil),
		(*SystemUndeleteTransactionBody_ContractID)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_system_undelete_proto_rawDesc), len(file_system_undelete_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_system_undelete_proto_goTypes,
		DependencyIndexes: file_system_undelete_proto_depIdxs,
		MessageInfos:      file_system_undelete_proto_msgTypes,
	}.Build()
	File_system_undelete_proto = out.File
	file_system_undelete_proto_goTypes = nil
	file_system_undelete_proto_depIdxs = nil
}
