//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protowire"
	protobuf "google.golang.org/protobuf/proto"
)

func TestUnitTransactionBodyHashPendingSelectedID(t *testing.T) {
	tx := bodyHashNewAppend(t, "abcd", []AccountID{{Account: 3}})
	tx.SetTransactionID(NewTransactionIDWithValidStart(AccountID{Account: 1001}, time.Unix(1700000000, 999)))
	before := bodyHashSnapshot(t, tx)
	for range 2 {
		actual, err := tx.GetTransactionBodyHash()
		require.Error(t, err, "pending selected ID must not attest stale prepared bodies")
		require.Empty(t, actual)
		bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), false)
	}
}

func TestUnitTransactionBodyHashSelectedMembership(t *testing.T) {
	for name, edit := range map[string]func(*Transaction[TransactionInterface]){
		"missing-id": func(tx *Transaction[TransactionInterface]) { tx.transactionIDs.slice = tx.transactionIDs.slice[:1] },
		"extra-id": func(tx *Transaction[TransactionInterface]) {
			tx.transactionIDs.slice = append(tx.transactionIDs.slice, bodyHashTransactionID())
		},
		"duplicate-id": func(tx *Transaction[TransactionInterface]) { tx.transactionIDs.slice[1] = tx.transactionIDs.slice[0] },
		"changed-id": func(tx *Transaction[TransactionInterface]) {
			tx.transactionIDs.slice[1] = NewTransactionIDWithValidStart(AccountID{Account: 1001}, time.Unix(1700000000, 999))
		},
		"missing-node": func(tx *Transaction[TransactionInterface]) { tx.nodeAccountIDs.slice = tx.nodeAccountIDs.slice[:1] },
		"extra-node": func(tx *Transaction[TransactionInterface]) {
			tx.nodeAccountIDs.slice = append(tx.nodeAccountIDs.slice, AccountID{Account: 5})
		},
		"duplicate-node": func(tx *Transaction[TransactionInterface]) { tx.nodeAccountIDs.slice[1] = tx.nodeAccountIDs.slice[0] },
		"changed-node":   func(tx *Transaction[TransactionInterface]) { tx.nodeAccountIDs.slice[1] = AccountID{Account: 5} },
	} {
		t.Run(name, func(t *testing.T) {
			tx := bodyHashNewAppend(t, "abcdefgh", []AccountID{{Account: 3}, {Account: 4}})
			edit(tx.getBaseTransaction())
			before := bodyHashSnapshot(t, tx)
			actual, err := tx.GetTransactionBodyHash()
			require.Error(t, err)
			require.Empty(t, actual)
			bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), false)
		})
	}
	t.Run("selected-order-does-not-reorder-hash", func(t *testing.T) {
		vector := bodyHashVectorNamed(t, "later-node-order-normalized")
		tx := bodyHashImport(t, bodyHashUnhex(t, vector.TransactionListHex))
		base := tx.getBaseTransaction()
		base.transactionIDs.slice[0], base.transactionIDs.slice[1] = base.transactionIDs.slice[1], base.transactionIDs.slice[0]
		base.nodeAccountIDs.slice[0], base.nodeAccountIDs.slice[1] = base.nodeAccountIDs.slice[1], base.nodeAccountIDs.slice[0]
		actual, err := TransactionGetTransactionBodyHash(tx)
		require.NoError(t, err)
		require.Equal(t, bodyHashUnhex(t, vector.ExpectedHashHex), actual)
	})
}

