//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitHookCreationDetailsExtensionPoint(t *testing.T) {
	t.Parallel()

	hcd := NewHookCreationDetails()

	// Test default value
	assert.Equal(t, HookExtensionPoint(0), hcd.GetExtensionPoint())

	// Test setting extension point
	hcd.SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK)
	assert.Equal(t, ACCOUNT_ALLOWANCE_HOOK, hcd.GetExtensionPoint())
}

func TestUnitHookCreationDetailsHookId(t *testing.T) {
	t.Parallel()

	hcd := NewHookCreationDetails()

	// Test default value
	assert.Equal(t, int64(0), hcd.GetHookId())

	// Test setting hook ID
	hookId := int64(12345)
	hcd.SetHookId(hookId)
	assert.Equal(t, hookId, hcd.GetHookId())

	// Test negative hook ID
	negativeHookId := int64(-1)
	hcd.SetHookId(negativeHookId)
	assert.Equal(t, negativeHookId, hcd.GetHookId())
}

func TestUnitHookCreationDetailsLambdaEvmHook(t *testing.T) {
	t.Parallel()

	hcd := NewHookCreationDetails()

	// Test default value
	assert.Equal(t, EvmHook{}, hcd.GetLambdaEvmHook())

	// Test setting lambda EVM hook
	lambdaHook := NewEvmHook()
	hcd.SetEvmHook(*lambdaHook)
	assert.Equal(t, *lambdaHook, hcd.GetLambdaEvmHook())
}

func TestUnitHookCreationDetailsAdminKey(t *testing.T) {
	t.Parallel()

	hcd := NewHookCreationDetails()

	// Test default value
	assert.Nil(t, hcd.GetAdminKey())

	// Test setting admin key with Ed25519 public key
	ed25519PrivateKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ed25519PublicKey := ed25519PrivateKey.PublicKey()
	hcd.SetAdminKey(ed25519PublicKey)
	assert.Equal(t, ed25519PublicKey, hcd.GetAdminKey())

	// Test setting admin key with ECDSA public key
	ecdsaPrivateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	ecdsaPublicKey := ecdsaPrivateKey.PublicKey()
	hcd.SetAdminKey(ecdsaPublicKey)
	assert.Equal(t, ecdsaPublicKey, hcd.GetAdminKey())

	// Test setting admin key with ContractID
	contractID, err := ContractIDFromString("0.0.123")
	require.NoError(t, err)
	hcd.SetAdminKey(contractID)
	assert.Equal(t, contractID, hcd.GetAdminKey())
}

func TestUnitHookCreationDetailsMethodChaining(t *testing.T) {
	t.Parallel()

	ed25519PrivateKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ed25519PublicKey := ed25519PrivateKey.PublicKey()

	lambdaHook := NewEvmHook()

	// Test method chaining
	hcd := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(12345).
		SetEvmHook(*lambdaHook).
		SetAdminKey(ed25519PublicKey)

	assert.Equal(t, ACCOUNT_ALLOWANCE_HOOK, hcd.GetExtensionPoint())
	assert.Equal(t, int64(12345), hcd.GetHookId())
	assert.Equal(t, *lambdaHook, hcd.GetLambdaEvmHook())
	assert.Equal(t, ed25519PublicKey, hcd.GetAdminKey())
}

func TestUnitHookCreationDetailsToProtobuf(t *testing.T) {
	t.Parallel()

	// Test with all fields set
	ed25519PrivateKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ed25519PublicKey := ed25519PrivateKey.PublicKey()

	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)
	lambdaHook := NewEvmHook().SetContractId(&contractID)

	hcd := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(789).
		SetEvmHook(*lambdaHook).
		SetAdminKey(ed25519PublicKey)

	pb := hcd.toProtobuf()
	require.NotNil(t, pb)
	assert.Equal(t, services.HookExtensionPoint(ACCOUNT_ALLOWANCE_HOOK), pb.ExtensionPoint)
	assert.Equal(t, int64(789), pb.HookId)
	assert.NotNil(t, pb.Hook)
	assert.NotNil(t, pb.Hook.(*services.HookCreationDetails_EvmHook).EvmHook)
	assert.NotNil(t, pb.AdminKey)
}

func TestUnitHookCreationDetailsToProtobufWithNilAdminKey(t *testing.T) {
	t.Parallel()

	// Test with nil admin key - we need to create a valid LambdaEvmHook to avoid nil pointer dereference
	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)
	lambdaHook := NewEvmHook().SetContractId(&contractID)

	hcd := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(789).
		SetEvmHook(*lambdaHook)

	lambdaEvmHook := &services.EvmHook{
		Spec: &services.EvmHookSpec{
			BytecodeSource: &services.EvmHookSpec_ContractId{
				ContractId: contractID._ToProtobuf(),
			},
		},
	}

	pb := hcd.toProtobuf()
	require.NotNil(t, pb)
	assert.Equal(t, services.HookExtensionPoint(ACCOUNT_ALLOWANCE_HOOK), pb.ExtensionPoint)
	assert.Equal(t, int64(789), pb.HookId)
	assert.Equal(t, lambdaEvmHook, pb.Hook.(*services.HookCreationDetails_EvmHook).EvmHook)
	assert.Nil(t, pb.AdminKey)
}

