//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"testing"

	"github.com/hashgraph/hedera-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitTokenCancelAirdropTransactionSetPendingAirdropIds(t *testing.T) {
	t.Parallel()

	pendingAirdropId1 := &PendingAirdropId{tokenID: &TokenID{Token: 1}}
	pendingAirdropId2 := &PendingAirdropId{tokenID: &TokenID{Token: 2}}

	transaction := NewTokenCancelAirdropTransaction().
		SetPendingAirdropIds([]*PendingAirdropId{pendingAirdropId1, pendingAirdropId2})

	assert.Equal(t, []*PendingAirdropId{pendingAirdropId1, pendingAirdropId2}, transaction.GetPendingAirdropIds())
}

func TestUnitTokenCancelAirdropTransactionAddPendingAirdropId(t *testing.T) {
	t.Parallel()

	pendingAirdropId1 := PendingAirdropId{tokenID: &TokenID{Token: 1}}
	pendingAirdropId2 := PendingAirdropId{tokenID: &TokenID{Token: 2}}

	transaction := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(pendingAirdropId1).
		AddPendingAirdropId(pendingAirdropId2)

	assert.Equal(t, []*PendingAirdropId{&pendingAirdropId1, &pendingAirdropId2}, transaction.GetPendingAirdropIds())
}

func TestUnitTokenCancelAirdropTransactionFreeze(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}}

	pendingAirdropId := PendingAirdropId{tokenID: &TokenID{Token: 1}, sender: &AccountID{Account: 3}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})
	transaction := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(pendingAirdropId).
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID)

	_, err := transaction.Freeze()
	require.NoError(t, err)
}

func TestUnitTokenCancelAirdropTransactionToBytes(t *testing.T) {
	t.Parallel()

	pendingAirdropId := PendingAirdropId{tokenID: &TokenID{Token: 1}}

	transaction := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(pendingAirdropId)

	bytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.NotNil(t, bytes)
}

func TestUnitTokenCancelAirdropTransactionFromBytes(t *testing.T) {
	t.Parallel()

	pendingAirdropId := PendingAirdropId{tokenID: &TokenID{Token: 1}, sender: &AccountID{Account: 3}, receiver: &AccountID{Account: 4}}

	transaction := NewTokenCancelAirdropTransaction().
		AddPendingAirdropId(pendingAirdropId)

	bytes, err := transaction.ToBytes()
	require.NoError(t, err)
	require.NotNil(t, bytes)

	deserializedTransaction, err := TransactionFromBytes(bytes)
	require.NoError(t, err)

	switch tx := deserializedTransaction.(type) {
	case TokenCancelAirdropTransaction:
		assert.Equal(t, transaction.GetPendingAirdropIds(), tx.GetPendingAirdropIds())
	default:
		t.Fatalf("expected TokenCancelAirdropTransaction, got %T", deserializedTransaction)
	}
}

func TestUnitTokenCancelAirdropTransactionScheduleProtobuf(t *testing.T) {
	t.Parallel()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	pendingAirdropId1 := PendingAirdropId{tokenID: &TokenID{Token: 1}}
	pendingAirdropId2 := PendingAirdropId{tokenID: &TokenID{Token: 2}}

	tx, err := NewTokenCancelAirdropTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		AddPendingAirdropId(pendingAirdropId1).
		AddPendingAirdropId(pendingAirdropId2).
		Freeze()
	require.NoError(t, err)

	expected := &services.SchedulableTransactionBody{
		TransactionFee: 100000000,
		Data: &services.SchedulableTransactionBody_TokenCancelAirdrop{
			TokenCancelAirdrop: &services.TokenCancelAirdropTransactionBody{
				PendingAirdrops: []*services.PendingAirdropId{
					pendingAirdropId1._ToProtobuf(),
					pendingAirdropId2._ToProtobuf(),
				},
			},
		},
	}

	actual, err := tx.buildScheduled()
	require.NoError(t, err)
	require.Equal(t, expected.String(), actual.String())
}

func TestUnitTokenCancelAirdropTransactionValidateNetworkOnIDs(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	checksum := "dmqui"
	pendingAirdropId := &PendingAirdropId{
		tokenID:  &TokenID{Token: 3, checksum: &checksum},
		sender:   &AccountID{Account: 3, checksum: &checksum},
		receiver: &AccountID{Account: 3, checksum: &checksum},
	}

	transaction := NewTokenCancelAirdropTransaction().
		SetPendingAirdropIds([]*PendingAirdropId{pendingAirdropId})

	err = transaction.validateNetworkOnIDs(client)
	require.NoError(t, err)
}
