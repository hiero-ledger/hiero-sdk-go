//go:build all || e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

// Limited max auto association tests
func TestLimitedMaxAutoAssociationsFungibleTokensFlow(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create token1
	tokenID1, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create token2
	tokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	// account create with 1 max auto associations
	receiver, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(1)
	})
	require.NoError(t, err)

	// transfer token1 to receiver account
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer token2 to the receiver should fail with NO_REMAINING_AUTOMATIC_ASSOCIATIONS
	tokenTransferTransaction2, err := NewTransferTransaction().
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID2, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction2.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "NO_REMAINING_AUTOMATIC_ASSOCIATIONS")
}

func TestLimitedMaxAutoAssociationsNFTsFlow(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create 2 NFT collections and mint 10 NFTs for each collection
	nftID1, err := createNft(&env)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(mintMetadata).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftID2, err := createNft(&env)
	require.NoError(t, err)

	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(mintMetadata).Execute(env.Client)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// account create with 1 max auto associations
	receiver, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(1)
	})
	require.NoError(t, err)

	// transfer nftID1 nfts to receiver account
	tokenTransferTransaction, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[2]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[3]), env.OperatorID, receiver).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nftID2 nft to receiver should fail with NO_REMAINING_AUTOMATIC_ASSOCIATIONS
	tokenTransferTransaction2, err := NewTransferTransaction().
		AddNftTransfer(nftID2.Nft(serials[0]), env.OperatorID, receiver).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction2.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "NO_REMAINING_AUTOMATIC_ASSOCIATIONS")
}

func TestLimitedMaxAutoAssociationsFungibleTokensWithManualAssociate(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create token1
	tokenID1, err := createFungibleToken(&env)

	// account create
	receiver, key, err := createAccount(&env)
	require.NoError(t, err)

	frozenAssociateTxn, err := NewTokenAssociateTransaction().SetAccountID(receiver).AddTokenID(tokenID1).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenAssociateTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer token1 to receiver account
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -10).
		AddTokenTransfer(tokenID1, receiver, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 10
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(receiver).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(10), tokenBalance.Tokens.Get(tokenID1))
}

func TestLimitedMaxAutoAssociationsNFTsManualAssociate(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create NFT collection and mint 10
	nftID1, err := createNft(&env)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(mintMetadata).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// account create
	receiver, key, err := createAccount(&env)
	require.NoError(t, err)

	frozenAssociateTxn, err := NewTokenAssociateTransaction().SetAccountID(receiver).AddTokenID(nftID1).FreezeWith(env.Client)
	require.NoError(t, err)
	resp, err := frozenAssociateTxn.Sign(key).Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nftID1 nfts to receiver account
	tokenTransferTransaction, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[2]), env.OperatorID, receiver).
		AddNftTransfer(nftID1.Nft(serials[3]), env.OperatorID, receiver).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

