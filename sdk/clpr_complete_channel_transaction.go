package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprCompleteChannelTransaction reveals the ownership key of a registered CLPR Channel and supplies
// its verifier contract and peer proofs, the reveal phase of the Channel commit-reveal (HIP-1535).
// Any payer may submit it.
type ClprCompleteChannelTransaction struct {
	*Transaction[*ClprCompleteChannelTransaction]
	fieldErrors                _ClprFieldErrors
	channelID                  []byte
	verifierContract           *ContractID
	configProofBytes           []byte
	endpointManifestProofBytes []byte
	ownershipProof             *_ClprOwnershipProof
}

// NewClprCompleteChannelTransaction creates a ClprCompleteChannelTransaction.
func NewClprCompleteChannelTransaction() *ClprCompleteChannelTransaction {
	tx := &ClprCompleteChannelTransaction{}
	tx.Transaction = _NewTransaction(tx)

	return tx
}

func _ClprCompleteChannelTransactionFromProtobuf(tx Transaction[*ClprCompleteChannelTransaction], pb *services.TransactionBody) ClprCompleteChannelTransaction {
	body := pb.GetClprCompleteChannel()

	clprCompleteChannelTransaction := ClprCompleteChannelTransaction{
		channelID:                  _ClprBytesFromProtobuf(body.GetChannelId()),
		configProofBytes:           _ClprBytesFromProtobuf(body.GetConfigProofBytes()),
		endpointManifestProofBytes: _ClprCopyBytes(body.GetEndpointManifestProofBytes()),
		ownershipProof:             _ClprOwnershipProofFromProtobuf(body.GetPublicKey(), body.GetSignature(), body.GetSignatureScheme()),
	}
	if body.GetVerifierContract() != nil {
		clprCompleteChannelTransaction.verifierContract = _ContractIDFromProtobuf(body.GetVerifierContract())
	}

	tx.childTransaction = &clprCompleteChannelTransaction
	clprCompleteChannelTransaction.Transaction = &tx
	return clprCompleteChannelTransaction
}

// SetChannelID sets the 32-byte Channel ID. Changing it clears the ownership proof.
// A value of another length is not stored; the previous value is kept and freezing fails.
func (tx *ClprCompleteChannelTransaction) SetChannelID(channelID []byte) *ClprCompleteChannelTransaction {
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
func (tx *ClprCompleteChannelTransaction) GetChannelID() []byte {
	return _ClprCopyBytes(tx.channelID)
}

// SetVerifierContract sets the contract that verifies the peer ledger's proofs.
func (tx *ClprCompleteChannelTransaction) SetVerifierContract(contractID ContractID) *ClprCompleteChannelTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() {
		return tx
	}
	tx.verifierContract = _ClprCopyContractID(&contractID)
	return tx
}

// GetVerifierContract returns the verifier contract, or an empty ContractID if it is not set.
func (tx *ClprCompleteChannelTransaction) GetVerifierContract() ContractID {
	if tx.verifierContract == nil {
		return ContractID{}
	}

	return *_ClprCopyContractID(tx.verifierContract)
}

// SetConfigProofBytes sets the opaque proof of the peer ledger configuration.
// An empty value is not stored; the previous value is kept and freezing fails.
func (tx *ClprCompleteChannelTransaction) SetConfigProofBytes(configProofBytes []byte) *ClprCompleteChannelTransaction {
	var err error
	if len(configProofBytes) == 0 {
		err = errClprConfigProofEmpty
	}
	if tx._RequireNotFrozen(); tx.IsFrozen() || !tx.fieldErrors._Record("configProofBytes", err) {
		return tx
	}
	tx.configProofBytes = _ClprCopyBytes(configProofBytes)
	return tx
}

// GetConfigProofBytes returns a copy of the configuration proof, or nil if it is not set.
func (tx *ClprCompleteChannelTransaction) GetConfigProofBytes() []byte {
	return _ClprCopyBytes(tx.configProofBytes)
}

// SetEndpointManifestProofBytes sets the opaque proof of the peer endpoint manifest. It may be empty.
func (tx *ClprCompleteChannelTransaction) SetEndpointManifestProofBytes(endpointManifestProofBytes []byte) *ClprCompleteChannelTransaction {
	if tx._RequireNotFrozen(); tx.IsFrozen() {
		return tx
	}
	tx.endpointManifestProofBytes = _ClprCopyBytes(endpointManifestProofBytes)
	return tx
}