func TestUnitTransactionBodyHashLegacyOriginalSignerMetadata(t *testing.T) {
	body := bodyHashBodiesFromList(t, bodyHashUnhex(t, bodyHashVectorNamed(t, "ordinary-single").TransactionListHex))[0]
	oldKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	originalSignature := oldKey.PublicKey()._ToSignaturePairProtobuf(oldKey.Sign(body))
	for _, asList := range []bool{false, true} {
		t.Run(map[bool]string{false: "bare", true: "list"}[asList], func(t *testing.T) {
			envelope := &services.Transaction{BodyBytes: body, SigMap: &services.SignatureMap{SigPair: []*services.SignaturePair{originalSignature}}}
			input := bodyHashMarshal(t, envelope)
			if asList {
				input = bodyHashMarshal(t, &sdk.TransactionList{TransactionList: []*services.Transaction{envelope}})
			}
			tx := bodyHashImport(t, input).(TransferTransaction)
			want, err := tx.GetTransactionBodyHash()
			require.NoError(t, err)
			newKey, err := PrivateKeyGenerateEd25519()
			require.NoError(t, err)
			calls := 0
			tx.SignWith(newKey.PublicKey(), func(raw []byte) []byte { calls++; return newKey.Sign(raw) })
			actual, err := tx.GetTransactionBodyHash()
			require.NoError(t, err)
			require.Zero(t, calls)
			require.Equal(t, want, actual)
			serialized, err := tx.ToBytes()
			require.NoError(t, err)
			require.Equal(t, 1, calls, "the original legacy signature must not suppress the new signer")
			require.Equal(t, [][]byte{body}, bodyHashBodiesFromList(t, serialized))
			signed := tx.signedTransactions._Get(0).(*services.SignedTransaction)
			require.Len(t, signed.SigMap.SigPair, 2)
			require.True(t, protobuf.Equal(originalSignature, signed.SigMap.SigPair[0]))
			require.True(t, newKey.PublicKey().VerifySignedMessage(body, signed.SigMap.SigPair[1].GetEd25519()))
			actual, err = TransactionGetTransactionBodyHash(bodyHashImport(t, serialized))
			require.NoError(t, err)
			require.Equal(t, want, actual)
		})
	}
	t.Run("list-later-original-signer", func(t *testing.T) {
		bodies := bodyHashBodiesFromList(t, bodyHashUnhex(t, bodyHashVectorNamed(t, "ordinary-multi-node").TransactionListHex))
		original := oldKey.PublicKey()._ToSignaturePairProtobuf(oldKey.Sign(bodies[1]))
		input := bodyHashMarshal(t, &sdk.TransactionList{TransactionList: []*services.Transaction{
			{BodyBytes: bodies[0]},
			{BodyBytes: bodies[1], SigMap: &services.SignatureMap{SigPair: []*services.SignaturePair{original}}},
		}})
		tx := bodyHashImport(t, input).(TransferTransaction)
		want, err := tx.GetTransactionBodyHash()
		require.NoError(t, err)
		newKey, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)
		calls := 0
		tx.SignWith(newKey.PublicKey(), func(raw []byte) []byte { calls++; return newKey.Sign(raw) })
		serialized, err := tx.ToBytes()
		require.NoError(t, err)
		require.Equal(t, 2, calls)
		require.Equal(t, bodies, bodyHashBodiesFromList(t, serialized))
		for i, value := range tx.signedTransactions.slice {
			signed := value.(*services.SignedTransaction)
			require.Len(t, signed.SigMap.SigPair, i+1)
			if i == 1 {
				require.True(t, protobuf.Equal(original, signed.SigMap.SigPair[0]))
			}
			require.True(t, newKey.PublicKey().VerifySignedMessage(bodies[i], signed.SigMap.SigPair[i].GetEd25519()))
		}
		actual, err := TransactionGetTransactionBodyHash(bodyHashImport(t, serialized))
		require.NoError(t, err)
		require.Equal(t, want, actual)
	})
}

func TestUnitTransactionBodyHashUnreviewedRegeneration(t *testing.T) {
	client := bodyHashReviewOfflineClient(t)
	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 1001}, key)
	tx, err := NewTransferTransaction().SetNodeAccountIDs([]AccountID{{Account: 3}}).FreezeWith(client)
	require.NoError(t, err)
	tx.SetMinBackoff(0).SetMaxBackoff(0).SetMaxRetry(4)
	request := &bodyHashOfflineImported{TransactionInterface: tx}
	var ids []string
	request.method = _Method{transaction: func(_ context.Context, envelope *services.Transaction, _ ...grpc.CallOption) (*services.TransactionResponse, error) {
		var signed services.SignedTransaction
		require.NoError(t, protobuf.Unmarshal(envelope.SignedTransactionBytes, &signed))
		var body services.TransactionBody
		require.NoError(t, protobuf.Unmarshal(signed.BodyBytes, &body))
		require.True(t, key.PublicKey().VerifySignedMessage(signed.BodyBytes, signed.SigMap.SigPair[0].GetEd25519()))
		id := _TransactionIDFromProtobuf(body.TransactionID).String()
		require.NotContains(t, ids, id)
		ids = append(ids, id)
		code := services.ResponseCodeEnum_OK
		if len(ids) < 3 {
			code = services.ResponseCodeEnum_TRANSACTION_EXPIRED
		}
		return &services.TransactionResponse{NodeTransactionPrecheckCode: code}, nil
	}}
	_, err = TransactionExecute(request, client)
	require.NoError(t, err)
	require.Len(t, ids, 3)
	require.False(t, tx.bodyReviewed)
}

