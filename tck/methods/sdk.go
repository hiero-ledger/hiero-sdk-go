package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type SDKService struct {
	clients *utils.SafeClientMap
}

func NewSdkService() *SDKService {
	return &SDKService{
		clients: utils.NewSafeClientMap(),
	}
}

func (s *SDKService) GetClient(sessionId string) *hiero.Client {
	return s.clients.Get(sessionId)
}

// Setup function for the SDK
func (s *SDKService) Setup(_ context.Context, params param.SetupParams) (response.SetupResponse, error) {
	var clientType string
	var client *hiero.Client

	if params.NodeIp != nil && params.NodeAccountId != nil && params.MirrorNetworkIp != nil {
		// Custom client setup
		nodeId, err := hiero.AccountIDFromString(*params.NodeAccountId)
		if err != nil {
			return response.SetupResponse{}, err
		}
		node := map[string]hiero.AccountID{
			*params.NodeIp: nodeId,
		}
		client = hiero.ClientForNetwork(node)
		clientType = "custom"
		client.SetMirrorNetwork([]string{*params.MirrorNetworkIp})
	} else {
		// Default to testnet
		client = hiero.ClientForTestnet()
		clientType = "testnet"
	}

	// Set operator (adjustments may be needed based on the Hiero SDK)
	operatorId, _ := hiero.AccountIDFromString(params.OperatorAccountId)
	operatorKey, _ := hiero.PrivateKeyFromString(params.OperatorPrivateKey)
	client.SetOperator(operatorId, operatorKey)
	s.clients.Set(params.SessionId, client)

	return response.SetupResponse{
		Message: "Successfully setup " + clientType + " client.",
		Status:  "SUCCESS",
	}, nil
}

func (s *SDKService) SetOperator(_ context.Context, params param.SetupParams) response.SetupResponse {
	client := s.clients.Get(params.SessionId)
	operatorId, _ := hiero.AccountIDFromString(params.OperatorAccountId)
	operatorKey, _ := hiero.PrivateKeyFromString(params.OperatorPrivateKey)
	client.SetOperator(operatorId, operatorKey)
	return response.SetupResponse{
		Status: "SUCCESS",
	}
}

// Reset function for the SDK
func (s *SDKService) Reset(_ context.Context, params param.BaseParams) response.SetupResponse {
	s.clients.Delete(params.SessionId)
	return response.SetupResponse{
		Status: "SUCCESS",
	}
}
