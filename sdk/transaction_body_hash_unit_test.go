//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
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

// Candidate for hub PR #288, pending maintainer agreement on framing, imports,
// and rejection rules. The shared fixture bytes and expectations are copied verbatim.
func TestUnitTransactionBodyHashSharedVectors(t *testing.T) {
	vectors := bodyHashVectors(t)
	successes, rejections := 0, 0
	for _, vector := range vectors {
		if vector.ExpectedError == "" {
			successes++
		} else {
			rejections++
		}
		t.Run(vector.Name, func(t *testing.T) {
			input := bodyHashUnhex(t, vector.TransactionListHex)
			original := bytes.Clone(input)
			if vector.ExpectedError != "" {
				bodyHashRequireRejected(t, input)
				require.Equal(t, original, input)
				return
			}

			tx := bodyHashImport(t, input)
			wantBodies := bodyHashBodiesFromList(t, input)
			require.Equal(t, wantBodies, bodyHashPreparedBodies(tx))
			before := bodyHashSnapshot(t, tx)
			want := bodyHashUnhex(t, vector.ExpectedHashHex)
			for range 3 {
				actual, err := TransactionGetTransactionBodyHash(tx)
				require.NoError(t, err)
				require.Len(t, actual, sha512.Size384)
				require.Equal(t, want, actual)
				bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), true)
			}

			// Compare with an identical import that has never been body-hashed.
			control := bodyHashImport(t, original)
			controlBytes, err := TransactionToBytes(control)
			require.NoError(t, err)
			serialized, err := TransactionToBytes(tx)
			require.NoError(t, err)
			require.Equal(t, controlBytes, serialized, "hashing must not change network serialization")
			require.Equal(t, wantBodies, bodyHashBodiesFromList(t, serialized), "roundtrip must retain original body bytes")
			roundtrip := bodyHashImport(t, serialized)
			actual, err := TransactionGetTransactionBodyHash(roundtrip)
			require.NoError(t, err)
			require.Equal(t, want, actual)
			require.Equal(t, original, input)
		})
	}
	require.Equal(t, 20, successes)
	require.Equal(t, 15, rejections)
}

func TestUnitTransactionBodyHashFreshOrdinary(t *testing.T) {
	nodes := []AccountID{{Account: 3}, {Account: 4}}
	single := bodyHashNewTransfer(t, nodes[:1])
	multi := bodyHashNewTransfer(t, nodes)
	reversed := bodyHashNewTransfer(t, []AccountID{nodes[1], nodes[0]})

	var hashes [][]byte
	for _, tx := range []*TransferTransaction{single, multi, reversed} {
		before := bodyHashSnapshot(t, tx)
		actual, err := tx.GetTransactionBodyHash()
		require.NoError(t, err)
		require.Len(t, actual, sha512.Size384)
		require.Equal(t, bodyHashOrdinaryExpected(bodyHashPreparedBodies(tx)), actual)
		wrapped, err := TransactionGetTransactionBodyHash(tx)
		require.NoError(t, err)
		require.Equal(t, actual, wrapped)
		bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), true)
		hashes = append(hashes, actual)
	}
	inner := sha512.Sum384(bodyHashPreparedBodies(single)[0])
	require.NotEqual(t, inner[:], hashes[0], "even one entry requires the final SHA-384")
	require.NotEqual(t, hashes[0], hashes[1], "every prepared node is committed")
	require.NotEqual(t, hashes[1], hashes[2], "the first chunk determines node order")
	laterNodeChanged := bodyHashNewTransfer(t, []AccountID{{Account: 3}, {Account: 5}})
	require.Equal(t, bodyHashPreparedBodies(multi)[0], bodyHashPreparedBodies(laterNodeChanged)[0])
	changedHash, err := laterNodeChanged.GetTransactionBodyHash()
	require.NoError(t, err)
	require.NotEqual(t, hashes[1], changedHash, "a change only in the final node must be detected")
}

func TestUnitTransactionBodyHashFreshChunks(t *testing.T) {
	nodes := []AccountID{{Account: 3}, {Account: 4}}
	tx := bodyHashNewAppend(t, "abcdefgh", nodes)
	changed := bodyHashNewAppend(t, "abcdefg!", nodes)
	reversed := bodyHashNewAppend(t, "abcdefgh", []AccountID{nodes[1], nodes[0]})
	originalBodies := bodyHashPreparedBodies(tx)
	require.Len(t, originalBodies, 4)
	changedBodies := bodyHashPreparedBodies(changed)
	require.Equal(t, originalBodies[:2], changedBodies[:2], "only the final chunk changes")
	require.NotEqual(t, originalBodies[2:], changedBodies[2:])

	actual, err := tx.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Equal(t, bodyHashOrdinaryExpected(originalBodies), actual)
	changedHash, err := changed.GetTransactionBodyHash()
	require.NoError(t, err)
	require.NotEqual(t, actual, changedHash)
	reversedHash, err := reversed.GetTransactionBodyHash()
	require.NoError(t, err)
	require.NotEqual(t, actual, reversedHash)

	for _, order := range [][]int{{0, 2, 1, 3}, {0, 1, 3, 2}} {
		imported := bodyHashImport(t, bodyHashList(t,
			originalBodies[order[0]], originalBodies[order[1]],
			originalBodies[order[2]], originalBodies[order[3]]))
		hash, err := TransactionGetTransactionBodyHash(imported)
		require.NoError(t, err)
		require.Equal(t, actual, hash, "node-major and later-node permutations canonicalize hash traversal")
	}
	chunkReversed := bodyHashImport(t, bodyHashList(t,
		originalBodies[2], originalBodies[3], originalBodies[0], originalBodies[1]))
	hash, err := TransactionGetTransactionBodyHash(chunkReversed)
	require.NoError(t, err)
	require.NotEqual(t, actual, hash, "chunk first appearance is committed")
}

