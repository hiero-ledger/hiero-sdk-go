//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTransactionSerializationDeserialization(t *testing.T) {
	t.Parallel()

	transaction, err := _NewMockTransaction()
	require.NoError(t, err)

	_, err = transaction.Freeze()
	require.NoError(t, err)

	_, err = transaction.GetSignatures()
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.
		SetTransactionMemo("memo").
		SetMaxTransactionFee(NewHbar(5))

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	deserializedTX, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	var deserializedTXTyped TransferTransaction
	switch tx := deserializedTX.(type) {
	case TransferTransaction:
		deserializedTXTyped = tx
	default:
		panic("Transaction was not TransferTransaction")
	}

	require.Equal(t, "memo", deserializedTXTyped.memo)
	require.Equal(t, NewHbar(5), deserializedTXTyped.GetMaxTransactionFee())
	assert.Equal(t, transaction.String(), deserializedTXTyped.String())
}

func TestUnitTransactionValidateBodiesEqual(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyFromString(mockPrivateKey)
	require.NoError(t, err)
	transaction := services.TransactionBody{
		TransactionID:            testTransactionID._ToProtobuf(),
		NodeAccountID:            AccountID{Account: 3}._ToProtobuf(),
		TransactionFee:           0,
		TransactionValidDuration: nil,
		GenerateRecord:           false,
		Memo:                     "",
		Data: &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: &services.CryptoCreateTransactionBody{
				Key:                           key._ToProtoKey(),
				InitialBalance:                0,
				ProxyAccountID:                AccountID{Account: 123}._ToProtobuf(),
				SendRecordThreshold:           0,
				ReceiveRecordThreshold:        0,
				ReceiverSigRequired:           false,
				AutoRenewPeriod:               nil,
				ShardID:                       nil,
				RealmID:                       nil,
				NewRealmAdminKey:              nil,
				Memo:                          "",
				MaxAutomaticTokenAssociations: 0,
			},
		},
	}

	transactionBody, err := protobuf.Marshal(&transaction)
	require.NoError(t, err)

	signed, err := protobuf.Marshal(&services.SignedTransaction{
		BodyBytes: transactionBody,
	})
	require.NoError(t, err)
	list, err := protobuf.Marshal(&sdk.TransactionList{
		TransactionList: []*services.Transaction{
			{
				SignedTransactionBytes: signed,
			},
			{
				SignedTransactionBytes: signed,
			},
			{
				SignedTransactionBytes: signed,
			},
		},
	})

	deserializedTX, err := TransactionFromBytes(list)
	require.NoError(t, err)

	var deserializedTXTyped AccountCreateTransaction
	switch tx := deserializedTX.(type) {
	case AccountCreateTransaction:
		deserializedTXTyped = tx
	default:
		panic("Transaction was not AccountCreateTransaction")
	}

	assert.Equal(t, uint64(transaction.TransactionID.AccountID.GetAccountNum()), deserializedTXTyped.GetTransactionID().AccountID.Account)
}

func DisabledTestUnitTransactionValidateBodiesNotEqual(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyFromString(mockPrivateKey)
	require.NoError(t, err)
	transaction := services.TransactionBody{
		TransactionID:            testTransactionID._ToProtobuf(),
		NodeAccountID:            AccountID{Account: 3}._ToProtobuf(),
		TransactionFee:           0,
		TransactionValidDuration: nil,
		GenerateRecord:           false,
		Memo:                     "",
		Data: &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: &services.CryptoCreateTransactionBody{
				Key:                           key._ToProtoKey(),
				InitialBalance:                0,
				ProxyAccountID:                AccountID{Account: 123}._ToProtobuf(),
				SendRecordThreshold:           0,
				ReceiveRecordThreshold:        0,
				ReceiverSigRequired:           false,
				AutoRenewPeriod:               nil,
				ShardID:                       nil,
				RealmID:                       nil,
				NewRealmAdminKey:              nil,
				Memo:                          "",
				MaxAutomaticTokenAssociations: 0,
			},
		},
	}

	transaction2 := services.TransactionBody{
		TransactionID:            testTransactionID._ToProtobuf(),
		NodeAccountID:            AccountID{Account: 3}._ToProtobuf(),
		TransactionFee:           0,
		TransactionValidDuration: nil,
		GenerateRecord:           false,
		Memo:                     "",
		Data: &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: &services.CryptoCreateTransactionBody{
				Key:                           key._ToProtoKey(),
				InitialBalance:                0,
				ProxyAccountID:                AccountID{Account: 1}._ToProtobuf(),
				SendRecordThreshold:           0,
				ReceiveRecordThreshold:        0,
				ReceiverSigRequired:           false,
				AutoRenewPeriod:               nil,
				ShardID:                       nil,
				RealmID:                       nil,
				NewRealmAdminKey:              nil,
				Memo:                          "",
				MaxAutomaticTokenAssociations: 0,
			},
		},
	}

	transactionBody, err := protobuf.Marshal(&transaction)
	require.NoError(t, err)

	signed, err := protobuf.Marshal(&services.SignedTransaction{
		BodyBytes: transactionBody,
	})

	transactionBody2, err := protobuf.Marshal(&transaction2)
	require.NoError(t, err)

	signed2, err := protobuf.Marshal(&services.SignedTransaction{
		BodyBytes: transactionBody2,
	})

	require.NoError(t, err)
	list, err := protobuf.Marshal(&sdk.TransactionList{
		TransactionList: []*services.Transaction{
			{
				SignedTransactionBytes: signed,
			},
			{
				SignedTransactionBytes: signed2,
			},
			{
				SignedTransactionBytes: signed2,
			},
		},
	})

	_, err = TransactionFromBytes(list)
	require.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("failed to validate transaction bodies"), err.Error())
	}
}

