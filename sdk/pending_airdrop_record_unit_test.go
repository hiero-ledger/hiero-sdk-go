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

	assert.Equal(t, original.GetPendingAirdropId(), roundTripped.GetPendingAirdropId())
	assert.Equal(t, uint64(789), roundTripped.GetPendingAirdropAmount())
	assert.Contains(t, original.String(), "789")
}
