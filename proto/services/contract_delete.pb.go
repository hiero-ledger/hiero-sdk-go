//*
// # Contract Delete
// Delete a smart contract, transferring any remaining balance to a
// designated account.
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
// source: contract_delete.proto

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
// Delete a smart contract, and transfer any remaining HBAR balance to a
// designated account.
//
// If this call succeeds then all subsequent calls to that smart contract
// SHALL execute the `0x0` opcode, as required for EVM equivalence.
//
// ### Requirements
//   - An account or smart contract MUST be designated to receive all remaining
//     account balances.
//   - The smart contract MUST have an admin key set. If the contract does not
//     have `admin_key` set, then this transaction SHALL fail and response code
//     `MODIFYING_IMMUTABLE_CONTRACT` SHALL be set.
//   - If `admin_key` is, or contains, an empty `KeyList` key, it SHALL be
//     treated the same as an admin key that is not set.
//   - The `Key` set for `admin_key` on the smart contract MUST have a valid
//     signature set on this transaction.
//   - The designated receiving account MAY have `receiver_sig_required` set. If
//     that field is set, the receiver account MUST also sign this transaction.
//   - The field `permanent_removal` MUST NOT be set. That field is reserved for
//     internal system use when purging the smart contract from state. Any user
//     transaction with that field set SHALL be rejected and a response code
//     `PERMANENT_REMOVAL_REQUIRES_SYSTEM_INITIATION` SHALL be set.
//
// ### Block Stream Effects
// None
type ContractDeleteTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The id of the contract to be deleted.
	// <p>
	// This field is REQUIRED.
	ContractID *ContractID `protobuf:"bytes,1,opt,name=contractID,proto3" json:"contractID,omitempty"`
	// Types that are assignable to Obtainers:
	//
	//	*ContractDeleteTransactionBody_TransferAccountID
	//	*ContractDeleteTransactionBody_TransferContractID
	Obtainers isContractDeleteTransactionBody_Obtainers `protobuf_oneof:"obtainers"`
	// *
	// A flag indicating that this transaction is "synthetic"; initiated by the
	// node software.
	// <p>
	// The consensus nodes create such "synthetic" transactions to both to
	// properly manage state changes and to communicate those changes to other
	// systems via the Block Stream.<br/>
	// A user-initiated transaction MUST NOT set this flag.
	PermanentRemoval bool `protobuf:"varint,4,opt,name=permanent_removal,json=permanentRemoval,proto3" json:"permanent_removal,omitempty"`
}

func (x *ContractDeleteTransactionBody) Reset() {
	*x = ContractDeleteTransactionBody{}
	mi := &file_contract_delete_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContractDeleteTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractDeleteTransactionBody) ProtoMessage() {}