func TestUnitTransactionBodyHashFreshBatchSignaturesAndContent(t *testing.T) {
	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	makeBatch := func(signChildren bool, memo string, reverse bool, count int) *BatchTransaction {
		children := make([]TransactionInterface, 0, 2)
		for i, childMemo := range []string{"first child", "second child"} {
			child := NewTransferTransaction().
				SetTransactionID(NewTransactionIDWithValidStart(AccountID{Account: 1001}, time.Unix(1700000000, int64(123+i)))).
				SetNodeAccountIDs([]AccountID{{Account: 3}}).
				SetBatchKey(key.PublicKey()).
				SetTransactionMemo(childMemo).
				AddHbarTransfer(AccountID{Account: 1001}, NewHbar(-1)).
				AddHbarTransfer(AccountID{Account: 1002}, NewHbar(1))
			_, err := child.Freeze()
			require.NoError(t, err)
			if signChildren {
				child.Sign(key)
			}
			children = append(children, child)
		}
		if reverse {
			children[0], children[1] = children[1], children[0]
		}
		batch := NewBatchTransaction().
			SetTransactionID(bodyHashTransactionID()).
			SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
			SetTransactionMemo(memo).
			SetInnerTransactions(children[:count])
		_, err := batch.Freeze()
		require.NoError(t, err)
		return batch
	}

	unsigned := makeBatch(false, "outer", false, 2)
	signed := makeBatch(true, "outer", false, 2)
	want, err := unsigned.GetTransactionBodyHash()
	require.NoError(t, err)
	before := bodyHashSnapshot(t, signed)
	actual, err := signed.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Equal(t, want, actual, "child signatures are excluded")
	bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, signed), true)

	// Actual outer signing after review must remain possible for every node.
	signed.Sign(key)
	signable, err := signed.GetSignableNodeBodyBytesList()
	require.NoError(t, err)
	for _, entry := range signable {
		_, err := signed.AddSignatureV2(key.PublicKey(), key.Sign(entry.Body), entry.TransactionID, entry.NodeID)
		require.NoError(t, err)
	}
	actual, err = signed.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Equal(t, want, actual, "outer signatures are excluded")
	for _, variant := range []*BatchTransaction{
		makeBatch(false, "changed outer", false, 2),
		makeBatch(false, "outer", true, 2),
		makeBatch(false, "outer", false, 1),
	} {
		actual, err := variant.GetTransactionBodyHash()
		require.NoError(t, err)
		require.NotEqual(t, want, actual)
	}
}

func TestUnitTransactionBodyHashLegacyImports(t *testing.T) {
	vector := bodyHashVectorNamed(t, "ordinary-single")
	body := bodyHashBodiesFromList(t, bodyHashUnhex(t, vector.TransactionListHex))[0]
	modernSigned := bodyHashMarshal(t, &services.SignedTransaction{BodyBytes: body})
	var signedList sdk.TransactionList
	signedVector := bodyHashVectorNamed(t, "ordinary-signed")
	require.NoError(t, protobuf.Unmarshal(bodyHashUnhex(t, signedVector.TransactionListHex), &signedList))
	var signed services.SignedTransaction
	require.NoError(t, protobuf.Unmarshal(signedList.TransactionList[0].SignedTransactionBytes, &signed))
	for name, envelope := range map[string]*services.Transaction{
		"modern":        {SignedTransactionBytes: modernSigned},
		"legacy":        {BodyBytes: body},
		"legacy-signed": {BodyBytes: body, SigMap: signed.SigMap},
	} {
		for _, asList := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "/single", true: "/list"}[asList], func(t *testing.T) {
				input := bodyHashMarshal(t, envelope)
				if asList {
					input = bodyHashMarshal(t, &sdk.TransactionList{TransactionList: []*services.Transaction{envelope}})
				}
				control := bodyHashImport(t, input)
				before := bodyHashSnapshot(t, control)
				controlBytes, err := TransactionToBytes(control)
				require.NoError(t, err)
				require.Equal(t, before.SignedTransactions, bodyHashSnapshot(t, control).SignedTransactions)
				tx := bodyHashImport(t, input)
				require.Equal(t, [][]byte{body}, bodyHashPreparedBodies(tx))
				actual, err := TransactionGetTransactionBodyHash(tx)
				require.NoError(t, err)
				require.Equal(t, bodyHashUnhex(t, vector.ExpectedHashHex), actual)
				serialized, err := TransactionToBytes(tx)
				require.NoError(t, err)
				require.Equal(t, controlBytes, serialized)
				require.Equal(t, [][]byte{body}, bodyHashBodiesFromList(t, serialized))
				actual, err = TransactionGetTransactionBodyHash(bodyHashImport(t, serialized))
				require.NoError(t, err)
				require.Equal(t, bodyHashUnhex(t, vector.ExpectedHashHex), actual)
			})
		}
	}
}

