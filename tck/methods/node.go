package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"encoding/hex"
	"strconv"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// NodeService provides node-related TCK operations
type NodeService struct {
	sdkService *SDKService
}

// SetSdkService sets the SDK service reference
func (n *NodeService) SetSdkService(service *SDKService) {
	n.sdkService = service
}

// CreateNode jRPC method for createNode
func (n *NodeService) CreateNode(_ context.Context, params param.CreateNodeParams) (*response.NodeResponse, error) {
	transaction := hiero.NewNodeCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	// Set description
	if params.Description != nil {
		transaction.SetDescription(*params.Description)
	}

	// Set gossip endpoints
	if params.GossipEndpoints != nil {
		for _, endpointParam := range *params.GossipEndpoints {
			endpoint, err := convertEndpointParam(endpointParam)
			if err != nil {
				return nil, err
			}
			transaction.AddGossipEndpoint(endpoint)
		}
	}

	// Set service endpoints
	if params.ServiceEndpoints != nil {
		for _, endpointParam := range *params.ServiceEndpoints {
			endpoint, err := convertEndpointParam(endpointParam)
			if err != nil {
				return nil, err
			}
			transaction.AddServiceEndpoint(endpoint)
		}
	}

	// Set gossip CA certificate
	if params.GossipCaCertificate != nil {
		certBytes, err := hex.DecodeString(*params.GossipCaCertificate)
		if err != nil {
			return nil, err
		}
		transaction.SetGossipCaCertificate(certBytes)
	}

	// Set gRPC certificate hash
	if params.GrpcCertificateHash != nil {
		hashBytes, err := hex.DecodeString(*params.GrpcCertificateHash)
		if err != nil {
			return nil, err
		}
		transaction.SetGrpcCertificateHash(hashBytes)
	}

	// Set admin key
	if err := utils.SetKeyIfPresent(params.AdminKey, transaction.SetAdminKey); err != nil {
		return nil, err
	}

	// Set decline reward
	if params.DeclineReward != nil {
		transaction.SetDeclineReward(*params.DeclineReward)
	}

	// Set gRPC web proxy endpoint
	if params.GrpcWebProxyEndpoint != nil {
		endpoint, err := convertEndpointParam(*params.GrpcWebProxyEndpoint)
		if err != nil {
			return nil, err
		}
		transaction.SetGrpcWebProxyEndpoint(endpoint)
	}

	// Set common transaction parameters
	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, n.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	// Execute transaction
	txResponse, err := transaction.Execute(n.sdkService.Client)
	if err != nil {
		return nil, err
	}

	// Get receipt
	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(n.sdkService.Client)
	if err != nil {
		return nil, err
	}

	var nodeId = strconv.FormatUint(receipt.NodeID, 10)

	return &response.NodeResponse{NodeId: nodeId, Status: receipt.Status.String()}, nil
}

func (n *NodeService) DeleteNode(_ context.Context, params param.DeleteNodeParams) (*response.NodeResponse, error) {
	transaction := hiero.NewNodeDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set node ID
	if params.NodeId != nil {
		nodeId, err := strconv.ParseUint(*params.NodeId, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetNodeID(nodeId)
	}

	// Set common transaction parameters
	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, n.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	// Execute transaction
	txResponse, err := transaction.Execute(n.sdkService.Client)
	if err != nil {
		return nil, err
	}

	// Get receipt
	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(n.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.NodeResponse{Status: receipt.Status.String()}, nil
}

// convertEndpointParam converts a param.EndpointParams to hiero.Endpoint
func convertEndpointParam(endpointParam param.EndpointParams) (hiero.Endpoint, error) {
	endpoint := hiero.Endpoint{}

	// Set address
	if endpointParam.Address != nil {
		decoded, err := hex.DecodeString(*endpointParam.Address)
		if err != nil {
			return endpoint, err
		}
		endpoint.SetAddress(decoded)
	}

	// Set port
	if endpointParam.Port != nil {
		endpoint.SetPort(*endpointParam.Port)
	}

	// Set domain name (this will override if both address and domainName are provided)
	if endpointParam.DomainName != nil {
		endpoint.SetDomainName(*endpointParam.DomainName)
	}

	return endpoint, nil
}
