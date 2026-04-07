//go:build unit
// +build unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitDilithiumKeyGenerate(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyGenerateDilithium()
	require.NoError(t, err)

	pub := key.PublicKey()
	assert.NotEmpty(t, key.BytesRaw())
	assert.NotEmpty(t, pub.BytesRaw())
	assert.Equal(t, 4000, len(key.BytesRaw()), "Dilithium Mode3 private key should be 4000 bytes")
	assert.Equal(t, 1952, len(pub.BytesRaw()), "Dilithium Mode3 public key should be 1952 bytes")
}

func TestUnitDilithiumSignAndVerify(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyGenerateDilithium()
	require.NoError(t, err)

	message := []byte("test transaction body bytes for PQC signing")
	sig := key.Sign(message)
	assert.NotEmpty(t, sig)
	assert.Equal(t, 3293, len(sig), "Dilithium Mode3 signature should be 3293 bytes")

	pub := key.PublicKey()
	assert.True(t, pub.VerifySignedMessage(message, sig))
}

func TestUnitDilithiumVerifyWrongMessage(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyGenerateDilithium()
	require.NoError(t, err)

	message := []byte("original message")
	sig := key.Sign(message)

	pub := key.PublicKey()
	assert.False(t, pub.VerifySignedMessage([]byte("tampered message"), sig))
}

func TestUnitDilithiumVerifyWrongKey(t *testing.T) {
	t.Parallel()

	key1, err := PrivateKeyGenerateDilithium()
	require.NoError(t, err)

	key2, err := PrivateKeyGenerateDilithium()
	require.NoError(t, err)

	message := []byte("test message")
	sig := key1.Sign(message)

	// Verify with wrong public key should fail
	assert.False(t, key2.PublicKey().VerifySignedMessage(message, sig))
}

func TestUnitDilithiumKeyRoundTrip(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyGenerateDilithium()
	require.NoError(t, err)

	// Private key round-trip
	privBytes := key.BytesRaw()
	restored, err := PrivateKeyFromBytesDilithium(privBytes)
	require.NoError(t, err)
	assert.Equal(t, key.BytesRaw(), restored.BytesRaw())

	// Public key round-trip
	pub := key.PublicKey()
	pubBytes := pub.BytesRaw()
	restoredPub, err := PublicKeyFromBytesDilithium(pubBytes)
	require.NoError(t, err)
	assert.Equal(t, pub.BytesRaw(), restoredPub.BytesRaw())

	// Sign with restored key, verify with restored public key
	message := []byte("round-trip test message")
	sig := restored.Sign(message)
	assert.True(t, restoredPub.VerifySignedMessage(message, sig))
}

func TestUnitDilithiumKeyString(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyGenerateDilithium()
	require.NoError(t, err)

	// String representation should be non-empty hex
	assert.NotEmpty(t, key.String())
	assert.NotEmpty(t, key.StringRaw())
	assert.NotEmpty(t, key.PublicKey().String())
	assert.NotEmpty(t, key.PublicKey().StringRaw())

	// String should be hex-encoded raw bytes (8000 hex chars for 4000-byte private key)
	assert.Equal(t, 8000, len(key.StringRaw()))
	// Public key: 3904 hex chars for 1952 bytes
	assert.Equal(t, 3904, len(key.PublicKey().StringRaw()))
}

func TestUnitDilithiumKeyFromBytesInvalidLength(t *testing.T) {
	t.Parallel()

	_, err := PrivateKeyFromBytesDilithium([]byte("too short"))
	assert.Error(t, err)

	_, err = PublicKeyFromBytesDilithium([]byte("too short"))
	assert.Error(t, err)
}