func TestUnitTransactionBodyHashLegacyRejectsIncompleteEntry(t *testing.T) {
	vector := bodyHashVectorNamed(t, "ordinary-multi-node")
	bodies := bodyHashBodiesFromList(t, bodyHashUnhex(t, vector.TransactionListHex))
	missingNode := bodyHashReplaceField(t, bodies[1], 2, nil)
	input := bodyHashMarshal(t, &sdk.TransactionList{TransactionList: []*services.Transaction{
		{BodyBytes: bodies[0]}, {BodyBytes: missingNode},
	}})
	bodyHashRequireRejected(t, input)
}

func TestUnitTransactionBodyHashOriginalWireBytes(t *testing.T) {
	vector := bodyHashVectorNamed(t, "ordinary-single")
	body := bodyHashBodiesFromList(t, bodyHashUnhex(t, vector.TransactionListHex))[0]
	for name, suffix := range map[string][]byte{
		"unknown-body-content": bodyHashBytesField(19000, []byte("opaque body proof")),
		"explicit-default":     {0x28, 0x00}, // generateRecord=false; protobuf-equivalent, different original bytes.
	} {
		t.Run(name, func(t *testing.T) {
			original := append(bytes.Clone(body), suffix...)
			tx := bodyHashImport(t, bodyHashList(t, original))
			actual, err := TransactionGetTransactionBodyHash(tx)
			require.NoError(t, err)
			require.Equal(t, bodyHashOrdinaryExpected([][]byte{original}), actual)
			require.NotEqual(t, bodyHashUnhex(t, vector.ExpectedHashHex), actual)
			serialized, err := TransactionToBytes(tx)
			require.NoError(t, err)
			require.Equal(t, [][]byte{original}, bodyHashBodiesFromList(t, serialized))
		})
	}

	batch := bodyHashVectorNamed(t, "batch-unsigned")
	batchBody := bodyHashBodiesFromList(t, bodyHashUnhex(t, batch.TransactionListHex))[0]
	var decoded services.TransactionBody
	require.NoError(t, protobuf.Unmarshal(batchBody, &decoded))
	children := decoded.GetAtomicBatch().GetTransactions()
	var wrapper []byte
	for _, child := range children {
		// Unknown SignedTransaction metadata, like the hash algorithm flag, is outside the commitment.
		child = append(bytes.Clone(child), bodyHashBytesField(19000, []byte("envelope metadata"))...)
		wrapper = append(wrapper, bodyHashBytesField(1, child)...)
	}
	modified := bodyHashReplaceField(t, batchBody, 74, bodyHashBytesField(74, wrapper))
	actual, err := TransactionGetTransactionBodyHash(bodyHashImport(t, bodyHashList(t, modified)))
	require.NoError(t, err)
	require.Equal(t, bodyHashUnhex(t, batch.ExpectedHashHex), actual)
}

func TestUnitTransactionBodyHashMalformedImports(t *testing.T) {
	ordinary := bodyHashVectorNamed(t, "ordinary-single")
	body := bodyHashBodiesFromList(t, bodyHashUnhex(t, ordinary.TransactionListHex))[0]
	batch := bodyHashVectorNamed(t, "batch-unsigned")
	batchBody := bodyHashBodiesFromList(t, bodyHashUnhex(t, batch.TransactionListHex))[0]
	var decoded services.TransactionBody
	require.NoError(t, protobuf.Unmarshal(batchBody, &decoded))
	child := decoded.GetAtomicBatch().Transactions[0]
	var signed services.SignedTransaction
	require.NoError(t, protobuf.Unmarshal(child, &signed))
	batchWithWrapper := func(wrapper []byte) []byte {
		return bodyHashReplaceField(t, batchBody, 74, bodyHashBytesField(74, wrapper))
	}
	batchWithChild := func(child []byte) []byte { return batchWithWrapper(bodyHashBytesField(1, child)) }

	cases := map[string][]byte{
		"truncated-body-varint":           append(bytes.Clone(body), 0x80),
		"invalid-field-number":            append(bytes.Clone(body), 0x00),
		"invalid-wire-type":               append(bytes.Clone(body), 0x7e),
		"truncated-fixed64":               append(bytes.Clone(body), 0x79, 0x01),
		"truncated-fixed32":               append(bytes.Clone(body), 0x7d, 0x01),
		"truncated-length-delimited":      append(bytes.Clone(body), 0x7a, 0x7f, 0x01),
		"overflow-varint":                 append(bytes.Clone(body), 0x78, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x02),
		"conflicting-ordinary-operations": append(bytes.Clone(body), bodyHashBytesField(16, nil)...),
		"conflicting-batch-operations":    append(bytes.Clone(batchBody), bodyHashBytesField(14, nil)...),
		"mistyped-batch-operation":        bodyHashReplaceField(t, batchBody, 74, protowire.AppendVarint(protowire.AppendTag(nil, 74, protowire.VarintType), 1)),
		"mistyped-wrapper-child":          batchWithWrapper([]byte{0x08, 0x01}),
		"invalid-wrapper-wire":            batchWithWrapper([]byte{0x0a, 0xff}),
		"mistyped-child-body":             batchWithChild([]byte{0x08, 0x01}),
		"invalid-child-envelope":          batchWithChild(append(bytes.Clone(child), 0x80)),
		"invalid-child-body":              batchWithChild(bodyHashBytesField(1, append(bytes.Clone(signed.BodyBytes), 0x80))),
		"child-without-operation":         batchWithChild(bodyHashBytesField(1, bodyHashReplaceField(t, signed.BodyBytes, 14, nil))),
		"child-conflicting-operations":    batchWithChild(bodyHashBytesField(1, append(bytes.Clone(signed.BodyBytes), bodyHashBytesField(16, nil)...))),
		"child-without-transaction-id":    batchWithChild(bodyHashBytesField(1, bodyHashReplaceField(t, signed.BodyBytes, 1, nil))),
		"child-without-node-id":           batchWithChild(bodyHashBytesField(1, bodyHashReplaceField(t, signed.BodyBytes, 2, nil))),
	}
	for name, malformed := range cases {
		t.Run(name, func(t *testing.T) { bodyHashRequireRejected(t, bodyHashList(t, malformed)) })
	}
	for name, input := range map[string][]byte{
		"empty-list":                {},
		"truncated-list":            {0x0a, 0x7f},
		"truncated-signed-envelope": bodyHashMarshal(t, &sdk.TransactionList{TransactionList: []*services.Transaction{{SignedTransactionBytes: []byte{0x0a, 0xff}}}}),
	} {
		t.Run(name, func(t *testing.T) { bodyHashRequireRejected(t, input) })
	}
}

