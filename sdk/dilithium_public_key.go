package hiero

// SPDX-License-Identifier: Apache-2.0

// PQC POC: CRYSTALS-Dilithium Mode3 (ML-DSA-65) public key implementation.
// This follows the same internal patterns as _Ed25519PublicKey and _ECDSAPublicKey
// to demonstrate that the SDK key abstraction can accommodate post-quantum keys.

import (
	"bytes"
	"encoding/hex"

	"github.com/cloudflare/circl/sign/dilithium/mode3"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// _DilithiumPublicKey wraps a CRYSTALS-Dilithium Mode3 public key.
type _DilithiumPublicKey struct {
	keyData *mode3.PublicKey
}

func _DilithiumPublicKeyFromBytes(byt []byte) (*_DilithiumPublicKey, error) {
	if len(byt) != mode3.PublicKeySize {
		return nil, _NewErrBadKeyf("invalid dilithium public key length: %d bytes, expected %d", len(byt), mode3.PublicKeySize)
	}

	pk := new(mode3.PublicKey)
	if err := pk.UnmarshalBinary(byt); err != nil {
		return nil, _NewErrBadKeyf("failed to unmarshal dilithium public key: %v", err)
	}

	return &_DilithiumPublicKey{
		keyData: pk,
	}, nil
}

func _DilithiumPublicKeyFromString(s string) (*_DilithiumPublicKey, error) {
	byt, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return _DilithiumPublicKeyFromBytes(byt)
}

func (pk _DilithiumPublicKey) String() string {
	return pk._StringRaw()
}

func (pk _DilithiumPublicKey) _StringRaw() string {
	b, _ := pk.keyData.MarshalBinary()
	return hex.EncodeToString(b)
}

func (pk _DilithiumPublicKey) _StringDer() string {
	// PQC POC: DER encoding for Dilithium is not yet standardized.
	return pk._StringRaw()
}

func (pk _DilithiumPublicKey) _Bytes() []byte {
	b, _ := pk.keyData.MarshalBinary()
	return b
}

func (pk _DilithiumPublicKey) _BytesRaw() []byte {
	return pk._Bytes()
}

func (pk _DilithiumPublicKey) _BytesDer() []byte {
	// PQC POC: DER encoding not yet standardized for Dilithium.
	return pk._BytesRaw()
}

// _ToProtoKey converts to a protobuf Key.
// PQC POC: The Hedera protobuf schema does not yet have a Dilithium key variant.
// In production, a new oneof field (e.g., `bytes dilithium_mode3 = 9;`) would be added
// to the Key message in basic_types.proto. For this POC, we store the raw public key
// bytes using the Ed25519 field as a placeholder — this is NOT valid for network use.
func (pk _DilithiumPublicKey) _ToProtoKey() *services.Key {
	b, _ := pk.keyData.MarshalBinary()
	return &services.Key{Key: &services.Key_Ed25519{Ed25519: b}}
}

// _ToSignaturePairProtobuf creates a signature pair for transaction signing.
// PQC POC: Uses Ed25519 signature field as placeholder. In production, a new oneof
// variant (e.g., `bytes dilithium_mode3 = 7;`) would be added to SignaturePair.
func (pk _DilithiumPublicKey) _ToSignaturePairProtobuf(signature []byte) *services.SignaturePair {
	b, _ := pk.keyData.MarshalBinary()
	return &services.SignaturePair{
		PubKeyPrefix: b,
		Signature: &services.SignaturePair_Ed25519{
			Ed25519: signature,
		},
	}
}

func (pk _DilithiumPublicKey) _VerifySignedMessage(message []byte, signature []byte) bool {
	return mode3.Verify(pk.keyData, message, signature)
}

func (pk _DilithiumPublicKey) _VerifyTransaction(tx *Transaction[TransactionInterface]) bool {
	if tx.signedTransactions._Length() == 0 {
		return false
	}

	for _, signedKey := range tx.publicKeys {
		if bytes.Equal(signedKey.BytesRaw(), pk._BytesRaw()) {
			return true
		}
	}

	for _, value := range tx.signedTransactions.slice {
		signedTx := value.(*services.SignedTransaction)
		found := false
		for _, sigPair := range signedTx.SigMap.GetSigPair() {
			if bytes.Equal(sigPair.GetPubKeyPrefix(), pk._Bytes()) {
				found = true
				if !pk._VerifySignedMessage(signedTx.BodyBytes, sigPair.GetEd25519()) {
					return false
				}
			}
		}

		if !found {
			return false
		}
	}

	return true
}
