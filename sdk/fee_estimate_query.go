package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"io"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FeeEstimateQuery allows users to query expected transaction fees without submitting transactions to the network
type FeeEstimateQuery struct {
	mode        FeeEstimateMode
	transaction TransactionInterface
	attempt     uint64
	maxAttempts uint64
}

// NewFeeEstimateQuery creates a new FeeEstimateQuery
func NewFeeEstimateQuery() *FeeEstimateQuery {
	return &FeeEstimateQuery{
		mode:        FeeEstimateModeState, // Default mode is STATE
		maxAttempts: maxAttempts,
	}
}

// SetMode sets the estimation mode (optional, defaults to STATE)
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

// SetMaxAttempts sets the maximum number of retry attempts
func (q *FeeEstimateQuery) SetMaxAttempts(maxAttempts uint64) *FeeEstimateQuery {
	q.maxAttempts = maxAttempts
	return q
}

// GetMaxAttempts returns the maximum number of retry attempts
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

	// Automatically freeze the transaction if not already frozen
	baseTx := q.transaction.getBaseTransaction()
	if !baseTx.IsFrozen() {
		_, err := baseTx.FreezeWith(client)
		if err != nil {
			return FeeEstimateResponse{}, errors.Wrap(err, "failed to freeze transaction")
		}
	}

	// Check if this is a chunked transaction
	if fileAppendTx, ok := q.transaction.(*FileAppendTransaction); ok {
		return q.executeChunkedTransaction(client, fileAppendTx)
	}

	if topicMessageTx, ok := q.transaction.(*TopicMessageSubmitTransaction); ok {
		return q.executeChunkedTransaction(client, topicMessageTx)
	}

	// For non-chunked transactions, estimate fees directly
	return q.estimateSingleTransaction(client, q.transaction)
}

// executeChunkedTransaction handles fee estimation for chunked transactions
func (q *FeeEstimateQuery) executeChunkedTransaction(client *Client, tx TransactionInterface) (FeeEstimateResponse, error) {
	baseTx := tx.getBaseTransaction()
	numChunks := baseTx.signedTransactions._Length() / baseTx.nodeAccountIDs._Length()
	if numChunks == 0 {
		return FeeEstimateResponse{}, errors.New("transaction has no chunks")
	}

	// Aggregate responses from all chunks
	var aggregatedResponse FeeEstimateResponse
	aggregatedResponse.Mode = q.mode
	aggregatedResponse.NodeFee = FeeEstimate{Base: 0, Extras: []FeeExtra{}}
	aggregatedResponse.ServiceFee = FeeEstimate{Base: 0, Extras: []FeeExtra{}}
	aggregatedResponse.NetworkFee = NetworkFee{Multiplier: 0, Subtotal: 0}
	aggregatedResponse.Notes = []string{}

	var totalNodeSubtotal uint64
	var totalServiceSubtotal uint64

	// Estimate fees for each chunk
	for i := 0; i < numChunks; i++ {
		// Build the transaction for this specific chunk
		chunkTx, err := baseTx._BuildTransaction(i)
		if err != nil {
			return FeeEstimateResponse{}, errors.Wrapf(err, "failed to build chunk %d", i)
		}

		// Create a temporary transaction list with just this chunk
		chunkResponse, err := q.callGetFeeEstimate(client, chunkTx)
		if err != nil {
			return FeeEstimateResponse{}, errors.Wrapf(err, "failed to estimate chunk %d", i)
		}

		// Aggregate fees
		totalNodeSubtotal += chunkResponse.NodeFee.Subtotal()
		totalServiceSubtotal += chunkResponse.ServiceFee.Subtotal()

		// Store network multiplier from first chunk (should be same for all)
		if i == 0 {
			aggregatedResponse.NetworkFee.Multiplier = chunkResponse.NetworkFee.Multiplier
		}

		// Aggregate notes
		aggregatedResponse.Notes = append(aggregatedResponse.Notes, chunkResponse.Notes...)
	}

	// Calculate aggregated totals
	aggregatedResponse.NodeFee.Base = totalNodeSubtotal
	aggregatedResponse.ServiceFee.Base = totalServiceSubtotal
	aggregatedResponse.NetworkFee.Subtotal = totalNodeSubtotal * uint64(aggregatedResponse.NetworkFee.Multiplier)
	aggregatedResponse.Total = aggregatedResponse.NetworkFee.Subtotal + totalNodeSubtotal + totalServiceSubtotal

	return aggregatedResponse, nil
}

// estimateSingleTransaction estimates fees for a single transaction
func (q *FeeEstimateQuery) estimateSingleTransaction(client *Client, tx TransactionInterface) (FeeEstimateResponse, error) {
	baseTx := tx.getBaseTransaction()

	// Build the transaction (use first chunk if multiple exist)
	protoTx, err := baseTx._BuildTransaction(0)
	if err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to build transaction")
	}

	return q.callGetFeeEstimate(client, protoTx)
}

// callGetFeeEstimate calls the GetFeeEstimate gRPC method
func (q *FeeEstimateQuery) callGetFeeEstimate(client *Client, protoTx *services.Transaction) (FeeEstimateResponse, error) {
	// Get mirror node client
	mirrorNode, err := client.mirrorNetwork._GetNextMirrorNode()
	if err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to get mirror node")
	}

	channel, err := mirrorNode._GetNetworkServiceClient()
	if err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to get network service client")
	}

	ctx := client.networkUpdateContext

	// Build the fee estimate query
	query := &mirror.FeeEstimateQuery{
		Mode:        q.mode.toProto(),
		Transaction: protoTx,
	}

	// Call the mirror node gRPC endpoint with retry logic
	var resp *mirror.FeeEstimateResponse
	var lastErr error

	for q.attempt < q.maxAttempts {
		resp, err = channel.GetFeeEstimate(ctx, query)
		if err == nil {
			break
		}

		lastErr = err

		// Check if we should retry
		if !q.shouldRetry(err) {
			return FeeEstimateResponse{}, errors.Wrap(err, "failed to call GetFeeEstimate")
		}

		// Calculate delay with exponential backoff
		delayMs := 250.0 * float64(uint64(1)<<q.attempt) // 250ms, 500ms, 1000ms, etc.
		if delayMs > 8000 {
			delayMs = 8000
		}

		// Wait before retry
		select {
		case <-ctx.Done():
			return FeeEstimateResponse{}, ctx.Err()
		case <-time.After(time.Duration(delayMs) * time.Millisecond):
		}

		q.attempt++
	}

	if lastErr != nil {
		return FeeEstimateResponse{}, errors.Wrapf(lastErr, "failed to call GetFeeEstimate after %d attempts", q.maxAttempts)
	}

	// Convert proto response to SDK response
	return feeEstimateResponseFromProto(resp), nil
}

// shouldRetry determines if an error should be retried
func (q *FeeEstimateQuery) shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	if err == io.EOF {
		return false
	}

	grpcStatus, ok := status.FromError(err)
	if !ok {
		return false
	}

	code := grpcStatus.Code()

	// Retry on transient transport errors
	switch code {
	case codes.Unavailable, codes.DeadlineExceeded:
		return true
	case codes.InvalidArgument:
		// Do not retry on malformed transaction
		return false
	default:
		return false
	}
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