func TestUnitTransactionBodyHashFreshCustomFeePreparation(t *testing.T) {
	for _, chunkSize := range []uint64{100, 4} {
		t.Run(map[uint64]string{100: "single-chunk", 4: "multi-chunk"}[chunkSize], func(t *testing.T) {
			limit := NewCustomFeeLimit().SetPayerId(AccountID{Account: 1001}).
				SetCustomFees([]*CustomFixedFee{NewCustomFixedFee().SetAmount(7)})
			tx, err := NewTopicMessageSubmitTransaction().SetTransactionID(bodyHashTransactionID()).
				SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).SetTopicID(TopicID{Topic: 1002}).
				SetMessage([]byte("abcdefgh")).SetChunkSize(chunkSize).SetCustomFeeLimits([]*CustomFeeLimit{limit}).FreezeWith(nil)
			require.NoError(t, err)
			key, err := PrivateKeyGenerateEd25519()
			require.NoError(t, err)
			calls := 0
			tx.SignWith(key.PublicKey(), func(body []byte) []byte { calls++; return key.Sign(body) })
			original := bodyHashPreparedBodies(tx)
			require.Len(t, original, map[uint64]int{100: 2, 4: 4}[chunkSize])
			before := bodyHashSnapshot(t, tx)
			want, err := tx.GetTransactionBodyHash()
			require.NoError(t, err, "Freeze must prepare custom fees without serialization")
			require.Zero(t, calls)
			bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), true)
			for _, raw := range original {
				var body services.TransactionBody
				require.NoError(t, protobuf.Unmarshal(raw, &body))
				require.Len(t, body.MaxCustomFees, 1)
				require.Equal(t, int64(7), body.MaxCustomFees[0].Fees[0].Amount)
			}
			serialized, err := tx.ToBytes()
			require.NoError(t, err)
			require.Equal(t, len(original), calls)
			require.Equal(t, original, bodyHashBodiesFromList(t, serialized))
			actual, err := TransactionGetTransactionBodyHash(bodyHashImport(t, serialized))
			require.NoError(t, err)
			require.Equal(t, want, actual)
			for _, signed := range tx.signedTransactions.slice {
				envelope := signed.(*services.SignedTransaction)
				require.True(t, key.PublicKey().VerifySignedMessage(envelope.BodyBytes, envelope.SigMap.SigPair[0].GetEd25519()))
			}
		})
	}
}

func TestUnitTransactionBodyHashImportedUnknownCommonFields(t *testing.T) {
	// Exact public probe from the supervising session; expected hashes were calculated independently.
	for _, probe := range []struct{ name, wire, hash string }{
		{"original-key", "0a5e2a5c0a5a0a150a080880e2cfaa06107b12070800100018e907180012060800100018001880c2d72f2202087832076368696c64204172020a00ca042212208a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c", "ea59778ddedc13c9b4fe32ac81a57b4834dd9a49f51c8a6343a0d3f9fea547bcb19eaa9575f3870afa89152b225f27f4"},
		{"unknown-key", "0a612a5f0a5d0a150a080880e2cfaa06107b12070800100018e907180012060800100018001880c2d72f2202087832076368696c64204172020a00ca042512208a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5ca00601", "809644d6d295c59087e79c4a6bee352be7b703645d96064983a054ca28045990b3a56ffe111f9234cf935cca23ea3113"},
	} {
		t.Run(probe.name, func(t *testing.T) {
			bodyHashCommonRoundtrip(t, bodyHashUnhex(t, probe.wire), bodyHashUnhex(t, probe.hash))
		})
	}

	fee := &services.CustomFeeLimit{AccountId: (&AccountID{Account: 1001})._ToProtobuf(), Fees: []*services.FixedFee{{Amount: 7,
		DenominatingTokenId: (&TokenID{Token: 1002})._ToProtobuf()}}}
	for _, message := range []protobuf.Message{fee, fee.AccountId, fee.Fees[0], fee.Fees[0].DenominatingTokenId} {
		message.ProtoReflect().SetUnknown(bodyHashBytesField(19000, []byte("nested unknown")))
	}
	tx, err := NewTopicMessageSubmitTransaction().SetTransactionID(bodyHashTransactionID()).SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTopicID(TopicID{Topic: 1002}).SetMessage([]byte("message")).FreezeWith(nil)
	require.NoError(t, err)
	var body services.TransactionBody
	require.NoError(t, protobuf.Unmarshal(bodyHashPreparedBodies(tx)[0], &body))
	body.MaxCustomFees = []*services.CustomFeeLimit{fee}
	withUnknown := bodyHashMarshal(t, &body)
	t.Run("unknown-fee-limit-payer-fee-token", func(t *testing.T) {
		bodyHashCommonRoundtrip(t, bodyHashList(t, withUnknown), bodyHashOrdinaryExpected([][]byte{withUnknown}))
	})
	for _, field := range []string{"amount", "payer", "token"} {
		t.Run("edited-"+field, func(t *testing.T) {
			imported := bodyHashImport(t, bodyHashList(t, withUnknown)).(TopicMessageSubmitTransaction)
			limit := imported.GetCustomFeeLimits()[0]
			switch field {
			case "amount":
				limit.CustomFees[0].SetAmount(8)
			case "payer":
				limit.SetPayerId(AccountID{Account: 1003})
			case "token":
				limit.CustomFees[0].SetDenominatingTokenID(TokenID{Token: 1003})
			}
			before := bodyHashSnapshot(t, imported)
			_, err := TransactionGetTransactionBodyHash(imported)
			require.Error(t, err)
			bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, imported), false)
		})
	}
}

