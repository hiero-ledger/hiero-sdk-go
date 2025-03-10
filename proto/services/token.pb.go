//*
// ## Token
// Tokens represent both fungible and non-fungible units of exchange.
// The `Token` here represents a token within the network state.
//
// ### Keywords
// The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
// "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this
// document are to be interpreted as described in [RFC2119](https://www.ietf.org/rfc/rfc2119)
// and clarified in [RFC8174](https://www.ietf.org/rfc/rfc8174).

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: token.proto

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
// An Hedera Token Service(HTS) token.
//
// A token SHALL represent a fungible or non-fungible unit of exchange.<br/>
// The specified Treasury Account SHALL receive the initial supply of tokens and
// SHALL determine distribution of all tokens once minted.
type Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// A unique identifier for this token.
	TokenId *TokenID `protobuf:"bytes,1,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
	// *
	// A human-readable name for this token.
	// <p>
	// This value MAY NOT be unique.<br/>
	// This value SHALL NOT exceed 100 bytes when encoded as UTF-8.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// *
	// A human-readable symbol for the token.
	// <p>
	// This value SHALL NOT be unique.<br/>
	// This value SHALL NOT exceed 100 bytes when encoded as UTF-8.
	Symbol string `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
	// *
	// A number of decimal places for this token.
	// <p>
	// If decimals are 8 or 11, then the number of whole tokens can be at most
	// billions or millions, respectively. More decimals allows for a more
	// finely-divided token, but also limits the maximum total supply.
	// <p>
	// Examples
	// <ul>
	//
	//	<li>Bitcoin satoshis (21 million whole tokens with 8 decimals).</li>
	//	<li>Hedera tinybar (50 billion whole tokens with 8 decimals).</li>
	//	<li>Bitcoin milli-satoshis (21 million whole tokens with 11 decimals).</li>
	//	<li>Theoretical limit is roughly 92.2 billion with 8 decimals, or
	//	    92.2 million with 11 decimals.</li>
	//
	// </ul>
	// All token amounts in the network are stored as integer amounts, with each
	// unit representing 10<sup>-decimals</sup> whole tokens.
	// <p>
	// For tokens with `token_type` set to `NON_FUNGIBLE_UNIQUE` this MUST be 0.
	Decimals int32 `protobuf:"varint,4,opt,name=decimals,proto3" json:"decimals,omitempty"`
	// *
	// A _current_ total supply of this token, expressed in the smallest unit
	// of the token.
	// <p>
	// The number of _whole_ tokens this represents is (total_supply /
	// 10<sup>decimals</sup>). The value of total supply, MUST be within the
	// positive range of a twos-compliment signed 64-bit integer.
	// The `total_supply`, therefore MUST be between 1, and
	// 9,223,372,036,854,775,807, inclusive.
	// <p>
	// This value SHALL be reduced when a `token_burn` or `token_wipe_account`
	// operation is executed, and SHALL be increased when a `token_mint`
	// operation is executed.
	TotalSupply int64 `protobuf:"varint,5,opt,name=total_supply,json=totalSupply,proto3" json:"total_supply,omitempty"`
	// *
	// A treasury account identifier for this token.
	// <p>
	// When the token is created, the initial supply given in the token create
	// transaction SHALL be minted and deposited in the treasury account.<br/>
	// All token mint transactions for this token SHALL deposit the new minted
	// tokens in the treasury account.<br/>
	// All token burn transactions for this token SHALL remove the tokens to be
	// burned from the treasury account.
	TreasuryAccountId *AccountID `protobuf:"bytes,6,opt,name=treasury_account_id,json=treasuryAccountId,proto3" json:"treasury_account_id,omitempty"`
	// *
	// Access control for general modification of this token.
	// <p>
	// This key MUST sign any `token_update` transaction that
	// changes any attribute of the token other than expiration_time.
	// Other attributes of this token MAY be changed by transactions other than
	// `token_update`, and MUST be signed by one of the other purpose-specific
	// keys assigned to the token.<br/>
	// This value can be set during token creation, and SHALL NOT be
	// modified thereafter, unless the update transaction is signed by both
	// the existing `admin_key` and the new `admin_key`.<br/>
	// If the `admin_key` is not set for a token, that token SHALL be immutable.
	AdminKey *Key `protobuf:"bytes,7,opt,name=admin_key,json=adminKey,proto3" json:"admin_key,omitempty"`
	// *
	// Access control for KYC for this token.
	// <p>
	// Know Your Customer (KYC) status may be granted for an account by a token
	// grant kyc transaction signed by this key.<br/>
	// If this key is not set, then KYC status cannot be granted to an account
	// for this token, and any `TokenGrantKyc` transaction attempting to grant
	// kyc to an account for this token SHALL NOT succeed.<br/>
	// This key MAY be set when the token is created, and MAY be set or modified
	// via a token update transaction signed by the `admin_key`.<br/>
	// If `admin_key` is not set, this value, whether set or unset,
	// SHALL be immutable.
	KycKey *Key `protobuf:"bytes,8,opt,name=kyc_key,json=kycKey,proto3" json:"kyc_key,omitempty"`
	// *
	// Access control to freeze this token.
	// <p>
	// A token may be frozen for an account, preventing any transaction from
	// transferring that token for that specified account, by a token freeze
	// account transaction signed by this key.<br/>
	// If this key is not set, the token cannot be frozen, and any transaction
	// attempting to freeze the token for an account SHALL NOT succeed.<br/>
	// This key MAY be set when the token is created, and MAY be set or modified
	// via a token update transaction signed by the `admin_key`.<br/>
	// If `admin_key` is not set, this value, whether set or unset,
	// SHALL be immutable.
	FreezeKey *Key `protobuf:"bytes,9,opt,name=freeze_key,json=freezeKey,proto3" json:"freeze_key,omitempty"`
	// *
	// Access control of account wipe for this token.
	// <p>
	// A token may be wiped, removing and burning tokens from a specific
	// account, by a token wipe transaction, which MUST be signed by this key.
	// The `treasury_account` cannot be subjected to a token wipe. A token burn
	// transaction, signed by the `supply_key`, serves to burn tokens held by
	// the `treasury_account` instead.<br/>
	// If this key is not set, the token cannot be wiped, and any transaction
	// attempting to wipe the token from an account SHALL NOT succeed.<br/>
	// This key MAY be set when the token is created, and MAY be set or modified
	// via a token update transaction signed by the `admin_key`.<br/>
	// If `admin_key` is not set, this value, whether set or unset,
	// SHALL be immutable.
	WipeKey *Key `protobuf:"bytes,10,opt,name=wipe_key,json=wipeKey,proto3" json:"wipe_key,omitempty"`
	// *
	// Access control of token mint/burn for this token.
	// <p>
	// A token mint transaction MUST be signed by this key, and any token mint
	// transaction not signed by the current `supply_key` for that token
	// SHALL NOT succeed.<br/>
	// A token burn transaction MUST be signed by this key, and any token burn
	// transaction not signed by the current `supply_key` for that token
	// SHALL NOT succeed.<br/>
	// This key MAY be set when the token is created, and MAY be set or modified
	// via a token update transaction signed by the `admin_key`.<br/>
	// If `admin_key` is not set, this value, whether set or unset,
	// SHALL be immutable.
	SupplyKey *Key `protobuf:"bytes,11,opt,name=supply_key,json=supplyKey,proto3" json:"supply_key,omitempty"`
	// *
	// Access control of the `custom_fees` field for this token.
	// <p>
	// The token custom fee schedule may be changed, modifying the fees charged
	// for transferring that token, by a token update transaction, which MUST
	// be signed by this key.<br/>
	// If this key is not set, the token custom fee schedule cannot be changed,
	// and any transaction attempting to change the custom fee schedule for
	// this token SHALL NOT succeed.<br/>
	// This key MAY be set when the token is created, and MAY be set or modified
	// via a token update transaction signed by the `admin_key`.<br/>
	// If `admin_key` is not set, this value, whether set or unset,
	// SHALL be immutable.
	FeeScheduleKey *Key `protobuf:"bytes,12,opt,name=fee_schedule_key,json=feeScheduleKey,proto3" json:"fee_schedule_key,omitempty"`
	// *
	// Access control of pause/unpause for this token.
	// <p>
	// A token may be paused, preventing any transaction from transferring that
	// token, by a token update transaction signed by this key.<br/>
	// If this key is not set, the token cannot be paused, and any transaction
	// attempting to pause the token SHALL NOT succeed.<br/>
	// This key MAY be set when the token is created, and MAY be set or modified
	// via a token update transaction signed by the `admin_key`.<br/>
	// If `admin_key` is not set, this value, whether set or unset,
	// SHALL be immutable.
	PauseKey *Key `protobuf:"bytes,13,opt,name=pause_key,json=pauseKey,proto3" json:"pause_key,omitempty"`
	// *
	// A last used serial number for this token.
	// <p>
	// This SHALL apply only to non-fungible tokens.<br/>
	// When a new NFT is minted, the serial number to apply SHALL be calculated
	// from this value.
	LastUsedSerialNumber int64 `protobuf:"varint,14,opt,name=last_used_serial_number,json=lastUsedSerialNumber,proto3" json:"last_used_serial_number,omitempty"`
	// *
	// A flag indicating that this token is deleted.
	// <p>
	// A transaction involving a deleted token MUST NOT succeed.
	Deleted bool `protobuf:"varint,15,opt,name=deleted,proto3" json:"deleted,omitempty"`
	// *
	// A type for this token.
	// <p>
	// A token SHALL be either `FUNGIBLE_COMMON` or `NON_FUNGIBLE_UNIQUE`.<br/>
	// If this value was omitted during token creation, `FUNGIBLE_COMMON`
	// SHALL be used.
	TokenType TokenType `protobuf:"varint,16,opt,name=token_type,json=tokenType,proto3,enum=proto.TokenType" json:"token_type,omitempty"`
	// *
	// A supply type for this token.
	// <p>
	// A token SHALL have either `INFINITE` or `FINITE` supply type.<br/>
	// If this value was omitted during token creation, the value `INFINITE`
	// SHALL be used.
	SupplyType TokenSupplyType `protobuf:"varint,17,opt,name=supply_type,json=supplyType,proto3,enum=proto.TokenSupplyType" json:"supply_type,omitempty"`
	// *
	// An identifier for the account (if any) that the network will attempt
	// to charge for this token's auto-renewal upon expiration.
	// <p>
	// This field is OPTIONAL. If it is not set then renewal fees SHALL be
	// charged to the account identified by `treasury_account_id`.
	AutoRenewAccountId *AccountID `protobuf:"bytes,18,opt,name=auto_renew_account_id,json=autoRenewAccountId,proto3" json:"auto_renew_account_id,omitempty"`
	// *
	// A number of seconds by which the network should automatically extend
	// this token's expiration.
	// <p>
	// If the token has a valid auto-renew account, and is not deleted upon
	// expiration, the network SHALL attempt to automatically renew this
	// token.<br/>
	// If this is not provided in an allowed range on token creation, the
	// transaction SHALL fail with `INVALID_AUTO_RENEWAL_PERIOD`.<br/>
	// The default values for the minimum period and maximum period are 30 days
	// and 90 days, respectively.
	AutoRenewSeconds int64 `protobuf:"varint,19,opt,name=auto_renew_seconds,json=autoRenewSeconds,proto3" json:"auto_renew_seconds,omitempty"`
	// *
	// An expiration time for this token, in seconds since the epoch.
	// <p>
	// For this purpose, `epoch` SHALL be the
	// UNIX epoch with 0 at `1970-01-01T00:00:00.000Z`.
	ExpirationSecond int64 `protobuf:"varint,20,opt,name=expiration_second,json=expirationSecond,proto3" json:"expiration_second,omitempty"`
	// *
	// A short description of this token.
	// <p>
	// This value, if set, MUST NOT exceed `transaction.maxMemoUtf8Bytes`
	// (default 100) bytes when encoded as UTF-8.
	Memo string `protobuf:"bytes,21,opt,name=memo,proto3" json:"memo,omitempty"`
	// *
	// A maximum supply of this token.<br/>
	// This is the maximum number of tokens of this type that may be issued.
	// <p>
	// This limit SHALL apply regardless of `token_type`.<br/>
	// If `supply_type` is `INFINITE` then this value MUST be 0.<br/>
	// If `supply_type` is `FINITE`, then this value MUST be greater than 0.
	MaxSupply int64 `protobuf:"varint,22,opt,name=max_supply,json=maxSupply,proto3" json:"max_supply,omitempty"`
	// *
	// A flag indicating that this token is paused.
	// <p>
	// A transaction involving a paused token, other than token_unpause,
	// MUST NOT succeed.
	Paused bool `protobuf:"varint,23,opt,name=paused,proto3" json:"paused,omitempty"`
	// *
	// A flag indicating that accounts associated to this token are frozen by
	// default.
	// <p>
	// Accounts newly associated with this token CANNOT transact in the token
	// until unfrozen.<br/>
	// This SHALL NOT prevent a `tokenReject` transaction to return the tokens
	// from an account to the treasury account.
	AccountsFrozenByDefault bool `protobuf:"varint,24,opt,name=accounts_frozen_by_default,json=accountsFrozenByDefault,proto3" json:"accounts_frozen_by_default,omitempty"`
	// *
	// A flag indicating that accounts associated with this token are granted
	// KYC by default.
	AccountsKycGrantedByDefault bool `protobuf:"varint,25,opt,name=accounts_kyc_granted_by_default,json=accountsKycGrantedByDefault,proto3" json:"accounts_kyc_granted_by_default,omitempty"`
	// *
	// A custom fee schedule for this token.
	CustomFees []*CustomFee `protobuf:"bytes,26,rep,name=custom_fees,json=customFees,proto3" json:"custom_fees,omitempty"`
	// *
	// A Token "Metadata".
	// <p>
	// This value, if set, SHALL NOT exceed 100 bytes.
	Metadata []byte `protobuf:"bytes,27,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// *
	// Access Control of metadata update for this token.
	// <p>
	// A transaction to update the `metadata` field of this token MUST be
	// signed by this key.<br/>
	// If this token is a non-fungible/unique token type, a transaction to
	// update the `metadata` field of any individual serialized unique token
	// of this type MUST be signed by this key.<br/>
	// If this key is not set, the token metadata SHALL NOT be changed after it
	// is created.<br/>
	// If this key is not set, the metadata for any individual serialized token
	// of this type SHALL NOT be changed after it is created.<br/>
	// This key MAY be set when the token is created, and MAY be set or modified
	// via a token update transaction signed by the `admin_key`.<br/>
	// If `admin_key` is not set, this value, whether set or unset,
	// SHALL be immutable.
	MetadataKey *Key `protobuf:"bytes,28,opt,name=metadata_key,json=metadataKey,proto3" json:"metadata_key,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	mi := &file_token_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_token_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_token_proto_rawDescGZIP(), []int{0}
}

