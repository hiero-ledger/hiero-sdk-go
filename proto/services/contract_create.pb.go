//*
// # Smart Contract Create
//
// Create a new smart contract.
//
// ## General Comments
//  - A smart contract normally enforces rules, so "the code is law".<br/>
//    For example, an ERC-20 contract prevents a transfer from being undone
//    without a signature by the recipient of the transfer. This characteristic
//    is generally true if the contract instance was created without a value
//    for the `adminKey` field. For some uses, however, it may be desirable to
//    create something like an ERC-20 contract that has a specific group of
//    trusted individuals who can act as a "supreme court" with the ability to
//    override the normal operation, when a sufficient number of them agree to
//    do so. If `adminKey` is set to a valid Key (which MAY be complex), then a
//    transaction that can change the state of the smart contract in arbitrary
//    ways MAY be signed with enough signatures to activate that Key. Such
//    transactions might reverse a transaction, change the code to close an
//    unexpected loophole, remove an exploit, or adjust outputs in ways not
//    covered by the code itself. The admin key MAY also be used to change the
//    autoRenewPeriod, and change the adminKey field itself (for example, to
//    remove that key after a suitable testing period). The API currently does
//    not implement all relevant capabilities. But it does allow the `adminKey`
//    field to be set and queried, and MAY implement further administrative
//    capability in future releases.
//  - The current API ignores shardID, realmID, and newRealmAdminKey, and
//    creates everything in shard 0 and realm 0. Future versions of the system
//    MAY support additional shards and realms.
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
// source: contract_create.proto

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
// Create a new smart contract.
//
// If this transaction succeeds, the `ContractID` for the new smart contract
// SHALL be set in the transaction receipt.<br/>
// The contract is defined by the initial bytecode (or `initcode`). The
// `initcode` SHALL be stored either in a previously created file, or in the
// transaction body itself for very small contracts.
//
// As part of contract creation, the constructor defined for the new smart
// contract SHALL run with the parameters provided in the
// `constructorParameters` field.<br/>
// The gas to "power" that constructor MUST be provided via the `gas` field,
// and SHALL be charged to the payer for this transaction.<br/>
// If the contract _constructor_ stores information, it is charged gas for that
// storage. There is a separate fee in HBAR to maintain that storage until the
// expiration, and that fee SHALL be added to this transaction as part of the
// _transaction fee_, rather than gas.
//
// ### Block Stream Effects
// A `CreateContractOutput` message SHALL be emitted for each transaction.
type ContractCreateTransactionBody struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to InitcodeSource:
	//
	//	*ContractCreateTransactionBody_FileID
	//	*ContractCreateTransactionBody_Initcode
	InitcodeSource isContractCreateTransactionBody_InitcodeSource `protobuf_oneof:"initcodeSource"`
	// *
	// Access control for modification of the smart contract after
	// it is created.
	// <p>
	// If this field is set, that key MUST sign this transaction.<br/>
	// If this field is set, that key MUST sign each future transaction to
	// update or delete the contract.<br/>
	// An updateContract transaction that _only_ extends the topic
	// expirationTime (a "manual renewal" transaction) SHALL NOT require
	// admin key signature.
	// <p>
	// A contract without an admin key SHALL be immutable, except for
	// expiration and renewal.
	AdminKey *Key `protobuf:"bytes,3,opt,name=adminKey,proto3" json:"adminKey,omitempty"`
	// *
	// A maximum limit to the amount of gas to use for the constructor call.
	// <p>
	// The network SHALL charge the greater of the following, but SHALL NOT
	// charge more than the value of this field.
	// <ol>
	//
	//	<li>The actual gas consumed by the smart contract
	//	    constructor call.</li>
	//	<li>`80%` of this value.</li>
	//
	// </ol>
	// The `80%` factor encourages reasonable estimation, while allowing for
	// some overage to ensure successful execution.
	Gas int64 `protobuf:"varint,4,opt,name=gas,proto3" json:"gas,omitempty"`
	// *
	// The amount of HBAR to use as an initial balance for the account
	// representing the new smart contract.
	// <p>
	// This value is presented in tinybar
	// (10<sup><strong>-</strong>8</sup> HBAR).<br/>
	// The HBAR provided here will be withdrawn from the payer account that
	// signed this transaction.
	InitialBalance int64 `protobuf:"varint,5,opt,name=initialBalance,proto3" json:"initialBalance,omitempty"`
	// *
	// Proxy account staking is handled via `staked_id`.
	// <p>
	// Former field to designate a proxy account for HBAR staking.
	// This field MUST NOT be set.
	//
	// Deprecated: Marked as deprecated in contract_create.proto.
	ProxyAccountID *AccountID `protobuf:"bytes,6,opt,name=proxyAccountID,proto3" json:"proxyAccountID,omitempty"`
	// *
	// The initial lifetime, in seconds, for the smart contract, and the number
	// of seconds for which the smart contract will be automatically renewed
	// upon expiration.
	// <p>
	// This value MUST be set.<br/>
	// This value MUST be greater than the configured MIN_AUTORENEW_PERIOD.<br/>
	// This value MUST be less than the configured MAX_AUTORENEW_PERIOD.<br/>
	AutoRenewPeriod *Duration `protobuf:"bytes,8,opt,name=autoRenewPeriod,proto3" json:"autoRenewPeriod,omitempty"`
	// *
	// An array of bytes containing the EVM-encoded parameters to pass to
	// the smart contract constructor defined in the smart contract init
	// code provided.
	ConstructorParameters []byte `protobuf:"bytes,9,opt,name=constructorParameters,proto3" json:"constructorParameters,omitempty"`
	// *
	// <blockquote>Review Question<br/>
	// <blockquote>Should this be deprecated?<br/>
	// It's never been used and probably never should be used...<br/>
	// Shard should be determined by the node the transaction is submitted to.
	// </blockquote></blockquote>
	// <p>
	// The shard in which to create the new smart contract.<br/>
	// This value is currently ignored.
	ShardID *ShardID `protobuf:"bytes,10,opt,name=shardID,proto3" json:"shardID,omitempty"`
	// *
	// <blockquote>Review Question<br/>
	// <blockquote>Should this be deprecated?<br/>
	// It's never been used and probably never should be used...<br/>
	// Realm should be determined by node and network parameters.
	// </blockquote></blockquote>
	// <p>
	// The shard/realm in which to create the new smart contract.<br/>
	// This value is currently ignored.
	RealmID *RealmID `protobuf:"bytes,11,opt,name=realmID,proto3" json:"realmID,omitempty"`
	// *
	// <blockquote>Review Question<br/>
	// <blockquote>Should this be deprecated?<br/>
	// It's never been used and probably never should be used...<br/>
	// If a realm is used, it must already exist; we shouldn't be creating it
	// without a separate transaction.</blockquote></blockquote>
	// <p>
	// This was intended to provide an admin key for any new realm created
	// during the creation of the smart contract.<br/>
	// This value is currently ignored. a new realm SHALL NOT be created,
	// regardless of the value of `realmID`.
	NewRealmAdminKey *Key `protobuf:"bytes,12,opt,name=newRealmAdminKey,proto3" json:"newRealmAdminKey,omitempty"`
	// *
	// A short memo for this smart contract.
	// <p>
	// This value, if set, MUST NOT exceed `transaction.maxMemoUtf8Bytes`
	// (default 100) bytes when encoded as UTF-8.
	Memo string `protobuf:"bytes,13,opt,name=memo,proto3" json:"memo,omitempty"`
	// *
	// The maximum number of tokens that can be auto-associated with this
	// smart contract.
	// <p>
	// If this is less than or equal to `used_auto_associations` (or 0), then
	// this contract MUST manually associate with a token before transacting
	// in that token.<br/>
	// Following HIP-904 This value may also be `-1` to indicate no limit.<br/>
	// This value MUST NOT be less than `-1`.
	MaxAutomaticTokenAssociations int32 `protobuf:"varint,14,opt,name=max_automatic_token_associations,json=maxAutomaticTokenAssociations,proto3" json:"max_automatic_token_associations,omitempty"`
	// *
	// The id of an account, in the same shard and realm as this smart
	// contract, that has signed this transaction, allowing the network to use
	// its balance, when needed, to automatically extend this contract's
	// expiration time.
	// <p>
	// If this field is set, that key MUST sign this transaction.<br/>
	// If this field is set, then the network SHALL deduct the necessary fees
	// from the designated auto renew account, if that account has sufficient
	// balance. If the auto renew account does not have sufficient balance,
	// then the fees for contract renewal SHALL be deducted from the HBAR
	// balance held by the smart contract.<br/>
	// If this field is not set, then all renewal fees SHALL be deducted from
	// the HBAR balance held by this contract.
	AutoRenewAccountId *AccountID `protobuf:"bytes,15,opt,name=auto_renew_account_id,json=autoRenewAccountId,proto3" json:"auto_renew_account_id,omitempty"`
	// Types that are valid to be assigned to StakedId:
	//
	//	*ContractCreateTransactionBody_StakedAccountId
	//	*ContractCreateTransactionBody_StakedNodeId
	StakedId isContractCreateTransactionBody_StakedId `protobuf_oneof:"staked_id"`
	// *
	// A flag indicating that this smart contract declines to receive any
	// reward for staking its HBAR balance to help secure the network.
	// <p>
	// If set to true, this smart contract SHALL NOT receive any reward for
	// staking its HBAR balance to help secure the network, regardless of
	// staking configuration, but MAY stake HBAR to support the network
	// without reward.
	DeclineReward bool `protobuf:"varint,19,opt,name=decline_reward,json=declineReward,proto3" json:"decline_reward,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ContractCreateTransactionBody) Reset() {
	*x = ContractCreateTransactionBody{}
	mi := &file_contract_create_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContractCreateTransactionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractCreateTransactionBody) ProtoMessage() {}

func (x *ContractCreateTransactionBody) ProtoReflect() protoreflect.Message {
	mi := &file_contract_create_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractCreateTransactionBody.ProtoReflect.Descriptor instead.
func (*ContractCreateTransactionBody) Descriptor() ([]byte, []int) {
	return file_contract_create_proto_rawDescGZIP(), []int{0}
}

func (x *ContractCreateTransactionBody) GetInitcodeSource() isContractCreateTransactionBody_InitcodeSource {
	if x != nil {
		return x.InitcodeSource
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetFileID() *FileID {
	if x != nil {
		if x, ok := x.InitcodeSource.(*ContractCreateTransactionBody_FileID); ok {
			return x.FileID
		}
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetInitcode() []byte {
	if x != nil {
		if x, ok := x.InitcodeSource.(*ContractCreateTransactionBody_Initcode); ok {
			return x.Initcode
		}
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetAdminKey() *Key {
	if x != nil {
		return x.AdminKey
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetGas() int64 {
	if x != nil {
		return x.Gas
	}
	return 0
}

func (x *ContractCreateTransactionBody) GetInitialBalance() int64 {
	if x != nil {
		return x.InitialBalance
	}
	return 0
}

// Deprecated: Marked as deprecated in contract_create.proto.
func (x *ContractCreateTransactionBody) GetProxyAccountID() *AccountID {
	if x != nil {
		return x.ProxyAccountID
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetAutoRenewPeriod() *Duration {
	if x != nil {
		return x.AutoRenewPeriod
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetConstructorParameters() []byte {
	if x != nil {
		return x.ConstructorParameters
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetShardID() *ShardID {
	if x != nil {
		return x.ShardID
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetRealmID() *RealmID {
	if x != nil {
		return x.RealmID
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetNewRealmAdminKey() *Key {
	if x != nil {
		return x.NewRealmAdminKey
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetMemo() string {
	if x != nil {
		return x.Memo
	}
	return ""
}

func (x *ContractCreateTransactionBody) GetMaxAutomaticTokenAssociations() int32 {
	if x != nil {
		return x.MaxAutomaticTokenAssociations
	}
	return 0
}

func (x *ContractCreateTransactionBody) GetAutoRenewAccountId() *AccountID {
	if x != nil {
		return x.AutoRenewAccountId
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetStakedId() isContractCreateTransactionBody_StakedId {
	if x != nil {
		return x.StakedId
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetStakedAccountId() *AccountID {
	if x != nil {
		if x, ok := x.StakedId.(*ContractCreateTransactionBody_StakedAccountId); ok {
			return x.StakedAccountId
		}
	}
	return nil
}

func (x *ContractCreateTransactionBody) GetStakedNodeId() int64 {
	if x != nil {
		if x, ok := x.StakedId.(*ContractCreateTransactionBody_StakedNodeId); ok {
			return x.StakedNodeId
		}
	}
	return 0
}

func (x *ContractCreateTransactionBody) GetDeclineReward() bool {
	if x != nil {
		return x.DeclineReward
	}
	return false
}

type isContractCreateTransactionBody_InitcodeSource interface {
	isContractCreateTransactionBody_InitcodeSource()
}

type ContractCreateTransactionBody_FileID struct {
	// *
	// The source for the smart contract EVM bytecode.
	// <p>
	// The file containing the smart contract initcode.
	// A copy of the contents SHALL be made and held as `bytes`
	// in smart contract state.<br/>
	// The contract bytecode is limited in size only by the
	// network file size limit.
	FileID *FileID `protobuf:"bytes,1,opt,name=fileID,proto3,oneof"`
}

type ContractCreateTransactionBody_Initcode struct {
	// *
	// The source for the smart contract EVM bytecode.
	// <p>
	// The bytes of the smart contract initcode. A copy of the contents
	// SHALL be made and held as `bytes` in smart contract state.<br/>
	// This value is limited in length by the network transaction size
	// limit. This entire transaction, including all fields and signatures,
	// MUST be less than the network transaction size limit.
	Initcode []byte `protobuf:"bytes,16,opt,name=initcode,proto3,oneof"`
}

func (*ContractCreateTransactionBody_FileID) isContractCreateTransactionBody_InitcodeSource() {}

func (*ContractCreateTransactionBody_Initcode) isContractCreateTransactionBody_InitcodeSource() {}

type isContractCreateTransactionBody_StakedId interface {
	isContractCreateTransactionBody_StakedId()
}

type ContractCreateTransactionBody_StakedAccountId struct {
	// *
	// An account ID.
	// <p>
	// This smart contract SHALL stake its HBAR via this account as proxy.
	StakedAccountId *AccountID `protobuf:"bytes,17,opt,name=staked_account_id,json=stakedAccountId,proto3,oneof"`
}

type ContractCreateTransactionBody_StakedNodeId struct {
	// *
	// The ID of a network node.
	// <p>
	// This smart contract SHALL stake its HBAR to this node.
	// <p>
	// <blockquote>Note: node IDs do fluctuate as node operators change.
	// Most contracts are immutable, and a contract staking to an invalid
	// node ID SHALL NOT participate in staking. Immutable contracts MAY
	// find it more reliable to use a proxy account for staking
	// (via `staked_account_id`) to enable updating the _effective_ staking
	// node ID when necessary through updating the proxy
	// account.</blockquote>
	StakedNodeId int64 `protobuf:"varint,18,opt,name=staked_node_id,json=stakedNodeId,proto3,oneof"`
}

func (*ContractCreateTransactionBody_StakedAccountId) isContractCreateTransactionBody_StakedId() {}

func (*ContractCreateTransactionBody_StakedNodeId) isContractCreateTransactionBody_StakedId() {}

var File_contract_create_proto protoreflect.FileDescriptor

const file_contract_create_proto_rawDesc = "" +
	"\n" +
	"\x15contract_create.proto\x12\x05proto\x1a\x11basic_types.proto\x1a\x0eduration.proto\"\xd3\x06\n" +
	"\x1dContractCreateTransactionBody\x12'\n" +
	"\x06fileID\x18\x01 \x01(\v2\r.proto.FileIDH\x00R\x06fileID\x12\x1c\n" +
	"\binitcode\x18\x10 \x01(\fH\x00R\binitcode\x12&\n" +
	"\badminKey\x18\x03 \x01(\v2\n" +
	".proto.KeyR\badminKey\x12\x10\n" +
	"\x03gas\x18\x04 \x01(\x03R\x03gas\x12&\n" +
	"\x0einitialBalance\x18\x05 \x01(\x03R\x0einitialBalance\x12<\n" +
	"\x0eproxyAccountID\x18\x06 \x01(\v2\x10.proto.AccountIDB\x02\x18\x01R\x0eproxyAccountID\x129\n" +
	"\x0fautoRenewPeriod\x18\b \x01(\v2\x0f.proto.DurationR\x0fautoRenewPeriod\x124\n" +
	"\x15constructorParameters\x18\t \x01(\fR\x15constructorParameters\x12(\n" +
	"\ashardID\x18\n" +
	" \x01(\v2\x0e.proto.ShardIDR\ashardID\x12(\n" +
	"\arealmID\x18\v \x01(\v2\x0e.proto.RealmIDR\arealmID\x126\n" +
	"\x10newRealmAdminKey\x18\f \x01(\v2\n" +
	".proto.KeyR\x10newRealmAdminKey\x12\x12\n" +
	"\x04memo\x18\r \x01(\tR\x04memo\x12G\n" +
	" max_automatic_token_associations\x18\x0e \x01(\x05R\x1dmaxAutomaticTokenAssociations\x12C\n" +
	"\x15auto_renew_account_id\x18\x0f \x01(\v2\x10.proto.AccountIDR\x12autoRenewAccountId\x12>\n" +
	"\x11staked_account_id\x18\x11 \x01(\v2\x10.proto.AccountIDH\x01R\x0fstakedAccountId\x12&\n" +
	"\x0estaked_node_id\x18\x12 \x01(\x03H\x01R\fstakedNodeId\x12%\n" +
	"\x0edecline_reward\x18\x13 \x01(\bR\rdeclineRewardB\x10\n" +
	"\x0einitcodeSourceB\v\n" +
	"\tstaked_idB&\n" +
	"\"com.hederahashgraph.api.proto.javaP\x01b\x06proto3"

var (
	file_contract_create_proto_rawDescOnce sync.Once
	file_contract_create_proto_rawDescData []byte
)

func file_contract_create_proto_rawDescGZIP() []byte {
	file_contract_create_proto_rawDescOnce.Do(func() {
		file_contract_create_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_contract_create_proto_rawDesc), len(file_contract_create_proto_rawDesc)))
	})
	return file_contract_create_proto_rawDescData
}

var file_contract_create_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_contract_create_proto_goTypes = []any{
	(*ContractCreateTransactionBody)(nil), // 0: proto.ContractCreateTransactionBody
	(*FileID)(nil),                        // 1: proto.FileID
	(*Key)(nil),                           // 2: proto.Key
	(*AccountID)(nil),                     // 3: proto.AccountID
	(*Duration)(nil),                      // 4: proto.Duration
	(*ShardID)(nil),                       // 5: proto.ShardID
	(*RealmID)(nil),                       // 6: proto.RealmID
}
var file_contract_create_proto_depIdxs = []int32{
	1, // 0: proto.ContractCreateTransactionBody.fileID:type_name -> proto.FileID
	2, // 1: proto.ContractCreateTransactionBody.adminKey:type_name -> proto.Key
	3, // 2: proto.ContractCreateTransactionBody.proxyAccountID:type_name -> proto.AccountID
	4, // 3: proto.ContractCreateTransactionBody.autoRenewPeriod:type_name -> proto.Duration
	5, // 4: proto.ContractCreateTransactionBody.shardID:type_name -> proto.ShardID
	6, // 5: proto.ContractCreateTransactionBody.realmID:type_name -> proto.RealmID
	2, // 6: proto.ContractCreateTransactionBody.newRealmAdminKey:type_name -> proto.Key
	3, // 7: proto.ContractCreateTransactionBody.auto_renew_account_id:type_name -> proto.AccountID
	3, // 8: proto.ContractCreateTransactionBody.staked_account_id:type_name -> proto.AccountID
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_contract_create_proto_init() }
func file_contract_create_proto_init() {
	if File_contract_create_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_duration_proto_init()
	file_contract_create_proto_msgTypes[0].OneofWrappers = []any{
		(*ContractCreateTransactionBody_FileID)(nil),
		(*ContractCreateTransactionBody_Initcode)(nil),
		(*ContractCreateTransactionBody_StakedAccountId)(nil),
		(*ContractCreateTransactionBody_StakedNodeId)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_contract_create_proto_rawDesc), len(file_contract_create_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_contract_create_proto_goTypes,
		DependencyIndexes: file_contract_create_proto_depIdxs,
		MessageInfos:      file_contract_create_proto_msgTypes,
	}.Build()
	File_contract_create_proto = out.File
	file_contract_create_proto_goTypes = nil
	file_contract_create_proto_depIdxs = nil
}
