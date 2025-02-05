//*
// # Get Account Information
// A standard query to inspect the full detail of an account.
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
// source: crypto_get_info.proto

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
// A query to read information for an account.
//
// The returned information SHALL include balance.<br/>
// The returned information SHALL NOT include allowances.<br/>
// The returned information SHALL NOT include token relationships.<br/>
// The returned information SHALL NOT include account records.
type CryptoGetInfoQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// Standard information sent with every query operation.<br/>
	// This includes the signed payment and what kind of response is requested
	// (cost, state proof, both, or neither).
	Header *QueryHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// The account ID for which information is requested
	AccountID *AccountID `protobuf:"bytes,2,opt,name=accountID,proto3" json:"accountID,omitempty"`
}

func (x *CryptoGetInfoQuery) Reset() {
	*x = CryptoGetInfoQuery{}
	mi := &file_crypto_get_info_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CryptoGetInfoQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoGetInfoQuery) ProtoMessage() {}

func (x *CryptoGetInfoQuery) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_get_info_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoGetInfoQuery.ProtoReflect.Descriptor instead.
func (*CryptoGetInfoQuery) Descriptor() ([]byte, []int) {
	return file_crypto_get_info_proto_rawDescGZIP(), []int{0}
}

func (x *CryptoGetInfoQuery) GetHeader() *QueryHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *CryptoGetInfoQuery) GetAccountID() *AccountID {
	if x != nil {
		return x.AccountID
	}
	return nil
}

// *
// Response when the client sends the node CryptoGetInfoQuery
type CryptoGetInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// The standard response information for queries.<br/>
	// This includes the values requested in the `QueryHeader`
	// (cost, state proof, both, or neither).
	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// Details of the account.
	// <p>
	// A state proof MAY be generated for this field.
	AccountInfo *CryptoGetInfoResponse_AccountInfo `protobuf:"bytes,2,opt,name=accountInfo,proto3" json:"accountInfo,omitempty"`
}

func (x *CryptoGetInfoResponse) Reset() {
	*x = CryptoGetInfoResponse{}
	mi := &file_crypto_get_info_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CryptoGetInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoGetInfoResponse) ProtoMessage() {}

func (x *CryptoGetInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_get_info_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoGetInfoResponse.ProtoReflect.Descriptor instead.
func (*CryptoGetInfoResponse) Descriptor() ([]byte, []int) {
	return file_crypto_get_info_proto_rawDescGZIP(), []int{1}
}

func (x *CryptoGetInfoResponse) GetHeader() *ResponseHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *CryptoGetInfoResponse) GetAccountInfo() *CryptoGetInfoResponse_AccountInfo {
	if x != nil {
		return x.AccountInfo
	}
	return nil
}