func TestUnitTransactionToFromBytes(t *testing.T) {
	t.Parallel()

	duration := time.Second * 10
	operatorID := AccountID{Account: 5}
	recepientID := AccountID{Account: 4}
	node := []AccountID{{Account: 3}}
	transaction, err := NewTransferTransaction().
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs(node).
		AddHbarTransfer(operatorID, NewHbar(-1)).
		AddHbarTransfer(recepientID, NewHbar(1)).
		SetTransactionMemo("go sdk example multi_app_transfer/main.go").
		SetTransactionValidDuration(duration).
		Freeze()
	require.NoError(t, err)

	_ = transaction.GetTransactionID()
	nodeID := transaction.GetNodeAccountIDs()
	require.NotEmpty(t, nodeID)
	require.False(t, nodeID[0]._IsZero())

	var tx services.TransactionBody
	_ = protobuf.Unmarshal(transaction.signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.Equal(t, tx.TransactionID.String(), testTransactionID._ToProtobuf().String())
	require.Equal(t, tx.NodeAccountID.String(), node[0]._ToProtobuf().String())
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, duration, _DurationFromProtobuf(tx.TransactionValidDuration))
	require.Equal(t, tx.GetCryptoTransfer().Transfers.AccountAmounts,
		[]*services.AccountAmount{
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
				Amount:    -100000000,
			},
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
				Amount:    100000000,
			},
		})

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	newTransaction, err := TransactionFromBytes(txBytes)

	_ = protobuf.Unmarshal(newTransaction.(TransferTransaction).signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.Equal(t, tx.TransactionID.String(), testTransactionID._ToProtobuf().String())
	require.Equal(t, tx.NodeAccountID.String(), node[0]._ToProtobuf().String())
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, duration, _DurationFromProtobuf(tx.TransactionValidDuration))
	require.Equal(t, tx.GetCryptoTransfer().Transfers.AccountAmounts,
		[]*services.AccountAmount{
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
				Amount:    -100000000,
			},
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
				Amount:    100000000,
			},
		})
}

func TestUnitTransactionToFromBytesWithClient(t *testing.T) {
	t.Parallel()

	duration := time.Second * 10
	operatorID := AccountID{Account: 5}
	recepientID := AccountID{Account: 4}
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	client.SetOperator(AccountID{Account: 2}, privateKey)

	transaction, err := NewTransferTransaction().
		AddHbarTransfer(operatorID, NewHbar(-1)).
		AddHbarTransfer(recepientID, NewHbar(1)).
		SetTransactionMemo("go sdk example multi_app_transfer/main.go").
		SetTransactionValidDuration(duration).
		FreezeWith(client)
	require.NoError(t, err)

	var tx services.TransactionBody
	_ = protobuf.Unmarshal(transaction.signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.NotNil(t, tx.TransactionID, tx.NodeAccountID)
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, duration, _DurationFromProtobuf(tx.TransactionValidDuration))
	require.Equal(t, tx.GetCryptoTransfer().Transfers.AccountAmounts,
		[]*services.AccountAmount{
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
				Amount:    -100000000,
			},
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
				Amount:    100000000,
			},
		})

	initialTxID := tx.TransactionID
	initialNode := tx.NodeAccountID

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	newTransaction, err := TransactionFromBytes(txBytes)

	_ = protobuf.Unmarshal(newTransaction.(TransferTransaction).signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.NotNil(t, tx.TransactionID, tx.NodeAccountID)
	require.Equal(t, tx.TransactionID.String(), initialTxID.String())
	require.Equal(t, tx.NodeAccountID.String(), initialNode.String())
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, duration, _DurationFromProtobuf(tx.TransactionValidDuration))
	require.Equal(t, tx.GetCryptoTransfer().Transfers.AccountAmounts,
		[]*services.AccountAmount{
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
				Amount:    -100000000,
			},
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
				Amount:    100000000,
			},
		})
}

func TestUnitQueryRegression(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 5}
	node := []AccountID{{Account: 3}}
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	client.SetOperator(AccountID{Account: 2}, privateKey)

	query := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs(node).
		SetPaymentTransactionID(testTransactionID).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25))

	body := query.buildQuery()
	_, err = query.generatePayments(client, HbarFromTinybar(20))
	require.NoError(t, err)

	var paymentTx services.TransactionBody
	_ = protobuf.Unmarshal(query.Query.paymentTransactions[0].BodyBytes, &paymentTx)

	require.Equal(t, body.GetCryptoGetInfo().AccountID.String(), accountID._ToProtobuf().String())
	require.Equal(t, paymentTx.NodeAccountID.String(), node[0]._ToProtobuf().String())
	require.Equal(t, paymentTx.TransactionFee, uint64(NewHbar(1).tinybar))
	require.Equal(t, paymentTx.TransactionValidDuration, &services.Duration{Seconds: 120})
	require.Equal(t, paymentTx.Data, &services.TransactionBody_CryptoTransfer{
		CryptoTransfer: &services.CryptoTransferTransactionBody{
			Transfers: &services.TransferList{
				AccountAmounts: []*services.AccountAmount{
					{
						AccountID: node[0]._ToProtobuf(),
						Amount:    HbarFromTinybar(20).AsTinybar(),
					},
					{
						AccountID: client.GetOperatorAccountID()._ToProtobuf(),
						Amount:    -HbarFromTinybar(20).AsTinybar(),
					},
				},
			},
		},
	})
}
func TestUnitTransactionInitFeeMaxTransactionWithouthSettingFee(t *testing.T) {
	t.Parallel()

	//Default Max Fee for TransferTransaction
	fee := NewHbar(1)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	transaction, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		FreezeWith(client)
	require.NoError(t, err)
	require.Equal(t, uint64(fee.AsTinybar()), transaction.transactionFee)
}