func TestUnitHookCreationDetailsFromProtobufNoStorageUpdates(t *testing.T) {
	t.Parallel()

	// Create a protobuf message
	ed25519PrivateKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ed25519PublicKey := ed25519PrivateKey.PublicKey()

	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)
	lambdaHook := NewEvmHook().SetContractId(&contractID)

	pb := &services.HookCreationDetails{
		ExtensionPoint: services.HookExtensionPoint(ACCOUNT_ALLOWANCE_HOOK),
		HookId:         789,
		Hook: &services.HookCreationDetails_EvmHook{
			EvmHook: &services.EvmHook{
				Spec: &services.EvmHookSpec{
					BytecodeSource: &services.EvmHookSpec_ContractId{
						ContractId: contractID._ToProtobuf(),
					},
				},
			},
		},
		AdminKey: ed25519PublicKey._ToProtoKey(),
	}

	hcd := hookCreationDetailsFromProtobuf(pb)
	assert.Equal(t, *lambdaHook.GetContractId(), *hcd.GetLambdaEvmHook().GetContractId())
	assert.Equal(t, lambdaHook.GetStorageUpdates(), hcd.GetLambdaEvmHook().GetStorageUpdates())
}

func TestUnitHookCreationDetailsFromProtobuf(t *testing.T) {
	t.Parallel()

	// Create a protobuf message
	ed25519PrivateKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ed25519PublicKey := ed25519PrivateKey.PublicKey()

	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)
	storageUpdate := NewEvmHookStorageSlot().SetKey([]byte{1, 2, 3}).SetValue([]byte{4, 5, 6})
	lambdaHook := NewEvmHook().SetContractId(&contractID).AddStorageUpdate(storageUpdate)

	pb := &services.HookCreationDetails{
		ExtensionPoint: services.HookExtensionPoint(ACCOUNT_ALLOWANCE_HOOK),
		HookId:         789,
		Hook: &services.HookCreationDetails_EvmHook{
			EvmHook: &services.EvmHook{
				Spec: &services.EvmHookSpec{
					BytecodeSource: &services.EvmHookSpec_ContractId{
						ContractId: contractID._ToProtobuf(),
					},
				},
				StorageUpdates: []*services.EvmHookStorageUpdate{{
					Update: &services.EvmHookStorageUpdate_StorageSlot{
						StorageSlot: &services.EvmHookStorageSlot{
							Key:   storageUpdate.GetKey(),
							Value: storageUpdate.GetValue(),
						},
					},
				}},
			},
		},
		AdminKey: ed25519PublicKey._ToProtoKey(),
	}

	hcd := hookCreationDetailsFromProtobuf(pb)
	assert.Equal(t, ACCOUNT_ALLOWANCE_HOOK, hcd.GetExtensionPoint())
	assert.Equal(t, int64(789), hcd.GetHookId())
	assert.Equal(t, *lambdaHook.GetContractId(), *hcd.GetLambdaEvmHook().GetContractId())
	assert.Equal(t, *storageUpdate, hcd.GetLambdaEvmHook().GetStorageUpdates()[0])
	assert.NotNil(t, hcd.GetAdminKey())
}

func TestUnitHookCreationDetailsFromProtobufWithNilAdminKey(t *testing.T) {
	t.Parallel()

	// Create a protobuf message with nil admin key
	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)
	lambdaHook := NewEvmHook().SetContractId(&contractID)

	pb := &services.HookCreationDetails{
		ExtensionPoint: services.HookExtensionPoint(ACCOUNT_ALLOWANCE_HOOK),
		HookId:         789,
		Hook: &services.HookCreationDetails_EvmHook{
			EvmHook: lambdaHook.toProtobuf(),
		},
		AdminKey: nil,
	}

	hcd := hookCreationDetailsFromProtobuf(pb)
	assert.Equal(t, ACCOUNT_ALLOWANCE_HOOK, hcd.GetExtensionPoint())
	assert.Equal(t, int64(789), hcd.GetHookId())
	assert.Nil(t, hcd.GetAdminKey())
}

func TestUnitHookCreationDetailsRoundTrip(t *testing.T) {
	t.Parallel()

	// Test round trip conversion: struct -> protobuf -> struct
	ed25519PrivateKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	ed25519PublicKey := ed25519PrivateKey.PublicKey()

	contractID, err := ContractIDFromString("0.0.456")
	require.NoError(t, err)
	lambdaHook := NewEvmHook().SetContractId(&contractID)

	original := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(789).
		SetEvmHook(*lambdaHook).
		SetAdminKey(ed25519PublicKey)

	// Convert to protobuf and back
	pb := original.toProtobuf()
	converted := hookCreationDetailsFromProtobuf(pb)

	// Compare original and converted
	assert.Equal(t, original.GetExtensionPoint(), converted.GetExtensionPoint())
	assert.Equal(t, original.GetHookId(), converted.GetHookId())
	assert.Equal(t, original.GetAdminKey(), converted.GetAdminKey())
}

func TestUnitHookCreationDetailsEdgeCases(t *testing.T) {
	t.Parallel()

	// Test with zero values
	hcd := NewHookCreationDetails()
	assert.Equal(t, HookExtensionPoint(0), hcd.GetExtensionPoint())
	assert.Equal(t, int64(0), hcd.GetHookId())
	assert.Equal(t, EvmHook{}, hcd.GetLambdaEvmHook())
	assert.Nil(t, hcd.GetAdminKey())

	// Test with maximum int64 value
	hcd.SetHookId(9223372036854775807)
	assert.Equal(t, int64(9223372036854775807), hcd.GetHookId())

	// Test with minimum int64 value
	hcd.SetHookId(-9223372036854775808)
	assert.Equal(t, int64(-9223372036854775808), hcd.GetHookId())
}

func TestUnitHookCreationDetailsEmptyLambdaEvmHook(t *testing.T) {
	t.Parallel()

	// Test with empty LambdaEvmHook
	emptyLambdaHook := EvmHook{}
	hcd := NewHookCreationDetails().SetEvmHook(emptyLambdaHook)

	assert.Equal(t, emptyLambdaHook, hcd.GetLambdaEvmHook())

	pb := hcd.toProtobuf()
	hookCreationDetailsFromProtobuf(pb)
}
