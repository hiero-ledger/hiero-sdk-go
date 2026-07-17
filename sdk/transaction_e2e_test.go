//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	protobuf "google.golang.org/protobuf/proto"
)

func TestIntegrationTransactionAddSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	updateBytes, err := tx.ToBytes()
	require.NoError(t, err)

	sig1, err := newKey.SignTransaction(tx)
	require.NoError(t, err)

	tx2, err := TransactionFromBytes(updateBytes)
	require.NoError(t, err)

	if newTx, ok := tx2.(AccountDeleteTransaction); ok {
		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).Execute(env.Client)
		require.NoError(t, err)
	}

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTransactionSignTransaction(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = newKey.SignTransaction(tx)
	require.NoError(t, err)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTransactionGetHash(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err = tx.SignWithOperator(env.Client)
	require.NoError(t, err)

	hash, err := tx.GetTransactionHash()
	require.NoError(t, err)

	resp, err := tx.Execute(env.Client)
	require.NoError(t, err)

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)

	assert.Equal(t, hash, record.TransactionHash)

}

func DisabledTestTransactionFromBytes(t *testing.T) { // nolint
	id := TransactionIDGenerate(AccountID{0, 0, 542348, nil, nil, nil})

	TransactionBody := services.TransactionBody{
		TransactionID: &services.TransactionID{
			AccountID: &services.AccountID{
				Account: &services.AccountID_AccountNum{AccountNum: 542348},
			},
			TransactionValidStart: &services.Timestamp{
				Seconds: id.ValidStart.Unix(),
				Nanos:   int32(id.ValidStart.Nanosecond()),
			},
		},
		NodeAccountID: &services.AccountID{
			Account: &services.AccountID_AccountNum{AccountNum: 3},
		},
		TransactionFee: 200_000_000,
		TransactionValidDuration: &services.Duration{
			Seconds: 120,
		},
		GenerateRecord: false,
		Memo:           "",
		Data: &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: &services.CryptoTransferTransactionBody{
				Transfers: &services.TransferList{
					AccountAmounts: []*services.AccountAmount{
						{
							AccountID: &services.AccountID{
								Account: &services.AccountID_AccountNum{AccountNum: 47439},
							},
							Amount: 10,
						},
						{
							AccountID: &services.AccountID{
								Account: &services.AccountID_AccountNum{AccountNum: 542348},
							},
							Amount: -10,
						},
					},
				},
			},
		},
	}

	BodyBytes, err := protobuf.Marshal(&TransactionBody)
	require.NoError(t, err)

	key1, _ := PrivateKeyFromString("302e020100300506032b6570042204203e7fda6dde63c3cdb3cb5ecf5264324c5faad7c9847b6db093c088838b35a110")
	key2, _ := PrivateKeyFromString("302e020100300506032b65700422042032d3d5a32e9d06776976b39c09a31fbda4a4a0208223da761c26a2ae560c1755")
	key3, _ := PrivateKeyFromString("302e020100300506032b657004220420195a919056d1d698f632c228dbf248bbbc3955adf8a80347032076832b8299f9")
	key4, _ := PrivateKeyFromString("302e020100300506032b657004220420b9962f17f94ffce73a23649718a11638cac4b47095a7a6520e88c7563865be62")
	key5, _ := PrivateKeyFromString("302e020100300506032b657004220420fef68591819080cd9d48b0cbaa10f65f919752abb50ffb3e7411ac66ab22692e")

	publicKey1 := key1.PublicKey()
	publicKey2 := key2.PublicKey()
	publicKey3 := key3.PublicKey()
	publicKey4 := key4.PublicKey()
	publicKey5 := key5.PublicKey()

	signature1 := key1.Sign(BodyBytes)
	signature2 := key2.Sign(BodyBytes)
	signature3 := key3.Sign(BodyBytes)
	signature4 := key4.Sign(BodyBytes)
	signature5 := key5.Sign(BodyBytes)

	signed := services.SignedTransaction{
		BodyBytes: BodyBytes,
		SigMap: &services.SignatureMap{
			SigPair: []*services.SignaturePair{
				{
					PubKeyPrefix: key1.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature1,
					},
				},
				{
					PubKeyPrefix: key2.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature2,
					},
				},
				{
					PubKeyPrefix: key3.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature3,
					},
				},
				{
					PubKeyPrefix: key4.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature4,
					},
				},
				{
					PubKeyPrefix: key5.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature5,
					},
				},
			},
		},
	}

	bytes, err := protobuf.Marshal(&signed)
	require.NoError(t, err)

	bytes, err = protobuf.Marshal(&sdk.TransactionList{
		TransactionList: []*services.Transaction{{
			SignedTransactionBytes: bytes,
		}},
	})
	require.NoError(t, err)

	transaction, err := TransactionFromBytes(bytes)
	require.NoError(t, err)

	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	switch tx := transaction.(type) {
	case *TransferTransaction:
		assert.Equal(t, tx.GetHbarTransfers()[AccountID{0, 0, 542348, nil, nil, nil}].AsTinybar(), int64(-10))
		assert.Equal(t, tx.GetHbarTransfers()[AccountID{0, 0, 47439, nil, nil, nil}].AsTinybar(), int64(10))

		signatures, err := tx.GetSignatures()
		require.NoError(t, err)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey1)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey2)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey3)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey4)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey5)

		assert.Equal(t, len(tx.GetNodeAccountIDs()), 1)
		assert.True(t, tx.GetNodeAccountIDs()[0]._Equals(AccountID{0, 0, 3, nil, nil, nil}))

		resp, err := tx.Execute(env.Client)
		require.NoError(t, err)

		_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
		require.NoError(t, err)
	default:
		panic("Transaction was not a crypto transfer?")
	}
}