// HIP-904 Unlimited max auto association tests
func TestUnlimitedMaxAutoAssociationsExecutes(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// account create with unlimited max auto associations - verify it executes
	_, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	accountID, newKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// update the account with unlimited max auto associations
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFungibleTokens(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create token1
	tokenID1, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create token2
	tokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	// account create with unlimited max auto associations
	accountID1, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)
	// create account with 100 max auto associations
	accountID2, newKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// update the account with unlimited max auto associations
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer to both receivers some token1 tokens
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, accountID1, 1000).
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, accountID2, 1000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer to both receivers some token2 tokens
	tokenTransferTransaction, err = NewTransferTransaction().
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, accountID1, 1000).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, accountID2, 1000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receivers is 1000
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID1).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID2))

	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(accountID2).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID2))
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFungibleTokensWithDecimals(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create token1
	tokenID1, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetDecimals(10)
	})
	require.NoError(t, err)

	// create token2
	tokenID2, err := createFungibleToken(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetDecimals(10)
	})
	require.NoError(t, err)

	// account create with unlimited max auto associations
	accountID, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// transfer some token1 and token2 tokens
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransferWithDecimals(tokenID1, env.Client.GetOperatorAccountID(), -1000, 10).
		AddTokenTransferWithDecimals(tokenID1, accountID, 1000, 10).
		AddTokenTransferWithDecimals(tokenID2, env.Client.GetOperatorAccountID(), -1000, 10).
		AddTokenTransferWithDecimals(tokenID2, accountID, 1000, 10).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receiver is 1000
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID2))
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFromFungibleTokens(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create spender account which will be approved to spend
	spender, spenderKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(10)
	})
	require.NoError(t, err)

	// create token1
	tokenID1, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create token2
	tokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	// account create with unlimited max auto associations
	accountID, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// approve the spender
	approve, err := NewAccountAllowanceApproveTransaction().
		AddTokenApproval(tokenID1, spender, 2000).
		AddTokenApproval(tokenID2, spender, 2000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = approve.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transferFrom some token1 and token2 tokens
	env.Client.SetOperator(spender, spenderKey)
	tokenTransferTransactionFrozen, err := NewTransferTransaction().
		AddApprovedTokenTransfer(tokenID1, env.OperatorID, -1000, true).
		AddTokenTransfer(tokenID1, accountID, 1000).
		AddApprovedTokenTransfer(tokenID2, env.OperatorID, -1000, true).
		AddTokenTransfer(tokenID2, accountID, 1000).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokenTransferTransaction, err := tokenTransferTransactionFrozen.Sign(spenderKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// verify the balance of the receiver is 1000
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1000), tokenBalance.Tokens.Get(tokenID2))
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferNFTs(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create 2 NFT collections and mint 10 NFTs for each collection
	nftID1, err := createNft(&env)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(mintMetadata).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftID2, err := createNft(&env)
	require.NoError(t, err)

	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(mintMetadata).Execute(env.Client)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// account create with unlimited max auto associations
	accountID1, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	accountID2, newKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// account update with unlimited max auto associations
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1).
		SetAccountID(accountID2).
		FreezeWith(env.Client)
	require.NoError(t, err)

	accountUpdate, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = accountUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nft1 to both receivers, 2 for each
	tokenTransferTransaction, err := NewTransferTransaction().
		AddNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, accountID1).
		AddNftTransfer(nftID1.Nft(serials[2]), env.OperatorID, accountID2).
		AddNftTransfer(nftID1.Nft(serials[3]), env.OperatorID, accountID2).
		Execute(env.Client)

	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transfer nft2 to both receivers, 2 for each
	tokenTransferTransaction, err = NewTransferTransaction().
		AddNftTransfer(nftID2.Nft(serials[0]), env.OperatorID, accountID1).
		AddNftTransfer(nftID2.Nft(serials[1]), env.OperatorID, accountID1).
		AddNftTransfer(nftID2.Nft(serials[2]), env.OperatorID, accountID2).
		AddNftTransfer(nftID2.Nft(serials[3]), env.OperatorID, accountID2).
		Execute(env.Client)

	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the balance of the receivers is 2
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID1).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID2))

	tokenBalance, err = NewAccountBalanceQuery().SetAccountID(accountID2).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID2))
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFromNFTs(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// create spender account which will be approved to spend
	spender, spenderKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(10)
	})
	require.NoError(t, err)

	// create 2 NFT collections and mint 10 NFTs for each collection
	nftID1, err := createNft(&env)
	require.NoError(t, err)

	mint, err := NewTokenMintTransaction().SetTokenID(nftID1).SetMetadatas(mintMetadata).Execute(env.Client)
	receipt, err := mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	nftID2, err := createNft(&env)
	require.NoError(t, err)

	mint, err = NewTokenMintTransaction().SetTokenID(nftID2).SetMetadatas(mintMetadata).Execute(env.Client)
	receipt, err = mint.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	serials := receipt.SerialNumbers

	// account create with unlimited max auto associations
	accountID, _, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(-1)
	})
	require.NoError(t, err)

	// approve the spender
	approve, err := NewAccountAllowanceApproveTransaction().
		AddAllTokenNftApproval(nftID1, spender).
		AddAllTokenNftApproval(nftID2, spender).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = approve.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// transferFrom some nft1 nfts
	env.Client.SetOperator(spender, spenderKey)
	tokenTransferTransactionFrozen, err := NewTransferTransaction().
		AddApprovedNftTransfer(nftID1.Nft(serials[0]), env.OperatorID, accountID, true).
		AddApprovedNftTransfer(nftID1.Nft(serials[1]), env.OperatorID, accountID, true).
		AddApprovedNftTransfer(nftID2.Nft(serials[0]), env.OperatorID, accountID, true).
		AddApprovedNftTransfer(nftID2.Nft(serials[1]), env.OperatorID, accountID, true).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tokenTransferTransaction, err := tokenTransferTransactionFrozen.Sign(spenderKey).Execute(env.Client)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	env.Client.SetOperator(env.OperatorID, env.OperatorKey)

	// verify the balance of the receiver is 2
	tokenBalance, err := NewAccountBalanceQuery().SetAccountID(accountID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID1))
	assert.Equal(t, uint64(2), tokenBalance.Tokens.Get(nftID2))
}

