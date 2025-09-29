//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationLambdaSStoreTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	hexBytecode, err := hex.DecodeString(EMPTY_CONTRACT)
	require.NoError(t, err)

	resp, err := NewContractCreateTransaction().
		SetGas(300_000).
		SetBytecode(hexBytecode).
		Execute(env.Client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	contractID := receipt.ContractID

	hookDetail := NewHookCreationDetails().
		SetExtensionPoint(ACCOUNT_ALLOWANCE_HOOK).
		SetHookId(1).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetLambdaEvmHook(*NewLambdaEvmHook().SetEvmHookSpec(*NewEvmHookSpec().SetContractId(*contractID)))

	resp, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(2)).
		SetMaxAutomaticTokenAssociations(100).
		AddHook(*hookDetail).
		Execute(env.Client)

	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	accountId := *receipt.AccountID

	resp, err = NewLambdaSStoreTransaction().
		SetHookId(*NewHookId().SetEntityId(*NewHookEntityId().SetAccountId(accountId)).SetHookId(1)).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	require.NoError(t, err)
}
