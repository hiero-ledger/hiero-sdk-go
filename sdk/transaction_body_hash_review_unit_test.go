//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protowire"
	protobuf "google.golang.org/protobuf/proto"
)

func TestUnitTransactionBodyHashBatchNodeContentConsistency(t *testing.T) {
	vector := bodyHashVectorNamed(t, "batch-multi-node-unsigned")
	bodies := bodyHashBodiesFromList(t, bodyHashUnhex(t, vector.TransactionListHex))
	for _, variant := range []string{"algorithm-flag", "unknown-envelope-metadata", "changed-child-body"} {
		t.Run(variant, func(t *testing.T) {
			var outer services.TransactionBody
			require.NoError(t, protobuf.Unmarshal(bodies[1], &outer))
			var wrapper []byte
			for childIndex, envelope := range outer.GetAtomicBatch().Transactions {
				switch variant {
				case "algorithm-flag":
					envelope = protowire.AppendVarint(protowire.AppendTag(bytes.Clone(envelope), 3, protowire.VarintType), 1)
				case "unknown-envelope-metadata":
					envelope = append(bytes.Clone(envelope), bodyHashBytesField(19000, []byte("metadata"))...)
				case "changed-child-body":
					if childIndex != 0 {
						wrapper = append(wrapper, bodyHashBytesField(1, envelope)...)
						continue
					}
					var signed services.SignedTransaction
					require.NoError(t, protobuf.Unmarshal(envelope, &signed))
					signed.BodyBytes = bodyHashReplaceField(t, signed.BodyBytes, 6, bodyHashBytesField(6, []byte("changed content")))
					envelope = bodyHashMarshal(t, &signed)
				}
				wrapper = append(wrapper, bodyHashBytesField(1, envelope)...)
			}
			changed := bodyHashReplaceField(t, bodies[1], 74, bodyHashBytesField(74, wrapper))
			input := bodyHashList(t, bodies[0], changed)
			if variant == "changed-child-body" {
				bodyHashRequireRejected(t, input)
				return
			}
			tx := bodyHashImport(t, input)
			before := bodyHashSnapshot(t, tx)
			for range 2 {
				actual, err := TransactionGetTransactionBodyHash(tx)
				require.NoError(t, err)
				require.Equal(t, bodyHashUnhex(t, vector.ExpectedHashHex), actual)
				bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), true)
			}
			serialized, err := TransactionToBytes(tx)
			require.NoError(t, err)
			require.Equal(t, [][]byte{bodies[0], changed}, bodyHashBodiesFromList(t, serialized))
		})
	}
}