func TestUnitTransactionBodyHashLayoutValidation(t *testing.T) {
	vector := bodyHashVectorNamed(t, "chunk-major")
	bodies := bodyHashBodiesFromList(t, bodyHashUnhex(t, vector.TransactionListHex))
	var decoded services.TransactionBody
	require.NoError(t, protobuf.Unmarshal(bodies[0], &decoded))
	// Explicit zero shard/realm fields decode to the same node as omitted defaults.
	equivalentNode := append(bodyHashMarshal(t, decoded.NodeAccountID), 0x08, 0x00, 0x10, 0x00)
	equivalentID := append(bodyHashMarshal(t, decoded.TransactionID), 0x18, 0x00, 0x20, 0x00)
	duplicate := bodyHashReplaceField(t, bodies[0], 2, bodyHashBytesField(2, equivalentNode))
	duplicate = bodyHashReplaceField(t, duplicate, 1, bodyHashBytesField(1, equivalentID))
	extraNode := bodyHashReplaceField(t, bodies[2], 2, bodyHashBytesField(2, bodyHashMarshal(t, AccountID{Account: 5}._ToProtobuf())))
	changedType := bodyHashReplaceField(t, bodies[1], 16, bodyHashBytesField(14, nil))
	changedMemo := bodyHashReplaceField(t, bodies[1], 6, bodyHashBytesField(6, []byte("inconsistent node content")))
	changedLaterMemo := bodyHashReplaceField(t, bodies[3], 6, bodyHashBytesField(6, []byte("inconsistent chunk content")))
	for name, entries := range map[string][][]byte{
		"missing-transaction-id":       {bodyHashReplaceField(t, bodies[0], 1, nil)},
		"missing-node-id":              {bodyHashReplaceField(t, bodies[0], 2, nil)},
		"empty-transaction-id":         {bodyHashReplaceField(t, bodies[0], 1, bodyHashBytesField(1, nil))},
		"duplicate-decoded-ids":        {bodies[0], duplicate},
		"extra-node-in-later-chunk":    {bodies[0], bodies[1], bodies[2], bodies[3], extraNode},
		"replacement-node-in-chunk":    {bodies[0], bodies[1], bodies[2], extraNode},
		"inconsistent-operation-type":  {bodies[0], changedType, bodies[2], bodies[3]},
		"inconsistent-content-at-node": {bodies[0], changedMemo, bodies[2], bodies[3]},
		"inconsistent-chunk-content":   {bodies[0], bodies[1], bodies[2], changedLaterMemo},
	} {
		t.Run(name, func(t *testing.T) { bodyHashRequireRejected(t, bodyHashList(t, entries...)) })
	}
}

func TestUnitTransactionBodyHashUnfrozen(t *testing.T) {
	for _, tx := range []TransactionInterface{NewTransferTransaction(), NewFileAppendTransaction(), NewBatchTransaction()} {
		before := bodyHashSnapshot(t, tx)
		actual, err := TransactionGetTransactionBodyHash(tx)
		require.Error(t, err)
		require.Empty(t, actual)
		bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), false)
	}
	tx := NewTransferTransaction()
	actual, err := tx.GetTransactionBodyHash()
	require.Error(t, err)
	require.Empty(t, actual)
}

func TestUnitTransactionBodyHashFailedGetterPreservesState(t *testing.T) {
	for name, body := range map[string][]byte{
		"empty-prepared-body":     {},
		"malformed-prepared-body": {0x0a, 0xff},
		"missing-prepared-ids":    bodyHashBytesField(14, nil),
	} {
		t.Run(name, func(t *testing.T) {
			tx := bodyHashNewAppend(t, "abcdefgh", []AccountID{{Account: 3}, {Account: 4}})
			tx.transactionIDs.index, tx.nodeAccountIDs.index = 1, 1
			tx.signedTransactions.index, tx.transactions.index = 2, 1
			// Corrupt a later prepared entry so validation fails after valid entries.
			tx.signedTransactions._Get(3).(*services.SignedTransaction).BodyBytes = bytes.Clone(body)
			key, err := PrivateKeyGenerateEd25519()
			require.NoError(t, err)
			calls := 0
			tx.SignWith(key.PublicKey(), func([]byte) []byte { calls++; return nil })
			before := bodyHashSnapshot(t, tx)
			for range 2 {
				var actual []byte
				require.NotPanics(t, func() { actual, err = tx.GetTransactionBodyHash() })
				require.Error(t, err)
				require.Empty(t, actual)
				require.Zero(t, calls)
				bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), false)
			}
		})
	}
}