func (x *Token) GetTokenId() *TokenID {
	if x != nil {
		return x.TokenId
	}
	return nil
}

func (x *Token) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Token) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Token) GetDecimals() int32 {
	if x != nil {
		return x.Decimals
	}
	return 0
}

func (x *Token) GetTotalSupply() int64 {
	if x != nil {
		return x.TotalSupply
	}
	return 0
}

func (x *Token) GetTreasuryAccountId() *AccountID {
	if x != nil {
		return x.TreasuryAccountId
	}
	return nil
}

func (x *Token) GetAdminKey() *Key {
	if x != nil {
		return x.AdminKey
	}
	return nil
}

func (x *Token) GetKycKey() *Key {
	if x != nil {
		return x.KycKey
	}
	return nil
}

func (x *Token) GetFreezeKey() *Key {
	if x != nil {
		return x.FreezeKey
	}
	return nil
}

func (x *Token) GetWipeKey() *Key {
	if x != nil {
		return x.WipeKey
	}
	return nil
}

func (x *Token) GetSupplyKey() *Key {
	if x != nil {
		return x.SupplyKey
	}
	return nil
}

func (x *Token) GetFeeScheduleKey() *Key {
	if x != nil {
		return x.FeeScheduleKey
	}
	return nil
}

func (x *Token) GetPauseKey() *Key {
	if x != nil {
		return x.PauseKey
	}
	return nil
}

