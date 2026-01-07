//go:build all || unit

package hiero

import (
	"errors"
	"net/http"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestUnitFeeEstimateQueryCoverage(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewFeeEstimateQuery().
		SetMode(FeeEstimateModeState).
		SetMaxAttempts(5)

	require.Equal(t, FeeEstimateModeState, query.GetMode())
	require.Equal(t, uint64(5), query.GetMaxAttempts())

	// Test default mode
	query2 := NewFeeEstimateQuery()
	require.Equal(t, FeeEstimateModeState, query2.GetMode())

	// Test mode setting
	query.SetMode(FeeEstimateModeIntrinsic)
	require.Equal(t, FeeEstimateModeIntrinsic, query.GetMode())

	// Test mode string representation
	require.Equal(t, "STATE", FeeEstimateModeState.String())
	require.Equal(t, "INTRINSIC", FeeEstimateModeIntrinsic.String())
	require.Equal(t, "UNKNOWN", FeeEstimateMode(99).String())
}

func TestUnitFeeEstimateQuerySetTransaction(t *testing.T) {
	t.Parallel()

	query := NewFeeEstimateQuery()
	require.Nil(t, query.GetTransaction())

	tx := NewTransferTransaction()
	query.SetTransaction(tx)
	require.NotNil(t, query.GetTransaction())
	require.Equal(t, tx, query.GetTransaction())
}

func TestUnitFeeEstimateQueryValidateNetworkOnIDs(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	query := NewFeeEstimateQuery()
	err = query.validateNetworkOnIDs(client)
	require.NoError(t, err)

	// Test with transaction
	tx := NewTransferTransaction()
	query.SetTransaction(tx)
	err = query.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitFeeEstimateQueryExecuteWithoutTransaction(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	query := NewFeeEstimateQuery()
	_, err = query.Execute(client)
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction is required")
}

func TestUnitFeeEstimateQueryExecuteWithoutClient(t *testing.T) {
	t.Parallel()

	query := NewFeeEstimateQuery()
	tx := NewTransferTransaction()
	query.SetTransaction(tx)

	_, err := query.Execute(nil)
	require.Error(t, err)
	require.Equal(t, errNoClientProvided, err)
}

func TestUnitFeeEstimateModeToProto(t *testing.T) {
	t.Parallel()

	require.Equal(t, mirror.EstimateMode_STATE, FeeEstimateModeState.toProto())
	require.Equal(t, mirror.EstimateMode_INTRINSIC, FeeEstimateModeIntrinsic.toProto())
	require.Equal(t, mirror.EstimateMode_STATE, FeeEstimateMode(99).toProto()) // Unknown defaults to STATE
}

func TestUnitFeeEstimateModeFromProto(t *testing.T) {
	t.Parallel()

	require.Equal(t, FeeEstimateModeState, feeEstimateModeFromProto(mirror.EstimateMode_STATE))
	require.Equal(t, FeeEstimateModeIntrinsic, feeEstimateModeFromProto(mirror.EstimateMode_INTRINSIC))
	require.Equal(t, FeeEstimateModeState, feeEstimateModeFromProto(mirror.EstimateMode(99))) // Unknown defaults to STATE
}

func TestUnitFeeExtraFromProto(t *testing.T) {
	t.Parallel()

	// Test nil
	extra := feeExtraFromProto(nil)
	require.Equal(t, FeeExtra{}, extra)

	// Test with values
	pbExtra := &mirror.FeeExtra{
		Name:       "test",
		Included:   5,
		Count:      10,
		Charged:    5,
		FeePerUnit: 100,
		Subtotal:   500,
	}
	extra = feeExtraFromProto(pbExtra)
	require.Equal(t, "test", extra.Name)
	require.Equal(t, uint32(5), extra.Included)
	require.Equal(t, uint32(10), extra.Count)
	require.Equal(t, uint32(5), extra.Charged)
	require.Equal(t, uint64(100), extra.FeePerUnit)
	require.Equal(t, uint64(500), extra.Subtotal)
}

func TestUnitFeeEstimateFromProto(t *testing.T) {
	t.Parallel()

	// Test nil
	estimate := feeEstimateFromProto(nil)
	require.Equal(t, FeeEstimate{}, estimate)

	// Test with values
	pbEstimate := &mirror.FeeEstimate{
		Base: 1000,
		Extras: []*mirror.FeeExtra{
			{
				Name:     "extra1",
				Subtotal: 500,
			},
			{
				Name:     "extra2",
				Subtotal: 300,
			},
		},
	}
	estimate = feeEstimateFromProto(pbEstimate)
	require.Equal(t, uint64(1000), estimate.Base)
	require.Len(t, estimate.Extras, 2)
	require.Equal(t, "extra1", estimate.Extras[0].Name)
	require.Equal(t, "extra2", estimate.Extras[1].Name)

	// Test Subtotal calculation
	require.Equal(t, uint64(1800), estimate.Subtotal())
}

func TestUnitNetworkFeeFromProto(t *testing.T) {
	t.Parallel()

	// Test nil
	networkFee := networkFeeFromProto(nil)
	require.Equal(t, NetworkFee{}, networkFee)

	// Test with values
	pbNetworkFee := &mirror.NetworkFee{
		Multiplier: 3,
		Subtotal:   3000,
	}
	networkFee = networkFeeFromProto(pbNetworkFee)
	require.Equal(t, uint32(3), networkFee.Multiplier)
	require.Equal(t, uint64(3000), networkFee.Subtotal)
}

func TestUnitFeeEstimateResponseFromProto(t *testing.T) {
	t.Parallel()

	// Test nil
	response := feeEstimateResponseFromProto(nil)
	require.Equal(t, FeeEstimateResponse{}, response)

	// Test with values
	pbResponse := &mirror.FeeEstimateResponse{
		Mode: mirror.EstimateMode_STATE,
		Network: &mirror.NetworkFee{
			Multiplier: 3,
			Subtotal:   3000,
		},
		Node: &mirror.FeeEstimate{
			Base: 1000,
		},
		Service: &mirror.FeeEstimate{
			Base: 500,
		},
		Notes: []string{"note1", "note2"},
		Total: 4500,
	}
	response = feeEstimateResponseFromProto(pbResponse)
	require.Equal(t, FeeEstimateModeState, response.Mode)
	require.Equal(t, uint32(3), response.NetworkFee.Multiplier)
	require.Equal(t, uint64(1000), response.NodeFee.Base)
	require.Equal(t, uint64(500), response.ServiceFee.Base)
	require.Len(t, response.Notes, 2)
	require.Equal(t, "note1", response.Notes[0])
	require.Equal(t, "note2", response.Notes[1])
	require.Equal(t, uint64(4500), response.Total)
}

func TestUnitFeeEstimateQueryShouldRetry(t *testing.T) {
	t.Parallel()

	query := NewFeeEstimateQuery()

	// Test nil error and nil response
	require.False(t, query.shouldRetry(nil, nil))

	// Test nil error with 200 response (should not retry)
	resp200 := &http.Response{StatusCode: http.StatusOK}
	require.False(t, query.shouldRetry(nil, resp200))

	// Test nil error with 500 response (should retry)
	resp500 := &http.Response{StatusCode: http.StatusInternalServerError}
	require.True(t, query.shouldRetry(nil, resp500))

	// Test nil error with 503 response (should retry)
	resp503 := &http.Response{StatusCode: http.StatusServiceUnavailable}
	require.True(t, query.shouldRetry(nil, resp503))

	// Test nil error with 429 response (should retry)
	resp429 := &http.Response{StatusCode: http.StatusTooManyRequests}
	require.True(t, query.shouldRetry(nil, resp429))

	// Test nil error with 400 response (should not retry)
	resp400 := &http.Response{StatusCode: http.StatusBadRequest}
	require.False(t, query.shouldRetry(nil, resp400))

	// Test nil error with 404 response (should not retry)
	resp404 := &http.Response{StatusCode: http.StatusNotFound}
	require.False(t, query.shouldRetry(nil, resp404))

	// Test network error (should retry)
	err := errors.New("connection refused")
	require.True(t, query.shouldRetry(err, nil))

	// Test any error (should retry)
	err = errors.New("timeout")
	require.True(t, query.shouldRetry(err, nil))
}

func TestUnitFeeEstimateModeFromString(t *testing.T) {
	t.Parallel()

	require.Equal(t, FeeEstimateModeState, feeEstimateModeFromString("STATE"))
	require.Equal(t, FeeEstimateModeIntrinsic, feeEstimateModeFromString("INTRINSIC"))
	require.Equal(t, FeeEstimateModeState, feeEstimateModeFromString("UNKNOWN")) // Unknown defaults to STATE
	require.Equal(t, FeeEstimateModeState, feeEstimateModeFromString(""))        // Empty defaults to STATE
}

func TestUnitFeeEstimateResponseFromREST(t *testing.T) {
	t.Parallel()

	// Test valid JSON response
	jsonData := `{
		"mode": "STATE",
		"network": {
			"multiplier": 3,
			"subtotal": 3000
		},
		"node": {
			"base": 1000,
			"extras": [
				{
					"name": "extra1",
					"included": 0,
					"count": 1,
					"charged": 1,
					"feePerUnit": 100,
					"subtotal": 100
				}
			]
		},
		"service": {
			"base": 500
		},
		"notes": ["note1", "note2"],
		"total": 4600
	}`

	response, err := feeEstimateResponseFromREST([]byte(jsonData))
	require.NoError(t, err)
	require.Equal(t, FeeEstimateModeState, response.Mode)
	require.Equal(t, uint32(3), response.NetworkFee.Multiplier)
	require.Equal(t, uint64(3000), response.NetworkFee.Subtotal)
	require.Equal(t, uint64(1000), response.NodeFee.Base)
	require.Len(t, response.NodeFee.Extras, 1)
	require.Equal(t, "extra1", response.NodeFee.Extras[0].Name)
	require.Equal(t, uint64(500), response.ServiceFee.Base)
	require.Len(t, response.Notes, 2)
	require.Equal(t, "note1", response.Notes[0])
	require.Equal(t, "note2", response.Notes[1])
	require.Equal(t, uint64(4600), response.Total)

	// Test INTRINSIC mode
	jsonData2 := `{
		"mode": "INTRINSIC",
		"network": {
			"multiplier": 2,
			"subtotal": 2000
		},
		"node": {
			"base": 500
		},
		"service": {
			"base": 250
		},
		"total": 2750
	}`

	response2, err := feeEstimateResponseFromREST([]byte(jsonData2))
	require.NoError(t, err)
	require.Equal(t, FeeEstimateModeIntrinsic, response2.Mode)
	require.Equal(t, uint32(2), response2.NetworkFee.Multiplier)
	require.Empty(t, response2.Notes) // Empty notes should be empty slice

	// Test invalid JSON
	_, err = feeEstimateResponseFromREST([]byte("invalid json"))
	require.Error(t, err)
}
