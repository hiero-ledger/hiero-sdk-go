package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprEndpoint is one endpoint of a CLPR endpoint manifest (HIP-1535). TLSCertificate is opaque DER
// bytes carried as published; it is slated for deprecation and does not by itself authenticate the
// endpoint. AccountID is the raw account bytes as published, not a parsed AccountID.
type ClprEndpoint struct {
	ServiceEndpoint ClprServiceEndpoint
	TLSCertificate  []byte
	AccountID       []byte
}

func _ClprEndpointFromProtobuf(pb *services.ClprEndpoint) ClprEndpoint {
	return ClprEndpoint{
		ServiceEndpoint: _ClprServiceEndpointFromProtobuf(pb.GetServiceEndpoint()),
		TLSCertificate:  _ClprCopyBytes(pb.GetTlsCertificate()),
		AccountID:       _ClprCopyBytes(pb.GetAccountId()),
	}
}
