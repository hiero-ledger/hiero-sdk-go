//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitAccountInfoQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	infoQuery := NewAccountInfoQuery().
		SetAccountID(accountID)

	err = infoQuery.validateNetworkOnIDs(client)
	require.NoError(t, err)
}
func TestAccountInfoQuery_Get(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	deadline := time.Duration(time.Minute)
	accountId := AccountID{Account: 3, checksum: &checksum}
	transactionID := TransactionIDGenerate(accountId)
	query := NewAccountInfoQuery().
		SetAccountID(accountId).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(10)).
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		SetMaxRetry(5).
		SetMaxBackoff(10 * time.Second).
		SetMinBackoff(1 * time.Second).
		SetPaymentTransactionID(transactionID).
		SetGrpcDeadline(&deadline)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	err = query.validateNetworkOnIDs(client)
	require.NoError(t, err)
	require.Equal(t, accountId, query.GetAccountID())
	require.Equal(t, NewHbar(2), query.GetQueryPayment())
	require.Equal(t, NewHbar(10), query.GetMaxQueryPayment())
	require.Equal(t, []AccountID{{Account: 3}, {Account: 4}}, query.GetNodeAccountIDs())
	require.Equal(t, 5, query.GetMaxRetryCount())
	require.Equal(t, 10*time.Second, query.GetMaxBackoff())
	require.Equal(t, 1*time.Second, query.GetMinBackoff())
	require.Equal(t, transactionID, query.GetPaymentTransactionID())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitAccountInfoQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	infoQuery := NewAccountInfoQuery().
		SetAccountID(accountID)

	err = infoQuery.validateNetworkOnIDs(client)
	require.Error(t, err)
	if err != nil {
		require.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountInfoQuerySetNothing(t *testing.T) {
	t.Parallel()

	balance := NewAccountInfoQuery()

	require.Equal(t, AccountID{}, balance.GetAccountID())
	require.Equal(t, []AccountID{}, balance.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, balance.GetMinBackoff())
	require.Equal(t, 8*time.Second, balance.GetMaxBackoff())
	require.Equal(t, 10, balance.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, balance.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, balance.GetQueryPayment())
	require.Equal(t, Hbar{}, balance.GetMaxQueryPayment())
}

func Test_AccountInfoQueryMapStatusError(t *testing.T) {
	t.Parallel()

	response := services.Response{
		Response: &services.Response_CryptoGetInfo{
			CryptoGetInfo: &services.CryptoGetInfoResponse{
				Header: &services.ResponseHeader{
					NodeTransactionPrecheckCode: services.ResponseCodeEnum(StatusInvalidAccountID),
					ResponseType:                services.ResponseType_COST_ANSWER,
				},
			},
		},
	}

	query := NewAccountInfoQuery()
	actualError := query.mapStatusError(query, &response)

	expectedError := ErrHederaPreCheckStatus{
		Status: StatusInvalidAccountID,
	}

	require.Equal(t, expectedError, actualError)
}

func TestUnitAccountInfoQueryMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_CryptoGetInfo{
				CryptoGetInfo: &services.CryptoGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewAccountInfoQuery().
		SetAccountID(AccountID{Account: 1234}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, HbarFromTinybar(2), cost)
}

// A paid query on a client without an operator must fail cleanly instead of panicking while
// building the payment.
func TestUnitAccountInfoQueryExecuteWithoutOperator(t *testing.T) {
	t.Parallel()

	client, server := NewMockClientAndServer([][]interface{}{{}})
	defer server.Close()

	client.operator = nil

	// An explicit payment skips the cost lookup and goes straight to payment signing.
	_, err := NewAccountInfoQuery().
		SetAccountID(AccountID{Account: 1800}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetQueryPayment(NewHbar(1)).
		Execute(client)
	require.ErrorIs(t, err, errNoClientProvided)
}

// A COST_ANSWER query is unpaid, so it must go out without a payment even on a query object a
// previous Execute already attached the client to - reaching for the absent operator there would
// panic while building the payment.
func TestUnitAccountInfoQueryGetCostAfterExecuteWithoutOperator(t *testing.T) {
	t.Parallel()

	var captured []*services.Query
	client, server := NewMockClientAndServer([][]interface{}{{_CostAnswer(services.ResponseCodeEnum_OK, &captured)}})
	defer server.Close()

	client.operator = nil

	query := NewAccountInfoQuery().
		SetAccountID(AccountID{Account: 1800}).
		SetNodeAccountIDs([]AccountID{{Account: 3}})

	// Fails for want of an operator, but leaves the client on the query.
	_, err := query.Execute(client)
	require.ErrorIs(t, err, errNoClientProvided)
	require.NotNil(t, query.client)

	cost, err := query.GetCost(client)
	require.NoError(t, err)
	require.Equal(t, HbarFromTinybar(25), cost)

	// The failed Execute must not have reached the network, so this is the cost lookup.
	require.Len(t, captured, 1)
	info, ok := captured[0].Query.(*services.Query_CryptoGetInfo)
	require.True(t, ok)
	require.Equal(t, services.ResponseType_COST_ANSWER, info.CryptoGetInfo.Header.ResponseType)
	require.Nil(t, info.CryptoGetInfo.Header.Payment, "a COST_ANSWER query must not attach a payment transaction")
}

// GetCost must fail cleanly without a client.
func TestUnitAccountInfoQueryGetCostNilClient(t *testing.T) {
	t.Parallel()

	_, err := NewAccountInfoQuery().
		SetAccountID(AccountID{Account: 1800}).
		GetCost(nil)
	require.ErrorIs(t, err, errNoClientProvided)
}
