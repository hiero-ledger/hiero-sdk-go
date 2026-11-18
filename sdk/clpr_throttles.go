package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"errors"
	"fmt"
	"math"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// _ClprSyncBytesHeadroom is the minimum gap between MaxSyncBytes and MaxMessagePayloadBytes.
const _ClprSyncBytesHeadroom = 65536

var errClprThrottleOutOfRange = errors.New("out of range")
var errClprSyncBytesBelowFloor = errors.New("maxSyncBytes must be at least maxMessagePayloadBytes + 65536")

// ClprThrottles are the CLPR ledger throttles (HIP-1535). Values read from the network are not validated.
type ClprThrottles struct {
	MaxMessagesPerBundle   uint32
	MaxMessagePayloadBytes uint32
	MaxGasPerMessage       uint64
	MaxQueueDepth          uint32
	MaxSyncBytes           uint64
	// MaxLocalEndpoints is the most endpoints the live manifest may hold.
	MaxLocalEndpoints uint32
	// MaxPeerEndpoints is the most peer endpoints cached per Channel; zero means no limit.
	MaxPeerEndpoints uint32
}

// Validate checks the limits a ledger configuration update must satisfy: the five non-endpoint
// limits are positive signed values and MaxSyncBytes is at least MaxMessagePayloadBytes + 65536.
func (throttles ClprThrottles) Validate() error {
	for _, limit := range []struct {
		name  string
		value uint64
		max   uint64
	}{
		{"maxMessagesPerBundle", uint64(throttles.MaxMessagesPerBundle), math.MaxInt32},
		{"maxMessagePayloadBytes", uint64(throttles.MaxMessagePayloadBytes), math.MaxInt32},
		{"maxGasPerMessage", throttles.MaxGasPerMessage, math.MaxInt64},
		{"maxQueueDepth", uint64(throttles.MaxQueueDepth), math.MaxInt32},
		{"maxSyncBytes", throttles.MaxSyncBytes, math.MaxInt64},
	} {
		if limit.value == 0 || limit.value > limit.max {
			return fmt.Errorf("%s must be between 1 and %d: %w", limit.name, limit.max, errClprThrottleOutOfRange)
		}
	}

	if throttles.MaxSyncBytes < uint64(throttles.MaxMessagePayloadBytes)+_ClprSyncBytesHeadroom {
		return errClprSyncBytesBelowFloor
	}

	return nil
}

func _ClprThrottlesFromProtobuf(pb *services.ClprThrottles) ClprThrottles {
	return ClprThrottles{
		MaxMessagesPerBundle:   pb.GetMaxMessagesPerBundle(),
		MaxMessagePayloadBytes: pb.GetMaxMessagePayloadBytes(),
		MaxGasPerMessage:       pb.GetMaxGasPerMessage(),
		MaxQueueDepth:          pb.GetMaxQueueDepth(),
		MaxSyncBytes:           pb.GetMaxSyncBytes(),
		MaxLocalEndpoints:      pb.GetMaxLocalEndpoints(),
		MaxPeerEndpoints:       pb.GetMaxPeerEndpoints(),
	}
}

func (throttles ClprThrottles) _ToProtobuf() *services.ClprThrottles {
	return &services.ClprThrottles{
		MaxMessagesPerBundle:   throttles.MaxMessagesPerBundle,
		MaxMessagePayloadBytes: throttles.MaxMessagePayloadBytes,
		MaxGasPerMessage:       throttles.MaxGasPerMessage,
		MaxQueueDepth:          throttles.MaxQueueDepth,
		MaxSyncBytes:           throttles.MaxSyncBytes,
		MaxLocalEndpoints:      throttles.MaxLocalEndpoints,
		MaxPeerEndpoints:       throttles.MaxPeerEndpoints,
	}
}
