package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprGetEndpointManifestQuery returns the local CLPR endpoint manifest and its state proof (HIP-1535).
type ClprGetEndpointManifestQuery struct {
	Query
}

// NewClprGetEndpointManifestQuery creates a ClprGetEndpointManifestQuery.
func NewClprGetEndpointManifestQuery() *ClprGetEndpointManifestQuery {
	header := services.QueryHeader{}
	return &ClprGetEndpointManifestQuery{
		Query: _NewQuery(true, &header),
	}
}

// SetGrpcDeadline When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *ClprGetEndpointManifestQuery) SetGrpcDeadline(deadline *time.Duration) *ClprGetEndpointManifestQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

func (q *ClprGetEndpointManifestQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *ClprGetEndpointManifestQuery) Execute(client *Client) (ClprEndpointManifestInfo, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return ClprEndpointManifestInfo{}, err
	}

	return _ClprEndpointManifestInfoFromProtobuf(resp.GetClprGetEndpointManifest()), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *ClprGetEndpointManifestQuery) SetMaxQueryPayment(maxPayment Hbar) *ClprGetEndpointManifestQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *ClprGetEndpointManifestQuery) SetQueryPayment(paymentAmount Hbar) *ClprGetEndpointManifestQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this ClprGetEndpointManifestQuery.
func (q *ClprGetEndpointManifestQuery) SetNodeAccountIDs(accountID []AccountID) *ClprGetEndpointManifestQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *ClprGetEndpointManifestQuery) SetMaxRetry(count int) *ClprGetEndpointManifestQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *ClprGetEndpointManifestQuery) SetMaxBackoff(max time.Duration) *ClprGetEndpointManifestQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *ClprGetEndpointManifestQuery) SetMinBackoff(min time.Duration) *ClprGetEndpointManifestQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *ClprGetEndpointManifestQuery) SetPaymentTransactionID(transactionID TransactionID) *ClprGetEndpointManifestQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *ClprGetEndpointManifestQuery) SetLogLevel(level LogLevel) *ClprGetEndpointManifestQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *ClprGetEndpointManifestQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetClpr().GetEndpointManifest,
	}
}

func (q *ClprGetEndpointManifestQuery) getName() string {
	return "ClprGetEndpointManifestQuery"
}

func (q *ClprGetEndpointManifestQuery) buildQuery() *services.Query {
	pb := services.Query_ClprGetEndpointManifest{
		ClprGetEndpointManifest: &services.ClprGetEndpointManifestQuery{
			Header: q.pbHeader,
		},
	}

	return &services.Query{
		Query: &pb,
	}
}

func (q *ClprGetEndpointManifestQuery) validateNetworkOnIDs(*Client) error {
	return nil
}

func (q *ClprGetEndpointManifestQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetClprGetEndpointManifest()
}