func TestUnitTransactionInitFeeMaxTransactionFeeSetExplicitly(t *testing.T) {
	t.Parallel()

	clientMaxFee := NewHbar(14)
	explicitMaxFee := NewHbar(15)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetDefaultMaxTransactionFee(clientMaxFee)
	require.NoError(t, err)
	transaction, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		SetMaxTransactionFee(explicitMaxFee).
		FreezeWith(client)
	require.NoError(t, err)
	require.Equal(t, uint64(explicitMaxFee.AsTinybar()), transaction.transactionFee)
}

func TestUnitTransactionInitFeeMaxTransactionFromClientDefault(t *testing.T) {
	t.Parallel()

	fee := NewHbar(14)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetDefaultMaxTransactionFee(fee)
	require.NoError(t, err)
	transaction, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		FreezeWith(client)
	require.NoError(t, err)
	require.Equal(t, uint64(fee.AsTinybar()), transaction.transactionFee)
}

func TestUnitTransactionSignSwitchCases(t *testing.T) {
	t.Parallel()

	newKey, client, nodeAccountId := signSwitchCaseaSetup(t)

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {

		txVal, signature, transferTxBytes := signSwitchCaseaHelper(t, tx, newKey, client)

		signTests := signTestsForTransaction(txVal, newKey, signature, client)

		for _, tt := range signTests {
			t.Run(tt.name, func(t *testing.T) {
				transactionInterface, err := TransactionFromBytes(transferTxBytes)
				require.NoError(t, err)

				tx, err := tt.sign(transactionInterface, newKey)
				assert.NoError(t, err)
				assert.NotEmpty(t, tx)

				signs, err := TransactionGetSignatures(transactionInterface)
				assert.NoError(t, err)

				// verify with range because var signs = map[AccountID]map[*PublicKey][]byte, where *PublicKey is unknown memory address
				for key := range signs[nodeAccountId] {
					assert.Equal(t, signs[nodeAccountId][key], signature)
				}
			})
		}
	}
}

func TestUnitTransactionSignSwitchCasesPointers(t *testing.T) {
	t.Parallel()

	newKey, client, nodeAccountId := signSwitchCaseaSetup(t)

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {

		txVal, signature, transferTxBytes := signSwitchCaseaHelper(t, tx, newKey, client)
		signTests := signTestsForTransaction(txVal, newKey, signature, client)

		for _, tt := range signTests {
			t.Run(tt.name, func(t *testing.T) {
				transactionInterface, err := TransactionFromBytes(transferTxBytes)
				require.NoError(t, err)

				signs, err := TransactionGetSignatures(transactionInterface)
				assert.NoError(t, err)

				// verify with range because var signs = map[AccountID]map[*PublicKey][]byte, where *PublicKey is unknown memory address
				for key := range signs[nodeAccountId] {
					assert.Equal(t, signs[nodeAccountId][key], signature)
				}
			})
		}
	}
}

func TestUnitTransactionAttributes(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	nodeAccountIds := client.network._GetNodeAccountIDsForExecute()

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewContractCreateTransaction(),
		NewContractDeleteTransaction(),
		NewContractExecuteTransaction(),
		NewContractUpdateTransaction(),
		NewFileAppendTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewScheduleCreateTransaction(),
		NewScheduleDeleteTransaction(),
		NewScheduleSignTransaction(),
		NewSystemDeleteTransaction(),
		NewSystemUndeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {
		txName := reflect.TypeOf(tx).Elem().Name()

		tests := createTransactionTests(txName, nodeAccountIds)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				txSet, err := tt.set(tx)
				require.NoError(t, err)

				txGet, err := tt.get(txSet)
				require.NoError(t, err)

				tt.assert(t, txGet)
			})
		}
	}
}

func TestUnitTransactionAttributesDereferanced(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	nodeAccountIds := client.network._GetNodeAccountIDsForExecute()

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewContractCreateTransaction(),
		NewContractDeleteTransaction(),
		NewContractExecuteTransaction(),
		NewContractUpdateTransaction(),
		NewFileAppendTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewScheduleCreateTransaction(),
		NewScheduleDeleteTransaction(),
		NewScheduleSignTransaction(),
		NewSystemDeleteTransaction(),
		NewSystemUndeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {
		txName := reflect.TypeOf(tx).Elem().Name()

		tests := createTransactionTests(txName, nodeAccountIds)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				txSet, err := tt.set(tx)
				require.NoError(t, err)

				txGet, err := tt.get(txSet)
				require.NoError(t, err)

				tt.assert(t, txGet)
			})
		}
	}
}