func TestUnitTransactionBodyHashUnknownTransactionIDRoundtrip(t *testing.T) {
	vector := bodyHashVectorNamed(t, "ordinary-multi-node")
	bodies := bodyHashBodiesFromList(t, bodyHashUnhex(t, vector.TransactionListHex))
	oldKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	list := &sdk.TransactionList{}
	for i, body := range bodies {
		var decoded services.TransactionBody
		require.NoError(t, protobuf.Unmarshal(body, &decoded))
		id := append(bodyHashMarshal(t, decoded.TransactionID), bodyHashBytesField(19000, bytes.Repeat([]byte{byte(i + 1)}, i+1))...)
		bodies[i] = bodyHashReplaceField(t, body, 1, bodyHashBytesField(1, id))
		list.TransactionList = append(list.TransactionList, &services.Transaction{
			SignedTransactionBytes: bodyHashMarshal(t, &services.SignedTransaction{
				BodyBytes: bodies[i], SigMap: &services.SignatureMap{SigPair: []*services.SignaturePair{
					oldKey.PublicKey()._ToSignaturePairProtobuf(oldKey.Sign(bodies[i])),
				}},
			}),
		})
	}
	tx := bodyHashImport(t, bodyHashMarshal(t, list))
	want := bodyHashOrdinaryExpected(bodies)
	actual, err := TransactionGetTransactionBodyHash(tx)
	require.NoError(t, err)
	require.Equal(t, want, actual)
	require.NotEqual(t, bodyHashUnhex(t, vector.ExpectedHashHex), actual)
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	_, err = TransactionSign(tx, newKey)
	require.NoError(t, err)
	serialized, err := TransactionToBytes(tx)
	require.NoError(t, err)
	require.Equal(t, bodies, bodyHashBodiesFromList(t, serialized))
	var signedList sdk.TransactionList
	require.NoError(t, protobuf.Unmarshal(serialized, &signedList))
	for _, envelope := range signedList.TransactionList {
		var signed services.SignedTransaction
		require.NoError(t, protobuf.Unmarshal(envelope.SignedTransactionBytes, &signed))
		require.Len(t, signed.GetSigMap().GetSigPair(), 2)
		require.True(t, oldKey.PublicKey().VerifySignedMessage(signed.BodyBytes, signed.SigMap.SigPair[0].GetEd25519()))
		require.True(t, newKey.PublicKey().VerifySignedMessage(signed.BodyBytes, signed.SigMap.SigPair[1].GetEd25519()))
	}
	actual, err = TransactionGetTransactionBodyHash(bodyHashImport(t, serialized))
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func TestUnitTransactionBodyHashImportedExternalSigning(t *testing.T) {
	vector := bodyHashVectorNamed(t, "ordinary-single")
	body := bodyHashBodiesFromList(t, bodyHashUnhex(t, vector.TransactionListHex))[0]
	for _, legacy := range []bool{false, true} {
		for _, v2 := range []bool{false, true} {
			t.Run(map[bool]string{false: "modern", true: "legacy"}[legacy]+map[bool]string{false: "/AddSignature", true: "/AddSignatureV2"}[v2], func(t *testing.T) {
				input := bodyHashList(t, body)
				if legacy {
					input = bodyHashMarshal(t, &sdk.TransactionList{TransactionList: []*services.Transaction{{BodyBytes: body}}})
				}
				tx := bodyHashImport(t, input).(TransferTransaction)
				want, err := tx.GetTransactionBodyHash()
				require.NoError(t, err)
				entries, err := tx.GetSignableNodeBodyBytesList()
				require.NoError(t, err)
				key, err := PrivateKeyGenerateEd25519()
				require.NoError(t, err)
				require.NotPanics(t, func() {
					if v2 {
						_, err = tx.AddSignatureV2(key.PublicKey(), key.Sign(entries[0].Body), entries[0].TransactionID, entries[0].NodeID)
					} else {
						_, err = TransactionAddSignature(tx, key.PublicKey(), key.Sign(entries[0].Body))
					}
				})
				require.NoError(t, err)
				serialized, err := tx.ToBytes()
				require.NoError(t, err)
				require.Equal(t, [][]byte{body}, bodyHashBodiesFromList(t, serialized))
				signatures, err := tx.GetSignatures()
				require.NoError(t, err)
				require.Len(t, signatures[entries[0].NodeID], 1)
				actual, err := tx.GetTransactionBodyHash()
				require.NoError(t, err)
				require.Equal(t, want, actual)
			})
		}
	}
}

func TestUnitTransactionBodyHashStandaloneBatchChild(t *testing.T) {
	for _, name := range []string{"batch-unsigned", "batch-child-and-outer-signatures"} {
		t.Run(name, func(t *testing.T) {
			vector := bodyHashVectorNamed(t, name)
			var outer services.TransactionBody
			require.NoError(t, protobuf.Unmarshal(bodyHashBodiesFromList(t, bodyHashUnhex(t, vector.TransactionListHex))[0], &outer))
			envelope := outer.GetAtomicBatch().Transactions[0]
			input := bodyHashMarshal(t, &sdk.TransactionList{TransactionList: []*services.Transaction{{SignedTransactionBytes: envelope}}})
			tx := bodyHashImport(t, input)
			nodes, err := TransactionGetNodeAccountIDs(tx)
			require.NoError(t, err)
			require.Equal(t, []AccountID{{}}, nodes)
			before := bodyHashSnapshot(t, tx)
			want, err := TransactionGetTransactionBodyHash(tx)
			require.NoError(t, err)
			serialized, err := TransactionToBytes(tx)
			require.NoError(t, err)
			require.Equal(t, bodyHashBodiesFromList(t, input), bodyHashBodiesFromList(t, serialized))
			require.Equal(t, before.SignedTransactions, bodyHashSnapshot(t, tx).SignedTransactions)
			actual, err := TransactionGetTransactionBodyHash(bodyHashImport(t, serialized))
			require.NoError(t, err)
			require.Equal(t, want, actual)
		})
	}
}

func TestUnitTransactionBodyHashImportedDispatch(t *testing.T) {
	for _, name := range []string{"node-major", "later-node-order-normalized"} {
		for _, chunk := range []int{0, 1} {
			t.Run(name+map[int]string{0: "/first-chunk", 1: "/last-chunk"}[chunk], func(t *testing.T) {
				vector := bodyHashVectorNamed(t, name)
				tx := bodyHashImport(t, bodyHashUnhex(t, vector.TransactionListHex))
				base := tx.getBaseTransaction()
				base.transactionIDs.index = chunk
				wantID := base.GetTransactionID().String()
				wantBodies := bodyHashPreparedBodies(tx)
				before := bodyHashSnapshot(t, tx)
				for range 2 {
					_, err := TransactionGetTransactionBodyHash(tx)
					require.NoError(t, err)
					bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), true)
				}
				client := bodyHashReviewOfflineClient(t)
				base.SetMinBackoff(0)
				base.SetMaxBackoff(0)
				base.SetMaxRetry(3)
				calls := 0
				request := &bodyHashOfflineImported{TransactionInterface: tx}
				request.method = _Method{transaction: func(_ context.Context, envelope *services.Transaction, _ ...grpc.CallOption) (*services.TransactionResponse, error) {
					var signed services.SignedTransaction
					require.NoError(t, protobuf.Unmarshal(envelope.SignedTransactionBytes, &signed))
					var body services.TransactionBody
					require.NoError(t, protobuf.Unmarshal(signed.BodyBytes, &body))
					require.Equal(t, wantID, _TransactionIDFromProtobuf(body.TransactionID).String())
					require.Equal(t, int64(3+calls), body.NodeAccountID.GetAccountNum())
					require.Contains(t, wantBodies, signed.BodyBytes)
					calls++
					if calls == 1 {
						return nil, status.Error(codes.Unavailable, "retry next reviewed node")
					}
					return &services.TransactionResponse{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK}, nil
				}}
				_, err := TransactionExecute(request, client)
				require.NoError(t, err)
				require.Equal(t, 2, calls)
				require.Equal(t, wantBodies, bodyHashPreparedBodies(tx))
				afterDispatch := bodyHashSnapshot(t, tx)
				actual, err := TransactionGetTransactionBodyHash(tx)
				require.NoError(t, err)
				require.Equal(t, bodyHashUnhex(t, vector.ExpectedHashHex), actual)
				bodyHashAssertUnchanged(t, afterDispatch, bodyHashSnapshot(t, tx), true)
			})
		}
	}
}

