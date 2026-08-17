package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// AccountBalanceQuery gets the balance of a CryptoCurrency account. This returns only the balance, so it is a smaller
// and faster reply than AccountInfoQuery, which returns the balance plus additional information.
//
// Deprecated: AccountBalanceQuery is no longer supported. Use MirrorNodeAccountBalanceQuery or the mirror node REST API (GET /api/v1/accounts/{id}) to retrieve account balances.
//
// Consensus node v0.77 no longer serves CryptoGetAccountBalance, so Execute and GetCost return
// errAccountBalanceQueryDeprecated without contacting the network. The type and its setters are
// retained so existing code keeps compiling; none of them affect the outcome, and the query never
// reaches the _Execute loop (hence no getMethod/buildQuery/getQueryResponse implementations).
type AccountBalanceQuery struct {
	Query
	accountID  *AccountID
	contractID *ContractID
}

// NewAccountBalanceQuery creates an AccountBalanceQuery query which can be used to construct and execute
// an AccountBalanceQuery.
// It is recommended that you use this for creating new instances of an AccountBalanceQuery
// instead of manually creating an instance of the struct.
//
// Deprecated: AccountBalanceQuery is no longer supported. Use MirrorNodeAccountBalanceQuery or the mirror node REST API (GET /api/v1/accounts/{id}) to retrieve account balances.
func NewAccountBalanceQuery() *AccountBalanceQuery {
	header := services.QueryHeader{}
	return &AccountBalanceQuery{
		Query: _NewQuery(false, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *AccountBalanceQuery) SetGrpcDeadline(deadline *time.Duration) *AccountBalanceQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetAccountID sets the AccountID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (q *AccountBalanceQuery) SetAccountID(accountID AccountID) *AccountBalanceQuery {
	q.accountID = &accountID
	return q
}

// GetAccountID returns the AccountID for which you wish to query the balance.
func (q *AccountBalanceQuery) GetAccountID() AccountID {
	if q.accountID == nil {
		return AccountID{}
	}

	return *q.accountID
}

// SetContractID sets the ContractID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (q *AccountBalanceQuery) SetContractID(contractID ContractID) *AccountBalanceQuery {
	q.contractID = &contractID
	return q
}

// GetContractID returns the ContractID for which you wish to query the balance.
func (q *AccountBalanceQuery) GetContractID() ContractID {
	if q.contractID == nil {
		return ContractID{}
	}

	return *q.contractID
}

// GetCost returns an error immediately without any network request. The client is ignored.
//
// Deprecated: AccountBalanceQuery is no longer supported. Use MirrorNodeAccountBalanceQuery or the mirror node REST API (GET /api/v1/accounts/{id}) to retrieve account balances.
func (q *AccountBalanceQuery) GetCost(_ *Client) (Hbar, error) {
	return Hbar{}, errAccountBalanceQueryDeprecated
}

// Execute returns an error immediately without any network request. The client is ignored.
//
// Deprecated: AccountBalanceQuery is no longer supported. Use MirrorNodeAccountBalanceQuery or the mirror node REST API (GET /api/v1/accounts/{id}) to retrieve account balances.
func (q *AccountBalanceQuery) Execute(_ *Client) (AccountBalance, error) {
	return AccountBalance{}, errAccountBalanceQueryDeprecated
}

// SetMaxQueryPayment sets the maximum payment allowed for this query.
func (q *AccountBalanceQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountBalanceQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this query.
func (q *AccountBalanceQuery) SetQueryPayment(paymentAmount Hbar) *AccountBalanceQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountBalanceQuery.
func (q *AccountBalanceQuery) SetNodeAccountIDs(accountID []AccountID) *AccountBalanceQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *AccountBalanceQuery) SetMaxRetry(count int) *AccountBalanceQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries. Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *AccountBalanceQuery) SetMaxBackoff(max time.Duration) *AccountBalanceQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *AccountBalanceQuery) SetMinBackoff(min time.Duration) *AccountBalanceQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *AccountBalanceQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountBalanceQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *AccountBalanceQuery) SetLogLevel(level LogLevel) *AccountBalanceQuery {
	q.Query.SetLogLevel(level)
	return q
}
