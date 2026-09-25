package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// FeeEstimateQuery allows users to query expected transaction fees without submitting transactions to the network
type FeeEstimateQuery struct {
	mode               FeeEstimateMode
	transaction        TransactionInterface
	highVolumeThrottle uint16
	maxAttempts        uint64
}

// NewFeeEstimateQuery creates a new FeeEstimateQuery
func NewFeeEstimateQuery() *FeeEstimateQuery {
	return &FeeEstimateQuery{
		mode: FeeEstimateModeIntrinsic, // Default mode is INTRINSIC
	}
}

// SetMode sets the estimation mode (optional, defaults to INTRINSIC)
func (q *FeeEstimateQuery) SetMode(mode FeeEstimateMode) *FeeEstimateQuery {
	q.mode = mode
	return q
}

// GetMode returns the current estimation mode
func (q *FeeEstimateQuery) GetMode() FeeEstimateMode {
	return q.mode
}

// SetTransaction sets the transaction to estimate (required)
func (q *FeeEstimateQuery) SetTransaction(transaction TransactionInterface) *FeeEstimateQuery {
	q.transaction = transaction
	return q
}

// GetTransaction returns the current transaction
func (q *FeeEstimateQuery) GetTransaction() TransactionInterface {
	return q.transaction
}

// SetHighVolumeThrottle sets the high-volume throttle utilization in basis points (0–10000, where 10000 = 100%).
// A value of 0 (the default) indicates no high-volume pricing simulation. Maps to the
// `high_volume_throttle` query parameter
func (q *FeeEstimateQuery) SetHighVolumeThrottle(throttle uint16) *FeeEstimateQuery {
	q.highVolumeThrottle = throttle
	return q
}

// GetHighVolumeThrottle returns the current high-volume throttle utilization in basis points
func (q *FeeEstimateQuery) GetHighVolumeThrottle() uint16 {
	return q.highVolumeThrottle
}

// SetMaxAttempts sets the maximum number of retry attempts
func (q *FeeEstimateQuery) SetMaxAttempts(maxAttempts uint64) *FeeEstimateQuery {
	q.maxAttempts = maxAttempts
	return q
}

// GetMaxAttempts returns the maximum number of retry attempts, or 0 if unset.
func (q *FeeEstimateQuery) GetMaxAttempts() uint64 {
	return q.maxAttempts
}

// Execute executes the fee estimation query with the provided client
func (q *FeeEstimateQuery) Execute(client *Client) (FeeEstimateResponse, error) {
	if client == nil {
		return FeeEstimateResponse{}, errNoClientProvided
	}

	if q.transaction == nil {
		return FeeEstimateResponse{}, errors.New("transaction is required")
	}

	err := q.validateNetworkOnIDs(client)
	if err != nil {
		return FeeEstimateResponse{}, err
	}

	baseTx := q.transaction.getBaseTransaction()
	if !baseTx.IsFrozen() {
		_, err := baseTx.FreezeWith(client)
		if err != nil {
			return FeeEstimateResponse{}, errors.Wrap(err, "failed to freeze transaction")
		}
	}

	path, err := q.buildPath()
	if err != nil {
		return FeeEstimateResponse{}, err
	}
	// One client for all chunks, so they go to the same mirror node and share one deadline.
	restClient, err := client.mirrorRestClient(path, client.mirrorHttpPolicyForQuery(q.maxAttempts))
	if err != nil {
		return FeeEstimateResponse{}, err
	}

	if fileAppendTx, ok := q.transaction.(*FileAppendTransaction); ok {
		return q.executeChunkedTransaction(restClient, path, fileAppendTx)
	}

	if topicMessageTx, ok := q.transaction.(*TopicMessageSubmitTransaction); ok {
		return q.executeChunkedTransaction(restClient, path, topicMessageTx)
	}

	return q.estimateSingleTransaction(restClient, path, q.transaction)
}