func TestUnlimitedMaxAutoAssociationsFailsWithInvalid(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// account create with -2 and with -1000 max auto associations
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	_, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetMaxAutomaticTokenAssociations(-2).
		Execute(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTO_ASSOCIATIONS")

	_, err = NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(0)).
		SetMaxAutomaticTokenAssociations(-1000).
		Execute(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTO_ASSOCIATIONS")

	// create account with 100 max auto associations
	accountID, newKey, err := createAccount(&env, func(tx *AccountCreateTransaction) {
		tx.SetMaxAutomaticTokenAssociations(100)
	})
	require.NoError(t, err)

	// account update with -2 max auto associations - should fail
	accountUpdateFrozen, err := NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-2).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err := accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTO_ASSOCIATIONS")

	// account update with -1000 max auto associations - should fail
	accountUpdateFrozen, err = NewAccountUpdateTransaction().
		SetMaxAutomaticTokenAssociations(-1000).
		SetAccountID(accountID).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err = accountUpdateFrozen.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	_, err = tx.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_MAX_AUTO_ASSOCIATIONS")
}

// HIP-904 contract tests

func TestUnlimitedMaxAutoAssociationsContractAllowsToTransferFungibleTokens(t *testing.T) {
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// deploy the bytecode
	fileCreate, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetContents(testContractByteCode).
		Execute(env.Client)
	require.NoError(t, err)

	fileReceipt, err := fileCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := *fileReceipt.FileID

	// create a contract with unlimited max auto associations
	contractCreate, err := NewContractCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetGas(contractDeployGas).
		SetBytecodeFileID(fileID).
		SetMaxAutomaticTokenAssociations(-1).
		Execute(env.Client)
	require.NoError(t, err)

	contractReceipt, err := contractCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	require.NotNil(t, contractReceipt.ContractID)
	contractID := *contractReceipt.ContractID

	// the contract's own account is addressed by the same shard.realm.num
	contractAccountID := AccountID{Shard: contractID.Shard, Realm: contractID.Realm, Account: contractID.Contract}

	// create token1
	tokenID1, err := createFungibleToken(&env)
	require.NoError(t, err)

	// create token2
	tokenID2, err := createFungibleToken(&env)
	require.NoError(t, err)

	// transfer both tokens to the contract with no explicit association
	tokenTransferTransaction, err := NewTransferTransaction().
		AddTokenTransfer(tokenID1, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID1, contractAccountID, 1000).
		AddTokenTransfer(tokenID2, env.Client.GetOperatorAccountID(), -1000).
		AddTokenTransfer(tokenID2, contractAccountID, 1000).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = tokenTransferTransaction.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify both balances arrived on the contract
	contractBalance, err := NewAccountBalanceQuery().SetContractID(contractID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), contractBalance.Tokens.Get(tokenID1))
	assert.Equal(t, uint64(1000), contractBalance.Tokens.Get(tokenID2))
}

func TestUnlimitedMaxAutoAssociationsContractUpdate(t *testing.T) {
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// deploy the bytecode
	fileCreate, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetContents(testContractByteCode).
		Execute(env.Client)
	require.NoError(t, err)

	fileReceipt, err := fileCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := *fileReceipt.FileID

	// create a contract with a limited number of max auto associations
	contractCreate, err := NewContractCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetGas(contractDeployGas).
		SetBytecodeFileID(fileID).
		SetMaxAutomaticTokenAssociations(1).
		Execute(env.Client)
	require.NoError(t, err)

	contractReceipt, err := contractCreate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	require.NotNil(t, contractReceipt.ContractID)
	contractID := *contractReceipt.ContractID

	// update the contract to unlimited max auto associations
	contractUpdate, err := NewContractUpdateTransaction().
		SetContractID(contractID).
		SetMaxAutomaticTokenAssociations(-1).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = contractUpdate.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// verify the contract now reports unlimited (-1) max auto associations
	info, err := NewContractInfoQuery().SetContractID(contractID).Execute(env.Client)
	require.NoError(t, err)
	assert.Equal(t, int32(-1), info.MaxAutomaticTokenAssociations)
}
