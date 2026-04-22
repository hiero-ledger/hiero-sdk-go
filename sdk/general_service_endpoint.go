package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// GeneralServiceEndpoint a registered service endpoint for a general-purpose service.
type GeneralServiceEndpoint struct {
	registeredEndpointBase
	description string
}

// SetIPAddress sets the IP address for this endpoint.
func (e *GeneralServiceEndpoint) SetIPAddress(ip []byte) *GeneralServiceEndpoint {
	e.registeredEndpointBase.SetIPAddress(ip)
	return e
}

// SetDomainName sets the domain name for this endpoint.
func (e *GeneralServiceEndpoint) SetDomainName(name string) *GeneralServiceEndpoint {
	e.registeredEndpointBase.SetDomainName(name)
	return e
}

// SetPort sets the port number for this endpoint.
func (e *GeneralServiceEndpoint) SetPort(port uint32) *GeneralServiceEndpoint {
	e.registeredEndpointBase.SetPort(port)
	return e
}

// SetRequiresTls sets whether this endpoint requires TLS.
func (e *GeneralServiceEndpoint) SetRequiresTls(tls bool) *GeneralServiceEndpoint {
	e.registeredEndpointBase.SetRequiresTls(tls)
	return e
}

// SetDescription sets a short description of the service provided by this endpoint.
func (e *GeneralServiceEndpoint) SetDescription(description string) *GeneralServiceEndpoint {
	e.description = description
	return e
}

// GetDescription returns the short description of the service provided by this endpoint.
func (e *GeneralServiceEndpoint) GetDescription() string {
	return e.description
}

// _ToProtobuf converts this GeneralServiceEndpoint to its protobuf representation.
func (e *GeneralServiceEndpoint) _ToProtobuf() *services.RegisteredServiceEndpoint {
	pb := &services.RegisteredServiceEndpoint{
		Port:        e.port,
		RequiresTls: e.requiresTls,
		EndpointType: &services.RegisteredServiceEndpoint_GeneralService{
			GeneralService: &services.RegisteredServiceEndpoint_GeneralServiceEndpoint{
				Description: e.description,
			},
		},
	}
	e.addressToProtobuf(pb)
	return pb
}

// _GeneralServiceEndpointFromProtobuf converts a protobuf RegisteredServiceEndpoint to a GeneralServiceEndpoint.
func _GeneralServiceEndpointFromProtobuf(pb *services.RegisteredServiceEndpoint) *GeneralServiceEndpoint {
	endpoint := &GeneralServiceEndpoint{
		registeredEndpointBase: baseFromProtobuf(pb),
	}

	if gs, ok := pb.EndpointType.(*services.RegisteredServiceEndpoint_GeneralService); ok && gs.GeneralService != nil {
		endpoint.description = gs.GeneralService.GetDescription()
	}

	return endpoint
}
