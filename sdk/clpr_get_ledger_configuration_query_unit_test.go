//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protobuf "google.golang.org/protobuf/proto"
)

func clprThrottleLimits() ClprThrottles {
	return ClprThrottles{
		MaxMessagesPerBundle:   8,
		MaxMessagePayloadBytes: 1024,
		MaxGasPerMessage:       9007199254740993,
		MaxQueueDepth:          16,
		MaxSyncBytes:           66560,
		MaxLocalEndpoints:      2,
		MaxPeerEndpoints:       3,
	}
}

func clprLedgerConfigurationResponse(header *services.ResponseHeader, proof []byte) *services.Response {
	return &services.Response{
		Response: &services.Response_ClprGetLedgerConfiguration{
			ClprGetLedgerConfiguration: &services.ClprGetLedgerConfigurationResponse{
				Header: header,
				Configuration: &services.ClprLedgerConfiguration{
					ProtocolVersion:      1,
					ChainId:              "hedera:testnet",
					ServiceAddress:       _ClprServiceAddress[:],
					Timestamp:            &services.Timestamp{Seconds: 1700000000, Nanos: 123456789},
					Throttles:            clprThrottleLimits()._ToProtobuf(),
					InitialTrustAnchor:   []byte{1, 2, 3, 4},
					InitialTrustAnchorId: []byte{5, 6, 7, 8},
				},
				ConfigurationStateProof: proof,
			},
		},
	}
}

func clprCostHeader(cost uint64) *services.ResponseHeader {
	return &services.ResponseHeader{ResponseType: services.ResponseType_COST_ANSWER, Cost: cost}
}

func clprAnswerHeader() *services.ResponseHeader {
	return &services.ResponseHeader{ResponseType: services.ResponseType_ANSWER_ONLY}
}

// clprPaymentAmount returns the tinybars the query payment transfers to the node, or -1 without a payment.
func clprPaymentAmount(t *testing.T, header *services.QueryHeader) int64 {
	t.Helper()
	if header.GetPayment() == nil {
		return -1
	}
	var body services.TransactionBody
	require.NoError(t, protobuf.Unmarshal(header.GetPayment().GetBodyBytes(), &body)) // nolint
	for _, amount := range body.GetCryptoTransfer().GetTransfers().GetAccountAmounts() {
		if amount.GetAmount() >= 0 {
			return amount.GetAmount()
		}
	}
	return -1
}

func TestUnitClprGetLedgerConfigurationQueryBuildQuery(t *testing.T) {
	t.Parallel()

	query := NewClprGetLedgerConfigurationQuery()
	pb := query.buildQuery()

	require.NotNil(t, pb.GetClprGetLedgerConfiguration())
	assert.NotNil(t, pb.GetClprGetLedgerConfiguration().GetHeader())
	assert.Equal(t, "ClprGetLedgerConfigurationQuery", query.getName())
}

func TestUnitClprGetLedgerConfigurationQueryMapsResponse(t *testing.T) {
	t.Parallel()

	info := _ClprLedgerConfigurationInfoFromProtobuf(
		clprLedgerConfigurationResponse(clprAnswerHeader(), []byte{9, 10, 11, 12}).GetClprGetLedgerConfiguration())

	assert.Equal(t, uint32(1), info.ProtocolVersion)
	assert.Equal(t, "hedera:testnet", info.ChainID)
	assert.Equal(t, _ClprServiceAddress[:], info.ServiceAddress)
	assert.Equal(t, time.Unix(1700000000, 123456789), info.Timestamp)
	assert.Equal(t, clprThrottleLimits(), info.Throttles)
	assert.Equal(t, []byte{1, 2, 3, 4}, info.InitialTrustAnchor)
	assert.Equal(t, []byte{5, 6, 7, 8}, info.InitialTrustAnchorID)
	assert.Equal(t, []byte{9, 10, 11, 12}, info.ConfigurationStateProof)

	empty := _ClprLedgerConfigurationInfoFromProtobuf(&services.ClprGetLedgerConfigurationResponse{})
	assert.Equal(t, ClprThrottles{}, empty.Throttles)
	assert.Empty(t, empty.ConfigurationStateProof)
}

