package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type SDKService struct {
	clients map[string]*hiero.Client
}

func NewSdkService() *SDKService {
	return &SDKService{
		clients: make(map[string]*hiero.Client, 0),
	}
}

func (s *SDKService) GetClient(sessionId string) *hiero.Client {
	return s.clients[sessionId]
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
	s.clients[params.SessionId] = client

	return response.SetupResponse{
		Message: "Successfully setup " + clientType + " client.",
		Status:  "SUCCESS",
	}, nil
}

// Reset function for the SDK
func (s *SDKService) Reset(_ context.Context, params param.BaseParams) response.SetupResponse {
	delete(s.clients, params.SessionId)
	return response.SetupResponse{
		Status: "SUCCESS",
	}
}