type bodyHashOfflineImported struct {
	TransactionInterface
	method _Method
}

func (tx *bodyHashOfflineImported) getMethod(*_Channel) _Method { return tx.method }
func (tx *bodyHashOfflineImported) getBaseTransaction() *Transaction[TransactionInterface] {
	base := tx.TransactionInterface.getBaseTransaction()
	base.childTransaction = tx
	return base
}

func bodyHashReviewOfflineClient(t *testing.T) *Client {
	t.Helper()
	client, err := ClientForNetworkV2(map[string]AccountID{"127.0.0.1:50211": {Account: 3}, "127.0.0.1:50212": {Account: 4}})
	require.NoError(t, err)
	for _, managed := range client.network.nodes {
		managed.(*_Node).channel = &_Channel{}
	}
	t.Cleanup(func() {
		for _, managed := range client.network.nodes {
			managed.(*_Node).channel = nil
		}
		require.NoError(t, client.Close())
	})
	return client
}

func TestUnitTransactionBodyHashReviewedCommonSetters(t *testing.T) {
	for name, edit := range map[string]func(*TransferTransaction){
		"memo":        func(tx *TransferTransaction) { tx.SetTransactionMemo("after") },
		"fee":         func(tx *TransferTransaction) { tx.SetMaxTransactionFee(NewHbar(99)) },
		"high-volume": func(tx *TransferTransaction) { tx.SetHighVolume(true) },
		"duration":    func(tx *TransferTransaction) { tx.SetTransactionValidDuration(time.Minute) },
		"batch-key": func(tx *TransferTransaction) {
			key, err := PrivateKeyGenerateEd25519()
			require.NoError(t, err)
			tx.SetBatchKey(key.PublicKey())
		},
	} {
		t.Run(name, func(t *testing.T) {
			tx := bodyHashNewTransfer(t, []AccountID{{Account: 3}})
			beforeBodies := bodyHashPreparedBodies(tx)
			want, err := tx.GetTransactionBodyHash()
			require.NoError(t, err)
			edit(tx)
			require.ErrorIs(t, tx.freezeError, errTransactionIsFrozen)
			actual, err := tx.GetTransactionBodyHash()
			require.NoError(t, err)
			require.Equal(t, want, actual)
			_, err = tx.ToBytes()
			require.NoError(t, err)
			require.Equal(t, beforeBodies, bodyHashPreparedBodies(tx))
		})
	}
}

