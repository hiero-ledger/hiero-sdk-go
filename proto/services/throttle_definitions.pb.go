//*
// # Throttle Definitions
// A set of messages that support maintaining throttling limits on network
// transactions to ensure no one transaction type consumes the entirety of
// network resources. Also used to charge congestion fees when network load
// is exceptionally high, as an incentive to delay transactions that are
// not time-sensitive.
//
// For details behind this throttling design, please see the
// `docs/throttle-design.md` document in the
// [Hedera Services](https://github.com/hashgraph/hedera-services) repository.
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
// source: throttle_definitions.proto

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
// A single throttle limit applied to one or more operations.
//
// The list of operations MUST contain at least one entry.<br/>
// The throttle limit SHALL be specified in thousandths of an operation
// per second; one operation per second for the network would be `1000`.<br/>
// The throttle limit MUST be greater than zero (`0`).
type ThrottleGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A list of operations to be throttled.
	// <p>
	// This list MUST contain at least one item.<br/>
	// This list SHOULD NOT contain any item included in any other
	// active `ThrottleGroup`.
	Operations []HederaFunctionality `protobuf:"varint,1,rep,packed,name=operations,proto3,enum=proto.HederaFunctionality" json:"operations,omitempty"`
	// *
	// A throttle limit for this group.<br/>
	// This is a total number of operations, in thousandths, the network may
	// perform each second for this group. Every node executes every transaction,
	// so this limit effectively applies individually to each node as well.<br/>
	// <p>
	// This value MUST be greater than zero (`0`).<br/>
	// This value SHOULD be less than `9,223,372`.<br/>
	MilliOpsPerSec uint64 `protobuf:"varint,2,opt,name=milliOpsPerSec,proto3" json:"milliOpsPerSec,omitempty"`
}

func (x *ThrottleGroup) Reset() {
	*x = ThrottleGroup{}
	mi := &file_throttle_definitions_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ThrottleGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ThrottleGroup) ProtoMessage() {}

func (x *ThrottleGroup) ProtoReflect() protoreflect.Message {
	mi := &file_throttle_definitions_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ThrottleGroup.ProtoReflect.Descriptor instead.
func (*ThrottleGroup) Descriptor() ([]byte, []int) {
	return file_throttle_definitions_proto_rawDescGZIP(), []int{0}
}

func (x *ThrottleGroup) GetOperations() []HederaFunctionality {
	if x != nil {
		return x.Operations
	}
	return nil
}

func (x *ThrottleGroup) GetMilliOpsPerSec() uint64 {
	if x != nil {
		return x.MilliOpsPerSec
	}
	return 0
}

// *
// A "bucket" of performance allocated across one or more throttle groups.<br/>
// This entry combines one or more throttle groups into a single unit to
// calculate limitations and congestion. Each "bucket" "fills" as operations
// are completed, then "drains" over a period of time defined for each bucket.
// This fill-and-drain characteristic enables the network to process sudden
// bursts of heavy traffic while still observing throttle limits over longer
// timeframes.
//
// The value of `burstPeriodMs` is combined with the `milliOpsPerSec`
// values for the individual throttle groups to determine the total
// bucket "capacity". This combination MUST be less than the maximum
// value of a signed long integer (`9223372036854775807`), when scaled to
// a nanosecond measurement resolution.
//
// > Note
// >> There is some question regarding the mechanism of calculating the
// >> combination of `burstPeriodMs` and `milliOpsPerSec`. The calculation
// >> Is implemented in difficult-to-find code, and very likely does not
// >> match the approach described here.
type ThrottleBucket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A name for this bucket.<br/>
	// This is used for log entries.
	// <p>
	// This value SHOULD NOT exceed 20 characters.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// *
	// A burst duration limit, in milliseconds.<br/>
	// This value determines the total "capacity" of the bucket. The rate
	// at which the bucket "drains" is set by the throttles, and this duration
	// sets how long that rate must be sustained to empty a "full" bucket.
	// That combination (calculated as the product of this value and the least
	// common multiple of the `milliOpsPerSec` values for all throttle groups)
	// determines the maximum amount of operations this bucket can "hold".
	// <p>
	// The calculated capacity of this bucket MUST NOT exceed `9,223,372,036,854`.
	BurstPeriodMs uint64 `protobuf:"varint,2,opt,name=burstPeriodMs,proto3" json:"burstPeriodMs,omitempty"`
	// *
	// A list of throttle groups.<br/>
	// These throttle groups combined define the effective throttle
	// rate for the bucket.
	// <p>
	// This list MUST contain at least one entry.
	ThrottleGroups []*ThrottleGroup `protobuf:"bytes,3,rep,name=throttleGroups,proto3" json:"throttleGroups,omitempty"`
}

func (x *ThrottleBucket) Reset() {
	*x = ThrottleBucket{}
	mi := &file_throttle_definitions_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ThrottleBucket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ThrottleBucket) ProtoMessage() {}

func (x *ThrottleBucket) ProtoReflect() protoreflect.Message {
	mi := &file_throttle_definitions_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ThrottleBucket.ProtoReflect.Descriptor instead.
func (*ThrottleBucket) Descriptor() ([]byte, []int) {
	return file_throttle_definitions_proto_rawDescGZIP(), []int{1}
}

func (x *ThrottleBucket) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ThrottleBucket) GetBurstPeriodMs() uint64 {
	if x != nil {
		return x.BurstPeriodMs
	}
	return 0
}

