package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// MirrorNodeServiceEndpoint a registered service endpoint for a mirror node.
type MirrorNodeServiceEndpoint struct {
	registeredEndpointBase
}

// SetIPAddress sets the IP address for this endpoint.
func (e *MirrorNodeServiceEndpoint) SetIPAddress(ip []byte) *MirrorNodeServiceEndpoint {
	e.registeredEndpointBase.SetIPAddress(ip)
	return e
}

// SetDomainName sets the domain name for this endpoint.
func (e *MirrorNodeServiceEndpoint) SetDomainName(name string) *MirrorNodeServiceEndpoint {
	e.registeredEndpointBase.SetDomainName(name)
	return e
}

// SetPort sets the port number for this endpoint.
func (e *MirrorNodeServiceEndpoint) SetPort(port uint32) *MirrorNodeServiceEndpoint {
	e.registeredEndpointBase.SetPort(port)
	return e
}

// SetRequiresTls sets whether this endpoint requires TLS.
func (e *MirrorNodeServiceEndpoint) SetRequiresTls(tls bool) *MirrorNodeServiceEndpoint {
	e.registeredEndpointBase.SetRequiresTls(tls)
	return e
}

// _ToProtobuf converts this MirrorNodeServiceEndpoint to its protobuf representation.
func (e *MirrorNodeServiceEndpoint) _ToProtobuf() *services.RegisteredServiceEndpoint {
	pb := &services.RegisteredServiceEndpoint{
		Port:        e.port,
		RequiresTls: e.requiresTls,
		EndpointType: &services.RegisteredServiceEndpoint_MirrorNode{
			MirrorNode: &services.RegisteredServiceEndpoint_MirrorNodeEndpoint{},
		},
	}
	e.addressToProtobuf(pb)
	return pb
}

// _MirrorNodeServiceEndpointFromProtobuf converts a protobuf RegisteredServiceEndpoint to a MirrorNodeServiceEndpoint.
func _MirrorNodeServiceEndpointFromProtobuf(pb *services.RegisteredServiceEndpoint) *MirrorNodeServiceEndpoint {
	return &MirrorNodeServiceEndpoint{
		registeredEndpointBase: baseFromProtobuf(pb),
	}
}