func TestIntegrationTransactionFailsWhenSigningWithoutFreezing(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs)

	_, err = tx.Sign(newKey).Execute(env.Client)
	require.ErrorContains(t, err, "transaction is not frozen")

}

func TestIntegrationTransactionAddSignatureV2(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	signableBodyList, err := tx.GetSignableNodeBodyBytesList()
	require.NoError(t, err)
	for _, signableBody := range signableBodyList {
		signature := newKey.Sign(signableBody.Body)
		tx, err = tx.AddSignatureV2(newKey.PublicKey(), signature, signableBody.TransactionID, signableBody.NodeID)
		require.NoError(t, err)
	}

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTransactionStaticMethodsFileAppend(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	frozenFileCreateTx, err := NewFileCreateTransaction().
		SetKeys(newKey.PublicKey()).
		SetContents([]byte("Initial content")).
		FreezeWith(env.Client)
	require.NoError(t, err)

	fileCreateResp, err := frozenFileCreateTx.Sign(newKey).Execute(env.Client)
	require.NoError(t, err)

	fileReceipt, err := fileCreateResp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	fileID := *fileReceipt.FileID

	fileAppendTx := NewFileAppendTransaction().
		SetFileID(fileID).
		SetContents([]byte(bigContents))

	require.False(t, fileAppendTx.IsFrozen())

	frozenFileTxInterface, err := TransactionFreezeWith(fileAppendTx, env.Client)
	require.NoError(t, err)
	frozenFileTx, ok := frozenFileTxInterface.(*FileAppendTransaction)
	require.True(t, ok)

	require.True(t, frozenFileTx.IsFrozen())

	frozenFileTxInterface, err = TransactionSign(frozenFileTx, newKey)
	require.NoError(t, err)

	resp, err := TransactionExecute(frozenFileTxInterface, env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

func TestIntegrationTransactionStaticMethodsTopicMessageSubmit(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	frozenTopicCreateTx, err := NewTopicCreateTransaction().
		SetAdminKey(adminKey.PublicKey()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	resp, err := frozenTopicCreateTx.Sign(adminKey).Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	topicID := *receipt.TopicID
	assert.NotNil(t, topicID)

	topicMessageSubmitTx := NewTopicMessageSubmitTransaction().
		SetMessage([]byte(bigContents)).
		SetTopicID(topicID)

	frozenTopicMessageSubmitTx, err := TransactionFreezeWith(topicMessageSubmitTx, env.Client)
	require.NoError(t, err)
	topicMessageSubmitTx, ok := frozenTopicMessageSubmitTx.(*TopicMessageSubmitTransaction)
	require.True(t, ok)

	topicMessageSubmitTx.Sign(adminKey)

	resp, err = TransactionExecute(topicMessageSubmitTx, env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

// _SignaturesContainKey reports whether GetSignatures records a signature for publicKey against
// any node (GetSignatures keys the inner map by *PublicKey, so we compare by String()).
func _SignaturesContainKey(signatures map[AccountID]map[*PublicKey][]byte, publicKey PublicKey) bool {
	for _, nodeSigs := range signatures {
		for key := range nodeSigs {
			if key.String() == publicKey.String() {
				return true
			}
		}
	}
	return false
}

func TestIntegrationTransactionRemoveSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	// Sign with the account key, then confirm the signature is present.
	_, err = newKey.SignTransaction(tx)
	require.NoError(t, err)

	signatures, err := tx.GetSignatures()
	require.NoError(t, err)
	require.True(t, _SignaturesContainKey(signatures, newKey.PublicKey()))

	// Remove it: the returned bytes are the signatures that were stripped off the wire.
	removed, err := tx.RemoveSignature(newKey.PublicKey())
	require.NoError(t, err)
	require.NotEmpty(t, removed)

	signatures, err = tx.GetSignatures()
	require.NoError(t, err)
	require.False(t, _SignaturesContainKey(signatures, newKey.PublicKey()))

	// Removing again errors, since the key is no longer on the transaction.
	_, err = tx.RemoveSignature(newKey.PublicKey())
	require.ErrorIs(t, err, errPublicKeyHasNotSigned)

	// Re-sign with the account key and execute: the delete requires the account's signature, so a
	// successful receipt proves the signature was genuinely re-applied after removal.
	_, err = newKey.SignTransaction(tx)
	require.NoError(t, err)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}

// TestIntegrationTransactionRemoveSignatureIsAppliedOnExecute proves the removal actually reaches
// the wire: after removing the account key we execute WITHOUT re-signing, and the network rejects
// the delete with INVALID_SIGNATURE. This is stronger than the GetSignatures check (which only
// reads the in-memory sig map) because it confirms the removed signature is genuinely absent from
// the bytes transmitted at build/execute time and is not resurrected.
func TestIntegrationTransactionRemoveSignatureIsAppliedOnExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	// Sign with the account key, confirm it is present, then remove it.
	_, err = newKey.SignTransaction(tx)
	require.NoError(t, err)

	signatures, err := tx.GetSignatures()
	require.NoError(t, err)
	require.True(t, _SignaturesContainKey(signatures, newKey.PublicKey()))

	removed, err := tx.RemoveSignature(newKey.PublicKey())
	require.NoError(t, err)
	require.NotEmpty(t, removed)

	// Execute WITHOUT re-signing. The operator (payer) signature is re-applied automatically, but the
	// account's required signature was removed, so the delete must be rejected with INVALID_SIGNATURE.
	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTransactionRemoveAllSignatures(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKeyWithoutAlias(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	// Sign with the account key plus a second, unrelated key so more than one signature is present on
	// the wire. Both are applied via SignTransaction, which materializes the signature into the sig
	// map immediately (unlike the operator signer, which is deferred until build/execute time and so
	// would not yet appear in GetSignatures/RemoveAllSignatures).
	secondKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	_, err = newKey.SignTransaction(tx)
	require.NoError(t, err)
	_, err = secondKey.SignTransaction(tx)
	require.NoError(t, err)

	// Both keys must be present before removal.
	signatures, err := tx.GetSignatures()
	require.NoError(t, err)
	require.True(t, _SignaturesContainKey(signatures, newKey.PublicKey()))
	require.True(t, _SignaturesContainKey(signatures, secondKey.PublicKey()))

	// Strip everything: the returned map is keyed by every public key that had signed, so both keys
	// must appear in it.
	removed, err := tx.RemoveAllSignatures()
	require.NoError(t, err)
	require.Len(t, removed, 2)

	var newKeyFound, secondKeyFound bool
	for key := range removed {
		switch key.String() {
		case newKey.PublicKey().String():
			newKeyFound = true
		case secondKey.PublicKey().String():
			secondKeyFound = true
		}
	}
	require.True(t, newKeyFound)
	require.True(t, secondKeyFound)

	// Every node's signature set is now empty.
	signatures, err = tx.GetSignatures()
	require.NoError(t, err)
	for _, nodeSigs := range signatures {
		require.Empty(t, nodeSigs)
	}

	// Re-sign with the account key and execute. The operator signature is re-applied automatically
	// at execution time, so re-adding the account key is enough for the delete to succeed.
	_, err = newKey.SignTransaction(tx)
	require.NoError(t, err)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
}