// *
// Information describing A single Account in the Hedera distributed ledger.
//
// #### Attributes
// Each Account may have a unique three-part identifier, a Key, and one or
// more token balances. Accounts also have an alias, which has multiple
// forms, and may be set automatically. Several additional items are
// associated with the Account to enable full functionality.
//
// #### Expiration
// Accounts, as most items in the network, have an expiration time, recorded
// as a `Timestamp`, and must be "renewed" for a small fee at expiration.
// This helps to reduce the amount of inactive accounts retained in state.
// Another account may be designated to pay any renewal fees and
// automatically renew the account for (by default) 30-90 days at a time as
// a means to optionally ensure important accounts remain active.
//
// ### Staking
// Accounts may participate in securing the network by "staking" the account
// balances to a particular network node, and receive a portion of network
// fees as a reward. An account may optionally decline these rewards but
// still stake its balances.
//
// #### Transfer Restrictions
// An account may optionally require that inbound transfer transactions be
// signed by that account as receiver (in addition to any other signatures
// required, including sender).
type CryptoGetInfoResponse_AccountInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// a unique identifier for this account.
	// <p>
	// An account identifier, when assigned to this field, SHALL be of
	// the form `shard.realm.number`.
	AccountID *AccountID `protobuf:"bytes,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	// *
	// A Solidity ID.
	// <p>
	// This SHALL be populated if this account is a smart contract, and
	// SHALL NOT be populated otherwise.<br/>
	// This SHALL be formatted as a string according to Solidity ID
	// standards.
	ContractAccountID string `protobuf:"bytes,2,opt,name=contractAccountID,proto3" json:"contractAccountID,omitempty"`
	// *
	// A boolean indicating that this account is deleted.
	// <p>
	// Any transaction involving a deleted account SHALL fail.
	Deleted bool `protobuf:"varint,3,opt,name=deleted,proto3" json:"deleted,omitempty"`
	// *
	// Replaced by StakingInfo.<br/>
	// ID of the account to which this account is staking its balances. If
	// this account is not currently staking its balances, then this field,
	// if set, SHALL be the sentinel value of `0.0.0`.
	//
	// Deprecated: Marked as deprecated in crypto_get_info.proto.
	ProxyAccountID *AccountID `protobuf:"bytes,4,opt,name=proxyAccountID,proto3" json:"proxyAccountID,omitempty"`
	// *
	// Replaced by StakingInfo.<br/>
	// The total amount of tinybar proxy staked to this account.
	//
	// Deprecated: Marked as deprecated in crypto_get_info.proto.
	ProxyReceived int64 `protobuf:"varint,6,opt,name=proxyReceived,proto3" json:"proxyReceived,omitempty"`
	// *
	// The key to be used to sign transactions from this account, if any.
	// <p>
	// This key SHALL NOT be set for hollow accounts until the account
	// is finalized.<br/>
	// This key SHALL be set on all other accounts, except for certain
	// immutable accounts (0.0.800 and 0.0.801) necessary for network
	// function and otherwise secured by the governing council.
	Key *Key `protobuf:"bytes,7,opt,name=key,proto3" json:"key,omitempty"`
	// *
	// The HBAR balance of this account, in tinybar (10<sup>-8</sup> HBAR).
	// <p>
	// This value SHALL always be a whole number.
	Balance uint64 `protobuf:"varint,8,opt,name=balance,proto3" json:"balance,omitempty"`
	// *
	// Obsolete and unused.<br/>
	// The threshold amount, in tinybars, at which a record was created for
	// any transaction that decreased the balance of this account.
	//
	// Deprecated: Marked as deprecated in crypto_get_info.proto.
	GenerateSendRecordThreshold uint64 `protobuf:"varint,9,opt,name=generateSendRecordThreshold,proto3" json:"generateSendRecordThreshold,omitempty"`
	// *
	// Obsolete and unused.<br/>
	// The threshold amount, in tinybars, at which a record was created for
	// any transaction that increased the balance of this account.
	//
	// Deprecated: Marked as deprecated in crypto_get_info.proto.
	GenerateReceiveRecordThreshold uint64 `protobuf:"varint,10,opt,name=generateReceiveRecordThreshold,proto3" json:"generateReceiveRecordThreshold,omitempty"`
	// *
	// A boolean indicating that the account requires a receiver signature
	// for inbound token transfer transactions.
	// <p>
	// If this value is `true` then a transaction to transfer tokens to this
	// account SHALL NOT succeed unless this account has signed the
	// transfer transaction.
	ReceiverSigRequired bool `protobuf:"varint,11,opt,name=receiverSigRequired,proto3" json:"receiverSigRequired,omitempty"`
	// *
	// The current expiration time for this account.
	// <p>
	// This account SHALL be due standard renewal fees when the network
	// consensus time exceeds this time.<br/>
	// If rent and expiration are enabled for the network, and automatic
	// renewal is enabled for this account, renewal fees SHALL be charged
	// after this time, and, if charged, the expiration time SHALL be
	// extended for another renewal period.<br/>
	// This account MAY be expired and removed from state at any point
	// after this time if not renewed.<br/>
	// An account holder MAY extend this time by submitting an account
	// update transaction to modify expiration time, subject to the current
	// maximum expiration time for the network.
	ExpirationTime *Timestamp `protobuf:"bytes,12,opt,name=expirationTime,proto3" json:"expirationTime,omitempty"`
	// *
	// A duration to extend this account's expiration.
	// <p>
	// The network SHALL extend the account's expiration by this
	// duration, if funds are available, upon automatic renewal.<br/>
	// This SHALL NOT apply if the account is already deleted
	// upon expiration.<br/>
	// If this is not provided in an allowed range on account creation, the
	// transaction SHALL fail with INVALID_AUTO_RENEWAL_PERIOD. The default
	// values for the minimum period and maximum period are currently
	// 30 days and 90 days, respectively.
	AutoRenewPeriod *Duration `protobuf:"bytes,13,opt,name=autoRenewPeriod,proto3" json:"autoRenewPeriod,omitempty"`
	// *
	// All of the livehashes attached to the account (each of which is a
	// hash along with the keys that authorized it and can delete it)
	LiveHashes []*LiveHash `protobuf:"bytes,14,rep,name=liveHashes,proto3" json:"liveHashes,omitempty"`
	// *
	// As of `HIP-367`, which enabled unlimited token associations, the
	// potential scale for this value requires that users consult a mirror
	// node for this information.
	//
	// Deprecated: Marked as deprecated in crypto_get_info.proto.
	TokenRelationships []*TokenRelationship `protobuf:"bytes,15,rep,name=tokenRelationships,proto3" json:"tokenRelationships,omitempty"`
	// *
	// A short description of this account.
	// <p>
	// This value, if set, MUST NOT exceed `transaction.maxMemoUtf8Bytes`
	// (default 100) bytes when encoded as UTF-8.
	Memo string `protobuf:"bytes,16,opt,name=memo,proto3" json:"memo,omitempty"`
	// *
	// The total number of non-fungible/unique tokens owned by this account.
	OwnedNfts int64 `protobuf:"varint,17,opt,name=ownedNfts,proto3" json:"ownedNfts,omitempty"`
	// *
	// The maximum number of tokens that can be auto-associated with the
	// account.
	// <p>
	// If this is less than or equal to `used_auto_associations` (or 0),
	// then this account MUST manually associate with a token before
	// transacting in that token.<br/>
	// Following HIP-904 This value may also be `-1` to indicate no
	// limit.<br/>
	// This value MUST NOT be less than `-1`.
	MaxAutomaticTokenAssociations int32 `protobuf:"varint,18,opt,name=max_automatic_token_associations,json=maxAutomaticTokenAssociations,proto3" json:"max_automatic_token_associations,omitempty"`
	// *
	// An account alias.<br/>
	// This is a value used in some contexts to reference an account when
	// the tripartite account identifier is not available.
	// <p>
	// This field, when set to a non-default value, is immutable and
	// SHALL NOT be changed.
	Alias []byte `protobuf:"bytes,19,opt,name=alias,proto3" json:"alias,omitempty"`
	// *
	// The ledger ID of the network that generated this response.
	// <p>
	// This value SHALL identify the distributed ledger that responded to
	// this query.
	LedgerId []byte `protobuf:"bytes,20,opt,name=ledger_id,json=ledgerId,proto3" json:"ledger_id,omitempty"`
	// *
	// The ethereum transaction nonce associated with this account.
	EthereumNonce int64 `protobuf:"varint,21,opt,name=ethereum_nonce,json=ethereumNonce,proto3" json:"ethereum_nonce,omitempty"`
	// *
	// Staking information for this account.
	StakingInfo *StakingInfo `protobuf:"bytes,22,opt,name=staking_info,json=stakingInfo,proto3" json:"staking_info,omitempty"`
}

func (x *CryptoGetInfoResponse_AccountInfo) Reset() {
	*x = CryptoGetInfoResponse_AccountInfo{}
	mi := &file_crypto_get_info_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CryptoGetInfoResponse_AccountInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoGetInfoResponse_AccountInfo) ProtoMessage() {}

func (x *CryptoGetInfoResponse_AccountInfo) ProtoReflect() protoreflect.Message {
	mi := &file_crypto_get_info_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoGetInfoResponse_AccountInfo.ProtoReflect.Descriptor instead.
func (*CryptoGetInfoResponse_AccountInfo) Descriptor() ([]byte, []int) {
	return file_crypto_get_info_proto_rawDescGZIP(), []int{1, 0}
}

func (x *CryptoGetInfoResponse_AccountInfo) GetAccountID() *AccountID {
	if x != nil {
		return x.AccountID
	}
	return nil
}

func (x *CryptoGetInfoResponse_AccountInfo) GetContractAccountID() string {
	if x != nil {
		return x.ContractAccountID
	}
	return ""
}

func (x *CryptoGetInfoResponse_AccountInfo) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

// Deprecated: Marked as deprecated in crypto_get_info.proto.
func (x *CryptoGetInfoResponse_AccountInfo) GetProxyAccountID() *AccountID {
	if x != nil {
		return x.ProxyAccountID
	}
	return nil
}

// Deprecated: Marked as deprecated in crypto_get_info.proto.
func (x *CryptoGetInfoResponse_AccountInfo) GetProxyReceived() int64 {
	if x != nil {
		return x.ProxyReceived
	}
	return 0
}

func (x *CryptoGetInfoResponse_AccountInfo) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *CryptoGetInfoResponse_AccountInfo) GetBalance() uint64 {
	if x != nil {
		return x.Balance
	}
	return 0
}

// Deprecated: Marked as deprecated in crypto_get_info.proto.
func (x *CryptoGetInfoResponse_AccountInfo) GetGenerateSendRecordThreshold() uint64 {
	if x != nil {
		return x.GenerateSendRecordThreshold
	}
	return 0
}

// Deprecated: Marked as deprecated in crypto_get_info.proto.
func (x *CryptoGetInfoResponse_AccountInfo) GetGenerateReceiveRecordThreshold() uint64 {
	if x != nil {
		return x.GenerateReceiveRecordThreshold
	}
	return 0
}

func (x *CryptoGetInfoResponse_AccountInfo) GetReceiverSigRequired() bool {
	if x != nil {
		return x.ReceiverSigRequired
	}
	return false
}

func (x *CryptoGetInfoResponse_AccountInfo) GetExpirationTime() *Timestamp {
	if x != nil {
		return x.ExpirationTime
	}
	return nil
}

func (x *CryptoGetInfoResponse_AccountInfo) GetAutoRenewPeriod() *Duration {
	if x != nil {
		return x.AutoRenewPeriod
	}
	return nil
}

func (x *CryptoGetInfoResponse_AccountInfo) GetLiveHashes() []*LiveHash {
	if x != nil {
		return x.LiveHashes
	}
	return nil
}

// Deprecated: Marked as deprecated in crypto_get_info.proto.
func (x *CryptoGetInfoResponse_AccountInfo) GetTokenRelationships() []*TokenRelationship {
	if x != nil {
		return x.TokenRelationships
	}
	return nil
}

func (x *CryptoGetInfoResponse_AccountInfo) GetMemo() string {
	if x != nil {
		return x.Memo
	}
	return ""
}

func (x *CryptoGetInfoResponse_AccountInfo) GetOwnedNfts() int64 {
	if x != nil {
		return x.OwnedNfts
	}
	return 0
}

func (x *CryptoGetInfoResponse_AccountInfo) GetMaxAutomaticTokenAssociations() int32 {
	if x != nil {
		return x.MaxAutomaticTokenAssociations
	}
	return 0
}

func (x *CryptoGetInfoResponse_AccountInfo) GetAlias() []byte {
	if x != nil {
		return x.Alias
	}
	return nil
}

func (x *CryptoGetInfoResponse_AccountInfo) GetLedgerId() []byte {
	if x != nil {
		return x.LedgerId
	}
	return nil
}

func (x *CryptoGetInfoResponse_AccountInfo) GetEthereumNonce() int64 {
	if x != nil {
		return x.EthereumNonce
	}
	return 0
}

func (x *CryptoGetInfoResponse_AccountInfo) GetStakingInfo() *StakingInfo {
	if x != nil {
		return x.StakingInfo
	}
	return nil
}

var File_crypto_get_info_proto protoreflect.FileDescriptor

var file_crypto_get_info_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x5f, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0e, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x12, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x63,
	0x72, 0x79, 0x70, 0x74, 0x6f, 0x5f, 0x61, 0x64, 0x64, 0x5f, 0x6c, 0x69, 0x76, 0x65, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x70, 0x0a, 0x12, 0x43, 0x72, 0x79,
	0x70, 0x74, 0x6f, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12,
	0x2a, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x09, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44,
	0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x22, 0x84, 0x09, 0x0a, 0x15,
	0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x12, 0x4a, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x1a, 0xef, 0x07, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x2e, 0x0a, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44,
	0x12, 0x2c, 0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x18,
	0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x3c, 0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x78,
	0x79, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x49, 0x44, 0x42, 0x02, 0x18, 0x01, 0x52, 0x0e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x28, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x42, 0x02, 0x18,
	0x01, 0x52, 0x0d, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64,
	0x12, 0x1c, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x18,
	0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x1b, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x68,
	0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x18,
	0x01, 0x52, 0x1b, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x4a,
	0x0a, 0x1e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x18, 0x01, 0x52, 0x1e, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x12, 0x30, 0x0a, 0x13, 0x72, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x53, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x72, 0x53, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x12, 0x38, 0x0a, 0x0e,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x39, 0x0a, 0x0f, 0x61, 0x75, 0x74, 0x6f, 0x52, 0x65,
	0x6e, 0x65, 0x77, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0f, 0x61, 0x75, 0x74, 0x6f, 0x52, 0x65, 0x6e, 0x65, 0x77, 0x50, 0x65, 0x72, 0x69, 0x6f,
	0x64, 0x12, 0x2f, 0x0a, 0x0a, 0x6c, 0x69, 0x76, 0x65, 0x48, 0x61, 0x73, 0x68, 0x65, 0x73, 0x18,
	0x0e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69,
	0x76, 0x65, 0x48, 0x61, 0x73, 0x68, 0x52, 0x0a, 0x6c, 0x69, 0x76, 0x65, 0x48, 0x61, 0x73, 0x68,
	0x65, 0x73, 0x12, 0x4c, 0x0a, 0x12, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x73, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x42, 0x02, 0x18, 0x01, 0x52, 0x12, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x6d, 0x65, 0x6d, 0x6f, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6d, 0x65, 0x6d, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x77, 0x6e, 0x65, 0x64, 0x4e, 0x66, 0x74,
	0x73, 0x18, 0x11, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6f, 0x77, 0x6e, 0x65, 0x64, 0x4e, 0x66,
	0x74, 0x73, 0x12, 0x47, 0x0a, 0x20, 0x6d, 0x61, 0x78, 0x5f, 0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61,
	0x74, 0x69, 0x63, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x12, 0x20, 0x01, 0x28, 0x05, 0x52, 0x1d, 0x6d, 0x61,
	0x78, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x41,
	0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x61,
	0x6c, 0x69, 0x61, 0x73, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x14,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x49, 0x64, 0x12, 0x25,
	0x0a, 0x0e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x5f, 0x6e, 0x6f, 0x6e, 0x63, 0x65,
	0x18, 0x15, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d,
	0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x6b, 0x69, 0x6e, 0x67,
	0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x0b, 0x73, 0x74, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x4a, 0x04, 0x08, 0x05,
	0x10, 0x06, 0x42, 0x26, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61,
	0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_crypto_get_info_proto_rawDescOnce sync.Once
	file_crypto_get_info_proto_rawDescData = file_crypto_get_info_proto_rawDesc
)

func file_crypto_get_info_proto_rawDescGZIP() []byte {
	file_crypto_get_info_proto_rawDescOnce.Do(func() {
		file_crypto_get_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_crypto_get_info_proto_rawDescData)
	})
	return file_crypto_get_info_proto_rawDescData
}

var file_crypto_get_info_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_crypto_get_info_proto_goTypes = []any{
	(*CryptoGetInfoQuery)(nil),                // 0: proto.CryptoGetInfoQuery
	(*CryptoGetInfoResponse)(nil),             // 1: proto.CryptoGetInfoResponse
	(*CryptoGetInfoResponse_AccountInfo)(nil), // 2: proto.CryptoGetInfoResponse.AccountInfo
	(*QueryHeader)(nil),                       // 3: proto.QueryHeader
	(*AccountID)(nil),                         // 4: proto.AccountID
	(*ResponseHeader)(nil),                    // 5: proto.ResponseHeader
	(*Key)(nil),                               // 6: proto.Key
	(*Timestamp)(nil),                         // 7: proto.Timestamp
	(*Duration)(nil),                          // 8: proto.Duration
	(*LiveHash)(nil),                          // 9: proto.LiveHash
	(*TokenRelationship)(nil),                 // 10: proto.TokenRelationship
	(*StakingInfo)(nil),                       // 11: proto.StakingInfo
}
var file_crypto_get_info_proto_depIdxs = []int32{
	3,  // 0: proto.CryptoGetInfoQuery.header:type_name -> proto.QueryHeader
	4,  // 1: proto.CryptoGetInfoQuery.accountID:type_name -> proto.AccountID
	5,  // 2: proto.CryptoGetInfoResponse.header:type_name -> proto.ResponseHeader
	2,  // 3: proto.CryptoGetInfoResponse.accountInfo:type_name -> proto.CryptoGetInfoResponse.AccountInfo
	4,  // 4: proto.CryptoGetInfoResponse.AccountInfo.accountID:type_name -> proto.AccountID
	4,  // 5: proto.CryptoGetInfoResponse.AccountInfo.proxyAccountID:type_name -> proto.AccountID
	6,  // 6: proto.CryptoGetInfoResponse.AccountInfo.key:type_name -> proto.Key
	7,  // 7: proto.CryptoGetInfoResponse.AccountInfo.expirationTime:type_name -> proto.Timestamp
	8,  // 8: proto.CryptoGetInfoResponse.AccountInfo.autoRenewPeriod:type_name -> proto.Duration
	9,  // 9: proto.CryptoGetInfoResponse.AccountInfo.liveHashes:type_name -> proto.LiveHash
	10, // 10: proto.CryptoGetInfoResponse.AccountInfo.tokenRelationships:type_name -> proto.TokenRelationship
	11, // 11: proto.CryptoGetInfoResponse.AccountInfo.staking_info:type_name -> proto.StakingInfo
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_crypto_get_info_proto_init() }
func file_crypto_get_info_proto_init() {
	if File_crypto_get_info_proto != nil {
		return
	}
	file_timestamp_proto_init()
	file_duration_proto_init()
	file_basic_types_proto_init()
	file_query_header_proto_init()
	file_response_header_proto_init()
	file_crypto_add_live_hash_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_crypto_get_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_crypto_get_info_proto_goTypes,
		DependencyIndexes: file_crypto_get_info_proto_depIdxs,
		MessageInfos:      file_crypto_get_info_proto_msgTypes,
	}.Build()
	File_crypto_get_info_proto = out.File
	file_crypto_get_info_proto_rawDesc = nil
	file_crypto_get_info_proto_goTypes = nil
	file_crypto_get_info_proto_depIdxs = nil
}
