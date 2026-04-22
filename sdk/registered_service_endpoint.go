package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// RegisteredServiceEndpoint a service endpoint published by a registered node. Each endpoint
// declares an address (IP or FQDN), port, TLS requirement, and the type of node service it provides.
// The concrete subtypes are BlockNodeServiceEndpoint, MirrorNodeServiceEndpoint,
// RpcRelayServiceEndpoint and GeneralServiceEndpoint.
type RegisteredServiceEndpoint interface {
	GetIPAddress() []byte
	GetDomainName() string
	GetPort() uint32
	GetRequiresTls() bool
	Validate() error
	_ToProtobuf() *services.RegisteredServiceEndpoint
}

type registeredEndpointBase struct {
	ipAddress   []byte
	domainName  string
	port        uint32
	requiresTls bool
}

// SetIPAddress sets the IP address for this endpoint.
func (e *registeredEndpointBase) SetIPAddress(ip []byte) *registeredEndpointBase {
	e.ipAddress = ip
	return e
}

// GetIPAddress returns the IP address of this endpoint.
func (e *registeredEndpointBase) GetIPAddress() []byte {
	return e.ipAddress
}

// SetDomainName sets the domain name for this endpoint.
func (e *registeredEndpointBase) SetDomainName(name string) *registeredEndpointBase {
	e.domainName = name
	return e
}

// GetDomainName returns the domain name of this endpoint.
func (e *registeredEndpointBase) GetDomainName() string {
	return e.domainName
}

// SetPort sets the port number for this endpoint.
func (e *registeredEndpointBase) SetPort(port uint32) *registeredEndpointBase {
	e.port = port
	return e
}

// GetPort returns the port number of this endpoint.
func (e *registeredEndpointBase) GetPort() uint32 {
	return e.port
}

// SetRequiresTls sets whether this endpoint requires TLS.
func (e *registeredEndpointBase) SetRequiresTls(tls bool) *registeredEndpointBase {
	e.requiresTls = tls
	return e
}

// GetRequiresTls returns whether this endpoint requires TLS.
func (e *registeredEndpointBase) GetRequiresTls() bool {
	return e.requiresTls
}

// Validate checks that this endpoint has a valid address configuration.
func (e *registeredEndpointBase) Validate() error {
	if e.ipAddress != nil && e.domainName != "" {
		return errEndpointCannotHaveBothAddressAndDomainName
	}
	if e.ipAddress == nil && e.domainName == "" {
		return errEndpointMustHaveAddressOrDomainName
	}
	return nil
}

// addressToProtobuf sets the address oneof on the given protobuf message.
func (e *registeredEndpointBase) addressToProtobuf(pb *services.RegisteredServiceEndpoint) {
	if e.ipAddress != nil {
		pb.Address = &services.RegisteredServiceEndpoint_IpAddress{
			IpAddress: e.ipAddress,
		}
	} else if e.domainName != "" {
		pb.Address = &services.RegisteredServiceEndpoint_DomainName{
			DomainName: e.domainName,
		}
	}
}

// baseFromProtobuf populates the base fields from a protobuf message.
func baseFromProtobuf(pb *services.RegisteredServiceEndpoint) registeredEndpointBase {
	base := registeredEndpointBase{
		port:        pb.GetPort(),
		requiresTls: pb.GetRequiresTls(),
	}

	switch addr := pb.Address.(type) {
	case *services.RegisteredServiceEndpoint_IpAddress:
		base.ipAddress = addr.IpAddress
	case *services.RegisteredServiceEndpoint_DomainName:
		base.domainName = addr.DomainName
	}

	return base
}

// _RegisteredServiceEndpointFromProtobuf dispatches on the protobuf endpoint_type
// oneof and returns the appropriate concrete SDK type.
func _RegisteredServiceEndpointFromProtobuf(pb *services.RegisteredServiceEndpoint) RegisteredServiceEndpoint {
	switch pb.EndpointType.(type) {
	case *services.RegisteredServiceEndpoint_BlockNode:
		return _BlockNodeServiceEndpointFromProtobuf(pb)
	case *services.RegisteredServiceEndpoint_MirrorNode:
		return _MirrorNodeServiceEndpointFromProtobuf(pb)
	case *services.RegisteredServiceEndpoint_RpcRelay:
		return _RpcRelayServiceEndpointFromProtobuf(pb)
	case *services.RegisteredServiceEndpoint_GeneralService:
		return _GeneralServiceEndpointFromProtobuf(pb)
	}
	return nil
}