func TestUnitTransactionBodyHashFailedLayoutPreservesState(t *testing.T) {
	for _, name := range []string{"duplicate", "missing-node", "inconsistent-content"} {
		t.Run(name, func(t *testing.T) {
			tx := bodyHashNewAppend(t, "abcdefgh", []AccountID{{Account: 3}, {Account: 4}})
			switch name {
			case "duplicate":
				tx.signedTransactions._Get(3).(*services.SignedTransaction).BodyBytes = bytes.Clone(tx.signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes)
			case "missing-node":
				tx.signedTransactions.slice = tx.signedTransactions.slice[:3]
			case "inconsistent-content":
				last := tx.signedTransactions._Get(3).(*services.SignedTransaction)
				last.BodyBytes = bodyHashReplaceField(t, last.BodyBytes, 6, bodyHashBytesField(6, []byte("changed memo")))
			}
			key, err := PrivateKeyGenerateEd25519()
			require.NoError(t, err)
			calls := 0
			tx.SignWith(key.PublicKey(), func([]byte) []byte { calls++; return nil })
			tx.transactionIDs.index, tx.nodeAccountIDs.index = 1, 1
			before := bodyHashSnapshot(t, tx)
			require.False(t, before.Locks[0])
			for range 2 {
				actual, err := tx.GetTransactionBodyHash()
				require.Error(t, err)
				require.Empty(t, actual)
				require.Zero(t, calls)
				bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), false)
			}
		})
	}
}

func TestUnitTransactionBodyHashDeferredSignerAndCursors(t *testing.T) {
	tx := bodyHashNewAppend(t, "abcdefghijkl", []AccountID{{Account: 3}, {Account: 4}})
	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	var signedBodies [][]byte
	tx.SignWith(key.PublicKey(), func(body []byte) []byte {
		signedBodies = append(signedBodies, bytes.Clone(body))
		return key.Sign(body)
	})
	tx.transactionIDs.index, tx.nodeAccountIDs.index = 1, 1
	tx.signedTransactions.index, tx.transactions.index = 3, 2
	before := bodyHashSnapshot(t, tx)
	wantBodies := bodyHashPreparedBodies(tx)
	for range 3 {
		actual, err := tx.GetTransactionBodyHash()
		require.NoError(t, err)
		require.Equal(t, bodyHashOrdinaryExpected(wantBodies), actual)
		require.Empty(t, signedBodies, "hashing must never invoke a deferred signer")
		bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), true)
	}

	// Add externally collected signatures to all reviewed chunks after the getter.
	entries, err := tx.GetSignableNodeBodyBytesList()
	require.NoError(t, err)
	for _, entry := range entries {
		_, err := tx.AddSignatureV2(key.PublicKey(), key.Sign(entry.Body), entry.TransactionID, entry.NodeID)
		require.NoError(t, err)
	}
	require.Empty(t, signedBodies)
	require.Equal(t, wantBodies, bodyHashPreparedBodies(tx))
	for _, value := range tx.signedTransactions.slice {
		signed := value.(*services.SignedTransaction)
		require.Len(t, signed.GetSigMap().GetSigPair(), 1)
		require.True(t, key.PublicKey().VerifySignedMessage(signed.BodyBytes, signed.SigMap.SigPair[0].GetEd25519()))
	}
	actual, err := tx.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Equal(t, bodyHashOrdinaryExpected(wantBodies), actual)

	// Deferred callbacks can subsequently sign the exact reviewed bytes too.
	ordinary := bodyHashNewTransfer(t, []AccountID{{Account: 3}, {Account: 4}})
	ordinaryBodies := bodyHashPreparedBodies(ordinary)
	ordinary.SignWith(key.PublicKey(), func(body []byte) []byte {
		signedBodies = append(signedBodies, bytes.Clone(body))
		return key.Sign(body)
	})
	_, err = ordinary.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Empty(t, signedBodies)
	serialized, err := ordinary.ToBytes()
	require.NoError(t, err)
	require.Equal(t, ordinaryBodies, signedBodies)
	require.Equal(t, ordinaryBodies, bodyHashBodiesFromList(t, serialized))
}