func TestUnitTransactionAttributesSerialization(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewContractCreateTransaction(),
		NewContractDeleteTransaction(),
		NewContractExecuteTransaction(),
		NewContractUpdateTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewScheduleCreateTransaction(),
		NewScheduleDeleteTransaction(),
		NewScheduleSignTransaction(),
		NewSystemDeleteTransaction(),
		NewSystemUndeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {
		txName := reflect.TypeOf(tx).Elem().Name()

		// Get the reflect.Value of the pointer to the Transaction
		txPtr := reflect.ValueOf(tx)
		txPtr.MethodByName("FreezeWith").Call([]reflect.Value{reflect.ValueOf(client)})

		tests := []struct {
			name string
			act  func(transactionInterface TransactionInterface)
		}{
			{
				name: "TransactionString/" + txName,
				act: func(transactionInterface TransactionInterface) {
					txString, err := TransactionString(transactionInterface)
					require.NoError(t, err)
					require.NotEmpty(t, txString)
				},
			},
			{
				name: "TransactionToBytes/" + txName,
				act: func(transactionInterface TransactionInterface) {
					txBytes, err := TransactionToBytes(transactionInterface)
					require.NoError(t, err)
					require.NotEmpty(t, txBytes)
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.act(tx)
				// txValue := reflect.ValueOf(tx).Elem().Interface()
				// tt.act(txValue)
			})
		}
	}
}

