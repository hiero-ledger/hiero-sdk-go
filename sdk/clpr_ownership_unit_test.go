//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HIP-1535 ownership vectors from sdk-collaboration-hub proposals/hips/hip-1535.md (d04df37c).
// Public test keys only. The digests and signatures use the ID-only completion formulas.
type clprOwnershipVector struct {
	name                string
	scheme              ClprSignatureScheme
	privateKey          string
	publicKey           string
	channelCommitment   string
	channelDigest       string
	channelSignature    string
	connectorID         string
	connectorCommitment string
	connectorDigest     string
	connectorSignature  string
}

const (
	clprVectorChannelID = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	clprVectorSalt      = "202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"
)

var clprOwnershipVectors = []clprOwnershipVector{
	{
		name:                "ECDSA_SECP256K1",
		scheme:              ClprSignatureSchemeEcdsaSecp256k1,
		privateKey:          "0000000000000000000000000000000000000000000000000000000000000001",
		publicKey:           "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8",
		channelCommitment:   "f476ca287a171dd8aa94165fb27fadee1fc61f051840382aed326e1fd41d41b3",
		channelDigest:       "8ae1aa597fa146ebd3aa2ceddf360668dea5e526567e92b0321816a4e895bd2d",
		channelSignature:    "cd52448f226a599c8d61d568a6186d50bb263a2559e08c1c8389233e632d2a5839232b91d7b8974b862ff5540bae836c08792936ebd708682413e0c8e2884086",
		connectorID:         "532a8331032028874f4d7fa1953636f8b076876c493606f93a86d8880058b4a6",
		connectorCommitment: "647e3f78cfdbceb3714e6c4d043f60e3169cabae3d90be168aeb58796ed6533d",
		connectorDigest:     "ba76a200dbda978d3be9120def504046d2d378bc85c8da3ac8a1c80710faee01",
		connectorSignature:  "4b133aadf1d8a8c8bc6509842e245e01f71639e6e09fa2ed8ddf315faa38926f29e87b0cf850345f9cf74e82de367c7e31c96dfde6cb4603afb42750a551197d",
	},
	{
		name:                "ED25519",
		scheme:              ClprSignatureSchemeEd25519,
		privateKey:          "9d61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f60",
		publicKey:           "d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a",
		channelCommitment:   "4816bb408545af511a22953026b40500890c7fa0a44c83682e251b114a0b6496",
		channelDigest:       "8ae1aa597fa146ebd3aa2ceddf360668dea5e526567e92b0321816a4e895bd2d",
		channelSignature:    "8ef8546b7a464125e9bda37b52203c8a8f1c59c51deae85c6178632b4050141226d6018ed52fcf8afd2a4c141a04109220eb8c7cd9d923a06a2fff23155cab0f",
		connectorID:         "488d7b7fdcc061f1de9fe98666c1337f0fc90275889f2450d59e1e6e837e74b1",
		connectorCommitment: "a36b6836b720019a9fb1846df6bb4b76d13a1f276da0ebc2371b4d0eb8e30618",
		connectorDigest:     "e4b51c24e102cf97db8b2c5320d5b2664693789b2de186a4e8cd67c2f2cb0fca",
		connectorSignature:  "7374599d239cc8300e8f8f99057656a0c4b3e42f7fbaa99ab8cd3907da4d9fbe1ef28f12b7af9cd3f126fa513e3d6661b24fe790c1bd8a80275e8eaf69f1c70f",
	},
}

func clprHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	require.NoError(t, err)
	return decoded
}

func (vector clprOwnershipVector) key(t *testing.T) PrivateKey {
	t.Helper()
	if vector.scheme == ClprSignatureSchemeEd25519 {
		key, err := PrivateKeyFromStringEd25519(vector.privateKey)
		require.NoError(t, err)
		return key
	}
	key, err := PrivateKeyFromStringECDSA(vector.privateKey)
	require.NoError(t, err)
	return key
}

func clprBytes(length int, value byte) []byte {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = value
	}
	return bytes
}

// clprOffline freezes a transaction for two nodes without a client.
func clprOffline[T TransactionInterface](t *testing.T, tx *Transaction[T]) {
	t.Helper()
	tx.SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}})
	tx.SetTransactionID(TransactionIDGenerate(AccountID{Account: 5006}))
	tx.SetMaxTransactionFee(NewHbar(1))
	tx.SetTransactionValidDuration(120 * time.Second)
}

func TestUnitClprSignatureSchemeString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "ECDSA_SECP256K1", ClprSignatureSchemeEcdsaSecp256k1.String())
	assert.Equal(t, "ED25519", ClprSignatureSchemeEd25519.String())
	assert.Equal(t, "UNKNOWN (7)", ClprSignatureScheme(7).String())
}

func TestUnitClprServiceAddress(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "000000000000000000000000000000000000016e", hex.EncodeToString(_ClprServiceAddress[:]))
}

