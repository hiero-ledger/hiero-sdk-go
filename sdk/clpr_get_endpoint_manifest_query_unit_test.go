//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clprEndpointManifestResponse(header *services.ResponseHeader, proof []byte) *services.Response {
	return &services.Response{
		Response: &services.Response_ClprGetEndpointManifest{
			ClprGetEndpointManifest: &services.ClprGetEndpointManifestResponse{
				Header: header,
				Manifest: &services.ClprEndpointManifest{
					Version:        9007199254740993,
					ServiceAddress: _ClprServiceAddress[:],
					Endpoints: []*services.ClprEndpoint{
						{
							ServiceEndpoint: &services.ClprServiceEndpoint{IpAddress: "192.0.2.1", Port: 50211},
							TlsCertificate:  []byte{1, 2, 3, 4},
							AccountId:       []byte{0xff, 0x00, 0x80, 0x7f},
						},
						{
							ServiceEndpoint: &services.ClprServiceEndpoint{IpAddress: "2001:db8::2", Port: 50212},
							TlsCertificate:  []byte{5, 6, 7, 8},
							AccountId:       []byte{0xaa, 0xbb, 0xcc, 0xdd},
						},
					},
				},
				ManifestStateProof: proof,
			},
		},
	}
}

func TestUnitClprGetEndpointManifestQueryBuildQuery(t *testing.T) {
	t.Parallel()

	query := NewClprGetEndpointManifestQuery()
	pb := query.buildQuery()

	require.NotNil(t, pb.GetClprGetEndpointManifest())
	assert.NotNil(t, pb.GetClprGetEndpointManifest().GetHeader())
	assert.Equal(t, "ClprGetEndpointManifestQuery", query.getName())
}

func TestUnitClprGetEndpointManifestQueryMapsResponse(t *testing.T) {
	t.Parallel()

	info := _ClprEndpointManifestInfoFromProtobuf(
		clprEndpointManifestResponse(clprAnswerHeader(), []byte{9, 10, 11, 12}).GetClprGetEndpointManifest())

	assert.Equal(t, uint64(9007199254740993), info.Version)
	assert.Equal(t, _ClprServiceAddress[:], info.ServiceAddress)
	assert.Equal(t, []ClprEndpoint{
		{ServiceEndpoint: ClprServiceEndpoint{IPAddress: "192.0.2.1", Port: 50211}, TLSCertificate: []byte{1, 2, 3, 4}, AccountID: []byte{0xff, 0x00, 0x80, 0x7f}},
		{ServiceEndpoint: ClprServiceEndpoint{IPAddress: "2001:db8::2", Port: 50212}, TLSCertificate: []byte{5, 6, 7, 8}, AccountID: []byte{0xaa, 0xbb, 0xcc, 0xdd}},
	}, info.Endpoints)
	assert.Equal(t, []byte{9, 10, 11, 12}, info.ManifestStateProof)

	empty := _ClprEndpointManifestInfoFromProtobuf(&services.ClprGetEndpointManifestResponse{})
	assert.NotNil(t, empty.Endpoints)
	assert.Empty(t, empty.Endpoints)
	assert.Empty(t, empty.ManifestStateProof)
}

func TestUnitClprGetEndpointManifestQueryCopiesDecodedBytes(t *testing.T) {
	t.Parallel()

	pb := clprEndpointManifestResponse(clprAnswerHeader(), []byte{9}).GetClprGetEndpointManifest()
	info := _ClprEndpointManifestInfoFromProtobuf(pb)

	pb.ManifestStateProof[0] = 0
	pb.Manifest.Endpoints[0].TlsCertificate[0] = 0
	pb.Manifest.Endpoints[0].AccountId[0] = 0

	assert.Equal(t, []byte{9}, info.ManifestStateProof)
	assert.Equal(t, []byte{1, 2, 3, 4}, info.Endpoints[0].TLSCertificate)
	assert.Equal(t, []byte{0xff, 0x00, 0x80, 0x7f}, info.Endpoints[0].AccountID)
}

func TestUnitClprGetEndpointManifestQueryCostAndPayment(t *testing.T) {
	t.Parallel()

	for _, cost := range []uint64{0, 100} {
		var responseTypes []services.ResponseType
		var payments []int64
		record := func(request *services.Query) *services.Response {
			header := request.GetClprGetEndpointManifest().GetHeader()
			responseTypes = append(responseTypes, header.GetResponseType())
			payments = append(payments, clprPaymentAmount(t, header))
			if header.GetResponseType() == services.ResponseType_COST_ANSWER {
				return clprEndpointManifestResponse(clprCostHeader(cost), nil)
			}
			return clprEndpointManifestResponse(clprAnswerHeader(), []byte{9})
		}

		client, server := NewMockClientAndServer([][]interface{}{{record, record, record}})
		query := NewClprGetEndpointManifestQuery().SetNodeAccountIDs([]AccountID{{Account: 3}})

		queryCost, err := query.GetCost(client)
		require.NoError(t, err)
		assert.Equal(t, HbarFromTinybar(int64(cost)), queryCost)

		info, err := query.Execute(client)
		server.Close()

		require.NoError(t, err)
		assert.Len(t, info.Endpoints, 2)
		assert.Equal(t, []services.ResponseType{
			services.ResponseType_COST_ANSWER, services.ResponseType_COST_ANSWER, services.ResponseType_ANSWER_ONLY,
		}, responseTypes)
		assert.Equal(t, int64(cost), payments[2])
	}
}

func TestUnitClprGetEndpointManifestQueryPrecheckStatus(t *testing.T) {
	t.Parallel()

	client, server := NewMockClientAndServer([][]interface{}{{
		clprEndpointManifestResponse(&services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_CLPR_NOT_ENABLED}, nil),
	}})
	defer server.Close()

	_, err := NewClprGetEndpointManifestQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		Execute(client)

	var precheck ErrHederaPreCheckStatus
	require.ErrorAs(t, err, &precheck)
	assert.Equal(t, CLPR_NOT_ENABLED, precheck.Status)
}

func TestUnitClprGetEndpointManifestQueryKeepsUint64VersionAndEmptyProof(t *testing.T) {
	t.Parallel()

	response := clprEndpointManifestResponse(clprAnswerHeader(), []byte{})
	response.GetClprGetEndpointManifest().GetManifest().Version = math.MaxUint64

	client, server := NewMockClientAndServer([][]interface{}{{response}})
	defer server.Close()

	info, err := NewClprGetEndpointManifestQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetQueryPayment(HbarFromTinybar(1)).
		Execute(client)
	require.NoError(t, err)
	assert.Equal(t, uint64(math.MaxUint64), info.Version)
	assert.Equal(t, "2001:db8::2", info.Endpoints[1].ServiceEndpoint.IPAddress)
	assert.Empty(t, info.ManifestStateProof)
}