func (x *Token) GetLastUsedSerialNumber() int64 {
	if x != nil {
		return x.LastUsedSerialNumber
	}
	return 0
}

func (x *Token) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

func (x *Token) GetTokenType() TokenType {
	if x != nil {
		return x.TokenType
	}
	return TokenType_FUNGIBLE_COMMON
}

func (x *Token) GetSupplyType() TokenSupplyType {
	if x != nil {
		return x.SupplyType
	}
	return TokenSupplyType_INFINITE
}

func (x *Token) GetAutoRenewAccountId() *AccountID {
	if x != nil {
		return x.AutoRenewAccountId
	}
	return nil
}

func (x *Token) GetAutoRenewSeconds() int64 {
	if x != nil {
		return x.AutoRenewSeconds
	}
	return 0
}

func (x *Token) GetExpirationSecond() int64 {
	if x != nil {
		return x.ExpirationSecond
	}
	return 0
}

func (x *Token) GetMemo() string {
	if x != nil {
		return x.Memo
	}
	return ""
}

func (x *Token) GetMaxSupply() int64 {
	if x != nil {
		return x.MaxSupply
	}
	return 0
}

func (x *Token) GetPaused() bool {
	if x != nil {
		return x.Paused
	}
	return false
}

func (x *Token) GetAccountsFrozenByDefault() bool {
	if x != nil {
		return x.AccountsFrozenByDefault
	}
	return false
}