func bodyHashCommonRoundtrip(t *testing.T, input, want []byte) {
	t.Helper()
	original := bytes.Clone(input)
	tx := bodyHashImport(t, input)
	bodies := bodyHashPreparedBodies(tx)
	oldKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	_, err = TransactionSign(tx, oldKey)
	require.NoError(t, err)
	for range 2 {
		before := bodyHashSnapshot(t, tx)
		actual, err := TransactionGetTransactionBodyHash(tx)
		require.NoError(t, err)
		require.Equal(t, want, actual)
		bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), true)
	}
	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	_, err = TransactionSign(tx, newKey)
	require.NoError(t, err)
	serialized, err := TransactionToBytes(tx)
	require.NoError(t, err)
	require.Equal(t, bodies, bodyHashBodiesFromList(t, serialized))
	for _, value := range tx.getBaseTransaction().signedTransactions.slice {
		signed := value.(*services.SignedTransaction)
		require.Len(t, signed.GetSigMap().GetSigPair(), 2)
		require.True(t, oldKey.PublicKey().VerifySignedMessage(signed.BodyBytes, signed.SigMap.SigPair[0].GetEd25519()))
		require.True(t, newKey.PublicKey().VerifySignedMessage(signed.BodyBytes, signed.SigMap.SigPair[1].GetEd25519()))
	}
	actual, err := TransactionGetTransactionBodyHash(bodyHashImport(t, serialized))
	require.NoError(t, err)
	require.Equal(t, want, actual)
	require.Equal(t, original, input)
}

func TestUnitTransactionBodyHashImportedNestedBatchKey(t *testing.T) {
	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	tx, err := NewTransferTransaction().SetTransactionID(bodyHashTransactionID()).SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetBatchKey(KeyListWithThreshold(1).Add(key.PublicKey())).Freeze()
	require.NoError(t, err)
	var body services.TransactionBody
	require.NoError(t, protobuf.Unmarshal(bodyHashPreparedBodies(tx)[0], &body))
	require.NotNil(t, body.BatchKey, "Freeze must prepare the batch key without serialization")
	threshold := body.BatchKey.GetThresholdKey()
	for _, message := range []protobuf.Message{body.BatchKey, threshold, threshold.Keys, threshold.Keys.Keys[0]} {
		message.ProtoReflect().SetUnknown(bodyHashBytesField(19000, []byte("unknown nested key")))
	}
	raw := bodyHashMarshal(t, &body)
	want := bodyHashOrdinaryExpected([][]byte{raw})
	bodyHashCommonRoundtrip(t, bodyHashList(t, raw), want)
	for _, reviewed := range []bool{false, true} {
		t.Run(map[bool]string{false: "pending-key-edit", true: "reviewed-key-edit"}[reviewed], func(t *testing.T) {
			imported := bodyHashImport(t, bodyHashList(t, raw))
			base := imported.getBaseTransaction()
			if reviewed {
				actual, err := TransactionGetTransactionBodyHash(imported)
				require.NoError(t, err)
				require.Equal(t, want, actual)
			}
			base.GetBatchKey().(*KeyList).SetThreshold(2)
			before := bodyHashSnapshot(t, imported)
			_, err := TransactionGetTransactionBodyHash(imported)
			require.Error(t, err)
			bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, imported), false)
			if reviewed {
				_, err = TransactionToBytes(imported)
				require.Error(t, err)
				require.Equal(t, before.SignedTransactions, bodyHashSnapshot(t, imported).SignedTransactions)
			}
		})
	}
}