func TestUnitTransactionBodyHashReviewedCommonWrappers(t *testing.T) {
	for name, edit := range map[string]func(TransactionInterface) (TransactionInterface, error){
		"memo": func(tx TransactionInterface) (TransactionInterface, error) {
			return TransactionSetTransactionMemo(tx, "after")
		},
		"fee": func(tx TransactionInterface) (TransactionInterface, error) {
			return TransactionSetMaxTransactionFee(tx, NewHbar(99))
		},
		"duration": func(tx TransactionInterface) (TransactionInterface, error) {
			return TransactionSetTransactionValidDuration(tx, time.Minute)
		},
		"high-volume": func(tx TransactionInterface) (TransactionInterface, error) { return TransactionSetHighVolume(tx, true) },
	} {
		t.Run(name, func(t *testing.T) {
			tx := bodyHashImport(t, bodyHashList(t, bodyHashPreparedBodies(bodyHashNewTransfer(t, []AccountID{{Account: 3}}))[0]))
			want, err := TransactionGetTransactionBodyHash(tx)
			require.NoError(t, err)
			_, err = edit(tx)
			require.ErrorIs(t, err, errTransactionIsFrozen)
			actual, err := TransactionGetTransactionBodyHash(tx)
			require.NoError(t, err)
			require.Equal(t, want, actual)
		})
	}
}

func TestUnitTransactionBodyHashPendingCommonUpdates(t *testing.T) {
	for name, edit := range map[string]func(*TransferTransaction){
		"memo":        func(tx *TransferTransaction) { tx.SetTransactionMemo("pending") },
		"fee":         func(tx *TransferTransaction) { tx.SetMaxTransactionFee(NewHbar(99)) },
		"high-volume": func(tx *TransferTransaction) { tx.SetHighVolume(true) },
		"duration":    func(tx *TransferTransaction) { tx.SetTransactionValidDuration(time.Minute) },
	} {
		t.Run(name, func(t *testing.T) {
			tx := bodyHashNewTransfer(t, []AccountID{{Account: 3}})
			edit(tx)
			before := bodyHashSnapshot(t, tx)
			actual, err := tx.GetTransactionBodyHash()
			require.Error(t, err)
			require.Empty(t, actual)
			bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), false)
		})
	}
}

func TestUnitTransactionBodyHashMutableCustomFeeReview(t *testing.T) {
	for _, reviewed := range []bool{false, true} {
		t.Run(map[bool]string{false: "pending", true: "reviewed"}[reviewed], func(t *testing.T) {
			fee := NewCustomFixedFee().SetAmount(1)
			limits := []*CustomFeeLimit{NewCustomFeeLimit().SetPayerId(AccountID{Account: 1001}).SetCustomFees([]*CustomFixedFee{fee})}
			tx, err := NewTopicMessageSubmitTransaction().SetTransactionID(bodyHashTransactionID()).
				SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).SetTopicID(TopicID{Topic: 1002}).
				SetMessage([]byte("message")).SetCustomFeeLimits(limits).FreezeWith(nil)
			require.NoError(t, err)
			if reviewed {
				_, err = tx.GetTransactionBodyHash()
				require.NoError(t, err)
			}
			key, err := PrivateKeyGenerateEd25519()
			require.NoError(t, err)
			calls := 0
			tx.SignWith(key.PublicKey(), func(body []byte) []byte { calls++; return key.Sign(body) })
			limits[0].CustomFees[0].SetAmount(99)
			before := bodyHashSnapshot(t, tx)
			actual, err := tx.GetTransactionBodyHash()
			require.Error(t, err)
			require.Empty(t, actual)
			bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), false)
			if reviewed {
				_, err = tx.ToBytes()
				require.Error(t, err)
				_, err = tx.GetTransactionHash()
				require.Error(t, err)
				request := &bodyHashOfflineImported{TransactionInterface: tx}
				request.method = _Method{transaction: func(context.Context, *services.Transaction, ...grpc.CallOption) (*services.TransactionResponse, error) {
					calls++
					return nil, nil
				}}
				_, err = TransactionExecute(request, bodyHashReviewOfflineClient(t))
				require.Error(t, err)
				require.Zero(t, calls)
				require.Equal(t, before.SignedTransactions, bodyHashSnapshot(t, tx).SignedTransactions)
			}
		})
	}
}