func TestUnitAddSignatureV2(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	fileID := FileID{File: 3}
	nodeAccountID1 := AccountID{Account: 3}
	nodeAccountID2 := AccountID{Account: 4}
	nodeAccountIDs := []AccountID{nodeAccountID1, nodeAccountID2}

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	require.NoError(t, err)

	// Mock signature bytes
	mockSignature := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Test Case 1: Single Node, Single Chunk
	t.Run("Single Node Single Chunk", func(t *testing.T) {
		transaction := NewFileAppendTransaction().
			SetFileID(fileID).
			SetContents([]byte("test content")).
			SetNodeAccountIDs([]AccountID{nodeAccountID1}).
			SetTransactionID(testTransactionID)

		transaction.SetMaxChunks(1)
		transaction.SetMaxChunkSize(2048)

		frozen, err := transaction.FreezeWith(client)
		require.NoError(t, err)

		frozen, err = frozen.AddSignatureV2(privateKey.PublicKey(), mockSignature, testTransactionID, nodeAccountID1)
		require.NoError(t, err)

		signs, err := frozen.GetSignatures()
		require.NoError(t, err)
		require.Len(t, signs, 1)
		require.Contains(t, signs, nodeAccountID1)

		// Verify the signature bytes
		for key := range signs[nodeAccountID1] {
			require.Equal(t, mockSignature, signs[nodeAccountID1][key])
		}
	})

	// Test Case 2: Multiple Nodes, Single Chunk
	t.Run("Multiple Nodes Single Chunk", func(t *testing.T) {
		transaction := NewFileAppendTransaction().
			SetFileID(fileID).
			SetContents([]byte("test content")).
			SetNodeAccountIDs(nodeAccountIDs).
			SetTransactionID(testTransactionID)

		transaction.SetMaxChunks(1)
		transaction.SetMaxChunkSize(2048)

		frozen, err := transaction.FreezeWith(client)
		require.NoError(t, err)

		// Add signatures for both nodes
		frozen, err = frozen.AddSignatureV2(privateKey.PublicKey(), mockSignature, testTransactionID, nodeAccountID1)
		require.NoError(t, err)
		frozen, err = frozen.AddSignatureV2(privateKey.PublicKey(), mockSignature, testTransactionID, nodeAccountID2)
		require.NoError(t, err)

		signs, err := frozen.GetSignatures()
		require.NoError(t, err)
		require.Len(t, signs, 2)
		require.Contains(t, signs, nodeAccountID1)
		require.Contains(t, signs, nodeAccountID2)

		// Verify signatures for both nodes
		for _, nodeID := range nodeAccountIDs {
			for key := range signs[nodeID] {
				require.Equal(t, mockSignature, signs[nodeID][key])
			}
		}
	})

	// Test Case 3: Multiple Nodes, Multiple Chunks
	t.Run("Multiple Nodes Multiple Chunks", func(t *testing.T) {
		content := make([]byte, 3000)
		for i := range content {
			content[i] = byte(i % 256)
		}

		transaction := NewFileAppendTransaction().
			SetFileID(fileID).
			SetContents(content).
			SetNodeAccountIDs(nodeAccountIDs).
			SetTransactionID(testTransactionID)

		transaction.SetMaxChunks(2)
		transaction.SetMaxChunkSize(2048)

		frozen, err := transaction.FreezeWith(client)
		require.NoError(t, err)

		// Add signatures for both nodes
		frozen, err = frozen.AddSignatureV2(privateKey.PublicKey(), mockSignature, testTransactionID, nodeAccountID1)
		require.NoError(t, err)
		frozen, err = frozen.AddSignatureV2(privateKey.PublicKey(), mockSignature, testTransactionID, nodeAccountID2)
		require.NoError(t, err)

		signs, err := frozen.GetSignatures()
		require.NoError(t, err)
		require.Len(t, signs, 2)
		require.Contains(t, signs, nodeAccountID1)
		require.Contains(t, signs, nodeAccountID2)

		// Verify each node has the correct signature
		for _, nodeID := range nodeAccountIDs {
			nodeSigs := signs[nodeID]
			require.Len(t, nodeSigs, 1)
			for key := range nodeSigs {
				require.NotNil(t, key)
				require.Equal(t, key.String(), privateKey.PublicKey().String())
				require.Equal(t, mockSignature, nodeSigs[key])
			}
		}
	})

	// Test Case 4: Wrong Node ID
	t.Run("Wrong Node ID", func(t *testing.T) {
		transaction := NewFileAppendTransaction().
			SetFileID(fileID).
			SetContents([]byte("test content")).
			SetNodeAccountIDs(nodeAccountIDs).
			SetTransactionID(testTransactionID)

		transaction.SetMaxChunks(1)
		transaction.SetMaxChunkSize(2048)

		frozen, err := transaction.FreezeWith(client)
		require.NoError(t, err)

		invalidNodeID := AccountID{Account: 999}
		frozen, err = frozen.AddSignatureV2(privateKey.PublicKey(), mockSignature, testTransactionID, invalidNodeID)
		require.NoError(t, err)

		signs, err := frozen.GetSignatures()
		require.NoError(t, err)
		require.NotContains(t, signs, invalidNodeID)
	})

	// Test Case 5: Wrong Transaction ID
	t.Run("Wrong Transaction ID", func(t *testing.T) {
		transaction := NewFileAppendTransaction().
			SetFileID(fileID).
			SetContents([]byte("test content")).
			SetNodeAccountIDs(nodeAccountIDs).
			SetTransactionID(testTransactionID)

		transaction.SetMaxChunks(1)
		transaction.SetMaxChunkSize(2048)

		frozen, err := transaction.FreezeWith(client)
		require.NoError(t, err)

		invalidTxID := TransactionID{
			AccountID:  &AccountID{Account: 999},
			ValidStart: &time.Time{},
		}

		frozen, err = frozen.AddSignatureV2(privateKey.PublicKey(), mockSignature, invalidTxID, nodeAccountID1)
		require.NoError(t, err)

		signs, err := frozen.GetSignatures()
		require.NoError(t, err)
		require.NotContains(t, signs[nodeAccountID1], privateKey.PublicKey())
	})

	// Test Case 6: Adding Same Signature Twice
	t.Run("Adding Same Signature Twice", func(t *testing.T) {
		transaction := NewFileAppendTransaction().
			SetFileID(fileID).
			SetContents([]byte("test content")).
			SetNodeAccountIDs([]AccountID{nodeAccountID1}).
			SetTransactionID(testTransactionID)

		transaction.SetMaxChunks(1)
		transaction.SetMaxChunkSize(2048)

		frozen, err := transaction.FreezeWith(client)
		require.NoError(t, err)

		// Add signature first time
		frozen, err = frozen.AddSignatureV2(privateKey.PublicKey(), mockSignature, testTransactionID, nodeAccountID1)
		require.NoError(t, err)

		// Add same signature second time
		frozen, err = frozen.AddSignatureV2(privateKey.PublicKey(), mockSignature, testTransactionID, nodeAccountID1)
		require.NoError(t, err)

		signs, err := frozen.GetSignatures()
		require.NoError(t, err)
		require.Len(t, signs, 1)
		require.Contains(t, signs, nodeAccountID1)

		// Verify there's only one signature
		nodeSigs := signs[nodeAccountID1]
		require.Len(t, nodeSigs, 1)
		for key := range nodeSigs {
			require.Equal(t, mockSignature, nodeSigs[key])
		}
	})
}

func TestUnitAddSignatureV2WithEmptySignedTransactions(t *testing.T) {
	t.Parallel()

	// Create a new transaction
	tx := NewFileAppendTransaction()

	// Generate a test key
	key, err := GeneratePrivateKey()
	require.NoError(t, err)

	// Create mock signature
	mockSignature := []byte{0, 1, 2, 3, 4}

	// Create test node and transaction IDs
	nodeID := AccountID{Account: 3}
	testTxID := TransactionID{
		AccountID:  &AccountID{Account: 5},
		ValidStart: &time.Time{},
	}

	// Add signature when signedTransactions is empty
	tx.AddSignatureV2(key.PublicKey(), mockSignature, testTxID, nodeID)

	// Verify no signatures were added
	signs, err := tx.GetSignatures()
	require.NoError(t, err)
	require.Empty(t, signs)
}

func TestUnitGetSignableNodeBodyBytesListUnfrozen(t *testing.T) {
	t.Parallel()

	// Test unfrozen transaction
	tx := NewTransferTransaction()
	list, err := tx.GetSignableNodeBodyBytesList()
	require.Error(t, err)
	require.Equal(t, errTransactionIsNotFrozen, err)
	require.Empty(t, list)
}

