//go:build all || unit

package hiero

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	require.Equal(t, FeeEstimateModeState, query2.GetMode())

	query.SetMode(FeeEstimateModeIntrinsic)
	require.Equal(t, FeeEstimateModeIntrinsic, query.GetMode())

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

func TestUnitFeeEstimateModeFromString(t *testing.T) {
	t.Parallel()

	require.Equal(t, FeeEstimateModeState, feeEstimateModeFromString("STATE"))
	require.Equal(t, FeeEstimateModeIntrinsic, feeEstimateModeFromString("INTRINSIC"))
	require.Equal(t, FeeEstimateModeState, feeEstimateModeFromString("UNKNOWN"))
	require.Equal(t, FeeEstimateModeState, feeEstimateModeFromString(""))
}

func TestUnitFeeEstimateResponseFromREST(t *testing.T) {
	t.Parallel()

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
	require.Empty(t, response2.Notes)
	_, err = feeEstimateResponseFromREST([]byte("invalid json"))
	require.Error(t, err)
}

func TestUnitFeeEstimateQueryRetriesOnUnavailableErrors(t *testing.T) {
	// Note: Not using t.Parallel() because these tests need to bind to port 5551
	// which is hardcoded for localhost in the mirror node code

	stub := newStubMirrorRestServer(t)
	defer stub.stop()

	client := newMockClientForREST(stub.getMirrorNetworkAddress())
	client.SetMaxBackoff(500 * time.Millisecond)

	tx := NewTransferTransaction()

	query := NewFeeEstimateQuery().
		SetTransaction(tx).
		SetMaxAttempts(3)

	stub.enqueue(stubResponse{status: 503, body: "transient error"})
	stub.enqueue(stubResponse{status: 200, body: newSuccessResponse(FeeEstimateModeState, 2, 6, 8)})

	response, err := query.Execute(client)
	require.NoError(t, err)

	assert.Equal(t, FeeEstimateModeState, response.Mode)
	assert.Equal(t, uint64(26), response.Total)
	assert.Equal(t, 2, stub.requestCount())
	stub.verify(t)
}

func TestUnitFeeEstimateQueryRetriesOnDeadlineExceededErrors(t *testing.T) {
	// Note: Not using t.Parallel() because these tests need to bind to port 5551
	// which is hardcoded for localhost in the mirror node code

	stub := newStubMirrorRestServer(t)
	defer stub.stop()

	client := newMockClientForREST(stub.getMirrorNetworkAddress())
	client.SetMaxBackoff(500 * time.Millisecond)

	tx := NewTransferTransaction()

	query := NewFeeEstimateQuery().
		SetTransaction(tx).
		SetMaxAttempts(3)

	stub.enqueue(stubResponse{status: 504, body: "gateway timeout"})
	stub.enqueue(stubResponse{status: 200, body: newSuccessResponse(FeeEstimateModeState, 4, 8, 20)})

	response, err := query.Execute(client)
	require.NoError(t, err)

	assert.Equal(t, FeeEstimateModeState, response.Mode)
	assert.Equal(t, uint64(60), response.Total)
	assert.Equal(t, 2, stub.requestCount())
	stub.verify(t)
}

func TestUnitFeeEstimateQuerySucceedsOnFirstAttempt(t *testing.T) {
	// Note: Not using t.Parallel() because these tests need to bind to port 5551
	// which is hardcoded for localhost in the mirror node code

	stub := newStubMirrorRestServer(t)
	defer stub.stop()

	client := newMockClientForREST(stub.getMirrorNetworkAddress())

	tx := NewTransferTransaction()

	query := NewFeeEstimateQuery().
		SetTransaction(tx).
		SetMode(FeeEstimateModeIntrinsic)

	stub.enqueue(stubResponse{status: 200, body: newSuccessResponse(FeeEstimateModeIntrinsic, 3, 10, 20)})

	response, err := query.Execute(client)
	require.NoError(t, err)

	assert.Equal(t, FeeEstimateModeIntrinsic, response.Mode)
	assert.Equal(t, uint64(60), response.Total)
	assert.Equal(t, 1, stub.requestCount())
	stub.verify(t)
}
