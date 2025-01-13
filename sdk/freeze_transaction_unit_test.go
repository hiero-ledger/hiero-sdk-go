//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitFreezeTransactionSetStartTime(t *testing.T) {
	t.Parallel()

	startTime := time.Now().Add(time.Hour)

	transaction := NewFreezeTransaction().
		SetStartTime(startTime)

	require.Equal(t, transaction.GetStartTime(), startTime)
}

func TestUnitFreezeTransactionSetEndTime(t *testing.T) {
	t.Parallel()

	endTime := time.Now().Add(2 * time.Hour)

	transaction := NewFreezeTransaction().
		SetEndTime(endTime)

	require.Equal(t, transaction.GetEndTime(), endTime)
}

func TestUnitFreezeTransactionSetFileID(t *testing.T) {
	t.Parallel()

	fileID := FileID{File: 1}

	transaction := NewFreezeTransaction().
		SetFileID(fileID)

	require.Equal(t, transaction.GetFileID(), &fileID)
}

func TestUnitFreezeTransactionSetFreezeType(t *testing.T) {
	t.Parallel()

	freezeType := FreezeTypeFreezeOnly

	transaction := NewFreezeTransaction().
		SetFreezeType(freezeType)

	require.Equal(t, transaction.GetFreezeType(), freezeType)
}

func TestUnitFreezeTransactionSetFileHash(t *testing.T) {
	t.Parallel()

	fileHash := []byte{1, 2, 3, 4, 5}

	transaction := NewFreezeTransaction().
		SetFileHash(fileHash)

	require.Equal(t, transaction.GetFileHash(), fileHash)
}

func TestUnitFreezeTransactionToBytes(t *testing.T) {
	t.Parallel()

	startTime := time.Now().Add(time.Hour)
	endTime := time.Now().Add(2 * time.Hour)
	fileID := FileID{File: 1}
	fileHash := []byte{1, 2, 3, 4, 5}
	freezeType := FreezeTypeFreezeOnly

	transaction := NewFreezeTransaction().
		SetStartTime(startTime).
		SetEndTime(endTime).
		SetFileID(fileID).
		SetFileHash(fileHash).
		SetFreezeType(freezeType)

	bytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.NotNil(t, bytes)
}

func TestUnitFreezeTransactionFromBytes(t *testing.T) {
	t.Parallel()

	startTime := time.Now().Add(time.Hour)
	fileID := FileID{File: 1}
	fileHash := []byte{1, 2, 3, 4, 5}
	freezeType := FreezeTypeFreezeOnly

	transaction := NewFreezeTransaction().
		SetStartTime(startTime).
		SetFileID(fileID).
		SetFileHash(fileHash).
		SetFreezeType(freezeType)

	bytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.NotNil(t, bytes)

	deserializedTransaction, err := TransactionFromBytes(bytes)
	require.NoError(t, err)

	switch tx := deserializedTransaction.(type) {
	case FreezeTransaction:
		assert.True(t, transaction.GetStartTime().Equal(tx.GetStartTime()))
		assert.Equal(t, transaction.GetFileID(), tx.GetFileID())
		assert.Equal(t, transaction.GetFileHash(), tx.GetFileHash())
		assert.Equal(t, transaction.GetFreezeType(), tx.GetFreezeType())
	default:
		t.Fatalf("expected FreezeTransaction, got %T", deserializedTransaction)
	}
}

func TestUnitFreezeTransactionScheduleProtobuf(t *testing.T) {
	t.Parallel()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})
	startTime := time.Now().Add(time.Hour)
	endTime := time.Now().Add(2 * time.Hour)
	fileID := FileID{File: 1}
	fileHash := []byte{1, 2, 3, 4, 5}
	freezeType := FreezeTypeFreezeOnly
	nodeAccountID := []AccountID{{Account: 10}}

	tx, err := NewFreezeTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetStartTime(startTime).
		SetEndTime(endTime).
		SetFileID(fileID).
		SetFileHash(fileHash).
		SetFreezeType(freezeType).
		Freeze()
	require.NoError(t, err)

	expected := &services.SchedulableTransactionBody{
		TransactionFee: 100000000,
		Data: &services.SchedulableTransactionBody_Freeze{
			Freeze: &services.FreezeTransactionBody{
				StartTime:  _TimeToProtobuf(startTime),
				UpdateFile: fileID._ToProtobuf(),
				FileHash:   fileHash,
				FreezeType: services.FreezeType(freezeType),
			},
		},
	}

	actual, err := tx.buildScheduled()
	require.NoError(t, err)
	require.Equal(t, expected.GetFreeze(), actual.GetFreeze())
}