func TestUnitTransactionBodyHashSignAfterReview(t *testing.T) {
	tx := bodyHashNewTransfer(t, []AccountID{{Account: 3}, {Account: 4}})
	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	wantBodies := bodyHashPreparedBodies(tx)
	want, err := tx.GetTransactionBodyHash()
	require.NoError(t, err)
	tx.Sign(key)
	serialized, err := tx.ToBytes()
	require.NoError(t, err)
	require.Equal(t, wantBodies, bodyHashBodiesFromList(t, serialized))
	for _, value := range tx.signedTransactions.slice {
		signed := value.(*services.SignedTransaction)
		require.Len(t, signed.GetSigMap().GetSigPair(), 1)
		require.True(t, key.PublicKey().VerifySignedMessage(signed.BodyBytes, signed.SigMap.SigPair[0].GetEd25519()))
	}
	actual, err := tx.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Equal(t, want, actual)
	secondKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	entries, err := tx.GetSignableNodeBodyBytesList()
	require.NoError(t, err)
	for _, entry := range entries {
		_, err := tx.AddSignatureV2(secondKey.PublicKey(), secondKey.Sign(entry.Body), entry.TransactionID, entry.NodeID)
		require.NoError(t, err)
	}
	for _, value := range tx.signedTransactions.slice {
		pairs := value.(*services.SignedTransaction).SigMap.SigPair
		require.Len(t, pairs, 2)
		pairs[0], pairs[1] = pairs[1], pairs[0]
	}
	before := bodyHashSnapshot(t, tx)
	actual, err = tx.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Equal(t, want, actual, "additional and reordered signatures are excluded")
	bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), true)
	removed, err := tx.RemoveAllSignatures()
	require.NoError(t, err)
	require.NotEmpty(t, removed)
	actual, err = tx.GetTransactionBodyHash()
	require.NoError(t, err)
	require.Equal(t, want, actual, "removing signatures must preserve the body hash")
}

func TestUnitTransactionBodyHashLocksPreparedSelection(t *testing.T) {
	// FileAppend.Freeze produces an unlocked ID collection, exercising a new lock.
	tx := bodyHashNewAppend(t, "abcdefgh", []AccountID{{Account: 3}, {Account: 4}})
	require.False(t, tx.transactionIDs.locked)
	before := bodyHashSnapshot(t, tx)
	_, err := tx.GetTransactionBodyHash()
	require.NoError(t, err)
	require.True(t, tx.transactionIDs.locked)
	require.True(t, tx.nodeAccountIDs.locked)
	require.False(t, tx.signedTransactions.locked, "signature collection must remain writable")
	require.False(t, tx.transactions.locked)

	// Exercise the same regeneration guard used on TRANSACTION_EXPIRED, offline.
	client, err := ClientForNetworkV2(map[string]AccountID{"127.0.0.1:50211": {Account: 3}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	client.SetOperator(AccountID{Account: 1001}, key)
	client.SetMaxNodesPerTransaction(1)
	require.False(t, tx.regenerateID(client))
	bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), true)
	require.Equal(t, []AccountID{{Account: 3}, {Account: 4}}, tx.GetNodeAccountIDs())
}

func TestUnitTransactionBodyHashNetworkHashUnchanged(t *testing.T) {
	// Existing network hash accessors are tested only on ordinary transactions.
	for _, nodes := range [][]AccountID{{{Account: 3}}, {{Account: 3}, {Account: 4}}} {
		tx := bodyHashNewTransfer(t, nodes)
		beforeBytes, err := tx.ToBytes()
		require.NoError(t, err)
		beforeNetworkHash, err := tx.GetTransactionHash()
		require.NoError(t, err)
		beforePerNode, err := tx.GetTransactionHashPerNode()
		require.NoError(t, err)
		bodyHash, err := tx.GetTransactionBodyHash()
		require.NoError(t, err)
		afterNetworkHash, err := tx.GetTransactionHash()
		require.NoError(t, err)
		afterPerNode, err := tx.GetTransactionHashPerNode()
		require.NoError(t, err)
		afterBytes, err := tx.ToBytes()
		require.NoError(t, err)
		require.Equal(t, beforeBytes, afterBytes)
		require.Equal(t, beforeNetworkHash, afterNetworkHash)
		require.Equal(t, beforePerNode, afterPerNode)

		key, err := PrivateKeyGenerateEd25519()
		require.NoError(t, err)
		tx.Sign(key)
		afterSigning, err := tx.GetTransactionHash()
		require.NoError(t, err)
		require.NotEqual(t, beforeNetworkHash, afterSigning, "network hashes still include signatures")
		actual, err := tx.GetTransactionBodyHash()
		require.NoError(t, err)
		require.Equal(t, bodyHash, actual)
	}
}