func TestUnitTransactionBodyHashIdentifierVectors(t *testing.T) {
	data, err := os.ReadFile("testdata/transaction-body-hash-identifier-vectors.json")
	require.NoError(t, err)
	var fixture struct {
		Vectors []bodyHashVector `json:"vectors"`
	}
	require.NoError(t, json.Unmarshal(data, &fixture))
	require.Len(t, fixture.Vectors, 5)
	successes, rejections := 0, 0
	for _, vector := range fixture.Vectors {
		t.Run(vector.Name, func(t *testing.T) {
			input := bodyHashUnhex(t, vector.TransactionListHex)
			if vector.ExpectedError != "" {
				rejections++
				bodyHashRequireRejected(t, input)
				return
			}
			successes++
			bodyHashCommonRoundtrip(t, input, bodyHashUnhex(t, vector.ExpectedHashHex))
		})
	}
	require.Equal(t, 1, successes)
	require.Equal(t, 4, rejections)
}

func TestUnitTransactionBodyHashNestedIdentifierWire(t *testing.T) {
	base := bodyHashBodiesFromList(t, bodyHashUnhex(t, bodyHashVectorNamed(t, "ordinary-single").TransactionListHex))[0]
	var body services.TransactionBody
	require.NoError(t, protobuf.Unmarshal(base, &body))
	for _, location := range []string{"transaction-id", "payer-account", "timestamp"} {
		for name, record := range map[string][]byte{
			"six-byte-tag":    {0xa0, 0x86, 0x80, 0x80, 0x80, 0x00, 0x01},
			"six-byte-length": {0xa2, 0x06, 0x81, 0x80, 0x80, 0x80, 0x80, 0x00, 'x'},
			"group":           {0xa3, 0x06, 0xa4, 0x06},
			"overflow":        append([]byte{0xa0, 0x06}, append(bytes.Repeat([]byte{0xff}, 9), 0x02)...),
			"truncated":       {0xa2, 0x06, 0x02, 'x'},
		} {
			t.Run(location+"/"+name, func(t *testing.T) {
				bodyHashRequireRejected(t, bodyHashList(t, bodyHashIdentifierRecord(t, base, body.TransactionID, location, record)))
			})
		}
		t.Run(location+"/opaque-unknown", func(t *testing.T) {
			// Unknown bytes are opaque even if their contents would be unsupported as wire records.
			record := bodyHashBytesField(100, []byte{0xa3, 0x06, 0xa4, 0x06})
			raw := bodyHashIdentifierRecord(t, base, body.TransactionID, location, record)
			bodyHashCommonRoundtrip(t, bodyHashList(t, raw), bodyHashOrdinaryExpected([][]byte{raw}))
		})
	}
}

func TestUnitTransactionBodyHashOperationWireScope(t *testing.T) {
	base := bodyHashBodiesFromList(t, bodyHashUnhex(t, bodyHashVectorNamed(t, "ordinary-single").TransactionListHex))[0]
	var body services.TransactionBody
	require.NoError(t, protobuf.Unmarshal(base, &body))
	// Operation-specific submessages retain existing protobuf decoding checks.
	operation := append(bodyHashMarshal(t, body.GetCryptoTransfer()), 0xa3, 0x06, 0xa4, 0x06)
	raw := bodyHashReplaceField(t, base, 14, bodyHashBytesField(14, operation))
	bodyHashCommonRoundtrip(t, bodyHashList(t, raw), bodyHashOrdinaryExpected([][]byte{raw}))
}

func bodyHashIdentifierRecord(t *testing.T, base []byte, id *services.TransactionID, location string, record []byte) []byte {
	t.Helper()
	rawID := bodyHashMarshal(t, id)
	switch location {
	case "transaction-id":
		rawID = append(rawID, record...)
	case "payer-account":
		rawID = bodyHashReplaceField(t, rawID, 2, bodyHashBytesField(2, append(bodyHashMarshal(t, id.AccountID), record...)))
	case "timestamp":
		rawID = bodyHashReplaceField(t, rawID, 1, bodyHashBytesField(1, append(bodyHashMarshal(t, id.TransactionValidStart), record...)))
	}
	return bodyHashReplaceField(t, base, protowire.Number(1), bodyHashBytesField(1, rawID))
}
