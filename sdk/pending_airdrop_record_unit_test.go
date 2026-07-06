//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitPendingAirdropRecordProtoRoundTrip(t *testing.T) {
	t.Parallel()

	accID, err := AccountIDFromString("0.0.123")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.456")
	require.NoError(t, err)

	original := PendingAirdropRecord{
		pendingAirdropId:     PendingAirdropId{&accID, &accID, &tokenID, nil},
		pendingAirdropAmount: 789,
	}

	roundTripped := _PendingAirdropRecordFromProtobuf(original._ToProtobuf())

	assert.Equal(t, original, roundTripped)
	assert.Equal(t, uint64(789), roundTripped.GetPendingAirdropAmount())
	assert.Equal(t, original.GetPendingAirdropId(), roundTripped.GetPendingAirdropId())
}

func TestUnitPendingAirdropRecordString(t *testing.T) {
	t.Parallel()

	accID, err := AccountIDFromString("0.0.123")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.456")
	require.NoError(t, err)

	record := PendingAirdropRecord{
		pendingAirdropId:     PendingAirdropId{&accID, &accID, &tokenID, nil},
		pendingAirdropAmount: 789,
	}

	expected := "PendingAirdropRecord{PendingAirdropId: Sender: 0.0.123, Receiver: 0.0.123, TokenID: 0.0.456, NftID: nil, PendingAirdropAmount: 789}"
	assert.Equal(t, expected, record.String())
}