func TestUnitTransactionBodyHashExecutionKeepsReviewedSelection(t *testing.T) {
	for _, status := range []services.ResponseCodeEnum{services.ResponseCodeEnum_OK, services.ResponseCodeEnum_TRANSACTION_EXPIRED} {
		t.Run(status.String(), func(t *testing.T) {
			client, err := ClientForNetworkV2(map[string]AccountID{
				"127.0.0.1:50211": {Account: 3},
				"127.0.0.1:50212": {Account: 4},
			})
			require.NoError(t, err)
			// Cached empty channels and an in-process method prevent sockets or RPCs.
			for _, managed := range client.network.nodes {
				managed.(*_Node).channel = &_Channel{}
			}
			t.Cleanup(func() {
				for _, managed := range client.network.nodes {
					managed.(*_Node).channel = nil
				}
				require.NoError(t, client.Close())
			})
			key, err := PrivateKeyGenerateEd25519()
			require.NoError(t, err)
			client.SetOperator(AccountID{Account: 1001}, key)
			client.SetMaxNodesPerTransaction(2)
			transfer, err := NewTransferTransaction().
				AddHbarTransfer(AccountID{Account: 1001}, NewHbar(-1)).
				AddHbarTransfer(AccountID{Account: 1002}, NewHbar(1)).
				FreezeWith(client)
			require.NoError(t, err)
			require.False(t, transfer.transactionIDs.locked)
			require.False(t, transfer.nodeAccountIDs.locked)
			reviewedNodes := transfer.GetNodeAccountIDs()
			require.Len(t, reviewedNodes, 2)
			reviewedID := transfer.GetTransactionID().String()
			reviewedBodies := bodyHashPreparedBodies(transfer)
			want, err := transfer.GetTransactionBodyHash()
			require.NoError(t, err)
			client.SetMaxNodesPerTransaction(1)
			transfer.nodeAccountIDs.index = 1
			calls := 0
			tx := &bodyHashOfflineTransfer{TransferTransaction: transfer}
			tx.method = _Method{transaction: func(_ context.Context, request *services.Transaction, _ ...grpc.CallOption) (*services.TransactionResponse, error) {
				calls++
				var signed services.SignedTransaction
				require.NoError(t, protobuf.Unmarshal(request.SignedTransactionBytes, &signed))
				require.Equal(t, reviewedBodies[1], signed.BodyBytes)
				return &services.TransactionResponse{NodeTransactionPrecheckCode: status}, nil
			}}
			response, err := TransactionExecute(tx, client)
			if status == services.ResponseCodeEnum_OK {
				require.NoError(t, err)
				require.Equal(t, reviewedNodes[1], response.NodeID)
			} else {
				require.Error(t, err, "reviewed IDs cannot regenerate on TRANSACTION_EXPIRED")
			}
			require.Equal(t, 1, calls)
			require.Equal(t, reviewedNodes, transfer.GetNodeAccountIDs())
			require.Equal(t, reviewedID, transfer.GetTransactionID().String())
			require.Equal(t, reviewedBodies, bodyHashPreparedBodies(transfer))
			actual, err := transfer.GetTransactionBodyHash()
			require.NoError(t, err)
			require.Equal(t, want, actual)
		})
	}
}

type bodyHashOfflineTransfer struct {
	*TransferTransaction
	method _Method
}

func (tx *bodyHashOfflineTransfer) getMethod(*_Channel) _Method { return tx.method }

func (tx *bodyHashOfflineTransfer) getBaseTransaction() *Transaction[TransactionInterface] {
	base := tx.TransferTransaction.getBaseTransaction()
	base.childTransaction = tx
	return base
}

type bodyHashVector struct {
	Name               string `json:"name"`
	TransactionListHex string `json:"transactionListHex"`
	ExpectedHashHex    string `json:"expectedHashHex"`
	ExpectedError      string `json:"expectedError"`
}

func bodyHashVectors(t *testing.T) []bodyHashVector {
	t.Helper()
	data, err := os.ReadFile("testdata/transaction-body-hash-vectors.json")
	require.NoError(t, err)
	var fixture struct {
		Profile string           `json:"profile"`
		Vectors []bodyHashVector `json:"vectors"`
	}
	require.NoError(t, json.Unmarshal(data, &fixture))
	require.Equal(t, "transaction-body-hash-proposed-v1", fixture.Profile)
	require.Equal(t, 35, len(fixture.Vectors))
	return fixture.Vectors
}

func bodyHashVectorNamed(t *testing.T, name string) bodyHashVector {
	t.Helper()
	for _, vector := range bodyHashVectors(t) {
		if vector.Name == name {
			return vector
		}
	}
	t.Fatalf("missing shared vector %q", name)
	return bodyHashVector{}
}

func bodyHashUnhex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	require.NoError(t, err)
	return decoded
}

func bodyHashMarshal(t *testing.T, value protobuf.Message) []byte {
	t.Helper()
	data, err := protobuf.Marshal(value)
	require.NoError(t, err)
	return data
}

func bodyHashImport(t *testing.T, input []byte) TransactionInterface {
	t.Helper()
	tx, err := TransactionFromBytes(input)
	require.NoError(t, err)
	tx, err = TransactionFreezeWith(tx, nil)
	require.NoError(t, err)
	require.True(t, tx.getBaseTransaction().IsFrozen())
	return tx
}

func bodyHashRequireRejected(t *testing.T, input []byte) {
	t.Helper()
	original := bytes.Clone(input)
	var err error
	var hash []byte
	require.NotPanics(t, func() {
		var tx TransactionInterface
		tx, err = TransactionFromBytes(input)
		if err != nil {
			return
		}
		tx, err = TransactionFreezeWith(tx, nil)
		if err == nil {
			before := bodyHashSnapshot(t, tx)
			hash, err = TransactionGetTransactionBodyHash(tx)
			if err != nil {
				bodyHashAssertUnchanged(t, before, bodyHashSnapshot(t, tx), false)
			}
		}
	})
	require.Error(t, err, "import, freeze, or getter must reject invalid input")
	require.Empty(t, hash)
	require.Equal(t, original, input)
}

func bodyHashBodiesFromList(t *testing.T, input []byte) [][]byte {
	t.Helper()
	var list sdk.TransactionList
	require.NoError(t, protobuf.Unmarshal(input, &list))
	bodies := make([][]byte, 0, len(list.TransactionList))
	for _, transaction := range list.TransactionList {
		var signed services.SignedTransaction
		require.NoError(t, protobuf.Unmarshal(transaction.SignedTransactionBytes, &signed))
		bodies = append(bodies, bytes.Clone(signed.BodyBytes))
	}
	return bodies
}