func (x *ThrottleBucket) GetThrottleGroups() []*ThrottleGroup {
	if x != nil {
		return x.ThrottleGroups
	}
	return nil
}

// *
// A list of throttle buckets.<br/>
// This list, simultaneously enforced, defines a complete throttling policy.
//
//  1. When an operation appears in more than one throttling bucket,
//     that operation SHALL be throttled unless all of the buckets where
//     the operation appears have "capacity" available.
//  1. An operation assigned to no buckets is SHALL be throttled in every
//     instance.  The _effective_ throttle for this case is `0`.
type ThrottleDefinitions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A list of throttle buckets.
	// <p>
	// This list MUST be set, and SHOULD NOT be empty.<br/>
	// An empty list SHALL have the effect of setting all operations to
	// a single group with throttle limit of `0` operations per second for the
	// entire network.
	ThrottleBuckets []*ThrottleBucket `protobuf:"bytes,1,rep,name=throttleBuckets,proto3" json:"throttleBuckets,omitempty"`
}

func (x *ThrottleDefinitions) Reset() {
	*x = ThrottleDefinitions{}
	mi := &file_throttle_definitions_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ThrottleDefinitions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ThrottleDefinitions) ProtoMessage() {}

func (x *ThrottleDefinitions) ProtoReflect() protoreflect.Message {
	mi := &file_throttle_definitions_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ThrottleDefinitions.ProtoReflect.Descriptor instead.
func (*ThrottleDefinitions) Descriptor() ([]byte, []int) {
	return file_throttle_definitions_proto_rawDescGZIP(), []int{2}
}

func (x *ThrottleDefinitions) GetThrottleBuckets() []*ThrottleBucket {
	if x != nil {
		return x.ThrottleBuckets
	}
	return nil
}

var File_throttle_definitions_proto protoreflect.FileDescriptor

var file_throttle_definitions_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x74, 0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c, 0x65, 0x5f, 0x64, 0x65, 0x66, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a, 0x0d, 0x54, 0x68, 0x72, 0x6f, 0x74, 0x74,
	0x6c, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x3a, 0x0a, 0x0a, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x64, 0x65, 0x72, 0x61, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x0a, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x26, 0x0a, 0x0e, 0x6d, 0x69, 0x6c, 0x6c, 0x69, 0x4f, 0x70, 0x73, 0x50,
	0x65, 0x72, 0x53, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x6d, 0x69, 0x6c,
	0x6c, 0x69, 0x4f, 0x70, 0x73, 0x50, 0x65, 0x72, 0x53, 0x65, 0x63, 0x22, 0x88, 0x01, 0x0a, 0x0e,
	0x54, 0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c, 0x65, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x62, 0x75, 0x72, 0x73, 0x74, 0x50, 0x65, 0x72, 0x69, 0x6f,
	0x64, 0x4d, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x62, 0x75, 0x72, 0x73, 0x74,
	0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x4d, 0x73, 0x12, 0x3c, 0x0a, 0x0e, 0x74, 0x68, 0x72, 0x6f,
	0x74, 0x74, 0x6c, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c,
	0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x0e, 0x74, 0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c, 0x65,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x22, 0x56, 0x0a, 0x13, 0x54, 0x68, 0x72, 0x6f, 0x74, 0x74,
	0x6c, 0x65, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3f, 0x0a,
	0x0f, 0x74, 0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c, 0x65, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54,
	0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c, 0x65, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x0f, 0x74,
	0x68, 0x72, 0x6f, 0x74, 0x74, 0x6c, 0x65, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x42, 0x26,
	0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68,
	0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_throttle_definitions_proto_rawDescOnce sync.Once
	file_throttle_definitions_proto_rawDescData = file_throttle_definitions_proto_rawDesc
)

func file_throttle_definitions_proto_rawDescGZIP() []byte {
	file_throttle_definitions_proto_rawDescOnce.Do(func() {
		file_throttle_definitions_proto_rawDescData = protoimpl.X.CompressGZIP(file_throttle_definitions_proto_rawDescData)
	})
	return file_throttle_definitions_proto_rawDescData
}

var file_throttle_definitions_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_throttle_definitions_proto_goTypes = []any{
	(*ThrottleGroup)(nil),       // 0: proto.ThrottleGroup
	(*ThrottleBucket)(nil),      // 1: proto.ThrottleBucket
	(*ThrottleDefinitions)(nil), // 2: proto.ThrottleDefinitions
	(HederaFunctionality)(0),    // 3: proto.HederaFunctionality
}
var file_throttle_definitions_proto_depIdxs = []int32{
	3, // 0: proto.ThrottleGroup.operations:type_name -> proto.HederaFunctionality
	0, // 1: proto.ThrottleBucket.throttleGroups:type_name -> proto.ThrottleGroup
	1, // 2: proto.ThrottleDefinitions.throttleBuckets:type_name -> proto.ThrottleBucket
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_throttle_definitions_proto_init() }
func file_throttle_definitions_proto_init() {
	if File_throttle_definitions_proto != nil {
		return
	}
	file_basic_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_throttle_definitions_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_throttle_definitions_proto_goTypes,
		DependencyIndexes: file_throttle_definitions_proto_depIdxs,
		MessageInfos:      file_throttle_definitions_proto_msgTypes,
	}.Build()
	File_throttle_definitions_proto = out.File
	file_throttle_definitions_proto_rawDesc = nil
	file_throttle_definitions_proto_goTypes = nil
	file_throttle_definitions_proto_depIdxs = nil
}
