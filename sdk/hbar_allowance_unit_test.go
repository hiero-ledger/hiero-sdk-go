//go:build all || unit

package hiero

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestUnitHbarAllowanceConstructor(t *testing.T) {
	t.Parallel()

	ownerAccountID := *_MockHbarAllowance().OwnerAccountID
	spenderAccountID := *_MockHbarAllowance().SpenderAccountID
	amount := int64(100)

	allowance := NewHbarAllowance(ownerAccountID, spenderAccountID, amount)

	require.Equal(t, ownerAccountID, *allowance.OwnerAccountID)
	require.Equal(t, spenderAccountID, *allowance.SpenderAccountID)
	require.Equal(t, amount, allowance.Amount)
}

func TestUnitHbarAllowanceProtobufRoundTrip(t *testing.T) {
	t.Parallel()

	allowance := _MockHbarAllowance()

	pb := allowance._ToProtobuf()
	require.NotNil(t, pb)

	allowanceFromProtobuf := _HbarAllowanceFromProtobuf(pb)
	require.Equal(t, allowance, allowanceFromProtobuf)
}

func TestUnitHbarAllowanceProtobufRoundTripNilOwner(t *testing.T) {
	t.Parallel()

	allowance := _MockHbarAllowance()
	allowance.OwnerAccountID = nil

	pb := allowance._ToProtobuf()
	require.NotNil(t, pb)
	require.Nil(t, pb.Owner)

	allowanceFromProtobuf := _HbarAllowanceFromProtobuf(pb)
	require.Equal(t, allowance, allowanceFromProtobuf)
}

func TestUnitHbarAllowanceProtobufRoundTripNilSpender(t *testing.T) {
	t.Parallel()

	allowance := _MockHbarAllowance()
	allowance.SpenderAccountID = nil

	pb := allowance._ToProtobuf()
	require.NotNil(t, pb)
	require.Nil(t, pb.Spender)

	allowanceFromProtobuf := _HbarAllowanceFromProtobuf(pb)
	require.Equal(t, allowance, allowanceFromProtobuf)
}

func TestUnitHbarAllowanceString(t *testing.T) {
	t.Parallel()

	allowance := _MockHbarAllowance()
	require.Equal(t, "OwnerAccountID: 0.0.123, SpenderAccountID: 0.0.456, Amount: 100 tℏ", allowance.String())

	allowance.OwnerAccountID = nil
	require.Equal(t, "SpenderAccountID: 0.0.456, Amount: 100 tℏ", allowance.String())

	allowance = _MockHbarAllowance()
	allowance.SpenderAccountID = nil
	require.Equal(t, "OwnerAccountID: 0.0.123, Amount: 100 tℏ", allowance.String())

	allowance.OwnerAccountID = nil
	require.Equal(t, "", allowance.String())
}

func _MockHbarAllowance() HbarAllowance {
	ownerAccountID, _ := AccountIDFromString("0.0.123-esxsf")
	ownerAccountID.checksum = nil

	spenderAccountID, _ := AccountIDFromString("0.0.123-esxsf")
	spenderAccountID.Account = 456
	spenderAccountID.checksum = nil

	return HbarAllowance{
		OwnerAccountID:   &ownerAccountID,
		SpenderAccountID: &spenderAccountID,
		Amount:           100,
	}
}
