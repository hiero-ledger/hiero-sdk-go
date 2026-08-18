//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mirrorHbarBalanceRetryAttempts = 10
	mirrorHbarBalanceRetryDelay    = 2 * time.Second
)

// mirrorHbarBalanceEventually polls until ready accepts the result, absorbing mirror node ingestion
// lag. On timeout it returns the last balance and error so the caller's assertion reports it.
func mirrorHbarBalanceEventually(env *IntegrationTestEnv, query *MirrorNodeAccountBalanceQuery, ready func(MirrorNodeAccountBalance) bool) (MirrorNodeAccountBalance, error) {
	var balance MirrorNodeAccountBalance
	var err error

	for attempt := 0; attempt < mirrorHbarBalanceRetryAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(mirrorHbarBalanceRetryDelay)
		}
		balance, err = query.Execute(env.Client)
		if err == nil && ready(balance) {
			return balance, nil
		}
	}

	return balance, err
}

func nonZeroHbars(balance MirrorNodeAccountBalance) bool {
	return balance.Hbars.AsTinybar() > 0
}

// Acceptance 1: a valid account ID returns the hbar balance the mirror node holds.
func TestIntegrationMirrorNodeAccountBalanceQueryCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	balance, err := mirrorHbarBalanceEventually(
		&env,
		NewMirrorNodeAccountBalanceQuery().SetAccountID(env.OperatorID),
		nonZeroHbars,
	)
	require.NoError(t, err)
	assert.Positive(t, balance.Hbars.AsTinybar(), "the operator account must hold hbars")
}

// The query is free and must work without an operator, unlike the consensus-node queries.
func TestIntegrationMirrorNodeAccountBalanceQueryWithoutOperator(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	operatorID := env.OperatorID
	env.Client.operator = nil

	balance, err := mirrorHbarBalanceEventually(
		&env,
		NewMirrorNodeAccountBalanceQuery().SetAccountID(operatorID),
		nonZeroHbars,
	)
	require.NoError(t, err)
	assert.Positive(t, balance.Hbars.AsTinybar())
}

// Acceptance 6: an unset account ID fails before any request is made.
func TestIntegrationMirrorNodeAccountBalanceQueryNoIDError(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewMirrorNodeAccountBalanceQuery().Execute(env.Client)
	require.ErrorIs(t, err, errMirrorNodeAccountBalanceQueryNoAccountID)
}

// Acceptance 2: an account addressed by its EVM address alias resolves and returns its balance.
func TestIntegrationMirrorNodeAccountBalanceQueryCanGetBalanceByEvmAddress(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Auto-create the account by funding a fresh ECDSA key's EVM address.
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	evmAddressAccount, err := AccountIDFromEvmPublicAddress(privateKey.PublicKey().ToEvmAddress())
	require.NoError(t, err)

	tx, err := NewTransferTransaction().
		AddHbarTransfer(evmAddressAccount, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetIncludeChildren(true).SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Query by the alias itself; no account number is known here.
	balance, err := mirrorHbarBalanceEventually(
		&env,
		NewMirrorNodeAccountBalanceQuery().SetAccountID(evmAddressAccount),
		nonZeroHbars,
	)
	require.NoError(t, err)
	assert.Equal(t, NewHbar(1).AsTinybar(), balance.Hbars.AsTinybar())
}

// Acceptance 3: an account addressed by an ED25519 public-key alias resolves and returns its balance.
func TestIntegrationMirrorNodeAccountBalanceQueryCanGetBalanceByPublicKeyAlias(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

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

	balance, err := mirrorHbarBalanceEventually(
		&env,
		NewMirrorNodeAccountBalanceQuery().SetAccountID(aliasAccountID),
		nonZeroHbars,
	)
	require.NoError(t, err)
	assert.Equal(t, NewHbar(1).AsTinybar(), balance.Hbars.AsTinybar())
}

// Acceptance 4: a contract's balance is read through the same account.id parameter.
func TestIntegrationMirrorNodeAccountBalanceQueryCanGetContractBalance(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents([]byte(SIMPLE_SMART_CONTRACT_BYTECODE)).
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

	// No SetContractID: a contract goes through account.id as the equivalent AccountID.
	contractAccountID := AccountID{Shard: contractID.Shard, Realm: contractID.Realm, Account: contractID.Contract}

	// Fund by transfer: the constructor is not payable, so an initial balance reverts.
	tx, err := NewTransferTransaction().
		AddHbarTransfer(contractAccountID, NewHbar(1)).
		AddHbarTransfer(env.OperatorID, NewHbar(-1)).
		Execute(env.Client)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	balance, err := mirrorHbarBalanceEventually(
		&env,
		NewMirrorNodeAccountBalanceQuery().SetAccountID(contractAccountID),
		nonZeroHbars,
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

// Acceptance 5: an unknown account yields a zero balance, not an error -- the endpoint returns an
// empty list rather than a 404, unlike the consensus-node query's INVALID_ACCOUNT_ID.
func TestIntegrationMirrorNodeAccountBalanceQueryNonExistentAccountIsZero(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	balance, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 999999999}).
		Execute(env.Client)

	require.NoError(t, err, "an unknown account is an empty balances array, not an error")
	assert.Zero(t, balance.Hbars.AsTinybar())
}
