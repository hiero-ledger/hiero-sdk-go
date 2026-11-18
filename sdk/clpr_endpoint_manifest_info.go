package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprEndpointManifestInfo is the local CLPR endpoint manifest returned by
// ClprGetEndpointManifestQuery (HIP-1535). ManifestStateProof is the state proof of the manifest in the
// latest signed snapshot. An empty proof is a successful read: no signed snapshot exists yet, so retry later.
type ClprEndpointManifestInfo struct {
	Version            uint64
	ServiceAddress     []byte
	Endpoints          []ClprEndpoint
	ManifestStateProof []byte
}

func _ClprEndpointManifestInfoFromProtobuf(pb *services.ClprGetEndpointManifestResponse) ClprEndpointManifestInfo {
	manifest := pb.GetManifest()

	endpoints := make([]ClprEndpoint, 0, len(manifest.GetEndpoints()))
	for _, endpoint := range manifest.GetEndpoints() {
		endpoints = append(endpoints, _ClprEndpointFromProtobuf(endpoint))
	}

	return ClprEndpointManifestInfo{
		Version:            manifest.GetVersion(),
		ServiceAddress:     _ClprCopyBytes(manifest.GetServiceAddress()),
		Endpoints:          endpoints,
		ManifestStateProof: _ClprCopyBytes(pb.GetManifestStateProof()),
	}
}