func TestUnitGetSignableNodeBodyBytesListBasic(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	nodeID := AccountID{Account: 3}
	tx := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetTransactionID(testTransactionID).
		AddHbarTransfer(AccountID{Account: 2}, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1))

	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	list, err := frozen.GetSignableNodeBodyBytesList()
	require.NoError(t, err)
	require.NotEmpty(t, list)
	require.Len(t, list, 1) // Should have one entry for our single node

	// Verify the basic contents of the list
	require.Equal(t, nodeID, list[0].NodeID)
	require.Equal(t, testTransactionID, list[0].TransactionID)
	require.NotEmpty(t, list[0].Body)
}

func TestUnitGetSignableNodeBodyBytesListContents(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	nodeID := AccountID{Account: 3}
	tx := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetTransactionID(testTransactionID).
		AddHbarTransfer(AccountID{Account: 2}, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1))

	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	list, err := frozen.GetSignableNodeBodyBytesList()
	require.NoError(t, err)

	// Verify the body bytes can be unmarshalled into a valid transaction body
	var body services.TransactionBody
	err = protobuf.Unmarshal(list[0].Body, &body)
	require.NoError(t, err)
	require.NotNil(t, body.GetCryptoTransfer())
	require.Equal(t, nodeID.String(), _AccountIDFromProtobuf(body.NodeAccountID).String())
	require.Equal(t, testTransactionID.String(), _TransactionIDFromProtobuf(body.TransactionID).String())
}

func TestUnitGetSignableNodeBodyBytesListMultipleNodeIDs(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	nodeID1 := AccountID{Account: 3}
	nodeID2 := AccountID{Account: 4}
	nodeIDs := []AccountID{nodeID1, nodeID2}

	tx := NewTransferTransaction().
		SetNodeAccountIDs(nodeIDs).
		SetTransactionID(testTransactionID).
		AddHbarTransfer(AccountID{Account: 2}, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1))

	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	list, err := frozen.GetSignableNodeBodyBytesList()
	require.NoError(t, err)
	require.Len(t, list, 2) // Should have two entries, one per node

	// Verify each node's entry
	for i, nodeID := range nodeIDs {
		require.Equal(t, nodeID, list[i].NodeID)
		require.Equal(t, testTransactionID, list[i].TransactionID)
		require.NotEmpty(t, list[i].Body)

		// Verify body contents
		var body services.TransactionBody
		err = protobuf.Unmarshal(list[i].Body, &body)
		require.NoError(t, err)
		require.NotNil(t, body.GetCryptoTransfer())
		require.Equal(t, nodeID.String(), _AccountIDFromProtobuf(body.NodeAccountID).String())
	}
}

func TestUnitGetSignableNodeBodyBytesListFileAppendMultipleChunks(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	nodeID1 := AccountID{Account: 3}
	nodeID2 := AccountID{Account: 4}
	nodeIDs := []AccountID{nodeID1, nodeID2}

	// Create content larger than chunk size to force multiple chunks
	content := make([]byte, 4096)
	for i := range content {
		content[i] = byte(i % 256)
	}

	tx := NewFileAppendTransaction().
		SetNodeAccountIDs(nodeIDs).
		SetTransactionID(testTransactionID).
		SetFileID(FileID{File: 5}).
		SetContents(content)

	// Set small chunk size to force multiple chunks
	tx.SetMaxChunkSize(2048)

	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	list, err := frozen.GetSignableNodeBodyBytesList()
	require.NoError(t, err)
	require.Len(t, list, 4) // Should have 4 entries: 2 nodes * 2 chunks

	// Map to track transaction IDs per node
	txIDsByNode := make(map[string]map[string]bool)
	for _, nodeID := range nodeIDs {
		txIDsByNode[nodeID.String()] = make(map[string]bool)
	}

	// Verify each entry
	for i := 0; i < len(list); i++ {
		require.Contains(t, nodeIDs, list[i].NodeID)
		require.NotEmpty(t, list[i].TransactionID)
		require.NotEmpty(t, list[i].Body)

		nodeIDStr := list[i].NodeID.String()
		txIDStr := list[i].TransactionID.String()

		// Each transaction ID should appear exactly once per node
		require.False(t, txIDsByNode[nodeIDStr][txIDStr], "Duplicate transaction ID found for the same node")
		txIDsByNode[nodeIDStr][txIDStr] = true

		// Verify body contents
		var body services.TransactionBody
		err = protobuf.Unmarshal(list[i].Body, &body)
		require.NoError(t, err)
		require.NotNil(t, body.GetFileAppend())
		require.Equal(t, list[i].NodeID.String(), _AccountIDFromProtobuf(body.NodeAccountID).String())
	}

	// Verify each node has the same number of unique transaction IDs
	for _, nodeID := range nodeIDs {
		require.Len(t, txIDsByNode[nodeID.String()], 2, "Each node should have exactly 2 unique transaction IDs")
	}

	// Verify that all nodes have the same set of transaction IDs
	firstNodeTxIDs := txIDsByNode[nodeID1.String()]
	for _, nodeID := range nodeIDs[1:] {
		nodeTxIDs := txIDsByNode[nodeID.String()]
		for txID := range firstNodeTxIDs {
			require.True(t, nodeTxIDs[txID], "All nodes should have the same set of transaction IDs")
		}
	}
}

