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

func TestUnitCustomFeeLimitSetPayerId(t *testing.T) {
	t.Parallel()

	payerId := AccountID{Account: 123}

	feeLimit := NewCustomFeeLimit().
		SetPayerId(payerId)

	require.Equal(t, feeLimit.GetPayerId(), payerId)
}

func TestUnitCustomFeeLimitSetCustomFees(t *testing.T) {
	t.Parallel()

	customFee := NewCustomFixedFee().SetAmount(100)
	customFees := []*CustomFixedFee{customFee}

	feeLimit := NewCustomFeeLimit().
		SetCustomFees(customFees)

	require.Equal(t, feeLimit.GetCustomFees(), customFees)
}

func TestUnitCustomFeeLimitAddCustomFee(t *testing.T) {
	t.Parallel()

	customFee := NewCustomFixedFee().SetAmount(100)

	feeLimit := NewCustomFeeLimit().
		AddCustomFee(customFee)

	require.Contains(t, feeLimit.GetCustomFees(), customFee)
}

func TestUnitCustomFeeLimitToProtobuf(t *testing.T) {
	t.Parallel()

	payerId := AccountID{Account: 123}
	customFee := NewCustomFixedFee().SetAmount(100)
	customFees := []*CustomFixedFee{customFee}

	feeLimit := NewCustomFeeLimit().
		SetPayerId(payerId).
		SetCustomFees(customFees)

	proto := feeLimit.toProtobuf()

	require.NotNil(t, proto)
	assert.Equal(t, payerId._ToProtobuf(), proto.AccountId)
	assert.Len(t, proto.Fees, 1)
	assert.Equal(t, customFee._ToProtobuf().GetFixedFee(), proto.Fees[0])
}

func TestUnitCustomFeeLimitFromProtobuf(t *testing.T) {
	t.Parallel()

	payerId := AccountID{Account: 123}
	customFee := NewCustomFixedFee().SetAmount(100)

	proto := &services.CustomFeeLimit{
		AccountId: payerId._ToProtobuf(),
		Fees:      []*services.FixedFee{customFee._ToProtobuf().GetFixedFee()},
	}

	feeLimit := customFeeLimitFromProtobuf(proto)

	require.NotNil(t, feeLimit)
	assert.Equal(t, payerId, feeLimit.GetPayerId())
	assert.Len(t, feeLimit.GetCustomFees(), 1)
	assert.Equal(t, customFee, feeLimit.GetCustomFees()[0])
}

func TestUnitCustomFeeLimitString(t *testing.T) {
	t.Parallel()

	payerId := AccountID{Account: 123}
	customFee := NewCustomFixedFee().SetAmount(100)
	customFees := []*CustomFixedFee{customFee}

	feeLimit := NewCustomFeeLimit().
		SetPayerId(payerId).
		SetCustomFees(customFees)

	expected := "CustomFeeLimit{PayerId: " + payerId.String() + ", CustomFees: [" + customFee.String() + "]}"
	assert.Equal(t, expected, feeLimit.String())
}