// GetEndpointManifestProofBytes returns a copy of the endpoint manifest proof.
func (tx *ClprCompleteChannelTransaction) GetEndpointManifestProofBytes() []byte {
	return _ClprCopyBytes(tx.endpointManifestProofBytes)
}

// SetOwnershipProof signs the Channel ownership reveal with the key committed by
// ClprRegisterChannelTransaction. The Channel ID must be set first. The private key is not retained.
func (tx *ClprCompleteChannelTransaction) SetOwnershipProof(privateKey PrivateKey) (*ClprCompleteChannelTransaction, error) {
	if tx._RequireNotFrozen(); tx.IsFrozen() {
		return tx, errTransactionIsFrozen
	}
	if err := _ClprRequireIDBefore(tx.getName(), "channelID", tx.channelID, "calling SetOwnershipProof"); err != nil {
		return tx, err
	}

	proof, err := _ClprSignOwnership(privateKey, _ClprChannelOwnershipPreimage(tx.channelID))
	if err != nil {
		return tx, err
	}

	tx.ownershipProof = proof
	return tx, nil
}

// GetOwnershipPublicKey returns a copy of the canonical ownership public key, or nil if no proof is set.
func (tx *ClprCompleteChannelTransaction) GetOwnershipPublicKey() []byte {
	return tx.ownershipProof._GetPublicKey()
}

// GetOwnershipSignature returns a copy of the ownership signature, or nil if no proof is set.
func (tx *ClprCompleteChannelTransaction) GetOwnershipSignature() []byte {
	return tx.ownershipProof._GetSignature()
}

// GetOwnershipSignatureScheme returns the ownership signature scheme and whether a proof is set.
func (tx *ClprCompleteChannelTransaction) GetOwnershipSignatureScheme() (ClprSignatureScheme, bool) {
	return tx.ownershipProof._GetScheme()
}

// _ClprChannelOwnershipPreimage is the message the Channel owner signs: the Channel ID.
func _ClprChannelOwnershipPreimage(channelID []byte) []byte {
	return _ClprCopyBytes(channelID)
}

// ----------- Overridden functions ----------------

func (tx ClprCompleteChannelTransaction) getName() string {
	return "ClprCompleteChannelTransaction"
}

func (tx ClprCompleteChannelTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums || tx.verifierContract == nil {
		return nil
	}

	return tx.verifierContract.ValidateChecksum(client)
}

func (tx ClprCompleteChannelTransaction) validateTransactionFields() error {
	if err := _ClprRequireID(tx.getName(), "channelID", tx.channelID, tx.fieldErrors); err != nil {
		return err
	}
	if err := tx.ownershipProof._Require(tx.getName()); err != nil {
		return err
	}
	if err := _ClprRequireSet(tx.getName(), "verifierContract", tx.verifierContract != nil, nil); err != nil {
		return err
	}
	if err := tx.fieldErrors["configProofBytes"]; err != nil {
		return fmt.Errorf("%s: %w", tx.getName(), err)
	}
	if len(tx.configProofBytes) == 0 {
		return fmt.Errorf("%s: %w", tx.getName(), errClprConfigProofEmpty)
	}

	return nil
}

func (tx ClprCompleteChannelTransaction) validateFrozenExecute() error {
	return tx.validateTransactionFields()
}

func (tx ClprCompleteChannelTransaction) build() *services.TransactionBody {
	body := tx.buildTransactionBody()
	body.Data = &services.TransactionBody_ClprCompleteChannel{
		ClprCompleteChannel: tx.buildProtoBody(),
	}

	return body
}

func (tx ClprCompleteChannelTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, _ClprCannotSchedule(tx.getName())
}

func (tx ClprCompleteChannelTransaction) buildProtoBody() *services.ClprCompleteChannelTransactionBody {
	body := &services.ClprCompleteChannelTransactionBody{
		ChannelId:                  _ClprCopyBytes(tx.channelID),
		ConfigProofBytes:           _ClprCopyBytes(tx.configProofBytes),
		EndpointManifestProofBytes: _ClprCopyBytes(tx.endpointManifestProofBytes),
	}
	if tx.verifierContract != nil {
		body.VerifierContract = tx.verifierContract._ToProtobuf()
	}
	if tx.ownershipProof != nil {
		body.PublicKey = tx.ownershipProof._GetPublicKey()
		body.Signature = tx.ownershipProof._GetSignature()
		body.SignatureScheme = services.ClprSignatureScheme(tx.ownershipProof.scheme)
	}

	return body
}

func (tx ClprCompleteChannelTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetClpr().CompleteChannel,
	}
}

func (tx ClprCompleteChannelTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
