//*
// # Crypto Update
// Modify a single account.
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
// source: crypto_update.proto

package services

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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
// Modify the current state of an account.
//
// ### Requirements
//   - The `key` for this account MUST sign all account update transactions.
//   - If the `key` field is set for this transaction, then _both_ the current
//     `key` and the new `key` MUST sign this transaction, for security and to
//     prevent setting the `key` field to an invalid value.
//   - If the `auto_renew_account` field is set for this transaction, the account
//     identified in that field MUST sign this transaction.
//   - Fields set to non-default values in this transaction SHALL be updated on
//     success. Fields not set to non-default values SHALL NOT be
//     updated on success.
//   - All fields that may be modified in this transaction SHALL have a
//     default value of unset (a.k.a. `null`).
//
// ### Block Stream Effects
// None
type CryptoUpdateTransactionBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// An account identifier.<br/>
	// This identifies the account which is to be modified in this transaction.
	// <p>
	// This field is REQUIRED.
	AccountIDToUpdate *AccountID `protobuf:"bytes,2,opt,name=accountIDToUpdate,proto3" json:"accountIDToUpdate,omitempty"`
	// *
	// An account key.<br/>
	// This may be a "primitive" key (a singly cryptographic key), or a
	// composite key.
	// <p>
	// If set, this key MUST be a valid key.<br/>
	// If set, the previous key and new key MUST both sign this transaction.
	Key *Key `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	// *
	// Removed in favor of the `staked_id` oneOf.<br/>
	// An account identifier for a "proxy" account. This account's HBAR are
	// staked to a node selected by the proxy account.
	//
	// Deprecated: Marked as deprecated in crypto_update.proto.
	ProxyAccountID *AccountID `protobuf:"bytes,4,opt,name=proxyAccountID,proto3" json:"proxyAccountID,omitempty"`
	// *
	// Removed prior to the first available history.<br/>
	// A fraction to split staking rewards between this account and the proxy
	// account.
	//
	// Deprecated: Marked as deprecated in crypto_update.proto.
	ProxyFraction int32 `protobuf:"varint,5,opt,name=proxyFraction,proto3" json:"proxyFraction,omitempty"`
	// This entire oneOf is deprecated, and the concept is not implemented.
	//
	// Types that are assignable to SendRecordThresholdField:
	//
	//	*CryptoUpdateTransactionBody_SendRecordThreshold
	//	*CryptoUpdateTransactionBody_SendRecordThresholdWrapper
	SendRecordThresholdField isCryptoUpdateTransactionBody_SendRecordThresholdField `protobuf_oneof:"sendRecordThresholdField"`
	// This entire oneOf is deprecated, and the concept is not implemented.
	//
	// Types that are assignable to ReceiveRecordThresholdField:
	//
	//	*CryptoUpdateTransactionBody_ReceiveRecordThreshold
	//	*CryptoUpdateTransactionBody_ReceiveRecordThresholdWrapper
	ReceiveRecordThresholdField isCryptoUpdateTransactionBody_ReceiveRecordThresholdField `protobuf_oneof:"receiveRecordThresholdField"`
	// *
	// A duration to extend account expiration.<br/>
	// An amount of time, in seconds, to extend the expiration date for this
	// account when _automatically_ renewed.
	// <p>
	// This duration MUST be between the current configured minimum and maximum
	// values defined for the network.<br/>
	// This duration SHALL be applied only when _automatically_ extending the
	// account expiration.
	AutoRenewPeriod *Duration `protobuf:"bytes,8,opt,name=autoRenewPeriod,proto3" json:"autoRenewPeriod,omitempty"`
	// *
	// A new account expiration time, in seconds since the epoch.
	// <p>
	// For this purpose, `epoch` SHALL be the UNIX epoch with 0
	// at `1970-01-01T00:00:00.000Z`.<br/>
	// If set, this value MUST be later than the current consensus time.<br/>
	// If set, this value MUST be earlier than the current consensus time
	// extended by the current maximum expiration time configured for the
	// network.
	ExpirationTime *Timestamp `protobuf:"bytes,9,opt,name=expirationTime,proto3" json:"expirationTime,omitempty"`
	// Types that are assignable to ReceiverSigRequiredField:
	//
	//	*CryptoUpdateTransactionBody_ReceiverSigRequired
	//	*CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper
	ReceiverSigRequiredField isCryptoUpdateTransactionBody_ReceiverSigRequiredField `protobuf_oneof:"receiverSigRequiredField"`
	// *
	// A short description of this Account.
	// <p>
	// This value, if set, MUST NOT exceed `transaction.maxMemoUtf8Bytes`
	// (default 100) bytes when encoded as UTF-8.
	Memo *wrapperspb.StringValue `protobuf:"bytes,14,opt,name=memo,proto3" json:"memo,omitempty"`
	// *
	// A maximum number of tokens that can be auto-associated
	// with this account.<br/>
	// By default this value is 0 for all accounts except for automatically
	// created accounts (i.e smart contracts) which default to -1.
	// <p>
	// If this value is `0`, then this account MUST manually associate to
	// a token before holding or transacting in that token.<br/>
	// This value MAY also be `-1` to indicate no limit.<br/>
	// If set, this value MUST NOT be less than `-1`.<br/>
	MaxAutomaticTokenAssociations *wrapperspb.Int32Value `protobuf:"bytes,15,opt,name=max_automatic_token_associations,json=maxAutomaticTokenAssociations,proto3" json:"max_automatic_token_associations,omitempty"`
	// Types that are assignable to StakedId:
	//
	//	*CryptoUpdateTransactionBody_StakedAccountId
	//	*CryptoUpdateTransactionBody_StakedNodeId
	StakedId isCryptoUpdateTransactionBody_StakedId `protobuf_oneof:"staked_id"`
	// *
	// A boolean indicating that this account has chosen to decline rewards for
	// staking its balances.
	// <p>
	// This account MAY still stake its balances, but SHALL NOT receive reward
	// payments for doing so, if this value is set, and `true`.
	DeclineReward *wrapperspb.BoolValue `protobuf:"bytes,18,opt,name=decline_reward,json=declineReward,proto3" json:"decline_reward,omitempty"`
}

func (x *CryptoUpdateTransactionBody) Reset() {
	*x = CryptoUpdateTransactionBody{}
	mi := &file_crypto_update_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CryptoUpdateTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoUpdateTransactionBody) ProtoMessage() {}

func (x *CryptoUpdateTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_update_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoUpdateTransactionBody.ProtoReflect.Descriptor instead.
func (*CryptoUpdateTransactionBody) Descriptor() ([]byte, []int) {
	return file_crypto_update_proto_rawDescGZIP(), []int{0}
}

func (x *CryptoUpdateTransactionBody) GetAccountIDToUpdate() *AccountID {
	if x != nil {
		return x.AccountIDToUpdate
	}
	return nil
}

func (x *CryptoUpdateTransactionBody) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

// Deprecated: Marked as deprecated in crypto_update.proto.
func (x *CryptoUpdateTransactionBody) GetProxyAccountID() *AccountID {
	if x != nil {
		return x.ProxyAccountID
	}
	return nil
}

// Deprecated: Marked as deprecated in crypto_update.proto.
func (x *CryptoUpdateTransactionBody) GetProxyFraction() int32 {
	if x != nil {
		return x.ProxyFraction
	}
	return 0
}

func (m *CryptoUpdateTransactionBody) GetSendRecordThresholdField() isCryptoUpdateTransactionBody_SendRecordThresholdField {
	if m != nil {
		return m.SendRecordThresholdField
	}
	return nil
}

// Deprecated: Marked as deprecated in crypto_update.proto.
func (x *CryptoUpdateTransactionBody) GetSendRecordThreshold() uint64 {
	if x, ok := x.GetSendRecordThresholdField().(*CryptoUpdateTransactionBody_SendRecordThreshold); ok {
		return x.SendRecordThreshold
	}
	return 0
}

// Deprecated: Marked as deprecated in crypto_update.proto.
func (x *CryptoUpdateTransactionBody) GetSendRecordThresholdWrapper() *wrapperspb.UInt64Value {
	if x, ok := x.GetSendRecordThresholdField().(*CryptoUpdateTransactionBody_SendRecordThresholdWrapper); ok {
		return x.SendRecordThresholdWrapper
	}
	return nil
}

func (m *CryptoUpdateTransactionBody) GetReceiveRecordThresholdField() isCryptoUpdateTransactionBody_ReceiveRecordThresholdField {
	if m != nil {
		return m.ReceiveRecordThresholdField
	}
	return nil
}

// Deprecated: Marked as deprecated in crypto_update.proto.
func (x *CryptoUpdateTransactionBody) GetReceiveRecordThreshold() uint64 {
	if x, ok := x.GetReceiveRecordThresholdField().(*CryptoUpdateTransactionBody_ReceiveRecordThreshold); ok {
		return x.ReceiveRecordThreshold
	}
	return 0
}

// Deprecated: Marked as deprecated in crypto_update.proto.
func (x *CryptoUpdateTransactionBody) GetReceiveRecordThresholdWrapper() *wrapperspb.UInt64Value {
	if x, ok := x.GetReceiveRecordThresholdField().(*CryptoUpdateTransactionBody_ReceiveRecordThresholdWrapper); ok {
		return x.ReceiveRecordThresholdWrapper
	}
	return nil
}

func (x *CryptoUpdateTransactionBody) GetAutoRenewPeriod() *Duration {
	if x != nil {
		return x.AutoRenewPeriod
	}
	return nil
}

func (x *CryptoUpdateTransactionBody) GetExpirationTime() *Timestamp {
	if x != nil {
		return x.ExpirationTime
	}
	return nil
}

func (m *CryptoUpdateTransactionBody) GetReceiverSigRequiredField() isCryptoUpdateTransactionBody_ReceiverSigRequiredField {
	if m != nil {
		return m.ReceiverSigRequiredField
	}
	return nil
}

// Deprecated: Marked as deprecated in crypto_update.proto.
func (x *CryptoUpdateTransactionBody) GetReceiverSigRequired() bool {
	if x, ok := x.GetReceiverSigRequiredField().(*CryptoUpdateTransactionBody_ReceiverSigRequired); ok {
		return x.ReceiverSigRequired
	}
	return false
}

func (x *CryptoUpdateTransactionBody) GetReceiverSigRequiredWrapper() *wrapperspb.BoolValue {
	if x, ok := x.GetReceiverSigRequiredField().(*CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper); ok {
		return x.ReceiverSigRequiredWrapper
	}
	return nil
}

func (x *CryptoUpdateTransactionBody) GetMemo() *wrapperspb.StringValue {
	if x != nil {
		return x.Memo
	}
	return nil
}

func (x *CryptoUpdateTransactionBody) GetMaxAutomaticTokenAssociations() *wrapperspb.Int32Value {
	if x != nil {
		return x.MaxAutomaticTokenAssociations
	}
	return nil
}

func (m *CryptoUpdateTransactionBody) GetStakedId() isCryptoUpdateTransactionBody_StakedId {
	if m != nil {
		return m.StakedId
	}
	return nil
}

func (x *CryptoUpdateTransactionBody) GetStakedAccountId() *AccountID {
	if x, ok := x.GetStakedId().(*CryptoUpdateTransactionBody_StakedAccountId); ok {
		return x.StakedAccountId
	}
	return nil
}

func (x *CryptoUpdateTransactionBody) GetStakedNodeId() int64 {
	if x, ok := x.GetStakedId().(*CryptoUpdateTransactionBody_StakedNodeId); ok {
		return x.StakedNodeId
	}
	return 0
}

func (x *CryptoUpdateTransactionBody) GetDeclineReward() *wrapperspb.BoolValue {
	if x != nil {
		return x.DeclineReward
	}
	return nil
}

type isCryptoUpdateTransactionBody_SendRecordThresholdField interface {
	isCryptoUpdateTransactionBody_SendRecordThresholdField()
}

type CryptoUpdateTransactionBody_SendRecordThreshold struct {
	// *
	// Removed prior to the first available history, and may be related
	// to an early design dead-end.<br/>
	// The new threshold amount (in tinybars) for which an account record is
	// created for any send/withdraw transaction
	//
	// Deprecated: Marked as deprecated in crypto_update.proto.
	SendRecordThreshold uint64 `protobuf:"varint,6,opt,name=sendRecordThreshold,proto3,oneof"`
}

type CryptoUpdateTransactionBody_SendRecordThresholdWrapper struct {
	// *
	// Removed prior to the first available history, and may be related
	// to an early design dead-end.<br/>
	// The new threshold amount (in tinybars) for which an account record is
	// created for any send/withdraw transaction
	//
	// Deprecated: Marked as deprecated in crypto_update.proto.
	SendRecordThresholdWrapper *wrapperspb.UInt64Value `protobuf:"bytes,11,opt,name=sendRecordThresholdWrapper,proto3,oneof"`
}

func (*CryptoUpdateTransactionBody_SendRecordThreshold) isCryptoUpdateTransactionBody_SendRecordThresholdField() {
}

func (*CryptoUpdateTransactionBody_SendRecordThresholdWrapper) isCryptoUpdateTransactionBody_SendRecordThresholdField() {
}

type isCryptoUpdateTransactionBody_ReceiveRecordThresholdField interface {
	isCryptoUpdateTransactionBody_ReceiveRecordThresholdField()
}

type CryptoUpdateTransactionBody_ReceiveRecordThreshold struct {
	// *
	// Removed prior to the first available history, and may be related
	// to an early design dead-end.<br/>
	// The new threshold amount (in tinybars) for which an account record is
	// created for any receive/deposit transaction.
	//
	// Deprecated: Marked as deprecated in crypto_update.proto.
	ReceiveRecordThreshold uint64 `protobuf:"varint,7,opt,name=receiveRecordThreshold,proto3,oneof"`
}

type CryptoUpdateTransactionBody_ReceiveRecordThresholdWrapper struct {
	// *
	// Removed prior to the first available history, and may be related
	// to an early design dead-end.<br/>
	// The new threshold amount (in tinybars) for which an account record is
	// created for any receive/deposit transaction.
	//
	// Deprecated: Marked as deprecated in crypto_update.proto.
	ReceiveRecordThresholdWrapper *wrapperspb.UInt64Value `protobuf:"bytes,12,opt,name=receiveRecordThresholdWrapper,proto3,oneof"`
}

func (*CryptoUpdateTransactionBody_ReceiveRecordThreshold) isCryptoUpdateTransactionBody_ReceiveRecordThresholdField() {
}

func (*CryptoUpdateTransactionBody_ReceiveRecordThresholdWrapper) isCryptoUpdateTransactionBody_ReceiveRecordThresholdField() {
}

type isCryptoUpdateTransactionBody_ReceiverSigRequiredField interface {
	isCryptoUpdateTransactionBody_ReceiverSigRequiredField()
}

type CryptoUpdateTransactionBody_ReceiverSigRequired struct {
	// *
	// Removed to distinguish between unset and a default value.<br/>
	// Do NOT use this field to set a false value because the server cannot
	// distinguish from the default value. Use receiverSigRequiredWrapper
	// field for this purpose.
	//
	// Deprecated: Marked as deprecated in crypto_update.proto.
	ReceiverSigRequired bool `protobuf:"varint,10,opt,name=receiverSigRequired,proto3,oneof"`
}

type CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper struct {
	// *
	// A flag indicating the account holder must authorize all incoming
	// token transfers.
	// <p>
	// If this flag is set then any transaction that would result in adding
	// hbar or other tokens to this account balance MUST be signed by the
	// identifying key of this account (that is, the `key` field).
	ReceiverSigRequiredWrapper *wrapperspb.BoolValue `protobuf:"bytes,13,opt,name=receiverSigRequiredWrapper,proto3,oneof"`
}

func (*CryptoUpdateTransactionBody_ReceiverSigRequired) isCryptoUpdateTransactionBody_ReceiverSigRequiredField() {
}

func (*CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper) isCryptoUpdateTransactionBody_ReceiverSigRequiredField() {
}

type isCryptoUpdateTransactionBody_StakedId interface {
	isCryptoUpdateTransactionBody_StakedId()
}

type CryptoUpdateTransactionBody_StakedAccountId struct {
	// *
	// ID of the account to which this account is staking its balances.
	// <p>
	// If this account is not currently staking its balances, then this
	// field, if set, MUST be the sentinel value of `0.0.0`.
	StakedAccountId *AccountID `protobuf:"bytes,16,opt,name=staked_account_id,json=stakedAccountId,proto3,oneof"`
}

type CryptoUpdateTransactionBody_StakedNodeId struct {
	// *
	// ID of the node this account is staked to.
	// <p>
	// If this account is not currently staking its balances, then this
	// field, if set, SHALL be the sentinel value of `-1`.<br/>
	// Wallet software SHOULD surface staking issues to users and provide a
	// simple mechanism to update staking to a new node ID in the event the
	// prior staked node ID ceases to be valid.
	StakedNodeId int64 `protobuf:"varint,17,opt,name=staked_node_id,json=stakedNodeId,proto3,oneof"`
}

func (*CryptoUpdateTransactionBody_StakedAccountId) isCryptoUpdateTransactionBody_StakedId() {}

func (*CryptoUpdateTransactionBody_StakedNodeId) isCryptoUpdateTransactionBody_StakedId() {}

var File_crypto_update_proto protoreflect.FileDescriptor

var file_crypto_update_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61,
	0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0e, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xd9, 0x09, 0x0a, 0x1b, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79,
	0x12, 0x3e, 0x0a, 0x11, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x54, 0x6f, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x11, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x54, 0x6f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x12, 0x1c, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3c,
	0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x42, 0x02, 0x18, 0x01, 0x52, 0x0e, 0x70, 0x72,
	0x6f, 0x78, 0x79, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x28, 0x0a, 0x0d,
	0x70, 0x72, 0x6f, 0x78, 0x79, 0x46, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x02, 0x18, 0x01, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x46, 0x72,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x36, 0x0a, 0x13, 0x73, 0x65, 0x6e, 0x64, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x04, 0x42, 0x02, 0x18, 0x01, 0x48, 0x00, 0x52, 0x13, 0x73, 0x65, 0x6e, 0x64, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x62,
	0x0a, 0x1a, 0x73, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65,
	0x73, 0x68, 0x6f, 0x6c, 0x64, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x42, 0x02, 0x18, 0x01, 0x48, 0x00, 0x52, 0x1a, 0x73, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x57, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x12, 0x3c, 0x0a, 0x16, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x04, 0x42, 0x02, 0x18, 0x01, 0x48, 0x01, 0x52, 0x16, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64,
	0x12, 0x68, 0x0a, 0x1d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x02, 0x18, 0x01, 0x48, 0x01, 0x52, 0x1d, 0x72, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68,
	0x6f, 0x6c, 0x64, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x0f, 0x61, 0x75,
	0x74, 0x6f, 0x52, 0x65, 0x6e, 0x65, 0x77, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x61, 0x75, 0x74, 0x6f, 0x52, 0x65, 0x6e, 0x65, 0x77, 0x50,
	0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x38, 0x0a, 0x0e, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0e, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x36, 0x0a, 0x13, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x53, 0x69, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x42, 0x02, 0x18, 0x01,
	0x48, 0x02, 0x52, 0x13, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x53, 0x69, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x12, 0x5c, 0x0a, 0x1a, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x53, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x57, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f,
	0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x02, 0x52, 0x1a, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x76, 0x65, 0x72, 0x53, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x57, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x30, 0x0a, 0x04, 0x6d, 0x65, 0x6d, 0x6f, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x04, 0x6d, 0x65, 0x6d, 0x6f, 0x12, 0x64, 0x0a, 0x20, 0x6d, 0x61, 0x78, 0x5f, 0x61,
	0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x61,
	0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0f, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x1d,
	0x6d, 0x61, 0x78, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3e, 0x0a,
	0x11, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x64, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x48, 0x03, 0x52, 0x0f, 0x73, 0x74,
	0x61, 0x6b, 0x65, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x26, 0x0a,
	0x0e, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x64, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x11, 0x20, 0x01, 0x28, 0x03, 0x48, 0x03, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x64, 0x4e,
	0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x41, 0x0a, 0x0e, 0x64, 0x65, 0x63, 0x6c, 0x69, 0x6e, 0x65,
	0x5f, 0x72, 0x65, 0x77, 0x61, 0x72, 0x64, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0d, 0x64, 0x65, 0x63, 0x6c, 0x69,
	0x6e, 0x65, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x42, 0x1a, 0x0a, 0x18, 0x73, 0x65, 0x6e, 0x64,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x42, 0x1d, 0x0a, 0x1b, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x42, 0x1a, 0x0a, 0x18, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x53,
	0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x42,
	0x0b, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x64, 0x5f, 0x69, 0x64, 0x42, 0x26, 0x0a, 0x22,
	0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61,
	0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_crypto_update_proto_rawDescOnce sync.Once
	file_crypto_update_proto_rawDescData = file_crypto_update_proto_rawDesc
)

func file_crypto_update_proto_rawDescGZIP() []byte {
	file_crypto_update_proto_rawDescOnce.Do(func() {
		file_crypto_update_proto_rawDescData = protoimpl.X.CompressGZIP(file_crypto_update_proto_rawDescData)
	})
	return file_crypto_update_proto_rawDescData
}

var file_crypto_update_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_crypto_update_proto_goTypes = []any{
	(*CryptoUpdateTransactionBody)(nil), // 0: proto.CryptoUpdateTransactionBody
	(*AccountID)(nil),                   // 1: proto.AccountID
	(*Key)(nil),                         // 2: proto.Key
	(*wrapperspb.UInt64Value)(nil),      // 3: google.protobuf.UInt64Value
	(*Duration)(nil),                    // 4: proto.Duration
	(*Timestamp)(nil),                   // 5: proto.Timestamp
	(*wrapperspb.BoolValue)(nil),        // 6: google.protobuf.BoolValue
	(*wrapperspb.StringValue)(nil),      // 7: google.protobuf.StringValue
	(*wrapperspb.Int32Value)(nil),       // 8: google.protobuf.Int32Value
}
var file_crypto_update_proto_depIdxs = []int32{
	1,  // 0: proto.CryptoUpdateTransactionBody.accountIDToUpdate:type_name -> proto.AccountID
	2,  // 1: proto.CryptoUpdateTransactionBody.key:type_name -> proto.Key
	1,  // 2: proto.CryptoUpdateTransactionBody.proxyAccountID:type_name -> proto.AccountID
	3,  // 3: proto.CryptoUpdateTransactionBody.sendRecordThresholdWrapper:type_name -> google.protobuf.UInt64Value
	3,  // 4: proto.CryptoUpdateTransactionBody.receiveRecordThresholdWrapper:type_name -> google.protobuf.UInt64Value
	4,  // 5: proto.CryptoUpdateTransactionBody.autoRenewPeriod:type_name -> proto.Duration
	5,  // 6: proto.CryptoUpdateTransactionBody.expirationTime:type_name -> proto.Timestamp
	6,  // 7: proto.CryptoUpdateTransactionBody.receiverSigRequiredWrapper:type_name -> google.protobuf.BoolValue
	7,  // 8: proto.CryptoUpdateTransactionBody.memo:type_name -> google.protobuf.StringValue
	8,  // 9: proto.CryptoUpdateTransactionBody.max_automatic_token_associations:type_name -> google.protobuf.Int32Value
	1,  // 10: proto.CryptoUpdateTransactionBody.staked_account_id:type_name -> proto.AccountID
	6,  // 11: proto.CryptoUpdateTransactionBody.decline_reward:type_name -> google.protobuf.BoolValue
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_crypto_update_proto_init() }
func file_crypto_update_proto_init() {
	if File_crypto_update_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_duration_proto_init()
	file_timestamp_proto_init()
	file_crypto_update_proto_msgTypes[0].OneofWrappers = []any{
		(*CryptoUpdateTransactionBody_SendRecordThreshold)(nil),
		(*CryptoUpdateTransactionBody_SendRecordThresholdWrapper)(nil),
		(*CryptoUpdateTransactionBody_ReceiveRecordThreshold)(nil),
		(*CryptoUpdateTransactionBody_ReceiveRecordThresholdWrapper)(nil),
		(*CryptoUpdateTransactionBody_ReceiverSigRequired)(nil),
		(*CryptoUpdateTransactionBody_ReceiverSigRequiredWrapper)(nil),
		(*CryptoUpdateTransactionBody_StakedAccountId)(nil),
		(*CryptoUpdateTransactionBody_StakedNodeId)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_crypto_update_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_crypto_update_proto_goTypes,
		DependencyIndexes: file_crypto_update_proto_depIdxs,
		MessageInfos:      file_crypto_update_proto_msgTypes,
	}.Build()
	File_crypto_update_proto = out.File
	file_crypto_update_proto_rawDesc = nil
	file_crypto_update_proto_goTypes = nil
	file_crypto_update_proto_depIdxs = nil
}
