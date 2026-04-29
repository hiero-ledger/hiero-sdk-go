//go:build all || unit

package hiero

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

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

	query2 := NewFeeEstimateQuery()
	require.Equal(t, FeeEstimateModeIntrinsic, query2.GetMode())
	require.Equal(t, uint16(0), query2.GetHighVolumeThrottle())

	query.SetMode(FeeEstimateModeIntrinsic)
	require.Equal(t, FeeEstimateModeIntrinsic, query.GetMode())

	query.SetHighVolumeThrottle(5000)
	require.Equal(t, uint16(5000), query.GetHighVolumeThrottle())

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

func TestUnitFeeEstimateQueryShouldRetry(t *testing.T) {
	t.Parallel()

	query := NewFeeEstimateQuery()

	require.False(t, query.shouldRetry(nil, nil))

	resp200 := &http.Response{StatusCode: http.StatusOK}
	require.False(t, query.shouldRetry(nil, resp200))

	resp500 := &http.Response{StatusCode: http.StatusInternalServerError}
	require.True(t, query.shouldRetry(nil, resp500))

	resp503 := &http.Response{StatusCode: http.StatusServiceUnavailable}
	require.True(t, query.shouldRetry(nil, resp503))

	resp429 := &http.Response{StatusCode: http.StatusTooManyRequests}
	require.True(t, query.shouldRetry(nil, resp429))

	resp400 := &http.Response{StatusCode: http.StatusBadRequest}
	require.False(t, query.shouldRetry(nil, resp400))

	resp404 := &http.Response{StatusCode: http.StatusNotFound}
	require.False(t, query.shouldRetry(nil, resp404))

	err := errors.New("connection refused")
	require.True(t, query.shouldRetry(err, nil))

	err = errors.New("timeout")
	require.True(t, query.shouldRetry(err, nil))
}

func TestUnitFeeEstimateResponseFromREST(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"high_volume_multiplier": 2,
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
					"fee_per_unit": 100,
					"subtotal": 100
				}
			]
		},
		"service": {
			"base": 500
		},
		"total": 4600
	}`

	var response FeeEstimateResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	require.NoError(t, err)
	require.Equal(t, uint64(2), response.HighVolumeMultiplier)
	require.Equal(t, uint32(3), response.NetworkFee.Multiplier)
	require.Equal(t, uint64(3000), response.NetworkFee.Subtotal)
	require.Equal(t, uint64(1000), response.NodeFee.Base)
	require.Len(t, response.NodeFee.Extras, 1)
	require.Equal(t, "extra1", response.NodeFee.Extras[0].Name)
	require.Equal(t, uint64(1), response.NodeFee.Extras[0].Count)
	require.Equal(t, uint64(1), response.NodeFee.Extras[0].Charged)
	require.Equal(t, uint64(100), response.NodeFee.Extras[0].FeePerUnit)
	require.Equal(t, uint64(500), response.ServiceFee.Base)
	require.Equal(t, uint64(4600), response.Total)

	jsonData2 := `{
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

	var response2 FeeEstimateResponse
	err = json.Unmarshal([]byte(jsonData2), &response2)
	require.NoError(t, err)
	require.Equal(t, uint32(2), response2.NetworkFee.Multiplier)
	require.Equal(t, uint64(0), response2.HighVolumeMultiplier)

	var response3 FeeEstimateResponse
	err = json.Unmarshal([]byte("invalid json"), &response3)
	require.Error(t, err)
}

func TestUnitFeeEstimateResponseTotalFormula(t *testing.T) {
	t.Parallel()

	response := FeeEstimateResponse{
		NetworkFee: NetworkFee{Multiplier: 3},
		NodeFee: FeeEstimate{
			Base: 1000,
			Extras: []FeeExtra{
				{Subtotal: 50},
				{Subtotal: 500},
			},
		},
		ServiceFee: FeeEstimate{
			Base:   200,
			Extras: []FeeExtra{{Subtotal: 75}},
		},
	}

	nodeSubtotal := response.NodeFee.Subtotal()
	serviceSubtotal := response.ServiceFee.Subtotal()
	response.NetworkFee.Subtotal = nodeSubtotal * uint64(response.NetworkFee.Multiplier)
	response.Total = response.NetworkFee.Subtotal + nodeSubtotal + serviceSubtotal

	require.Equal(t, uint64(1550), nodeSubtotal)
	require.Equal(t, uint64(275), serviceSubtotal)
	require.Equal(t, uint64(4650), response.NetworkFee.Subtotal)
	require.Equal(t, response.NetworkFee.Subtotal+nodeSubtotal+serviceSubtotal, response.Total)
}

func TestUnitTransactionEstimateFee(t *testing.T) {
	t.Parallel()

	tx := NewTransferTransaction()
	query := tx.EstimateFee()

	require.NotNil(t, query)
	require.Equal(t, FeeEstimateModeIntrinsic, query.GetMode())
	require.NotNil(t, query.GetTransaction())
	require.Equal(t, tx, query.GetTransaction())
}
