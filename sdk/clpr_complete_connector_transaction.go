package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprCompleteConnectorTransaction reveals the ownership key of a registered CLPR Connector and
// supplies its salt, contract, admin key and locked stake, the reveal phase of the Connector
// commit-reveal (HIP-1535).
type ClprCompleteConnectorTransaction struct {
	*Transaction[*ClprCompleteConnectorTransaction]
	fieldErrors       _ClprFieldErrors
	channelID         []byte
	connectorID       []byte
	salt              []byte
	connectorContract *ContractID
	adminKey          Key
	lockedStake       uint64
	ownershipProof    *_ClprOwnershipProof
}

// NewClprCompleteConnectorTransaction creates a ClprCompleteConnectorTransaction.
func NewClprCompleteConnectorTransaction() *ClprCompleteConnectorTransaction {
	tx := &ClprCompleteConnectorTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func _ClprCompleteConnectorTransactionFromProtobuf(tx Transaction[*ClprCompleteConnectorTransaction], pb *services.TransactionBody) ClprCompleteConnectorTransaction {
	body := pb.GetClprCompleteConnector()

	clprCompleteConnectorTransaction := ClprCompleteConnectorTransaction{
		channelID:      _ClprBytesFromProtobuf(body.GetChannelId()),
		connectorID:    _ClprBytesFromProtobuf(body.GetConnectorId()),
		salt:           _ClprBytesFromProtobuf(body.GetSalt()),
		lockedStake:    body.GetLockedStake(),
		ownershipProof: _ClprOwnershipProofFromProtobuf(body.GetPublicKey(), body.GetSignature(), body.GetSignatureScheme()),
	}
	if body.GetConnectorContract() != nil {
		clprCompleteConnectorTransaction.connectorContract = _ContractIDFromProtobuf(body.GetConnectorContract())
	}
	if body.GetAdminKey() != nil {
		adminKey, err := _KeyFromProtobuf(body.GetAdminKey())
		if err != nil {
			clprCompleteConnectorTransaction.fieldErrors._Record("adminKey", fmt.Errorf("%w: %w", errClprAdminKeyInvalid, err))
		} else {
			clprCompleteConnectorTransaction.adminKey = adminKey
		}
	}

	tx.childTransaction = &clprCompleteConnectorTransaction
	clprCompleteConnectorTransaction.Transaction = &tx
	return clprCompleteConnectorTransaction
}

// SetChannelID sets the 32-byte Channel ID. Changing it clears the ownership proof.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprCompleteConnectorTransaction) SetChannelID(channelID []byte) *ClprCompleteConnectorTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("channelID", _ClprCheckLength("channelID", channelID)) {
		return tx
	}
	if _ClprBytesChanged(tx.channelID, channelID) {
		tx.ownershipProof = nil
	}
	tx.channelID = _ClprCopyBytes(channelID)
	return tx
}

// GetChannelID returns a copy of the Channel ID, or nil if it is not set.
func (tx *ClprCompleteConnectorTransaction) GetChannelID() []byte {
	return _ClprCopyBytes(tx.channelID)
}

// SetConnectorID sets the 32-byte Connector ID from ComputeClprConnectorID. Changing it clears the ownership proof.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprCompleteConnectorTransaction) SetConnectorID(connectorID []byte) *ClprCompleteConnectorTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("connectorID", _ClprCheckLength("connectorID", connectorID)) {
		return tx
	}
	if _ClprBytesChanged(tx.connectorID, connectorID) {
		tx.ownershipProof = nil
	}
	tx.connectorID = _ClprCopyBytes(connectorID)
	return tx
}

// GetConnectorID returns a copy of the Connector ID, or nil if it is not set.
func (tx *ClprCompleteConnectorTransaction) GetConnectorID() []byte {
	return _ClprCopyBytes(tx.connectorID)
}

// SetSalt sets the 32-byte salt used to derive the Connector ID. Changing it clears the ownership proof.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprCompleteConnectorTransaction) SetSalt(salt []byte) *ClprCompleteConnectorTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("salt", _ClprCheckLength("salt", salt)) {
		return tx
	}
	if _ClprBytesChanged(tx.salt, salt) {
		tx.ownershipProof = nil
	}
	tx.salt = _ClprCopyBytes(salt)
	return tx
}

