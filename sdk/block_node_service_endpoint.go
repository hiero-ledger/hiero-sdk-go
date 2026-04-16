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
	endpointApis []BlockNodeApi
}

// SetIPAddress sets the IP address for this endpoint.
func (e *BlockNodeServiceEndpoint) SetIPAddress(ip []byte) *BlockNodeServiceEndpoint {
	e.registeredEndpointBase.SetIPAddress(ip)
	return e
}

// SetDomainName sets the domain name for this endpoint.
func (e *BlockNodeServiceEndpoint) SetDomainName(name string) *BlockNodeServiceEndpoint {
	e.registeredEndpointBase.SetDomainName(name)
	return e
}

// SetPort sets the port number for this endpoint.
func (e *BlockNodeServiceEndpoint) SetPort(port uint32) *BlockNodeServiceEndpoint {
	e.registeredEndpointBase.SetPort(port)
	return e
}

// SetRequiresTls sets whether this endpoint requires TLS.
func (e *BlockNodeServiceEndpoint) SetRequiresTls(tls bool) *BlockNodeServiceEndpoint {
	e.registeredEndpointBase.SetRequiresTls(tls)
	return e
}

// SetEndpointApis sets the block node API kinds for this endpoint.
func (e *BlockNodeServiceEndpoint) SetEndpointApis(apis []BlockNodeApi) *BlockNodeServiceEndpoint {
	e.endpointApis = apis
	return e
}

// AddEndpointApi adds a block node API kind to this endpoint.
func (e *BlockNodeServiceEndpoint) AddEndpointApi(api BlockNodeApi) *BlockNodeServiceEndpoint {
	e.endpointApis = append(e.endpointApis,api)
	return e
}

// GetEndpointApi returns the block node API kind for this endpoint.
func (e *BlockNodeServiceEndpoint) GetEndpointApis() []BlockNodeApi {
	return e.endpointApis
}

// _ToProtobuf converts this BlockNodeServiceEndpoint to its protobuf representation.
func (e *BlockNodeServiceEndpoint) _ToProtobuf() *services.RegisteredServiceEndpoint {
	apis := make([]services.RegisteredServiceEndpoint_BlockNodeEndpoint_BlockNodeApi, len(e.endpointApis))
	for i, api := range e.endpointApis {
		apis[i] = services.RegisteredServiceEndpoint_BlockNodeEndpoint_BlockNodeApi(api)
	}

	pb := &services.RegisteredServiceEndpoint{
		Port:        e.port,
		RequiresTls: e.requiresTls,
		EndpointType: &services.RegisteredServiceEndpoint_BlockNode{
			BlockNode: &services.RegisteredServiceEndpoint_BlockNodeEndpoint{
				EndpointApi: apis,
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
		pbApis := bn.BlockNode.GetEndpointApi()
		endpoint.endpointApis = make([]BlockNodeApi, len(pbApis))
		for i, api := range pbApis {
			endpoint.endpointApis[i] = BlockNodeApi(api)
		}
	}

	return endpoint
}