func TestUnitTransactionBodyHashReviewedCustomFeeSetters(t *testing.T) {
	for _, name := range []string{"set", "add", "clear"} {
		t.Run(name, func(t *testing.T) {
			limit := NewCustomFeeLimit().SetPayerId(AccountID{Account: 1001}).SetCustomFees([]*CustomFixedFee{NewCustomFixedFee().SetAmount(1)})
			tx, err := NewTopicMessageSubmitTransaction().SetTransactionID(bodyHashTransactionID()).SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
				SetTopicID(TopicID{Topic: 1002}).SetMessage([]byte("message")).SetCustomFeeLimits([]*CustomFeeLimit{limit}).FreezeWith(nil)
			require.NoError(t, err)
			want, err := tx.GetTransactionBodyHash()
			require.NoError(t, err)
			serialized, err := tx.ToBytes()
			require.NoError(t, err)
			imported, err := TransactionGetTransactionBodyHash(bodyHashImport(t, serialized))
			require.NoError(t, err)
			require.Equal(t, want, imported, "common fee limits must not duplicate per imported node")
			switch name {
			case "set":
				tx.SetCustomFeeLimits([]*CustomFeeLimit{limit})
			case "add":
				tx.AddCustomFeeLimit(limit)
			case "clear":
				tx.ClearCustomFeeLimits()
			}
			require.ErrorIs(t, tx.freezeError, errTransactionIsFrozen)
			require.Len(t, tx.GetCustomFeeLimits(), 1)
			actual, err := tx.GetTransactionBodyHash()
			require.NoError(t, err)
			require.Equal(t, want, actual)
		})
	}
}

func TestUnitTransactionBodyHashMutableBatchKey(t *testing.T) {
	for _, reviewed := range []bool{false, true} {
		t.Run(map[bool]string{false: "pending", true: "reviewed"}[reviewed], func(t *testing.T) {
			key, err := PrivateKeyGenerateEd25519()
			require.NoError(t, err)
			keyList := NewKeyList().Add(key.PublicKey())
			tx, err := NewTransferTransaction().SetTransactionID(bodyHashTransactionID()).SetNodeAccountIDs([]AccountID{{Account: 3}}).
				SetBatchKey(keyList).Freeze()
			require.NoError(t, err)
			if reviewed {
				_, err = tx.GetTransactionBodyHash()
				require.NoError(t, err)
			}
			keyList.Add(key.PublicKey())
			before := bodyHashSnapshot(t, tx)
			actual, err := tx.GetTransactionBodyHash()
			require.Error(t, err)
			require.Empty(t, actual)
			bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), false)
			if reviewed {
				_, err = tx.ToBytes()
				require.Error(t, err)
				require.Equal(t, before.SignedTransactions, bodyHashSnapshot(t, tx).SignedTransactions)
			}
		})
	}
}

func TestUnitTransactionBodyHashUnreviewedMutationCompatibility(t *testing.T) {
	tx := bodyHashNewTransfer(t, []AccountID{{Account: 3}})
	tx.SetTransactionMemo("allowed before review").SetMaxTransactionFee(NewHbar(5)).SetHighVolume(true)
	serialized, err := tx.ToBytes()
	require.NoError(t, err)
	var body services.TransactionBody
	require.NoError(t, protobuf.Unmarshal(bodyHashBodiesFromList(t, serialized)[0], &body))
	require.Equal(t, "allowed before review", body.Memo)
	require.Equal(t, uint64(NewHbar(5).AsTinybar()), body.TransactionFee)
	require.True(t, body.HighVolume)
}

func TestUnitTransactionBodyHashReviewedTopicRefreeze(t *testing.T) {
	tx, err := NewTopicMessageSubmitTransaction().SetTransactionID(bodyHashTransactionID()).SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTopicID(TopicID{Topic: 1002}).SetMessage([]byte("before review")).FreezeWith(nil)
	require.NoError(t, err)
	want, err := tx.GetTransactionBodyHash()
	require.NoError(t, err)
	before := bodyHashPreparedBodies(tx)
	tx.SetMessage([]byte("after review"))
	_, err = tx.FreezeWith(nil)
	require.ErrorIs(t, err, errTransactionIsFrozen)
	require.Equal(t, before, bodyHashPreparedBodies(tx))
	actual, err := tx.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Equal(t, want, actual)
}

func TestUnitTransactionBodyHashReviewedExecutionOptions(t *testing.T) {
	tx := bodyHashNewTransfer(t, []AccountID{{Account: 3}})
	want, err := tx.GetTransactionBodyHash()
	require.NoError(t, err)
	tx.SetMaxRetry(5).SetMinBackoff(time.Millisecond).SetMaxBackoff(time.Second).SetRegenerateTransactionID(true)
	require.NoError(t, tx.freezeError)
	actual, err := tx.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Equal(t, want, actual)
}
