//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitPrngTransactionSetRange(t *testing.T) {
	t.Parallel()

	rang := uint32(100)

	transaction := NewPrngTransaction().
		SetRange(rang)

	require.Equal(t, transaction.GetRange(), rang)
}

func TestUnitPrngTransactionFromBytes(t *testing.T) {
	t.Parallel()

	rang := uint32(100)

	transaction := NewPrngTransaction().
		SetRange(rang)

	bytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.NotNil(t, bytes)

	deserializedTransaction, err := TransactionFromBytes(bytes)
	require.NoError(t, err)

	switch tx := deserializedTransaction.(type) {
	case PrngTransaction:
		assert.Equal(t, transaction.GetRange(), tx.GetRange())
	default:
		t.Fatalf("expected PrngTransaction, got %T", deserializedTransaction)
	}
}

func TestUnitPrngTransactionScheduleProtobuf(t *testing.T) {
	t.Parallel()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})
	rang := uint32(100)
	nodeAccountID := []AccountID{{Account: 10}}

	tx, err := NewPrngTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetRange(rang).
		Freeze()
	require.NoError(t, err)

	expected := &services.SchedulableTransactionBody{
		TransactionFee: 100000000,
		Data: &services.SchedulableTransactionBody_UtilPrng{
			UtilPrng: &services.UtilPrngTransactionBody{
				Range: int32(rang),
			},
		},
	}

	actual, err := tx.buildScheduled()
	require.NoError(t, err)
	require.Equal(t, expected.GetUtilPrng(), actual.GetUtilPrng())
}
