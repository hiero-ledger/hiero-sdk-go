package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprServiceEndpoint is the numeric IPv4 or IPv6 address and port of a CLPR endpoint (HIP-1535).
type ClprServiceEndpoint struct {
	IPAddress string
	Port      uint32
}

func _ClprServiceEndpointFromProtobuf(pb *services.ClprServiceEndpoint) ClprServiceEndpoint {
	return ClprServiceEndpoint{
		IPAddress: pb.GetIpAddress(),
		Port:      pb.GetPort(),
	}
}