func (x *Token) GetAccountsKycGrantedByDefault() bool {
	if x != nil {
		return x.AccountsKycGrantedByDefault
	}
	return false
}

func (x *Token) GetCustomFees() []*CustomFee {
	if x != nil {
		return x.CustomFees
	}
	return nil
}

func (x *Token) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Token) GetMetadataKey() *Key {
	if x != nil {
		return x.MetadataKey
	}
	return nil
}

var File_token_proto protoreflect.FileDescriptor

var file_token_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f,
	0x66, 0x65, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb0, 0x09, 0x0a, 0x05, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x29, 0x0a, 0x08, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x44, 0x52, 0x07, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x64,
	0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64,
	0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x5f, 0x73, 0x75, 0x70, 0x70, 0x6c, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x53, 0x75, 0x70, 0x70, 0x6c, 0x79, 0x12, 0x40, 0x0a, 0x13, 0x74, 0x72,
	0x65, 0x61, 0x73, 0x75, 0x72, 0x79, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x11, 0x74, 0x72, 0x65, 0x61, 0x73,
	0x75, 0x72, 0x79, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x09,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x08, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x07, 0x6b, 0x79, 0x63, 0x5f, 0x6b, 0x65, 0x79,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b,
	0x65, 0x79, 0x52, 0x06, 0x6b, 0x79, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x0a, 0x66, 0x72,
	0x65, 0x65, 0x7a, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x09, 0x66, 0x72, 0x65, 0x65,
	0x7a, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x08, 0x77, 0x69, 0x70, 0x65, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4b, 0x65, 0x79, 0x52, 0x07, 0x77, 0x69, 0x70, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x0a,
	0x73, 0x75, 0x70, 0x70, 0x6c, 0x79, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x09, 0x73, 0x75,
	0x70, 0x70, 0x6c, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x34, 0x0a, 0x10, 0x66, 0x65, 0x65, 0x5f, 0x73,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x0e, 0x66,
	0x65, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x27, 0x0a,
	0x09, 0x70, 0x61, 0x75, 0x73, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x08, 0x70, 0x61,
	0x75, 0x73, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x35, 0x0a, 0x17, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75,
	0x73, 0x65, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x14, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x73, 0x65,
	0x64, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x2f, 0x0a, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x37, 0x0a, 0x0b, 0x73, 0x75, 0x70, 0x70,
	0x6c, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x75, 0x70, 0x70, 0x6c,
	0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x73, 0x75, 0x70, 0x70, 0x6c, 0x79, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x43, 0x0a, 0x15, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x72, 0x65, 0x6e, 0x65, 0x77, 0x5f,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x49, 0x44, 0x52, 0x12, 0x61, 0x75, 0x74, 0x6f, 0x52, 0x65, 0x6e, 0x65, 0x77, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x12, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x72,
	0x65, 0x6e, 0x65, 0x77, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x13, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x10, 0x61, 0x75, 0x74, 0x6f, 0x52, 0x65, 0x6e, 0x65, 0x77, 0x53, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x18, 0x14, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x10, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x65, 0x6d, 0x6f, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6d, 0x65, 0x6d, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x75, 0x70,
	0x70, 0x6c, 0x79, 0x18, 0x16, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6d, 0x61, 0x78, 0x53, 0x75,
	0x70, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x75, 0x73, 0x65, 0x64, 0x18, 0x17,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x70, 0x61, 0x75, 0x73, 0x65, 0x64, 0x12, 0x3b, 0x0a, 0x1a,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x5f, 0x66, 0x72, 0x6f, 0x7a, 0x65, 0x6e, 0x5f,
	0x62, 0x79, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18, 0x18, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x17, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x46, 0x72, 0x6f, 0x7a, 0x65, 0x6e,
	0x42, 0x79, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x12, 0x44, 0x0a, 0x1f, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x73, 0x5f, 0x6b, 0x79, 0x63, 0x5f, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x65,
	0x64, 0x5f, 0x62, 0x79, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18, 0x19, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x1b, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x4b, 0x79, 0x63, 0x47,
	0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x42, 0x79, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x12,
	0x31, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x66, 0x65, 0x65, 0x73, 0x18, 0x1a,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x46, 0x65, 0x65, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x46, 0x65,
	0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x1b,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2d,
	0x0a, 0x0c, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x1c,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79,
	0x52, 0x0b, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x42, 0x26, 0x0a,
	0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68, 0x67,
	0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6a,
	0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_token_proto_rawDescOnce sync.Once
	file_token_proto_rawDescData = file_token_proto_rawDesc
)