func TestUnitGetTransactionSizeComparison(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	tx := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(testTransactionID).
		SetKey(privateKeyECDSA.PublicKey())
	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	bodySize, err := frozen.GetTransactionBodySize()
	require.NoError(t, err)
	transactionSize, err := frozen.GetTransactionSize()
	require.NoError(t, err)
	require.Greater(t, transactionSize, bodySize)

	tx = NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		SetTransactionID(testTransactionID).
		SetKey(privateKeyECDSA.PublicKey())
	frozen, err = tx.FreezeWith(client)
	require.NoError(t, err)

	bodySizeMultiNode, err := frozen.GetTransactionBodySize()
	require.NoError(t, err)
	transactionSizeMultiNode, err := frozen.GetTransactionSize()
	require.NoError(t, err)
	require.Greater(t, transactionSizeMultiNode, bodySizeMultiNode)
	require.Equal(t, transactionSize, transactionSizeMultiNode)
	require.Equal(t, bodySizeMultiNode, bodySize)
}

func TestUnitGetTransactionSizeChunkedTransaction(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	content := make([]byte, 5000)
	for i := range content {
		content[i] = byte(i % 256)
	}

	tx := NewFileAppendTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(testTransactionID).
		SetFileID(FileID{File: 5}).
		SetContents(content).
		SetMaxChunkSize(1024)
	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	bodySize, err := frozen.GetTransactionBodySize()
	require.NoError(t, err)
	transactionSize, err := frozen.GetTransactionSize()
	require.NoError(t, err)
	require.Greater(t, transactionSize, bodySize)

	bodySizeList, err := frozen.GetTransactionBodySizeAllChunks()
	require.NoError(t, err)
	require.Len(t, bodySizeList, 5)

	require.Equal(t, bodySizeList[1], bodySizeList[0])
	require.Equal(t, bodySizeList[2], bodySizeList[0])
	require.Equal(t, bodySizeList[3], bodySizeList[0])
	require.Greater(t, bodySizeList[0], bodySizeList[4])
}

func TestUnitGetTransactionSizeUnfrozen(t *testing.T) {
	t.Parallel()

	tx := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(testTransactionID).
		SetKey(privateKeyECDSA.PublicKey())

	_, err := tx.GetTransactionSize()
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction is not froze")

	_, err = tx.GetTransactionBodySize()
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction is not froze")

	_, err = tx.GetTransactionBodySizeAllChunks()
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction is not froze")
}

func TestUnitGetTransactionSizeWithSignatures(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	tx := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(testTransactionID).
		SetKey(privateKeyECDSA.PublicKey())
	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	sizeBeforeSigning, err := frozen.GetTransactionSize()
	require.NoError(t, err)

	privateKeyECDSA.SignTransaction(frozen)
	sizeAfterSigning, err := frozen.GetTransactionSize()
	require.NoError(t, err)

	require.Greater(t, sizeAfterSigning, sizeBeforeSigning)
	bodySizeBeforeSigning, err := frozen.GetTransactionBodySize()
	require.NoError(t, err)
	bodySizeAfterSigning, err := frozen.GetTransactionBodySize()
	require.NoError(t, err)
	require.Equal(t, bodySizeBeforeSigning, bodySizeAfterSigning)
}

func TestUnitGetTransactionSizeLargeMemo(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	largeMemo := make([]byte, 100)
	for i := range largeMemo {
		largeMemo[i] = 'a'
	}

	tx := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(testTransactionID).
		AddHbarTransfer(AccountID{Account: 2}, NewHbar(1)).
		SetTransactionMemo(string(largeMemo))

	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	smallMemoTx := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(testTransactionID).
		AddHbarTransfer(AccountID{Account: 2}, NewHbar(1)).
		SetTransactionMemo("small memo")

	frozenSmall, err := smallMemoTx.FreezeWith(client)
	require.NoError(t, err)

	largeMemoSize, err := frozen.GetTransactionBodySize()
	require.NoError(t, err)
	smallMemoSize, err := frozenSmall.GetTransactionBodySize()
	require.NoError(t, err)

	require.Greater(t, largeMemoSize, smallMemoSize)
}

func TestUnitGetTransactionBodySizeAllChunksNonChunked(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	tx := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(testTransactionID).
		AddHbarTransfer(AccountID{Account: 2}, NewHbar(1))

	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	singleBodySize, err := frozen.GetTransactionBodySize()
	require.NoError(t, err)

	allBodySizes, err := frozen.GetTransactionBodySizeAllChunks()
	require.NoError(t, err)

	// For non-chunked transaction, should return array with single size
	require.Len(t, allBodySizes, 1)
	require.Equal(t, singleBodySize, allBodySizes[0])
}

func TestUnitGetTransactionBodySizeAllChunksTopicMessage(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	message := make([]byte, 2500)
	for i := range message {
		message[i] = byte(i % 256)
	}

	tx := NewTopicMessageSubmitTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(testTransactionID).
		SetTopicID(TopicID{Topic: 5}).
		SetMessage(message)
	frozen, err := tx.FreezeWith(client)
	require.NoError(t, err)

	allBodySizes, err := frozen.GetTransactionBodySizeAllChunks()
	require.NoError(t, err)

	require.Len(t, allBodySizes, 3)
	require.Equal(t, allBodySizes[0], allBodySizes[1])
	require.Greater(t, allBodySizes[0], allBodySizes[2])

	// Verify each chunk size is non-zero
	for _, size := range allBodySizes {
		require.Greater(t, size, int(0))
	}

	smallMessage := []byte("small message")
	txSmall := NewTopicMessageSubmitTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(testTransactionID).
		SetTopicID(TopicID{Topic: 5}).
		SetMessage(smallMessage)

	frozenSmall, err := txSmall.FreezeWith(client)
	require.NoError(t, err)

	allBodySizesSmall, err := frozenSmall.GetTransactionBodySizeAllChunks()
	require.NoError(t, err)

	require.Len(t, allBodySizesSmall, 1)
	require.Greater(t, allBodySizesSmall[0], int(0))
}

