package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
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

	baseTx := q.transaction.getBaseTransaction()
	if !baseTx.IsFrozen() {
		_, err := baseTx.FreezeWith(client)
		if err != nil {
			return FeeEstimateResponse{}, errors.Wrap(err, "failed to freeze transaction")
		}
	}

	if fileAppendTx, ok := q.transaction.(*FileAppendTransaction); ok {
		return q.executeChunkedTransaction(client, fileAppendTx)
	}

	if topicMessageTx, ok := q.transaction.(*TopicMessageSubmitTransaction); ok {
		return q.executeChunkedTransaction(client, topicMessageTx)
	}

	return q.estimateSingleTransaction(client, q.transaction)
}

// executeChunkedTransaction handles fee estimation for chunked transactions
func (q *FeeEstimateQuery) executeChunkedTransaction(client *Client, tx TransactionInterface) (FeeEstimateResponse, error) {
	baseTx := tx.getBaseTransaction()
	numChunks := baseTx.signedTransactions._Length() / baseTx.nodeAccountIDs._Length()
	if numChunks == 0 {
		return FeeEstimateResponse{}, errors.New("transaction has no chunks")
	}

	var aggregatedResponse FeeEstimateResponse
	aggregatedResponse.NodeFee = FeeEstimate{Base: 0, Extras: []FeeExtra{}}
	aggregatedResponse.ServiceFee = FeeEstimate{Base: 0, Extras: []FeeExtra{}}
	aggregatedResponse.NetworkFee = NetworkFee{Multiplier: 0, Subtotal: 0}
	aggregatedResponse.Notes = []string{}

	var totalNodeSubtotal uint64
	var totalServiceSubtotal uint64

	// Estimate fees for each chunk
	for i := 0; i < numChunks; i++ {
		chunkTx, err := baseTx._BuildTransaction(i)
		if err != nil {
			return FeeEstimateResponse{}, errors.Wrapf(err, "failed to build chunk %d", i)
		}

		chunkResponse, err := q.callGetFeeEstimate(client, chunkTx)
		if err != nil {
			return FeeEstimateResponse{}, errors.Wrapf(err, "failed to estimate chunk %d", i)
		}

		totalNodeSubtotal += chunkResponse.NodeFee.Subtotal()
		totalServiceSubtotal += chunkResponse.ServiceFee.Subtotal()

		if i == 0 {
			aggregatedResponse.NetworkFee.Multiplier = chunkResponse.NetworkFee.Multiplier
		}

		aggregatedResponse.Notes = append(aggregatedResponse.Notes, chunkResponse.Notes...)
	}

	aggregatedResponse.NodeFee.Base = totalNodeSubtotal
	aggregatedResponse.ServiceFee.Base = totalServiceSubtotal
	aggregatedResponse.NetworkFee.Subtotal = totalNodeSubtotal * uint64(aggregatedResponse.NetworkFee.Multiplier)
	aggregatedResponse.Total = aggregatedResponse.NetworkFee.Subtotal + totalNodeSubtotal + totalServiceSubtotal

	return aggregatedResponse, nil
}

// estimateSingleTransaction estimates fees for a single transaction
func (q *FeeEstimateQuery) estimateSingleTransaction(client *Client, tx TransactionInterface) (FeeEstimateResponse, error) {
	baseTx := tx.getBaseTransaction()

	protoTx, err := baseTx._BuildTransaction(0)
	if err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to build transaction")
	}

	return q.callGetFeeEstimate(client, protoTx)
}

// callGetFeeEstimate calls the fee estimate REST API endpoint
func (q *FeeEstimateQuery) callGetFeeEstimate(client *Client, protoTx *services.Transaction) (FeeEstimateResponse, error) {
	if client.mirrorNetwork == nil || len(client.GetMirrorNetwork()) == 0 {
		return FeeEstimateResponse{}, errors.New("mirror node is not set")
	}

	mirrorUrl, err := client.GetMirrorRestApiBaseUrl()
	if err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to get mirror REST API base URL")
	}

	isLocalHost := strings.Contains(mirrorUrl, "localhost") || strings.Contains(mirrorUrl, "127.0.0.1")
	if isLocalHost {
		mirrorUrl = "http://localhost:8084/api/v1"
	}

	txBytes, err := protobuf.Marshal(protoTx)
	if err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to marshal transaction")
	}

	url := fmt.Sprintf("%s/network/fees?mode=%s", mirrorUrl, q.mode.String())

	var lastErr error
	var resp *http.Response

	for q.attempt < q.maxAttempts {
		resp, err = http.Post(url, "application/protobuf", bytes.NewBuffer(txBytes)) // #nosec
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			break
		}

		switch {
		case err != nil:
			lastErr = err
		case resp != nil:
			body, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr == nil {
				lastErr = fmt.Errorf("received non-200 response: %d, details: %s", resp.StatusCode, body)
			} else {
				lastErr = fmt.Errorf("received non-200 response: %d", resp.StatusCode)
			}
		default:
			lastErr = fmt.Errorf("received nil response")
		}

		// Check if we should retry
		if !q.shouldRetry(err, resp) {
			return FeeEstimateResponse{}, errors.Wrap(lastErr, "failed to call fee estimate API")
		}

		// Calculate delay with exponential backoff
		delayMs := 250.0 * float64(uint64(1)<<q.attempt) // 250ms, 500ms, 1000ms, etc.
		if delayMs > 8000 {
			delayMs = 8000
		}

		// Wait before retry
		select {
		case <-client.networkUpdateContext.Done():
			return FeeEstimateResponse{}, client.networkUpdateContext.Err()
		case <-time.After(time.Duration(delayMs) * time.Millisecond):
		}

		q.attempt++
	}

	if resp == nil {
		return FeeEstimateResponse{}, errors.Wrapf(lastErr, "failed to call fee estimate API after %d attempts", q.maxAttempts)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to read response body")
	}

	var response FeeEstimateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return FeeEstimateResponse{}, errors.Wrap(err, "failed to unmarshal response")
	}

	return response, nil
}

// shouldRetry determines if an error should be retried
func (q *FeeEstimateQuery) shouldRetry(err error, resp *http.Response) bool {
	if err == nil && resp != nil {
		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			return true
		}
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return false
		}
	}

	if err == nil {
		return false
	}

	return true
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