func (x *ContractDeleteTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_contract_delete_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractDeleteTransactionBody.ProtoReflect.Descriptor instead.
func (*ContractDeleteTransactionBody) Descriptor() ([]byte, []int) {
	return file_contract_delete_proto_rawDescGZIP(), []int{0}
}

func (x *ContractDeleteTransactionBody) GetContractID() *ContractID {
	if x != nil {
		return x.ContractID
	}
	return nil
}

func (m *ContractDeleteTransactionBody) GetObtainers() isContractDeleteTransactionBody_Obtainers {
	if m != nil {
		return m.Obtainers
	}
	return nil
}

func (x *ContractDeleteTransactionBody) GetTransferAccountID() *AccountID {
	if x, ok := x.GetObtainers().(*ContractDeleteTransactionBody_TransferAccountID); ok {
		return x.TransferAccountID
	}
	return nil
}

func (x *ContractDeleteTransactionBody) GetTransferContractID() *ContractID {
	if x, ok := x.GetObtainers().(*ContractDeleteTransactionBody_TransferContractID); ok {
		return x.TransferContractID
	}
	return nil
}

func (x *ContractDeleteTransactionBody) GetPermanentRemoval() bool {
	if x != nil {
		return x.PermanentRemoval
	}
	return false
}

type isContractDeleteTransactionBody_Obtainers interface {
	isContractDeleteTransactionBody_Obtainers()
}

type ContractDeleteTransactionBody_TransferAccountID struct {
	// *
	// An Account ID recipient.
	// <p>
	// This account SHALL receive all HBAR and other tokens still owned by
	// the contract that is removed.
	TransferAccountID *AccountID `protobuf:"bytes,2,opt,name=transferAccountID,proto3,oneof"`
}

type ContractDeleteTransactionBody_TransferContractID struct {
	// *
	// A contract ID recipient.
	// <p>
	// This contract SHALL receive all HBAR and other tokens still owned by
	// the contract that is removed.
	TransferContractID *ContractID `protobuf:"bytes,3,opt,name=transferContractID,proto3,oneof"`
}

func (*ContractDeleteTransactionBody_TransferAccountID) isContractDeleteTransactionBody_Obtainers() {}

func (*ContractDeleteTransactionBody_TransferContractID) isContractDeleteTransactionBody_Obtainers() {
}

var File_contract_delete_proto protoreflect.FileDescriptor

var file_contract_delete_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11,
	0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x93, 0x02, 0x0a, 0x1d, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x6f, 0x64, 0x79, 0x12, 0x31, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x44, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x49, 0x44, 0x12, 0x40, 0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x65, 0x72, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x49, 0x44, 0x48, 0x00, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x43, 0x0a, 0x12, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x44, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x44, 0x48, 0x00, 0x52, 0x12, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x49, 0x44, 0x12, 0x2b, 0x0a,
	0x11, 0x70, 0x65, 0x72, 0x6d, 0x61, 0x6e, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x6d, 0x6f, 0x76,
	0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x70, 0x65, 0x72, 0x6d, 0x61, 0x6e,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x61, 0x6c, 0x42, 0x0b, 0x0a, 0x09, 0x6f, 0x62,
	0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x73, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68,
	0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contract_delete_proto_rawDescOnce sync.Once
	file_contract_delete_proto_rawDescData = file_contract_delete_proto_rawDesc
)

func file_contract_delete_proto_rawDescGZIP() []byte {
	file_contract_delete_proto_rawDescOnce.Do(func() {
		file_contract_delete_proto_rawDescData = protoimpl.X.CompressGZIP(file_contract_delete_proto_rawDescData)
	})
	return file_contract_delete_proto_rawDescData
}

var file_contract_delete_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_contract_delete_proto_goTypes = []any{
	(*ContractDeleteTransactionBody)(nil), // 0: proto.ContractDeleteTransactionBody
	(*ContractID)(nil),                    // 1: proto.ContractID
	(*AccountID)(nil),                     // 2: proto.AccountID
}
var file_contract_delete_proto_depIdxs = []int32{
	1, // 0: proto.ContractDeleteTransactionBody.contractID:type_name -> proto.ContractID
	2, // 1: proto.ContractDeleteTransactionBody.transferAccountID:type_name -> proto.AccountID
	1, // 2: proto.ContractDeleteTransactionBody.transferContractID:type_name -> proto.ContractID
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_contract_delete_proto_init() }
func file_contract_delete_proto_init() {
	if File_contract_delete_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_contract_delete_proto_msgTypes[0].OneofWrappers = []any{
		(*ContractDeleteTransactionBody_TransferAccountID)(nil),
		(*ContractDeleteTransactionBody_TransferContractID)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contract_delete_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_contract_delete_proto_goTypes,
		DependencyIndexes: file_contract_delete_proto_depIdxs,
		MessageInfos:      file_contract_delete_proto_msgTypes,
	}.Build()
	File_contract_delete_proto = out.File
	file_contract_delete_proto_rawDesc = nil
	file_contract_delete_proto_goTypes = nil
	file_contract_delete_proto_depIdxs = nil
}
