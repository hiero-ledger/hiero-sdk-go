package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ClprLedgerConfigurationInfo is the local CLPR ledger configuration returned by
// ClprGetLedgerConfigurationQuery (HIP-1535). The configuration fields are the live values;
// ConfigurationStateProof is the signed state proof of the latest snapshot. An empty proof is a
// successful read: no signed snapshot exists yet, so retry later.
type ClprLedgerConfigurationInfo struct {
	ProtocolVersion         uint32
	ChainID                 string
	ServiceAddress          []byte
	Timestamp               time.Time
	Throttles               ClprThrottles
	InitialTrustAnchor      []byte
	InitialTrustAnchorID    []byte
	ConfigurationStateProof []byte
}

func _ClprLedgerConfigurationInfoFromProtobuf(pb *services.ClprGetLedgerConfigurationResponse) ClprLedgerConfigurationInfo {
	configuration := pb.GetConfiguration()

	return ClprLedgerConfigurationInfo{
		ProtocolVersion:         configuration.GetProtocolVersion(),
		ChainID:                 configuration.GetChainId(),
		ServiceAddress:          _ClprCopyBytes(configuration.GetServiceAddress()),
		Timestamp:               _TimeFromProtobuf(configuration.GetTimestamp()),
		Throttles:               _ClprThrottlesFromProtobuf(configuration.GetThrottles()),
		InitialTrustAnchor:      _ClprCopyBytes(configuration.GetInitialTrustAnchor()),
		InitialTrustAnchorID:    _ClprCopyBytes(configuration.GetInitialTrustAnchorId()),
		ConfigurationStateProof: _ClprCopyBytes(pb.GetConfigurationStateProof()),
	}
}