// GetSalt returns a copy of the salt, or nil if it is not set.
func (tx *ClprCompleteConnectorTransaction) GetSalt() []byte {
	return _ClprCopyBytes(tx.salt)
}

// SetConnectorContract sets the contract that acts as the Connector.
func (tx *ClprCompleteConnectorTransaction) SetConnectorContract(contractID ContractID) *ClprCompleteConnectorTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() {
		return tx
	}
	tx.connectorContract = _ClprCopyContractID(&contractID)
	return tx
}

// GetConnectorContract returns the Connector contract, or an empty ContractID if it is not set.
func (tx *ClprCompleteConnectorTransaction) GetConnectorContract() ContractID {
	if tx.connectorContract == nil {
		return ContractID{}
	}

	return *_ClprCopyContractID(tx.connectorContract)
}

// SetAdminKey sets the key that administers the Connector.
// An empty key is not stored; the previous value is kept and freezing fails.
func (tx *ClprCompleteConnectorTransaction) SetAdminKey(adminKey Key) *ClprCompleteConnectorTransaction {
	var err error
	if _ClprKeyIsEmpty(adminKey) {
		err = errClprAdminKeyInvalid
	}
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("adminKey", err) {
		return tx
	}
	tx.adminKey = _ClprCopyKey(adminKey)
	return tx
}

// GetAdminKey returns a copy of the admin key, or nil if it is not set.
func (tx *ClprCompleteConnectorTransaction) GetAdminKey() Key {
	return _ClprCopyKey(tx.adminKey)
}

// SetLockedStake sets the stake locked by the Connector.
// A negative value is not stored; the previous value is kept and freezing fails.
func (tx *ClprCompleteConnectorTransaction) SetLockedStake(lockedStake Hbar) *ClprCompleteConnectorTransaction {
	tinybars := lockedStake.AsTinybar()
	var err error
	if tinybars < 0 {
		err = errClprLockedStakeNegative
	}
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("lockedStake", err) {
		return tx
	}
	tx.lockedStake = uint64(tinybars)
	return tx
}

// GetLockedStake returns the locked stake. A stake above 2^63-1 tinybars decoded from bytes saturates at
// that value; use GetLockedStakeTinybars for the exact wire value.
func (tx *ClprCompleteConnectorTransaction) GetLockedStake() Hbar {
	if tx.lockedStake > math.MaxInt64 {
		return HbarFromTinybar(math.MaxInt64)
	}

	return HbarFromTinybar(int64(tx.lockedStake))
}

// GetLockedStakeTinybars returns the exact locked stake in tinybars.
func (tx *ClprCompleteConnectorTransaction) GetLockedStakeTinybars() uint64 {
	return tx.lockedStake
}

// SetOwnershipProof signs the Connector ownership reveal with the key committed by
// ClprRegisterConnectorTransaction. The Channel ID, Connector ID and salt must be set first, and the
// Connector ID must derive from them and the key. The private key is not retained.
func (tx *ClprCompleteConnectorTransaction) SetOwnershipProof(privateKey PrivateKey) (*ClprCompleteConnectorTransaction, error) {
	if tx._RequireNotFrozen(); tx.IsFrozen() {
		return tx, errTransactionIsFrozen
	}
	if err := _ClprRequireIDBefore(tx.getName(), "channelID", tx.channelID, "calling SetOwnershipProof"); err != nil {
		return tx, err
	}
	if err := _ClprRequireIDBefore(tx.getName(), "connectorID", tx.connectorID, "calling SetOwnershipProof"); err != nil {
		return tx, err
	}
	if err := _ClprRequireIDBefore(tx.getName(), "salt", tx.salt, "calling SetOwnershipProof"); err != nil {
		return tx, err
	}
	if privateKey.ecdsaPrivateKey == nil && privateKey.ed25519PrivateKey == nil {
		return tx, errClprPrivateKeyRequired
	}

	expectedID, err := ComputeClprConnectorID(tx.channelID, privateKey.PublicKey(), tx.salt)
	if err != nil {
		return tx, err
	}
	if string(expectedID) != string(tx.connectorID) {
		return tx, fmt.Errorf("%s: %w", tx.getName(), errClprConnectorIDMismatch)
	}

	proof, err := _ClprSignOwnership(privateKey, _ClprConnectorOwnershipPreimage(tx.connectorID))
	if err != nil {
		return tx, err
	}

	tx.ownershipProof = proof
	return tx, nil
}

