// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: freeze.proto

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
// At consensus, sets the consensus time at which the platform should stop creating events and
// accepting transactions, and enter a maintenance window.
type FreezeTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// !! DEPRECATED and REJECTED by nodes
	// The start hour (in UTC time), a value between 0 and 23
	//
	// Deprecated: Marked as deprecated in freeze.proto.
	StartHour int32 `protobuf:"varint,1,opt,name=startHour,proto3" json:"startHour,omitempty"`
	// *
	// !! DEPRECATED and REJECTED by nodes
	// The start minute (in UTC time), a value between 0 and 59
	//
	// Deprecated: Marked as deprecated in freeze.proto.
	StartMin int32 `protobuf:"varint,2,opt,name=startMin,proto3" json:"startMin,omitempty"`
	// *
	// !! DEPRECATED and REJECTED by nodes
	// The end hour (in UTC time), a value between 0 and 23
	//
	// Deprecated: Marked as deprecated in freeze.proto.
	EndHour int32 `protobuf:"varint,3,opt,name=endHour,proto3" json:"endHour,omitempty"`
	// *
	// !! DEPRECATED and REJECTED by nodes
	// The end minute (in UTC time), a value between 0 and 59
	//
	// Deprecated: Marked as deprecated in freeze.proto.
	EndMin int32 `protobuf:"varint,4,opt,name=endMin,proto3" json:"endMin,omitempty"`
	// *
	// If set, the file whose contents should be used for a network software update during the
	// maintenance window.
	UpdateFile *FileID `protobuf:"bytes,5,opt,name=update_file,json=updateFile,proto3" json:"update_file,omitempty"`
	// *
	// If set, the expected hash of the contents of the update file (used to verify the update).
	FileHash []byte `protobuf:"bytes,6,opt,name=file_hash,json=fileHash,proto3" json:"file_hash,omitempty"`
	// *
	// The consensus time at which the maintenance window should begin.
	StartTime *Timestamp `protobuf:"bytes,7,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// *
	// The type of network freeze or upgrade operation to perform.
	FreezeType FreezeType `protobuf:"varint,8,opt,name=freeze_type,json=freezeType,proto3,enum=proto.FreezeType" json:"freeze_type,omitempty"`
}

func (x *FreezeTransactionBody) Reset() {
	*x = FreezeTransactionBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_freeze_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FreezeTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FreezeTransactionBody) ProtoMessage() {}

func (x *FreezeTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_freeze_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FreezeTransactionBody.ProtoReflect.Descriptor instead.
func (*FreezeTransactionBody) Descriptor() ([]byte, []int) {
	return file_freeze_proto_rawDescGZIP(), []int{0}
}

// Deprecated: Marked as deprecated in freeze.proto.
func (x *FreezeTransactionBody) GetStartHour() int32 {
	if x != nil {
		return x.StartHour
	}
	return 0
}

// Deprecated: Marked as deprecated in freeze.proto.
func (x *FreezeTransactionBody) GetStartMin() int32 {
	if x != nil {
		return x.StartMin
	}
	return 0
}

// Deprecated: Marked as deprecated in freeze.proto.
func (x *FreezeTransactionBody) GetEndHour() int32 {
	if x != nil {
		return x.EndHour
	}
	return 0
}

// Deprecated: Marked as deprecated in freeze.proto.
func (x *FreezeTransactionBody) GetEndMin() int32 {
	if x != nil {
		return x.EndMin
	}
	return 0
}

func (x *FreezeTransactionBody) GetUpdateFile() *FileID {
	if x != nil {
		return x.UpdateFile
	}
	return nil
}

func (x *FreezeTransactionBody) GetFileHash() []byte {
	if x != nil {
		return x.FileHash
	}
	return nil
}

func (x *FreezeTransactionBody) GetStartTime() *Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *FreezeTransactionBody) GetFreezeType() FreezeType {
	if x != nil {
		return x.FreezeType
	}
	return FreezeType_UNKNOWN_FREEZE_TYPE
}

var File_freeze_proto protoreflect.FileDescriptor

var file_freeze_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x66, 0x72, 0x65, 0x65, 0x7a, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x66, 0x72, 0x65, 0x65, 0x7a,
	0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc5, 0x02, 0x0a,
	0x15, 0x46, 0x72, 0x65, 0x65, 0x7a, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x20, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x48,
	0x6f, 0x75, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x02, 0x18, 0x01, 0x52, 0x09, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x48, 0x6f, 0x75, 0x72, 0x12, 0x1e, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x4d, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x02, 0x18, 0x01, 0x52, 0x08,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x4d, 0x69, 0x6e, 0x12, 0x1c, 0x0a, 0x07, 0x65, 0x6e, 0x64, 0x48,
	0x6f, 0x75, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x42, 0x02, 0x18, 0x01, 0x52, 0x07, 0x65,
	0x6e, 0x64, 0x48, 0x6f, 0x75, 0x72, 0x12, 0x1a, 0x0a, 0x06, 0x65, 0x6e, 0x64, 0x4d, 0x69, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x02, 0x18, 0x01, 0x52, 0x06, 0x65, 0x6e, 0x64, 0x4d,
	0x69, 0x6e, 0x12, 0x2e, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x2f, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x32, 0x0a, 0x0b, 0x66, 0x72, 0x65, 0x65, 0x7a, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x72,
	0x65, 0x65, 0x7a, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x66, 0x72, 0x65, 0x65, 0x7a, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65,
	0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_freeze_proto_rawDescOnce sync.Once
	file_freeze_proto_rawDescData = file_freeze_proto_rawDesc
)

func file_freeze_proto_rawDescGZIP() []byte {
	file_freeze_proto_rawDescOnce.Do(func() {
		file_freeze_proto_rawDescData = protoimpl.X.CompressGZIP(file_freeze_proto_rawDescData)
	})
	return file_freeze_proto_rawDescData
}

var file_freeze_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_freeze_proto_goTypes = []interface{}{
	(*FreezeTransactionBody)(nil), // 0: proto.FreezeTransactionBody
	(*FileID)(nil),                // 1: proto.FileID
	(*Timestamp)(nil),             // 2: proto.Timestamp
	(FreezeType)(0),               // 3: proto.FreezeType
}
var file_freeze_proto_depIdxs = []int32{
	1, // 0: proto.FreezeTransactionBody.update_file:type_name -> proto.FileID
	2, // 1: proto.FreezeTransactionBody.start_time:type_name -> proto.Timestamp
	3, // 2: proto.FreezeTransactionBody.freeze_type:type_name -> proto.FreezeType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_freeze_proto_init() }
func file_freeze_proto_init() {
	if File_freeze_proto != nil {
		return
	}
	file_timestamp_proto_init()
	file_basic_types_proto_init()
	file_freeze_type_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_freeze_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FreezeTransactionBody); i {
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
			RawDescriptor: file_freeze_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_freeze_proto_goTypes,
		DependencyIndexes: file_freeze_proto_depIdxs,
		MessageInfos:      file_freeze_proto_msgTypes,
	}.Build()
	File_freeze_proto = out.File
	file_freeze_proto_rawDesc = nil
	file_freeze_proto_goTypes = nil
	file_freeze_proto_depIdxs = nil
}
