//*
// # Freeze Type
// An enumeration to select the type of a network freeze.
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
// source: freeze_type.proto

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
// An enumeration of possible network freeze types.
//
// Each enumerated value SHALL be associated to a single network freeze
// scenario. Each freeze scenario defines the specific parameters
// REQUIRED for that freeze.
type FreezeType int32

const (
	// *
	// An invalid freeze type.
	// <p>
	// The first value in a protobuf enum is a default value. This default
	// is RECOMMENDED to be an invalid value to aid in detecting unset fields.
	FreezeType_UNKNOWN_FREEZE_TYPE FreezeType = 0
	// *
	// Freeze the network, and take no further action.
	// <p>
	// The `start_time` field is REQUIRED, MUST be strictly later than the
	// consensus time when this transaction is handled, and SHOULD be between
	// `300` and `3600` seconds after the transaction identifier
	// `transactionValidStart` field.<br/>
	// The fields `update_file` and `file_hash` SHALL be ignored.<br/>
	// A `FREEZE_ONLY` transaction SHALL NOT perform any network
	// changes or upgrades.<br/>
	// After this freeze is processed manual intervention is REQUIRED
	// to restart the network.
	FreezeType_FREEZE_ONLY FreezeType = 1
	// *
	// This freeze type does not freeze the network, but begins
	// "preparation" to upgrade the network.
	// <p>
	// The fields `update_file` and `file_hash` are REQUIRED
	// and MUST be valid.<br/>
	// The `start_time` field SHALL be ignored.<br/>
	// A `PREPARE_UPGRADE` transaction SHALL NOT freeze the network or
	// interfere with general transaction processing.<br/>
	// If this freeze type is initiated after a `TELEMETRY_UPGRADE`, the
	// prepared telemetry upgrade SHALL be reset and all telemetry upgrade
	// artifacts in the filesystem SHALL be deleted.<br/>
	// At some point after this freeze type completes (dependent on the size
	// of the upgrade file), the network SHALL be prepared to complete
	// a software upgrade of all nodes.
	FreezeType_PREPARE_UPGRADE FreezeType = 2
	// *
	// Freeze the network to perform a software upgrade.
	// <p>
	// The `start_time` field is REQUIRED, MUST be strictly later than the
	// consensus time when this transaction is handled, and SHOULD be between
	// `300` and `3600` seconds after the transaction identifier
	// `transactionValidStart` field.<br/>
	// A software upgrade file MUST be prepared prior to this transaction.<br/>
	// After this transaction completes, the network SHALL initiate an
	// upgrade and restart of all nodes at the start time specified.
	FreezeType_FREEZE_UPGRADE FreezeType = 3
	// *
	// Abort a pending network freeze operation.
	// <p>
	// All fields SHALL be ignored for this freeze type.<br/>
	// This freeze type MAY be submitted after a `FREEZE_ONLY`,
	// `FREEZE_UPGRADE`, or `TELEMETRY_UPGRADE` is initiated.<br/>
	// This freeze type MUST be submitted and reach consensus
	// before the `start_time` designated for the current pending
	// freeze to be effective.<br/>
	// After this freeze type is processed, the upgrade file hash
	// and pending freeze start time stored in the network SHALL
	// be reset to default (empty) values.
	FreezeType_FREEZE_ABORT FreezeType = 4
	// *
	// Prepare an upgrade of auxiliary services and containers
	// providing telemetry/metrics.
	// <p>
	// The `start_time` field is REQUIRED, MUST be strictly later than the
	// consensus time when this transaction is handled, and SHOULD be between
	// `300` and `3600` seconds after the transaction identifier
	// `transactionValidStart` field.<br/>
	// The `update_file` field is REQUIRED and MUST be valid.<br/>
	// A `TELEMETRY_UPGRADE` transaction SHALL NOT freeze the network or
	// interfere with general transaction processing.<br/>
	// This freeze type MUST NOT be initiated between a `PREPARE_UPGRADE`
	// and `FREEZE_UPGRADE`. If this freeze type is initiated after a
	// `PREPARE_UPGRADE`, the prepared upgrade SHALL be reset and all software
	// upgrade artifacts in the filesystem SHALL be deleted.<br/>
	// At some point after this freeze type completes (dependent on the
	// size of the upgrade file), the network SHALL automatically upgrade
	// the telemetry/metrics services and containers as directed in
	// the specified telemetry upgrade file.
	// <blockquote> The condition that `start_time` is REQUIRED is an
	// historical anomaly and SHOULD change in a future release.</blockquote>
	FreezeType_TELEMETRY_UPGRADE FreezeType = 5
)

