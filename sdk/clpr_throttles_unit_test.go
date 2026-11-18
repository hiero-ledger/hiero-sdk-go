//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitClprThrottlesValidate(t *testing.T) {
	t.Parallel()

	require.NoError(t, clprThrottleLimits().Validate())

	for _, tt := range []struct {
		field string
		zero  func(*ClprThrottles)
	}{
		{"maxMessagesPerBundle", func(l *ClprThrottles) { l.MaxMessagesPerBundle = 0 }},
		{"maxMessagePayloadBytes", func(l *ClprThrottles) { l.MaxMessagePayloadBytes = 0 }},
		{"maxGasPerMessage", func(l *ClprThrottles) { l.MaxGasPerMessage = 0 }},
		{"maxQueueDepth", func(l *ClprThrottles) { l.MaxQueueDepth = 0 }},
		{"maxSyncBytes", func(l *ClprThrottles) { l.MaxSyncBytes = 0 }},
	} {
		limits := clprThrottleLimits()
		tt.zero(&limits)
		err := limits.Validate()
		assert.ErrorIs(t, err, errClprThrottleOutOfRange, tt.field)
		assert.ErrorContains(t, err, tt.field)
	}

	for _, tt := range []struct {
		field string
		max   func(*ClprThrottles)
	}{
		{"maxMessagesPerBundle", func(l *ClprThrottles) { l.MaxMessagesPerBundle = math.MaxInt32 + 1 }},
		{"maxMessagePayloadBytes", func(l *ClprThrottles) { l.MaxMessagePayloadBytes = math.MaxInt32 + 1 }},
		{"maxGasPerMessage", func(l *ClprThrottles) { l.MaxGasPerMessage = math.MaxInt64 + 1 }},
		{"maxQueueDepth", func(l *ClprThrottles) { l.MaxQueueDepth = math.MaxInt32 + 1 }},
		{"maxSyncBytes", func(l *ClprThrottles) { l.MaxSyncBytes = math.MaxInt64 + 1 }},
	} {
		limits := clprThrottleLimits()
		tt.max(&limits)
		err := limits.Validate()
		assert.ErrorIs(t, err, errClprThrottleOutOfRange, tt.field)
		assert.ErrorContains(t, err, tt.field)
	}

	limits := clprThrottleLimits()
	limits.MaxLocalEndpoints = 0
	limits.MaxPeerEndpoints = 0
	assert.NoError(t, limits.Validate())
}

func TestUnitClprThrottlesSyncFloor(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		payload uint32
		floor   uint64
	}{{1024, 66560}, {math.MaxInt32, 2147549183}} {
		limits := clprThrottleLimits()
		limits.MaxMessagePayloadBytes = tt.payload

		limits.MaxSyncBytes = tt.floor
		assert.NoError(t, limits.Validate())

		limits.MaxSyncBytes = tt.floor - 1
		assert.ErrorIs(t, limits.Validate(), errClprSyncBytesBelowFloor)
	}
}

func TestUnitClprThrottlesProtobufRoundTrip(t *testing.T) {
	t.Parallel()

	limits := ClprThrottles{
		MaxMessagesPerBundle:   math.MaxUint32,
		MaxMessagePayloadBytes: math.MaxUint32,
		MaxGasPerMessage:       math.MaxUint64,
		MaxQueueDepth:          math.MaxUint32,
		MaxSyncBytes:           math.MaxUint64,
		MaxLocalEndpoints:      math.MaxUint32,
		MaxPeerEndpoints:       math.MaxUint32,
	}
	assert.Equal(t, limits, _ClprThrottlesFromProtobuf(limits._ToProtobuf()))

	unvalidated := ClprThrottles{MaxSyncBytes: 1}
	assert.Equal(t, unvalidated, _ClprThrottlesFromProtobuf(unvalidated._ToProtobuf()))
	assert.Equal(t, ClprThrottles{}, _ClprThrottlesFromProtobuf(nil))
}