func TestUnitClprOwnershipVectors(t *testing.T) {
	t.Parallel()

	channelID := clprHex(t, clprVectorChannelID)
	salt := clprHex(t, clprVectorSalt)

	for _, vector := range clprOwnershipVectors {
		t.Run(vector.name, func(t *testing.T) {
			publicKey := vector.key(t).PublicKey()

			scheme, encoded, err := _ClprEncodePublicKey(publicKey)
			require.NoError(t, err)
			assert.Equal(t, vector.scheme, scheme)
			assert.Equal(t, vector.publicKey, hex.EncodeToString(encoded))

			commitment, err := ComputeClprChannelOwnershipCommitment(channelID, publicKey)
			require.NoError(t, err)
			assert.Equal(t, vector.channelCommitment, hex.EncodeToString(commitment))

			connectorID, err := ComputeClprConnectorID(channelID, publicKey, salt)
			require.NoError(t, err)
			assert.Equal(t, vector.connectorID, hex.EncodeToString(connectorID))

			connectorCommitment, err := ComputeClprConnectorOwnershipCommitment(connectorID, publicKey)
			require.NoError(t, err)
			assert.Equal(t, vector.connectorCommitment, hex.EncodeToString(connectorCommitment))

			assert.Equal(t, vector.channelDigest, hex.EncodeToString(_ClprKeccak(_ClprChannelOwnershipPreimage(channelID))))
			assert.Equal(t, clprHex(t, clprVectorChannelID), channelID)
		})
	}
}

func TestUnitClprOwnershipSignChannelVectors(t *testing.T) {
	t.Parallel()

	channelID := clprHex(t, clprVectorChannelID)

	for _, vector := range clprOwnershipVectors {
		t.Run(vector.name, func(t *testing.T) {
			proof, err := _ClprSignOwnership(vector.key(t), _ClprChannelOwnershipPreimage(channelID))
			require.NoError(t, err)

			assert.True(t, proof._IsComplete())
			assert.Equal(t, vector.scheme, proof.scheme)
			assert.Equal(t, vector.publicKey, hex.EncodeToString(proof.publicKey))
			assert.Equal(t, vector.channelSignature, hex.EncodeToString(proof.signature))
			assertClprSignatureVerifies(t, vector, clprHex(t, vector.channelDigest), proof.signature)
		})
	}
}

// assertClprSignatureVerifies checks the signature against the vector digest, independently of the SDK signer.
func assertClprSignatureVerifies(t *testing.T, vector clprOwnershipVector, digest []byte, signature []byte) {
	t.Helper()
	require.Len(t, signature, 64)

	if vector.scheme == ClprSignatureSchemeEd25519 {
		assert.True(t, vector.key(t).PublicKey().VerifySignedMessage(digest, signature))
		return
	}

	recovered := false
	for _, v := range []byte{27, 28} {
		publicKey, _, err := ecdsa.RecoverCompact(append([]byte{v}, signature...), digest)
		if err == nil && hex.EncodeToString(publicKey.SerializeUncompressed()[1:]) == vector.publicKey {
			recovered = true
		}
	}
	assert.True(t, recovered, "ECDSA signature does not recover the vector public key")
}

func TestUnitClprOwnershipRejectsBadLengths(t *testing.T) {
	t.Parallel()

	publicKey := clprOwnershipVectors[1].key(t).PublicKey()
	valid := clprBytes(32, 1)

	for _, length := range []int{0, 31, 33} {
		bad := clprBytes(length, 1)

		_, err := ComputeClprChannelOwnershipCommitment(bad, publicKey)
		assert.ErrorContains(t, err, "channelID")
		_, err = ComputeClprConnectorID(bad, publicKey, valid)
		assert.ErrorContains(t, err, "channelID")
		_, err = ComputeClprConnectorID(valid, publicKey, bad)
		assert.ErrorContains(t, err, "salt")
		_, err = ComputeClprConnectorOwnershipCommitment(bad, publicKey)
		assert.ErrorContains(t, err, "connectorID")
	}

	_, err := ComputeClprChannelOwnershipCommitment(clprBytes(31, 1), PublicKey{})
	assert.ErrorContains(t, err, "channelID")
}

func TestUnitClprOwnershipRejectsEmptyKey(t *testing.T) {
	t.Parallel()

	valid := clprBytes(32, 1)

	_, err := ComputeClprChannelOwnershipCommitment(valid, PublicKey{})
	assert.ErrorIs(t, err, errClprOwnershipKeyUnsupported)
	_, err = ComputeClprConnectorID(valid, PublicKey{}, valid)
	assert.ErrorIs(t, err, errClprOwnershipKeyUnsupported)
	_, err = ComputeClprConnectorOwnershipCommitment(valid, PublicKey{})
	assert.ErrorIs(t, err, errClprOwnershipKeyUnsupported)
	_, err = _ClprSignOwnership(PrivateKey{}, valid)
	assert.ErrorIs(t, err, errClprPrivateKeyRequired)
}

func TestUnitClprOwnershipProofFromProtobuf(t *testing.T) {
	t.Parallel()

	assert.Nil(t, _ClprOwnershipProofFromProtobuf(nil, nil, 1))

	for _, tt := range []struct {
		keyLength, signatureLength int
		complete                   bool
	}{{32, 0, false}, {0, 64, false}, {32, 64, true}, {1, 1, true}} {
		proof := _ClprOwnershipProofFromProtobuf(clprBytes(tt.keyLength, 1), clprBytes(tt.signatureLength, 2), 0)
		require.NotNil(t, proof)
		assert.Equal(t, tt.complete, proof._IsComplete())
		assert.Equal(t, ClprSignatureSchemeEcdsaSecp256k1, proof.scheme)
	}

	unknown := _ClprOwnershipProofFromProtobuf([]byte{1}, []byte{2}, 7)
	assert.Equal(t, "UNKNOWN (7)", unknown.scheme.String())
}