// GetOwnershipPublicKey returns a copy of the canonical ownership public key, or nil if no proof is set.
func (tx *ClprCompleteConnectorTransaction) GetOwnershipPublicKey() []byte {
	return tx.ownershipProof._GetPublicKey()
}

// GetOwnershipSignature returns a copy of the ownership signature, or nil if no proof is set.
func (tx *ClprCompleteConnectorTransaction) GetOwnershipSignature() []byte {
	return tx.ownershipProof._GetSignature()
}

// GetOwnershipSignatureScheme returns the ownership signature scheme and whether a proof is set.
func (tx *ClprCompleteConnectorTransaction) GetOwnershipSignatureScheme() (ClprSignatureScheme, bool) {
	return tx.ownershipProof._GetScheme()
}

// _ClprConnectorOwnershipPreimage is the message the Connector owner signs: connectorID || 0x16e.
func _ClprConnectorOwnershipPreimage(connectorID []byte) []byte {
	return append(_ClprCopyBytes(connectorID), _ClprServiceAddress[:]...)
}

// ----------- Overridden functions ----------------

func (tx ClprCompleteConnectorTransaction) getName() string {
	return "ClprCompleteConnectorTransaction"
}

func (tx ClprCompleteConnectorTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.connectorContract != nil {
		if err := tx.connectorContract.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx ClprCompleteConnectorTransaction) validateTransactionFields() error {
	if err := _ClprRequireID(tx.getName(), "channelID", tx.channelID, tx.fieldErrors); err != nil {
		return err
	}
	if err := _ClprRequireID(tx.getName(), "connectorID", tx.connectorID, tx.fieldErrors); err != nil {
		return err
	}
	if err := _ClprRequireID(tx.getName(), "salt", tx.salt, tx.fieldErrors); err != nil {
		return err
	}
	if err := _ClprRequireSet(tx.getName(), "connectorContract", tx.connectorContract != nil, nil); err != nil {
		return err
	}
	if err := _ClprRequireSet(tx.getName(), "adminKey", tx.adminKey != nil, tx.fieldErrors); err != nil {
		return err
	}
	if err := tx.ownershipProof._Require(tx.getName()); err != nil {
		return err
	}
	if err := tx.fieldErrors["lockedStake"]; err != nil {
		return fmt.Errorf("%s: %w", tx.getName(), err)
	}

	return nil
}

func (tx ClprCompleteConnectorTransaction) validateFrozenExecute() error {
	return tx.validateTransactionFields()
}

func (tx ClprCompleteConnectorTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_ClprCompleteConnector{
		ClprCompleteConnector: tx.buildProtoBody(),
	}

	return body
}

func (tx ClprCompleteConnectorTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, _ClprCannotSchedule(tx.getName())
}

func (tx ClprCompleteConnectorTransaction) buildProtoBody() *services.ClprCompleteConnectorTransactionBody {
	body := &services.ClprCompleteConnectorTransactionBody{
		ChannelId:   _ClprCopyBytes(tx.channelID),
		ConnectorId: _ClprCopyBytes(tx.connectorID),
		Salt:        _ClprCopyBytes(tx.salt),
		LockedStake: tx.lockedStake,
	}
	if tx.connectorContract != nil {
		body.ConnectorContract = tx.connectorContract._ToProtobuf()
	}
	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}
	if tx.ownershipProof != nil {
		body.PublicKey = tx.ownershipProof._GetPublicKey()
		body.Signature = tx.ownershipProof._GetSignature()
		body.SignatureScheme = services.ClprSignatureScheme(tx.ownershipProof.scheme)
	}

	return body
}

func (tx ClprCompleteConnectorTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetClpr().CompleteConnector,
	}
}

func (tx ClprCompleteConnectorTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
