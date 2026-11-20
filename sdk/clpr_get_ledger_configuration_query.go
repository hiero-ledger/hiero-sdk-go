package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprGetLedgerConfigurationQuery returns the local CLPR ledger configuration and its state proof (HIP-1535).
type ClprGetLedgerConfigurationQuery struct {
	Query
}

// NewClprGetLedgerConfigurationQuery creates a ClprGetLedgerConfigurationQuery.
func NewClprGetLedgerConfigurationQuery() *ClprGetLedgerConfigurationQuery {
	header := services.QueryHeader{}
	return &ClprGetLedgerConfigurationQuery{
		Query: _NewQuery(true, &header),
	}
}

// SetGrpcDeadline When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *ClprGetLedgerConfigurationQuery) SetGrpcDeadline(deadline *time.Duration) *ClprGetLedgerConfigurationQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

func (q *ClprGetLedgerConfigurationQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *ClprGetLedgerConfigurationQuery) Execute(client *Client) (ClprLedgerConfigurationInfo, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return ClprLedgerConfigurationInfo{}, err
	}

	return _ClprLedgerConfigurationInfoFromProtobuf(resp.GetClprGetLedgerConfiguration()), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *ClprGetLedgerConfigurationQuery) SetMaxQueryPayment(maxPayment Hbar) *ClprGetLedgerConfigurationQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *ClprGetLedgerConfigurationQuery) SetQueryPayment(paymentAmount Hbar) *ClprGetLedgerConfigurationQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this ClprGetLedgerConfigurationQuery.
func (q *ClprGetLedgerConfigurationQuery) SetNodeAccountIDs(accountID []AccountID) *ClprGetLedgerConfigurationQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *ClprGetLedgerConfigurationQuery) SetMaxRetry(count int) *ClprGetLedgerConfigurationQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *ClprGetLedgerConfigurationQuery) SetMaxBackoff(max time.Duration) *ClprGetLedgerConfigurationQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *ClprGetLedgerConfigurationQuery) SetMinBackoff(min time.Duration) *ClprGetLedgerConfigurationQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *ClprGetLedgerConfigurationQuery) SetPaymentTransactionID(transactionID TransactionID) *ClprGetLedgerConfigurationQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *ClprGetLedgerConfigurationQuery) SetLogLevel(level LogLevel) *ClprGetLedgerConfigurationQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *ClprGetLedgerConfigurationQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetClpr().GetLedgerConfiguration,
	}
}

func (q *ClprGetLedgerConfigurationQuery) getName() string {
	return "ClprGetLedgerConfigurationQuery"
}

func (q *ClprGetLedgerConfigurationQuery) buildQuery() *services.Query {
	pb := services.Query_ClprGetLedgerConfiguration{
		ClprGetLedgerConfiguration: &services.ClprGetLedgerConfigurationQuery{
			Header: q.pbHeader,
		},
	}

	return &services.Query{
		Query: &pb,
	}
}

func (q *ClprGetLedgerConfigurationQuery) validateNetworkOnIDs(*Client) error {
	return nil
}

func (q *ClprGetLedgerConfigurationQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetClprGetLedgerConfiguration()
}