// executeChunkedTransaction handles fee estimation for chunked transactions
func (q *FeeEstimateQuery) executeChunkedTransaction(restClient *mirrorNodeHttpClient, path mirrorNodeRestPath, tx TransactionInterface) (FeeEstimateResponse, error) {
	baseTx := tx.getBaseTransaction()
	numChunks := baseTx.signedTransactions._Length() / baseTx.nodeAccountIDs._Length()
	if numChunks == 0 {
		return FeeEstimateResponse{}, errors.New("transaction has no chunks")
	}

	var aggregatedResponse FeeEstimateResponse
	aggregatedResponse.NodeFee = FeeEstimate{Base: 0, Extras: []FeeExtra{}}
	aggregatedResponse.ServiceFee = FeeEstimate{Base: 0, Extras: []FeeExtra{}}
	aggregatedResponse.NetworkFee = NetworkFee{Multiplier: 0, Subtotal: 0}

	var totalNodeSubtotal uint64
	var totalServiceSubtotal uint64

	// Estimate fees for each chunk
	for i := 0; i < numChunks; i++ {
		chunkTx, err := baseTx._BuildTransaction(i)
		if err != nil {
			return FeeEstimateResponse{}, errors.Wrapf(err, "failed to build chunk %d", i)
		}

		chunkResponse, err := q.callGetFeeEstimate(restClient, path, chunkTx)
		if err != nil {
			return FeeEstimateResponse{}, errors.Wrapf(err, "failed to estimate chunk %d", i)
		}

		totalNodeSubtotal += chunkResponse.NodeFee.Subtotal()
		totalServiceSubtotal += chunkResponse.ServiceFee.Subtotal()

		if i == 0 {
			aggregatedResponse.NetworkFee.Multiplier = chunkResponse.NetworkFee.Multiplier
			aggregatedResponse.HighVolumeMultiplier = chunkResponse.HighVolumeMultiplier
		}
	}

	aggregatedResponse.NodeFee.Base = totalNodeSubtotal
	aggregatedResponse.ServiceFee.Base = totalServiceSubtotal
	aggregatedResponse.NetworkFee.Subtotal = totalNodeSubtotal * uint64(aggregatedResponse.NetworkFee.Multiplier)
	aggregatedResponse.Total = aggregatedResponse.NetworkFee.Subtotal + totalNodeSubtotal + totalServiceSubtotal

	return aggregatedResponse, nil
}

// estimateSingleTransaction estimates fees for a single transaction
func (q *FeeEstimateQuery) estimateSingleTransaction(restClient *mirrorNodeHttpClient, path mirrorNodeRestPath, tx TransactionInterface) (FeeEstimateResponse, error) {
	baseTx := tx.getBaseTransaction()

	protoTx, err := baseTx._BuildTransaction(0)
	if err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to build transaction")
	}

	return q.callGetFeeEstimate(restClient, path, protoTx)
}

// callGetFeeEstimate POSTs one transaction to the fee estimate endpoint.
func (q *FeeEstimateQuery) callGetFeeEstimate(restClient *mirrorNodeHttpClient, path mirrorNodeRestPath, protoTx *services.Transaction) (FeeEstimateResponse, error) {
	txBytes, err := protobuf.Marshal(protoTx)
	if err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to marshal transaction")
	}

	resp, err := restClient.post(path, "application/protobuf", txBytes, CancellationNone())
	switch {
	case errors.Is(err, errMirrorHttpRetriesExhausted):
		return FeeEstimateResponse{}, errors.Wrap(mirrorNodeStatusError(resp), "failed to call fee estimate API")
	case err != nil:
		var transportErr *HttpTransportError
		if errors.As(err, &transportErr) && transportErr.IsRetryable() {
			return FeeEstimateResponse{}, errors.Wrapf(err, "failed to call fee estimate API after %d attempts", restClient.policy.GetMaxAttempts())
		}
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to call fee estimate API")
	case resp.statusCode != http.StatusOK:
		return FeeEstimateResponse{}, errors.Wrap(mirrorNodeStatusError(resp), "failed to call fee estimate API")
	}

	var response FeeEstimateResponse
	if err := json.Unmarshal(resp.body, &response); err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to unmarshal response")
	}

	return response, nil
}

// buildPath returns the fee estimate path for the query's mode and throttle.
func (q *FeeEstimateQuery) buildPath() (mirrorNodeRestPath, error) {
	path := "/network/fees?mode=" + q.mode.String()
	if q.highVolumeThrottle != 0 {
		path = fmt.Sprintf("%s&high_volume_throttle=%d", path, q.highVolumeThrottle)
	}

	return newMirrorNodeRestPath(path)
}

// validateNetworkOnIDs validates network and IDs on the query
func (q *FeeEstimateQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.transaction != nil {
		return q.transaction.validateNetworkOnIDs(client)
	}

	return nil
}
