//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mirrorHbarBalanceRetryAttempts = 10
	mirrorHbarBalanceRetryDelay    = 2 * time.Second
)

// mirrorHbarBalanceEventually polls until ready accepts the result, absorbing mirror node lag.
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

	if err == nil {
		err = fmt.Errorf("mirror node did not reach the expected state within %d attempts", mirrorHbarBalanceRetryAttempts)
	}

	return balance, err
}

func nonZeroHbars(balance MirrorNodeAccountBalance) bool {
	return balance.Hbars.AsTinybar() > 0
}

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

func TestIntegrationMirrorNodeAccountBalanceQueryNoIDError(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewMirrorNodeAccountBalanceQuery().Execute(env.Client)
	require.ErrorIs(t, err, errMirrorNodeAccountBalanceQueryNoAccountID)
}

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
		SetGas(400_000).
		SetBytecodeFileID(fileID).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := *receipt.ContractID

	defer func() {
		_, err := NewContractDeleteTransaction().
			SetContractID(contractID).
			SetTransferAccountID(env.Client.GetOperatorAccountID()).
			Execute(env.Client)
		require.NoError(t, err)
		_, err = NewFileDeleteTransaction().SetFileID(fileID).Execute(env.Client)
		require.NoError(t, err)
	}()

	contractAccountID := AccountID{Shard: contractID.Shard, Realm: contractID.Realm, Account: contractID.Contract}

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

}

func TestIntegrationMirrorNodeAccountBalanceQueryNonExistentAccountErrors(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	_, err := NewMirrorNodeAccountBalanceQuery().
		SetAccountID(AccountID{Account: 999999999}).
		Execute(env.Client)

	var status ErrHederaPreCheckStatus
	require.ErrorAs(t, err, &status)
	assert.Equal(t, StatusInvalidAccountID, status.Status)
}
