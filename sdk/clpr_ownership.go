package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"errors"
	"fmt"
)

const _ClprIDLength = 32

var errClprOwnershipKeyUnsupported = errors.New("CLPR ownership key must be an ECDSA secp256k1 or Ed25519 PublicKey")
var errClprPrivateKeyRequired = errors.New("CLPR ownership proof requires an ECDSA secp256k1 or Ed25519 PrivateKey")

// ComputeClprChannelOwnershipCommitment returns keccak256(channelID || publicKey), the commitment
// a ClprRegisterChannelTransaction registers before the Channel is completed.
func ComputeClprChannelOwnershipCommitment(channelID []byte, publicKey PublicKey) ([]byte, error) {
	if err := _ClprCheckLength("channelID", channelID); err != nil {
		return nil, err
	}

	_, key, err := _ClprEncodePublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return _ClprKeccak(channelID, key), nil
}

// ComputeClprConnectorID returns keccak256(channelID || publicKey || salt), the ID of a Connector.
func ComputeClprConnectorID(channelID []byte, publicKey PublicKey, salt []byte) ([]byte, error) {
	if err := _ClprCheckLength("channelID", channelID); err != nil {
		return nil, err
	}
	if err := _ClprCheckLength("salt", salt); err != nil {
		return nil, err
	}

	_, key, err := _ClprEncodePublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return _ClprKeccak(channelID, key, salt), nil
}

// ComputeClprConnectorOwnershipCommitment returns keccak256(connectorID || publicKey), the commitment
// a ClprRegisterConnectorTransaction registers before the Connector is completed.
func ComputeClprConnectorOwnershipCommitment(connectorID []byte, publicKey PublicKey) ([]byte, error) {
	if err := _ClprCheckLength("connectorID", connectorID); err != nil {
		return nil, err
	}

	_, key, err := _ClprEncodePublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return _ClprKeccak(connectorID, key), nil
}

func _ClprCheckLength(field string, value []byte) error {
	if len(value) != _ClprIDLength {
		return fmt.Errorf("%s must be exactly %d bytes, got %d", field, _ClprIDLength, len(value))
	}

	return nil
}

// _ClprEncodePublicKey returns the canonical ownership key: ECDSA x||y (64 bytes) or Ed25519 (32 bytes).
func _ClprEncodePublicKey(publicKey PublicKey) (ClprSignatureScheme, []byte, error) {
	switch {
	case publicKey.ecdsaPublicKey != nil:
		return ClprSignatureSchemeEcdsaSecp256k1, publicKey.ecdsaPublicKey._ToFullKey()[1:], nil
	case publicKey.ed25519PublicKey != nil:
		return ClprSignatureSchemeEd25519, publicKey.ed25519PublicKey._Bytes(), nil
	}

	return 0, nil, errClprOwnershipKeyUnsupported
}

func _ClprKeccak(parts ...[]byte) []byte {
	var data []byte
	for _, part := range parts {
		data = append(data, part...)
	}

	hash := Keccak256Hash(data)
	return hash.Bytes()
}