func TestUnitClprGetLedgerConfigurationQueryCostAndPayment(t *testing.T) {
	t.Parallel()

	for _, cost := range []uint64{0, 100} {
		var responseTypes []services.ResponseType
		var payments []int64
		record := func(request *services.Query) *services.Response {
			header := request.GetClprGetLedgerConfiguration().GetHeader()
			responseTypes = append(responseTypes, header.GetResponseType())
			payments = append(payments, clprPaymentAmount(t, header))
			if header.GetResponseType() == services.ResponseType_COST_ANSWER {
				return clprLedgerConfigurationResponse(clprCostHeader(cost), nil)
			}
			return clprLedgerConfigurationResponse(clprAnswerHeader(), []byte{9})
		}

		client, server := NewMockClientAndServer([][]interface{}{{record, record}})
		info, err := NewClprGetLedgerConfigurationQuery().
			SetNodeAccountIDs([]AccountID{{Account: 3}}).
			Execute(client)
		server.Close()

		require.NoError(t, err)
		assert.Equal(t, []byte{9}, info.ConfigurationStateProof)
		assert.Equal(t, []services.ResponseType{services.ResponseType_COST_ANSWER, services.ResponseType_ANSWER_ONLY}, responseTypes)
		assert.Equal(t, int64(cost), payments[1])
	}
}

func TestUnitClprGetLedgerConfigurationQueryExplicitPayment(t *testing.T) {
	t.Parallel()

	var payments []int64
	record := func(request *services.Query) *services.Response {
		header := request.GetClprGetLedgerConfiguration().GetHeader()
		assert.Equal(t, services.ResponseType_ANSWER_ONLY, header.GetResponseType())
		payments = append(payments, clprPaymentAmount(t, header))
		return clprLedgerConfigurationResponse(clprAnswerHeader(), nil)
	}

	client, server := NewMockClientAndServer([][]interface{}{{record}})
	defer server.Close()

	info, err := NewClprGetLedgerConfigurationQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetQueryPayment(HbarFromTinybar(7)).
		Execute(client)
	require.NoError(t, err)
	assert.Equal(t, []int64{7}, payments)
	assert.Empty(t, info.ConfigurationStateProof)
}

func TestUnitClprGetLedgerConfigurationQueryMaxPaymentExceeded(t *testing.T) {
	t.Parallel()

	client, server := NewMockClientAndServer([][]interface{}{{
		clprLedgerConfigurationResponse(clprCostHeader(100), nil),
	}})
	defer server.Close()

	_, err := NewClprGetLedgerConfigurationQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(HbarFromTinybar(99)).
		Execute(client)

	var exceeded ErrMaxQueryPaymentExceeded
	require.ErrorAs(t, err, &exceeded)
	assert.Equal(t, HbarFromTinybar(100), exceeded.QueryCost)
}

func TestUnitClprGetLedgerConfigurationQueryRequiresOperator(t *testing.T) {
	t.Parallel()

	client, server := NewMockClientAndServer([][]interface{}{{}})
	defer server.Close()
	client.operator = nil

	_, err := NewClprGetLedgerConfigurationQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetQueryPayment(HbarFromTinybar(7)).
		Execute(client)
	assert.ErrorIs(t, err, errNoClientProvided)
}

func TestUnitClprGetLedgerConfigurationQueryPrecheckStatus(t *testing.T) {
	t.Parallel()

	notEnabled := &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_CLPR_NOT_ENABLED}

	for _, payment := range []Hbar{{}, HbarFromTinybar(7)} {
		client, server := NewMockClientAndServer([][]interface{}{{
			clprLedgerConfigurationResponse(notEnabled, nil),
		}})

		_, err := NewClprGetLedgerConfigurationQuery().
			SetNodeAccountIDs([]AccountID{{Account: 3}}).
			SetQueryPayment(payment).
			Execute(client)
		server.Close()

		var precheck ErrHederaPreCheckStatus
		require.ErrorAs(t, err, &precheck)
		assert.Equal(t, CLPR_NOT_ENABLED, precheck.Status)
	}
}

func TestUnitClprGetLedgerConfigurationQueryKeepsUint64Throttles(t *testing.T) {
	t.Parallel()

	response := clprLedgerConfigurationResponse(clprAnswerHeader(), []byte{})
	response.GetClprGetLedgerConfiguration().GetConfiguration().GetThrottles().MaxGasPerMessage = math.MaxUint64
	response.GetClprGetLedgerConfiguration().GetConfiguration().GetThrottles().MaxSyncBytes = math.MaxUint64
	response.GetClprGetLedgerConfiguration().GetConfiguration().GetThrottles().MaxMessagesPerBundle = 0

	client, server := NewMockClientAndServer([][]interface{}{{response}})
	defer server.Close()

	info, err := NewClprGetLedgerConfigurationQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetQueryPayment(HbarFromTinybar(1)).
		Execute(client)
	require.NoError(t, err)
	assert.Equal(t, uint64(math.MaxUint64), info.Throttles.MaxGasPerMessage)
	assert.Equal(t, uint64(math.MaxUint64), info.Throttles.MaxSyncBytes)
	assert.Equal(t, uint32(0), info.Throttles.MaxMessagesPerBundle)
	assert.Empty(t, info.ConfigurationStateProof)
}
