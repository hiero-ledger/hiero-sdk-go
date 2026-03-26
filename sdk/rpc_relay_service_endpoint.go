package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// RpcRelayServiceEndpoint a registered service endpoint for an RPC relay.
type RpcRelayServiceEndpoint struct {
	registeredEndpointBase
}

// _ToProtobuf converts this RpcRelayServiceEndpoint to its protobuf representation.
func (e *RpcRelayServiceEndpoint) _ToProtobuf() *services.RegisteredServiceEndpoint {
	pb := &services.RegisteredServiceEndpoint{
		Port:        e.port,
		RequiresTls: e.requiresTls,
		EndpointType: &services.RegisteredServiceEndpoint_RpcRelay{
			RpcRelay: &services.RegisteredServiceEndpoint_RpcRelayEndpoint{},
		},
	}
	e.addressToProtobuf(pb)
	return pb
}

// _RpcRelayServiceEndpointFromProtobuf converts a protobuf RegisteredServiceEndpoint to a RpcRelayServiceEndpoint.
func _RpcRelayServiceEndpointFromProtobuf(pb *services.RegisteredServiceEndpoint) *RpcRelayServiceEndpoint {
	return &RpcRelayServiceEndpoint{
		registeredEndpointBase: baseFromProtobuf(pb),
	}
}
