//*
// # Schedule Create
// Message to create a schedule, which is an instruction to execute some other
// transaction (the scheduled transaction) at a future time, either when
// enough signatures are gathered (short term) or when the schedule expires
// (long term). In all cases the scheduled transaction is not executed if
// signature requirements are not met before the schedule expires.
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
// source: schedule_create.proto

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
// Create a new Schedule.
//
// #### Requirements
// This transaction SHALL create a new _schedule_ entity in network state.<br/>
// The schedule created SHALL contain the `scheduledTransactionBody` to be
// executed.<br/>
// If successful the receipt SHALL contain a `scheduleID` with the full
// identifier of the schedule created.<br/>
// When a schedule _executes_ successfully, the receipt SHALL include a
// `scheduledTransactionID` with the `TransactionID` of the transaction that
// executed.<br/>
// When a scheduled transaction is executed the network SHALL charge the
// regular _service_ fee for the transaction to the `payerAccountID` for
// that schedule, but SHALL NOT charge node or network fees.<br/>
// If the `payerAccountID` field is not set, the effective `payerAccountID`
// SHALL be the `payer` for this create transaction.<br/>
// If an `adminKey` is not specified, or is an empty `KeyList`, the schedule
// created SHALL be immutable.<br/>
// An immutable schedule MAY be signed, and MAY execute, but SHALL NOT be
// deleted.<br/>
// If two schedules have the same values for all fields except `payerAccountID`
// then those two schedules SHALL be deemed "identical".<br/>
// If a `scheduleCreate` requests a new schedule that is identical to an
// existing schedule, the transaction SHALL fail and SHALL return a status
// code of `IDENTICAL_SCHEDULE_ALREADY_CREATED` in the receipt.<br/>
// The receipt for a duplicate schedule SHALL include the `ScheduleID` of the
// existing schedule and the `TransactionID` of the earlier `scheduleCreate`
// so that the earlier schedule may be queried and/or referred to in a
// subsequent `scheduleSign`.
//
// #### Signature Requirements
// A `scheduleSign` transaction SHALL be used to add additional signatures
// to an existing schedule.<br/>
// Each signature SHALL "activate" the corresponding cryptographic("primitive")
// key for that schedule.<br/>
// Signature requirements SHALL be met when the set of active keys includes
// all keys required by the scheduled transaction.<br/>
// A scheduled transaction for a "long term" schedule SHALL NOT execute if
// the signature requirements for that transaction are not met when the
// network consensus time reaches the schedule `expiration_time`.<br/>
// A "short term" schedule SHALL execute immediately once signature
// requirements are met. This MAY be immediately when created.
//
// #### Long Term Schedules
// A "short term" schedule SHALL have the flag `wait_for_expiry` _unset_.<br/>
// A "long term" schedule SHALL have the flag  `wait_for_expiry` _set_.<br/>
// A "long term" schedule SHALL NOT be accepted if the network configuration
// `scheduling.longTermEnabled` is not enabled.<br/>
// A "long term" schedule SHALL execute when the current consensus time
// matches or exceeds the `expiration_time` for that schedule, if the
// signature requirements for the scheduled transaction
// are met at that instant.<br/>
// A "long term" schedule SHALL NOT execute before the current consensus time
// matches or exceeds the `expiration_time` for that schedule.<br/>
// A "long term" schedule SHALL expire, and be removed from state, after the
// network consensus time exceeds the schedule `expiration_time`.<br/>
// A short term schedule SHALL expire, and be removed from state,
// after the network consensus time exceeds the current network
// configuration for `ledger.scheduleTxExpiryTimeSecs`.
//
// > Note
// >> Long term schedules are not (as of release 0.56.0) enabled. Any schedule
// >> created currently MUST NOT set the `wait_for_expiry` flag.<br/>
// >> When long term schedules are not enabled, schedules SHALL NOT be
// >> executed at expiration, and MUST meet signature requirements strictly
// >> before expiration to be executed.
//
// ### Block Stream Effects
// If the scheduled transaction is executed immediately, the transaction
// record SHALL include a `scheduleRef` with the schedule identifier of the
// schedule created.
type ScheduleCreateTransactionBody struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// *
	// A scheduled transaction.
	// <p>
	// This value is REQUIRED.<br/>
	// This transaction body MUST be one of the types enabled in the
	// network configuration value `scheduling.whitelist`.
	ScheduledTransactionBody *SchedulableTransactionBody `protobuf:"bytes,1,opt,name=scheduledTransactionBody,proto3" json:"scheduledTransactionBody,omitempty"`
	// *
	// A short description of the schedule.
	// <p>
	// This value, if set, MUST NOT exceed `transaction.maxMemoUtf8Bytes`
	// (default 100) bytes when encoded as UTF-8.
	Memo string `protobuf:"bytes,2,opt,name=memo,proto3" json:"memo,omitempty"`
	// *
	// A `Key` required to delete this schedule.
	// <p>
	// If this is not set, or is an empty `KeyList`, this schedule SHALL be
	// immutable and SHALL NOT be deleted.
	AdminKey *Key `protobuf:"bytes,3,opt,name=adminKey,proto3" json:"adminKey,omitempty"`
	// *
	// An account identifier of a `payer` for the scheduled transaction.
	// <p>
	// This value MAY be unset. If unset, the `payer` for this `scheduleCreate`
	// transaction SHALL be the `payer` for the scheduled transaction.<br/>
	// If this is set, the identified account SHALL be charged the fees
	// required for the scheduled transaction when it is executed.<br/>
	// If the actual `payer` for the _scheduled_ transaction lacks
	// sufficient HBAR balance to pay service fees for the scheduled
	// transaction _when it executes_, the scheduled transaction
	// SHALL fail with `INSUFFICIENT_PAYER_BALANCE`.<br/>
	PayerAccountID *AccountID `protobuf:"bytes,4,opt,name=payerAccountID,proto3" json:"payerAccountID,omitempty"`
	// *
	// An expiration time.
	// <p>
	// If not set, the expiration SHALL default to the current consensus time
	// advanced by either the network configuration value
	// `scheduling.maxExpirationFutureSeconds`, if `wait_for_expiry` is set and
	// "long term" schedules are enabled, or the network configuration value
	// `ledger.scheduleTxExpiryTimeSecs` otherwise.
	ExpirationTime *Timestamp `protobuf:"bytes,5,opt,name=expiration_time,json=expirationTime,proto3" json:"expiration_time,omitempty"`
	// *
	// A flag to delay execution until expiration.
	// <p>
	// If this flag is set the scheduled transaction SHALL NOT be evaluated for
	// execution before the network consensus time matches or exceeds the
	// `expiration_time`.<br/>
	// If this flag is not set, the scheduled transaction SHALL be executed
	// immediately when all required signatures are received, whether in this
	// `scheduleCreate` transaction or a later `scheduleSign` transaction.<br/>
	// This value SHALL NOT be used and MUST NOT be set when the network
	// configuration value `scheduling.longTermEnabled` is not enabled.
	WaitForExpiry bool `protobuf:"varint,13,opt,name=wait_for_expiry,json=waitForExpiry,proto3" json:"wait_for_expiry,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScheduleCreateTransactionBody) Reset() {
	*x = ScheduleCreateTransactionBody{}
	mi := &file_schedule_create_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScheduleCreateTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleCreateTransactionBody) ProtoMessage() {}

func (x *ScheduleCreateTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_schedule_create_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleCreateTransactionBody.ProtoReflect.Descriptor instead.
func (*ScheduleCreateTransactionBody) Descriptor() ([]byte, []int) {
	return file_schedule_create_proto_rawDescGZIP(), []int{0}
}

func (x *ScheduleCreateTransactionBody) GetScheduledTransactionBody() *SchedulableTransactionBody {
	if x != nil {
		return x.ScheduledTransactionBody
	}
	return nil
}

func (x *ScheduleCreateTransactionBody) GetMemo() string {
	if x != nil {
		return x.Memo
	}
	return ""
}

func (x *ScheduleCreateTransactionBody) GetAdminKey() *Key {
	if x != nil {
		return x.AdminKey
	}
	return nil
}

func (x *ScheduleCreateTransactionBody) GetPayerAccountID() *AccountID {
	if x != nil {
		return x.PayerAccountID
	}
	return nil
}

func (x *ScheduleCreateTransactionBody) GetExpirationTime() *Timestamp {
	if x != nil {
		return x.ExpirationTime
	}
	return nil
}

func (x *ScheduleCreateTransactionBody) GetWaitForExpiry() bool {
	if x != nil {
		return x.WaitForExpiry
	}
	return false
}

var File_schedule_create_proto protoreflect.FileDescriptor

const file_schedule_create_proto_rawDesc = "" +
	"\n" +
	"\x15schedule_create.proto\x12\x05proto\x1a\x11basic_types.proto\x1a\x0ftimestamp.proto\x1a\"schedulable_transaction_body.proto\"\xd7\x02\n" +
	"\x1dScheduleCreateTransactionBody\x12]\n" +
	"\x18scheduledTransactionBody\x18\x01 \x01(\v2!.proto.SchedulableTransactionBodyR\x18scheduledTransactionBody\x12\x12\n" +
	"\x04memo\x18\x02 \x01(\tR\x04memo\x12&\n" +
	"\badminKey\x18\x03 \x01(\v2\n" +
	".proto.KeyR\badminKey\x128\n" +
	"\x0epayerAccountID\x18\x04 \x01(\v2\x10.proto.AccountIDR\x0epayerAccountID\x129\n" +
	"\x0fexpiration_time\x18\x05 \x01(\v2\x10.proto.TimestampR\x0eexpirationTime\x12&\n" +
	"\x0fwait_for_expiry\x18\r \x01(\bR\rwaitForExpiryB&\n" +
	"\"com.hederahashgraph.api.proto.javaP\x01b\x06proto3"

var (
	file_schedule_create_proto_rawDescOnce sync.Once
	file_schedule_create_proto_rawDescData []byte
)

func file_schedule_create_proto_rawDescGZIP() []byte {
	file_schedule_create_proto_rawDescOnce.Do(func() {
		file_schedule_create_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_schedule_create_proto_rawDesc), len(file_schedule_create_proto_rawDesc)))
	})
	return file_schedule_create_proto_rawDescData
}

var file_schedule_create_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_schedule_create_proto_goTypes = []any{
	(*ScheduleCreateTransactionBody)(nil), // 0: proto.ScheduleCreateTransactionBody
	(*SchedulableTransactionBody)(nil),    // 1: proto.SchedulableTransactionBody
	(*Key)(nil),                           // 2: proto.Key
	(*AccountID)(nil),                     // 3: proto.AccountID
	(*Timestamp)(nil),                     // 4: proto.Timestamp
}
var file_schedule_create_proto_depIdxs = []int32{
	1, // 0: proto.ScheduleCreateTransactionBody.scheduledTransactionBody:type_name -> proto.SchedulableTransactionBody
	2, // 1: proto.ScheduleCreateTransactionBody.adminKey:type_name -> proto.Key
	3, // 2: proto.ScheduleCreateTransactionBody.payerAccountID:type_name -> proto.AccountID
	4, // 3: proto.ScheduleCreateTransactionBody.expiration_time:type_name -> proto.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_schedule_create_proto_init() }
func file_schedule_create_proto_init() {
	if File_schedule_create_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_timestamp_proto_init()
	file_schedulable_transaction_body_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_schedule_create_proto_rawDesc), len(file_schedule_create_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_schedule_create_proto_goTypes,
		DependencyIndexes: file_schedule_create_proto_depIdxs,
		MessageInfos:      file_schedule_create_proto_msgTypes,
	}.Build()
	File_schedule_create_proto = out.File
	file_schedule_create_proto_goTypes = nil
	file_schedule_create_proto_depIdxs = nil
}