func file_token_proto_rawDescGZIP() []byte {
	file_token_proto_rawDescOnce.Do(func() {
		file_token_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_proto_rawDescData)
	})
	return file_token_proto_rawDescData
}

var file_token_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_token_proto_goTypes = []any{
	(*Token)(nil),        // 0: proto.Token
	(*TokenID)(nil),      // 1: proto.TokenID
	(*AccountID)(nil),    // 2: proto.AccountID
	(*Key)(nil),          // 3: proto.Key
	(TokenType)(0),       // 4: proto.TokenType
	(TokenSupplyType)(0), // 5: proto.TokenSupplyType
	(*CustomFee)(nil),    // 6: proto.CustomFee
}
var file_token_proto_depIdxs = []int32{
	1,  // 0: proto.Token.token_id:type_name -> proto.TokenID
	2,  // 1: proto.Token.treasury_account_id:type_name -> proto.AccountID
	3,  // 2: proto.Token.admin_key:type_name -> proto.Key
	3,  // 3: proto.Token.kyc_key:type_name -> proto.Key
	3,  // 4: proto.Token.freeze_key:type_name -> proto.Key
	3,  // 5: proto.Token.wipe_key:type_name -> proto.Key
	3,  // 6: proto.Token.supply_key:type_name -> proto.Key
	3,  // 7: proto.Token.fee_schedule_key:type_name -> proto.Key
	3,  // 8: proto.Token.pause_key:type_name -> proto.Key
	4,  // 9: proto.Token.token_type:type_name -> proto.TokenType
	5,  // 10: proto.Token.supply_type:type_name -> proto.TokenSupplyType
	2,  // 11: proto.Token.auto_renew_account_id:type_name -> proto.AccountID
	6,  // 12: proto.Token.custom_fees:type_name -> proto.CustomFee
	3,  // 13: proto.Token.metadata_key:type_name -> proto.Key
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_token_proto_init() }
func file_token_proto_init() {
	if File_token_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_custom_fees_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_token_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_token_proto_goTypes,
		DependencyIndexes: file_token_proto_depIdxs,
		MessageInfos:      file_token_proto_msgTypes,
	}.Build()
	File_token_proto = out.File
	file_token_proto_rawDesc = nil
	file_token_proto_goTypes = nil
	file_token_proto_depIdxs = nil
}