func signSwitchCaseaSetup(t *testing.T) (PrivateKey, *Client, AccountID) {
	newKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	nodeAccountIds := client.network._GetNodeAccountIDsForExecute()
	nodeAccountId := nodeAccountIds[0]

	return newKey, client, nodeAccountId
}

func signSwitchCaseaHelper(t *testing.T, tx TransactionInterface, newKey PrivateKey, client *Client) (txVal reflect.Value, signature []byte, transferTxBytes []byte) {
	// Get the reflect.Value of the pointer to the transaction
	txPtr := reflect.ValueOf(tx)
	txPtr.MethodByName("FreezeWith").Call([]reflect.Value{reflect.ValueOf(client)})

	// Get the reflect.Value of the transaction
	txVal = txPtr.Elem()

	// Get the transaction field by name
	// txField := txVal.FieldByName("Transaction")

	// Get the value of the Transaction field
	// txValue := txField.Interface().(Transaction[TransactionInterface])

	// refl_signature := reflect.ValueOf(newKey).MethodByName("SignTransaction").Call([]reflect.Value{reflect.ValueOf(&txValue)})
	signature, err := newKey.SignTransaction(tx)
	assert.NoError(t, err)

	transferTxBytes, err = TransactionToBytes(tx)
	assert.NoError(t, err)
	assert.NotEmpty(t, transferTxBytes)

	return txVal, signature, transferTxBytes
}

func signTestsForTransaction(txVal reflect.Value, newKey PrivateKey, signature []byte, client *Client) []struct {
	name string
	sign func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error)
} {
	return []struct {
		name string
		sign func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error)
	}{
		{
			name: "TransactionSign/" + txVal.Type().Name(),
			sign: func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error) {
				privateKey, ok := key.(PrivateKey)
				if !ok {
					panic("key is not a PrivateKey")
				}
				return TransactionSign(transactionInterface, privateKey)
			},
		},
		{
			name: "TransactionSignWith/" + txVal.Type().Name(),
			sign: func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error) {
				return TransactionSignWth(transactionInterface, newKey.PublicKey(), newKey.Sign)
			},
		},
		{
			name: "TransactionSignWithOperator/" + txVal.Type().Name(),
			sign: func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error) {
				return TransactionSignWithOperator(transactionInterface, client)
			},
		},
		{
			name: "TransactionAddSignature/" + txVal.Type().Name(),
			sign: func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error) {
				return TransactionAddSignature(transactionInterface, newKey.PublicKey(), signature)
			},
		},
	}
}

type transactionTest struct {
	name   string
	set    func(transactionInterface TransactionInterface) (TransactionInterface, error)
	get    func(transactionInterface TransactionInterface) (interface{}, error)
	assert func(t *testing.T, actual interface{})
}

func createTransactionTests(txName string, nodeAccountIds []AccountID) []transactionTest {
	return []transactionTest{
		{
			name: "TransactionTransactionID/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				transactionID := TransactionID{AccountID: &AccountID{Account: 9999}, ValidStart: &time.Time{}, scheduled: false, Nonce: nil}
				return TransactionSetTransactionID(transactionInterface, transactionID)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetTransactionID(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				transactionID := TransactionID{AccountID: &AccountID{Account: 9999}, ValidStart: &time.Time{}, scheduled: false, Nonce: nil}
				A := actual.(TransactionID)

				require.Equal(t, transactionID.AccountID, A.AccountID)
			},
		},
		{
			name: "TransactionTransactionMemo/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetTransactionMemo(transactionInterface, "test memo")
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetTransactionMemo(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, "test memo", actual)
			},
		},
		{
			name: "TransactionMaxTransactionFee/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetMaxTransactionFee(transactionInterface, NewHbar(1))
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetMaxTransactionFee(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, NewHbar(1), actual)
			},
		},
		{
			name: "TransactionTransactionValidDuration/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetTransactionValidDuration(transactionInterface, time.Second*10)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetTransactionValidDuration(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, time.Second*10, actual)
			},
		},
		{
			name: "TransactionNodeAccountIDs/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetNodeAccountIDs(transactionInterface, nodeAccountIds)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetNodeAccountIDs(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, nodeAccountIds, actual)
			},
		},
		{
			name: "TransactionMinBackoff/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				tx, _ := TransactionSetMaxBackoff(transactionInterface, time.Second*200)
				return TransactionSetMinBackoff(tx, time.Second*10)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetMinBackoff(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, time.Second*10, actual)
			},
		},
		{
			name: "TransactionMaxBackoff/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetMaxBackoff(transactionInterface, time.Second*200)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetMaxBackoff(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, time.Second*200, actual)
			},
		},
	}
}

// TransactionGetTransactionHash //needs to be tested in e2e tests
// TransactionGetTransactionHashPerNode //needs to be tested in e2e tests
// TransactionExecute //needs to be tested in e2e tests
