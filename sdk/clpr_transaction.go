package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"errors"
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// HIP-1535: shared helpers for the CLPR transactions.

var errClprFieldNotSet = errors.New("must be explicitly set")
var errClprOwnershipProofIncomplete = errors.New("a complete ownership proof must be set before freezing or executing")
var errClprConfigProofEmpty = errors.New("configProofBytes must not be empty")
var errClprBodyUnsupported = errors.New("the SDK does not build this CLPR transaction")

func _ClprCannotSchedule(name string) error {
	return fmt.Errorf("cannot schedule `%s`", name)
}

// _ClprFieldErrors keeps the error of a rejected setter value until the field is set validly. The
// rejected value is not stored, so the previous value and any ownership proof survive.
type _ClprFieldErrors map[string]error

// _Record stores or clears the field's error and reports whether the value may be stored.
func (errs *_ClprFieldErrors) _Record(field string, err error) bool {
	if err == nil {
		delete(*errs, field)
		return true
	}
	if *errs == nil {
		*errs = _ClprFieldErrors{}
	}
	(*errs)[field] = err
	return false
}

// _ClprRequireID checks that a 32-byte field was set validly before freezing or executing.
func _ClprRequireID(name string, field string, value []byte, errs _ClprFieldErrors) error {
	if err := errs[field]; err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}

	return _ClprRequireIDBefore(name, field, value, "freezing or executing")
}

func _ClprRequireIDBefore(name string, field string, value []byte, action string) error {
	if value == nil {
		return fmt.Errorf("%s: '%s' %w before %s", name, field, errClprFieldNotSet, action)
	}
	if err := _ClprCheckLength(field, value); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}

	return nil
}

func _ClprRequireSet(name string, field string, isSet bool, errs _ClprFieldErrors) error {
	if err := errs[field]; err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if !isSet {
		return fmt.Errorf("%s: '%s' %w before freezing or executing", name, field, errClprFieldNotSet)
	}

	return nil
}

// _ClprCopyBytes copies caller-owned bytes; nil stays nil so "unset" survives.
func _ClprCopyBytes(value []byte) []byte {
	if value == nil {
		return nil
	}

	return append([]byte{}, value...)
}

// _ClprBytesFromProtobuf maps an empty proto3 bytes field to unset.
func _ClprBytesFromProtobuf(value []byte) []byte {
	if len(value) == 0 {
		return nil
	}

	return _ClprCopyBytes(value)
}

func _ClprCopyContractID(id *ContractID) *ContractID {
	if id == nil {
		return nil
	}

	contractID := *id
	contractID.EvmAddress = _ClprCopyBytes(id.EvmAddress)
	return &contractID
}

// _ClprOwnershipProof is a signed ownership reveal. A decoded proof may be incomplete; the private
// key is never retained.
type _ClprOwnershipProof struct {
	scheme    ClprSignatureScheme
	publicKey []byte
	signature []byte
}

func _ClprSignOwnership(privateKey PrivateKey, preimage []byte) (*_ClprOwnershipProof, error) {
	scheme, publicKey, err := _ClprEncodePublicKey(privateKey.PublicKey())
	if err != nil {
		return nil, errClprPrivateKeyRequired
	}

	// ECDSA signing hashes with keccak256 internally; Ed25519 signs the keccak256 digest.
	message := preimage
	if scheme == ClprSignatureSchemeEd25519 {
		message = _ClprKeccak(preimage)
	}

	return &_ClprOwnershipProof{
		scheme:    scheme,
		publicKey: _ClprCopyBytes(publicKey),
		signature: privateKey.Sign(message),
	}, nil
}

// _ClprOwnershipProofFromProtobuf restores a proof from wire fields. An empty key and signature
// mean no proof, even though the scheme then decodes as ECDSA (0).
func _ClprOwnershipProofFromProtobuf(publicKey []byte, signature []byte, scheme services.ClprSignatureScheme) *_ClprOwnershipProof {
	if len(publicKey) == 0 && len(signature) == 0 {
		return nil
	}

	return &_ClprOwnershipProof{
		scheme:    ClprSignatureScheme(scheme),
		publicKey: _ClprCopyBytes(publicKey),
		signature: _ClprCopyBytes(signature),
	}
}

func (proof *_ClprOwnershipProof) _IsComplete() bool {
	return proof != nil && len(proof.publicKey) > 0 && len(proof.signature) > 0
}

func (proof *_ClprOwnershipProof) _GetPublicKey() []byte {
	if proof == nil {
		return nil
	}

	return _ClprCopyBytes(proof.publicKey)
}

func (proof *_ClprOwnershipProof) _GetSignature() []byte {
	if proof == nil {
		return nil
	}

	return _ClprCopyBytes(proof.signature)
}

func (proof *_ClprOwnershipProof) _GetScheme() (ClprSignatureScheme, bool) {
	if proof == nil {
		return 0, false
	}

	return proof.scheme, true
}

func (proof *_ClprOwnershipProof) _Require(name string) error {
	if !proof._IsComplete() {
		return fmt.Errorf("%s: %w", name, errClprOwnershipProofIncomplete)
	}

	return nil
}

// _ClprBytesChanged reports whether setting next over previous invalidates a proof bound to it.
// An unset previous value always counts as a change.
func _ClprBytesChanged(previous []byte, next []byte) bool {
	return previous == nil || string(previous) != string(next)
}
