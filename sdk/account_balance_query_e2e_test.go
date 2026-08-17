//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// AccountBalanceQuery is deprecated, so integration coverage targets MirrorNodeAccountBalanceQuery.

func TestIntegrationMirrorNodeAccountBalanceQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(env.OriginalOperatorID).
		Execute(env.Client)
	require.NoError(t, err)

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(env.OperatorID).
		Execute(env.Client)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, balance.Hbars.AsTinybar(), int64(0))
}

func TestIntegrationMirrorNodeAccountBalanceQueryCanGetTokenBalance(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	tokenID, err := createFungibleToken(&env)
	require.NoError(t, err)

	operatorID := env.Client.GetOperatorAccountID()

	// HIP-367 dropped token balances from the consensus response; assert via the mirror node.
	balance := requireMirrorTokenBalances(t, &env, operatorID, map[TokenID]uint64{tokenID: 1_000_000})
	assert.Equal(t, uint64(18), balance.TokenDecimals.Get(tokenID))

	err = CloseIntegrationTestEnv(env, &tokenID)
	require.NoError(t, err)
}

func TestIntegrationMirrorNodeAccountBalanceQueryNoIDError(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewMirrorNodeAccountBalanceQuery().
		Execute(env.Client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "either an account ID or a contract ID must be set")
}

// Acceptance #3: an account addressed by its EVM address alias resolves and returns its balance.
func TestIntegrationMirrorNodeAccountBalanceQueryCanGetBalanceByEvmAddress(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Auto-create an account by funding a fresh ECDSA public key's EVM address.
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	evmAddress := privateKey.PublicKey().ToEvmAddress()
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)

	tx, err := NewTransferTransaction().
		AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetIncludeChildren(true).SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Query the balance using the EVM-address alias itself (no account number known).
	balance, err := mirrorBalanceEventually(
		&env,
		NewMirrorNodeAccountBalanceQuery().SetAccountID(evmAddressAccount),
		func(balance AccountBalance) bool { return balance.Hbars.AsTinybar() > 0 },
	)
	require.NoError(t, err)

	assert.Equal(t, NewHbar(1).AsTinybar(), balance.Hbars.AsTinybar())
}

// Acceptance #4: an account addressed by an ED25519 public-key alias resolves and returns its balance.
func TestIntegrationMirrorNodeAccountBalanceQueryCanGetBalanceByPublicKeyAlias(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Auto-create an account by funding a fresh ED25519 public key's alias.
	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	aliasAccountID := *key.PublicKey().ToAccountID(0, 0)

	tx, err := NewTransferTransaction().
		AddHbarTransfer(aliasAccountID, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Query the balance using the public-key alias itself (no account number known).
	balance, err := mirrorBalanceEventually(
		&env,
		NewMirrorNodeAccountBalanceQuery().SetAccountID(aliasAccountID),
		func(balance AccountBalance) bool { return balance.Hbars.AsTinybar() > 0 },
	)
	require.NoError(t, err)

	assert.Equal(t, NewHbar(1).AsTinybar(), balance.Hbars.AsTinybar())
}

// Acceptance #5: a contract ID resolves through the accounts endpoint and returns its balance.
func TestIntegrationMirrorNodeAccountBalanceQueryCanGetContractBalance(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(testContractByteCode).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := *receipt.FileID

	resp, err = NewContractCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetGas(contractDeployGas).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hiero")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hiero-sdk-go::MirrorNodeAccountBalanceQuery").
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := *receipt.ContractID

	// Fund via crypto transfer: the greeter constructor is not payable, so an initial balance reverts.
	contractAccountID := AccountID{Shard: contractID.Shard, Realm: contractID.Realm, Account: contractID.Contract}
	tx, err := NewTransferTransaction().
		AddHbarTransfer(contractAccountID, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	balance, err := mirrorBalanceEventually(
		&env,
		NewMirrorNodeAccountBalanceQuery().SetContractID(contractID),
		func(balance AccountBalance) bool { return balance.Hbars.AsTinybar() > 0 },
	)
	require.NoError(t, err)

	assert.Equal(t, NewHbar(1).AsTinybar(), balance.Hbars.AsTinybar())

	_, err = NewContractDeleteTransaction().
		SetContractID(contractID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		Execute(env.Client)
	require.NoError(t, err)
}

// Acceptance #6: a non-existent account yields a clear error (mirror node 404) rather than a balance.
func TestIntegrationMirrorNodeAccountBalanceQueryNonExistentAccount(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// A syntactically valid but almost certainly non-existent account number.
	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 999999999}).
		Execute(env.Client)
	require.Error(t, err)
	// Assert the 404 specifically, so a transient 503 cannot masquerade as a pass.
	assert.Contains(t, err.Error(), "404")
}
