package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountStakersQuery gets all of the accounts that are proxy staking to this account. For each of  them, the amount
// currently staked will be given. This is not yet implemented, but will be in a future version of the API.
type AccountStakersQuery struct {
	query
	accountID *AccountID
}

// NewAccountStakersQuery creates an AccountStakersQuery query which can be used to construct and execute
// an AccountStakersQuery.
//
// It is recommended that you use this for creating new instances of an AccountStakersQuery
// instead of manually creating an instance of the struct.
func NewAccountStakersQuery() *AccountStakersQuery {
	header := services.QueryHeader{}
	result := AccountStakersQuery{
		query: _NewQuery(true, &header),
	}

	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *AccountStakersQuery) SetGrpcDeadline(deadline *time.Duration) *AccountStakersQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetAccountID sets the Account ID for which the stakers should be retrieved
func (this *AccountStakersQuery) SetAccountID(accountID AccountID) *AccountStakersQuery {
	this.accountID = &accountID
	return this
}

// GetAccountID returns the AccountID for this AccountStakersQuery.
func (this *AccountStakersQuery) GetAccountID() AccountID {
	if this.accountID == nil {
		return AccountID{}
	}

	return *this.accountID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *AccountStakersQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	err := this.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}
	if this.query.nodeAccountIDs.locked {
		for range this.nodeAccountIDs.slice {
			paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
			if err != nil {
				return Hbar{}, err
			}
			this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
		}
	}

	pb := this.build()
	pb.CryptoGetProxyStakers.Header = this.pbHeader

	this.pb = &services.Query{
		Query: pb,
	}

	this.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	this.paymentTransactionIDs._Advance()

	resp, err := _Execute(
		client,
		this.e,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetCryptoGetProxyStakers().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func (this *AccountStakersQuery) Execute(client *Client) ([]Transfer, error) {
	if client == nil || client.operator == nil {
		return []Transfer{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return []Transfer{}, err
	}

	if !this.paymentTransactionIDs.locked {
		this.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if this.queryPayment.tinybar != 0 {
		cost = this.queryPayment
	} else {
		if this.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = this.maxQueryPayment
		}

		actualCost, err := this.GetCost(client)
		if err != nil {
			return []Transfer{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []Transfer{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "AccountStakersQuery",
			}
		}

		cost = actualCost
	}

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			return []Transfer{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return []Transfer{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.CryptoGetProxyStakers.Header = this.pbHeader
	this.pb = &services.Query{
		Query: pb,
	}

	if this.isPaymentRequired && len(this.paymentTransactions) > 0 {
		this.paymentTransactionIDs._Advance()
	}
	this.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	resp, err := _Execute(
		client,
		this.e,
	)

	if err != nil {
		return []Transfer{}, err
	}

	var stakers = make([]Transfer, len(resp.(*services.Response).GetCryptoGetProxyStakers().Stakers.ProxyStaker))

	// TODO: This is wrong, this _Method shold return `[]ProxyStaker` not `[]Transfer`
	for i, element := range resp.(*services.Response).GetCryptoGetProxyStakers().Stakers.ProxyStaker {
		id := _AccountIDFromProtobuf(element.AccountID)
		accountID := AccountID{}

		if id == nil {
			accountID = *id
		}

		stakers[i] = Transfer{
			AccountID: accountID,
			Amount:    HbarFromTinybar(element.Amount),
		}
	}

	return stakers, err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *AccountStakersQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountStakersQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *AccountStakersQuery) SetQueryPayment(paymentAmount Hbar) *AccountStakersQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountStakersQuery.
func (this *AccountStakersQuery) SetNodeAccountIDs(accountID []AccountID) *AccountStakersQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *AccountStakersQuery) SetMaxRetry(count int) *AccountStakersQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *AccountStakersQuery) SetMaxBackoff(max time.Duration) *AccountStakersQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *AccountStakersQuery) GetMaxBackoff() time.Duration {
	return this.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *AccountStakersQuery) SetMinBackoff(min time.Duration) *AccountStakersQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *AccountStakersQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *AccountStakersQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("AccountStakersQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *AccountStakersQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountStakersQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *AccountStakersQuery) SetLogLevel(level LogLevel) *AccountStakersQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *AccountStakersQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetStakersByAccountID,
	}
}

func (this *AccountStakersQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetCryptoGetProxyStakers().Header.NodeTransactionPrecheckCode),
	}
}

func (this *AccountStakersQuery) getName() string {
	return "AccountStakersQuery"
}

func (this *AccountStakersQuery) build() *services.Query_CryptoGetProxyStakers {
	pb := services.Query_CryptoGetProxyStakers{
		CryptoGetProxyStakers: &services.CryptoGetStakersQuery{
			Header: &services.QueryHeader{},
		},
	}

	if this.accountID != nil {
		pb.CryptoGetProxyStakers.AccountID = this.accountID._ToProtobuf()
	}

	return &pb
}

func (this *AccountStakersQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.accountID != nil {
		if err := this.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *AccountStakersQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetCryptoGetProxyStakers().Header.NodeTransactionPrecheckCode)
}