package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// MirrorNodeServiceEndpoint a registered service endpoint for a mirror node.
type MirrorNodeServiceEndpoint struct {
	registeredEndpointBase
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
