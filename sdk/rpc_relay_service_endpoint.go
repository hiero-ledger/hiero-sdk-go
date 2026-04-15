package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// RpcRelayServiceEndpoint a registered service endpoint for an RPC relay.
type RpcRelayServiceEndpoint struct {
	registeredEndpointBase
}

// SetIPAddress sets the IP address for this endpoint.
func (e *RpcRelayServiceEndpoint) SetIPAddress(ip []byte) *RpcRelayServiceEndpoint {
	e.registeredEndpointBase.SetIPAddress(ip)
	return e
}

// SetDomainName sets the domain name for this endpoint.
func (e *RpcRelayServiceEndpoint) SetDomainName(name string) *RpcRelayServiceEndpoint {
	e.registeredEndpointBase.SetDomainName(name)
	return e
}

// SetPort sets the port number for this endpoint.
func (e *RpcRelayServiceEndpoint) SetPort(port uint32) *RpcRelayServiceEndpoint {
	e.registeredEndpointBase.SetPort(port)
	return e
}

// SetRequiresTls sets whether this endpoint requires TLS.
func (e *RpcRelayServiceEndpoint) SetRequiresTls(tls bool) *RpcRelayServiceEndpoint {
	e.registeredEndpointBase.SetRequiresTls(tls)
	return e
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