// Enum value maps for FreezeType.
var (
	FreezeType_name = map[int32]string{
		0: "UNKNOWN_FREEZE_TYPE",
		1: "FREEZE_ONLY",
		2: "PREPARE_UPGRADE",
		3: "FREEZE_UPGRADE",
		4: "FREEZE_ABORT",
		5: "TELEMETRY_UPGRADE",
	}
	FreezeType_value = map[string]int32{
		"UNKNOWN_FREEZE_TYPE": 0,
		"FREEZE_ONLY":         1,
		"PREPARE_UPGRADE":     2,
		"FREEZE_UPGRADE":      3,
		"FREEZE_ABORT":        4,
		"TELEMETRY_UPGRADE":   5,
	}
)

func (x FreezeType) Enum() *FreezeType {
	p := new(FreezeType)
	*p = x
	return p
}

func (x FreezeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FreezeType) Descriptor() protoreflect.EnumDescriptor {
	return file_freeze_type_proto_enumTypes[0].Descriptor()
}

func (FreezeType) Type() protoreflect.EnumType {
	return &file_freeze_type_proto_enumTypes[0]
}

func (x FreezeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FreezeType.Descriptor instead.
func (FreezeType) EnumDescriptor() ([]byte, []int) {
	return file_freeze_type_proto_rawDescGZIP(), []int{0}
}

var File_freeze_type_proto protoreflect.FileDescriptor

var file_freeze_type_proto_rawDesc = []byte{
	0x0a, 0x11, 0x66, 0x72, 0x65, 0x65, 0x7a, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x88, 0x01, 0x0a, 0x0a, 0x46,
	0x72, 0x65, 0x65, 0x7a, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x46, 0x52, 0x45, 0x45, 0x5a, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x46, 0x52, 0x45, 0x45, 0x5a, 0x45, 0x5f, 0x4f, 0x4e, 0x4c,
	0x59, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x52, 0x45, 0x50, 0x41, 0x52, 0x45, 0x5f, 0x55,
	0x50, 0x47, 0x52, 0x41, 0x44, 0x45, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x46, 0x52, 0x45, 0x45,
	0x5a, 0x45, 0x5f, 0x55, 0x50, 0x47, 0x52, 0x41, 0x44, 0x45, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c,
	0x46, 0x52, 0x45, 0x45, 0x5a, 0x45, 0x5f, 0x41, 0x42, 0x4f, 0x52, 0x54, 0x10, 0x04, 0x12, 0x15,
	0x0a, 0x11, 0x54, 0x45, 0x4c, 0x45, 0x4d, 0x45, 0x54, 0x52, 0x59, 0x5f, 0x55, 0x50, 0x47, 0x52,
	0x41, 0x44, 0x45, 0x10, 0x05, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64,
	0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_freeze_type_proto_rawDescOnce sync.Once
	file_freeze_type_proto_rawDescData = file_freeze_type_proto_rawDesc
)

func file_freeze_type_proto_rawDescGZIP() []byte {
	file_freeze_type_proto_rawDescOnce.Do(func() {
		file_freeze_type_proto_rawDescData = protoimpl.X.CompressGZIP(file_freeze_type_proto_rawDescData)
	})
	return file_freeze_type_proto_rawDescData
}

var file_freeze_type_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_freeze_type_proto_goTypes = []any{
	(FreezeType)(0), // 0: proto.FreezeType
}
var file_freeze_type_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_freeze_type_proto_init() }
func file_freeze_type_proto_init() {
	if File_freeze_type_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_freeze_type_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_freeze_type_proto_goTypes,
		DependencyIndexes: file_freeze_type_proto_depIdxs,
		EnumInfos:         file_freeze_type_proto_enumTypes,
	}.Build()
	File_freeze_type_proto = out.File
	file_freeze_type_proto_rawDesc = nil
	file_freeze_type_proto_goTypes = nil
	file_freeze_type_proto_depIdxs = nil
}