func bodyHashPreparedBodies(tx TransactionInterface) [][]byte {
	bodies := make([][]byte, 0, tx.getBaseTransaction().signedTransactions._Length())
	for _, value := range tx.getBaseTransaction().signedTransactions.slice {
		bodies = append(bodies, bytes.Clone(value.(*services.SignedTransaction).BodyBytes))
	}
	return bodies
}

func bodyHashOrdinaryExpected(bodies [][]byte) []byte {
	var entries []byte
	for _, body := range bodies {
		hash := sha512.Sum384(body)
		entries = append(entries, hash[:]...)
	}
	result := sha512.Sum384(entries)
	return result[:]
}

func bodyHashList(t *testing.T, bodies ...[]byte) []byte {
	t.Helper()
	list := &sdk.TransactionList{}
	for _, body := range bodies {
		list.TransactionList = append(list.TransactionList, &services.Transaction{
			SignedTransactionBytes: bodyHashMarshal(t, &services.SignedTransaction{BodyBytes: body}),
		})
	}
	return bodyHashMarshal(t, list)
}

func bodyHashBytesField(number protowire.Number, payload []byte) []byte {
	return protowire.AppendBytes(protowire.AppendTag(nil, number, protowire.BytesType), payload)
}

// Mutate test input records without re-encoding the other original body bytes.
func bodyHashReplaceField(t *testing.T, input []byte, number protowire.Number, replacement []byte) []byte {
	t.Helper()
	var result []byte
	for len(input) > 0 {
		field, wire, tagSize := protowire.ConsumeTag(input)
		require.Positive(t, tagSize)
		valueSize := protowire.ConsumeFieldValue(field, wire, input[tagSize:])
		require.GreaterOrEqual(t, valueSize, 0)
		size := tagSize + valueSize
		if field != number {
			result = append(result, input[:size]...)
		}
		input = input[size:]
	}
	return append(result, replacement...)
}

func bodyHashTransactionID() TransactionID {
	return NewTransactionIDWithValidStart(AccountID{Account: 1001}, time.Unix(1700000000, 123))
}

func bodyHashNewTransfer(t *testing.T, nodes []AccountID) *TransferTransaction {
	t.Helper()
	tx, err := NewTransferTransaction().
		SetTransactionID(bodyHashTransactionID()).
		SetNodeAccountIDs(nodes).
		AddHbarTransfer(AccountID{Account: 1001}, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 1002}, NewHbar(1)).
		Freeze()
	require.NoError(t, err)
	return tx
}

func bodyHashNewAppend(t *testing.T, contents string, nodes []AccountID) *FileAppendTransaction {
	t.Helper()
	tx, err := NewFileAppendTransaction().
		SetTransactionID(bodyHashTransactionID()).
		SetNodeAccountIDs(nodes).
		SetFileID(FileID{File: 1002}).
		SetContents([]byte(contents)).
		SetMaxChunkSize(4).
		FreezeWith(nil) // The embedded generic Freeze does not dispatch to chunk preparation.
	require.NoError(t, err)
	return tx
}

type bodyHashState struct {
	SignedTransactions [][]byte
	Transactions       [][]byte
	TransactionIDs     []string
	Nodes              []AccountID
	Indices            [4]int
	Locks              [4]bool
	PublicKeys         []string
	SignerCount        int
	FreezeError        error
	KeyError           error
	RegenerateID       bool
	BodyReviewed       bool
}

func bodyHashSnapshot(t *testing.T, transaction TransactionInterface) bodyHashState {
	t.Helper()
	tx := transaction.getBaseTransaction()
	state := bodyHashState{
		Indices:      [4]int{tx.transactionIDs.index, tx.nodeAccountIDs.index, tx.signedTransactions.index, tx.transactions.index},
		Locks:        [4]bool{tx.transactionIDs.locked, tx.nodeAccountIDs.locked, tx.signedTransactions.locked, tx.transactions.locked},
		SignerCount:  len(tx.transactionSigners),
		FreezeError:  tx.freezeError,
		KeyError:     tx.keyError,
		RegenerateID: tx.regenerateTransactionID,
		BodyReviewed: tx.bodyReviewed,
	}
	for _, value := range tx.signedTransactions.slice {
		state.SignedTransactions = append(state.SignedTransactions, bodyHashMarshal(t, value.(*services.SignedTransaction)))
	}
	for _, value := range tx.transactions.slice {
		state.Transactions = append(state.Transactions, bodyHashMarshal(t, value.(*services.Transaction)))
	}
	for _, value := range tx.transactionIDs.slice {
		state.TransactionIDs = append(state.TransactionIDs, value.(TransactionID).String())
	}
	for _, value := range tx.nodeAccountIDs.slice {
		state.Nodes = append(state.Nodes, value.(AccountID))
	}
	for _, key := range tx.publicKeys {
		state.PublicKeys = append(state.PublicKeys, key.String())
	}
	return state
}

func bodyHashAssertUnchanged(t *testing.T, before, after bodyHashState, allowReviewLock bool) {
	t.Helper()
	if allowReviewLock {
		require.True(t, after.Locks[0], "reviewed transaction IDs must be locked")
		require.True(t, after.Locks[1], "reviewed node membership must be locked")
		before.Locks[0], before.Locks[1] = true, true
		require.True(t, after.BodyReviewed)
		before.BodyReviewed = true
	}
	require.Equal(t, before, after)
}
