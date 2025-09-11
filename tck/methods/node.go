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

type NodeService struct {
	sdkService *SDKService
}

func (n *NodeService) SetSdkService(service *SDKService) {
	n.sdkService = service
}

// CreateNode jRPC method for createNode
func (n *NodeService) CreateNode(_ context.Context, params param.CreateNodeParams) (*response.NodeResponse, error) {
	transaction := hiero.NewNodeCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	if params.Description != nil {
		transaction.SetDescription(*params.Description)
	}

	if params.GossipEndpoints != nil {
		for _, endpointParam := range *params.GossipEndpoints {
			endpoint, err := convertEndpointParam(endpointParam)
			if err != nil {
				return nil, err
			}
			transaction.AddGossipEndpoint(endpoint)
		}
	}

	if params.ServiceEndpoints != nil {
		for _, endpointParam := range *params.ServiceEndpoints {
			endpoint, err := convertEndpointParam(endpointParam)
			if err != nil {
				return nil, err
			}
			transaction.AddServiceEndpoint(endpoint)
		}
	}

	if params.GossipCaCertificate != nil {
		certBytes, err := hex.DecodeString(*params.GossipCaCertificate)
		if err != nil {
			return nil, err
		}
		transaction.SetGossipCaCertificate(certBytes)
	}

	if params.GrpcCertificateHash != nil {
		hashBytes, err := hex.DecodeString(*params.GrpcCertificateHash)
		if err != nil {
			return nil, err
		}
		transaction.SetGrpcCertificateHash(hashBytes)
	}

	if err := utils.SetKeyIfPresent(params.AdminKey, transaction.SetAdminKey); err != nil {
		return nil, err
	}

	if params.DeclineReward != nil {
		transaction.SetDeclineReward(*params.DeclineReward)
	}

	if params.GrpcWebProxyEndpoint != nil {
		endpoint, err := convertEndpointParam(*params.GrpcWebProxyEndpoint)
		if err != nil {
			return nil, err
		}
		transaction.SetGrpcWebProxyEndpoint(endpoint)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, n.sdkService.Client)
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(n.sdkService.Client)
	if err != nil {
		return nil, err
	}

	receipt, err := txResponse.SetValidateStatus(true).GetReceipt(n.sdkService.Client)
	if err != nil {
		return nil, err
	}

	var nodeId = strconv.FormatUint(receipt.NodeID, 10)

	return &response.NodeResponse{NodeId: nodeId, Status: receipt.Status.String()}, nil
}

func convertEndpointParam(endpointParam param.EndpointParams) (hiero.Endpoint, error) {
	endpoint := hiero.Endpoint{}

	if endpointParam.Address != nil {
		decoded, err := hex.DecodeString(*endpointParam.Address)
		if err != nil {
			return endpoint, err
		}
		endpoint.SetAddress(decoded)
	}

	if endpointParam.Port != nil {
		endpoint.SetPort(*endpointParam.Port)
	}

	if endpointParam.DomainName != nil {
		endpoint.SetDomainName(*endpointParam.DomainName)
	}

	return endpoint, nil
}
