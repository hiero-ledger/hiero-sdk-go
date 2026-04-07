package hiero

// SPDX-License-Identifier: Apache-2.0

// PQC POC: CRYSTALS-Dilithium Mode3 (ML-DSA-65) private key implementation.
// This follows the same internal patterns as _Ed25519PrivateKey and _ECDSAPrivateKey
// to demonstrate that the SDK key abstraction can accommodate post-quantum keys.

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/cloudflare/circl/sign/dilithium/mode3"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
)

// _DilithiumPrivateKey wraps a CRYSTALS-Dilithium Mode3 private key.
type _DilithiumPrivateKey struct {
	keyData *mode3.PrivateKey
}

func _GenerateDilithiumPrivateKey() (*_DilithiumPrivateKey, error) {
	_, sk, err := mode3.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &_DilithiumPrivateKey{
		keyData: sk,
	}, nil
}

func _DilithiumPrivateKeyFromBytes(bytes []byte) (*_DilithiumPrivateKey, error) {
	if len(bytes) != mode3.PrivateKeySize {
		return nil, _NewErrBadKeyf("invalid dilithium private key length: %d bytes, expected %d", len(bytes), mode3.PrivateKeySize)
	}

	sk := new(mode3.PrivateKey)
	if err := sk.UnmarshalBinary(bytes); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal dilithium private key")
	}

	return &_DilithiumPrivateKey{
		keyData: sk,
	}, nil
}

func _DilithiumPrivateKeyFromString(s string) (*_DilithiumPrivateKey, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return _DilithiumPrivateKeyFromBytes(bytes)
}

func (sk _DilithiumPrivateKey) _PublicKey() *_DilithiumPublicKey {
	pk := sk.keyData.Public().(*mode3.PublicKey)
	return &_DilithiumPublicKey{
		keyData: pk,
	}
}

func (sk _DilithiumPrivateKey) String() string {
	return sk._StringRaw()
}

func (sk _DilithiumPrivateKey) _StringRaw() string {
	b, _ := sk.keyData.MarshalBinary()
	return hex.EncodeToString(b)
}

func (sk _DilithiumPrivateKey) _StringDer() string {
	// PQC POC: DER encoding for Dilithium is not yet standardized.
	// Using raw hex for now; FIPS 204 OIDs would be used in production.
	return sk._StringRaw()
}

func (sk _DilithiumPrivateKey) _BytesRaw() []byte {
	b, _ := sk.keyData.MarshalBinary()
	return b
}

func (sk _DilithiumPrivateKey) _BytesDer() []byte {
	// PQC POC: DER encoding not yet standardized for Dilithium.
	return sk._BytesRaw()
}

// _Sign signs the provided message with the Dilithium private key.
func (sk _DilithiumPrivateKey) _Sign(message []byte) []byte {
	sig := make([]byte, mode3.SignatureSize)
	mode3.SignTo(sk.keyData, message, sig)
	return sig
}

func (sk _DilithiumPrivateKey) _ToProtoKey() *services.Key {
	return sk._PublicKey()._ToProtoKey()
}

func (sk _DilithiumPrivateKey) _SignTransaction(tx *Transaction[TransactionInterface]) ([]byte, error) {
	tx._RequireOneNodeAccountID()

	if tx.signedTransactions._Length() == 0 {
		return make([]byte, 0), errTransactionRequiresSingleNodeAccountID
	}

	signature := sk._Sign(tx.signedTransactions._Get(0).(*services.SignedTransaction).GetBodyBytes())

	publicKey := sk._PublicKey()
	if publicKey == nil {
		return []byte{}, errors.New("public key is nil")
	}

	tx.AddSignature(PublicKey{
		dilithiumPublicKey: publicKey,
	}, signature)
	return signature, nil
}
