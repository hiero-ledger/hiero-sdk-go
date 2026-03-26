package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// BlockNodeServiceEndpoint a registered service endpoint for a block node.
// Extends the base registered endpoint with the specific block node API
// that the endpoint exposes.
type BlockNodeServiceEndpoint struct {
	registeredEndpointBase
	endpointApi BlockNodeApi
}

// SetEndpointApi sets the block node API kind for this endpoint.
func (e *BlockNodeServiceEndpoint) SetEndpointApi(api BlockNodeApi) *BlockNodeServiceEndpoint {
	e.endpointApi = api
	return e
}

// GetEndpointApi returns the block node API kind for this endpoint.
func (e *BlockNodeServiceEndpoint) GetEndpointApi() BlockNodeApi {
	return e.endpointApi
}

// _ToProtobuf converts this BlockNodeServiceEndpoint to its protobuf representation.
func (e *BlockNodeServiceEndpoint) _ToProtobuf() *services.RegisteredServiceEndpoint {
	pb := &services.RegisteredServiceEndpoint{
		Port:        e.port,
		RequiresTls: e.requiresTls,
		EndpointType: &services.RegisteredServiceEndpoint_BlockNode{
			BlockNode: &services.RegisteredServiceEndpoint_BlockNodeEndpoint{
				EndpointApi: services.RegisteredServiceEndpoint_BlockNodeEndpoint_BlockNodeApi(e.endpointApi),
			},
		},
	}
	e.addressToProtobuf(pb)
	return pb
}

// _BlockNodeServiceEndpointFromProtobuf converts a protobuf RegisteredServiceEndpoint to a BlockNodeServiceEndpoint.
func _BlockNodeServiceEndpointFromProtobuf(pb *services.RegisteredServiceEndpoint) *BlockNodeServiceEndpoint {
	endpoint := &BlockNodeServiceEndpoint{
		registeredEndpointBase: baseFromProtobuf(pb),
	}

	if bn, ok := pb.EndpointType.(*services.RegisteredServiceEndpoint_BlockNode); ok && bn.BlockNode != nil {
		endpoint.endpointApi = BlockNodeApi(bn.BlockNode.GetEndpointApi())
	}

	return endpoint
}
