//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnitMirrorNodeHttpRetryPolicyZeroValueIsDefault(t *testing.T) {
	t.Parallel()

	var policy MirrorNodeHttpRetryPolicy

	assert.Equal(t, DefaultMirrorNodeHttpRetryPolicy(), policy)
	assert.Equal(t, uint16(5), policy.GetMaxAttempts())
	assert.Equal(t, 30*time.Second, policy.GetPerAttemptTimeout())
	assert.Zero(t, policy.GetTotalDeadline())
	assert.Equal(t, 250*time.Millisecond, policy.GetInitialBackoff())
	assert.Equal(t, 8*time.Second, policy.GetMaxBackoff())
	assert.Equal(t, []uint16{408, 429, 500, 502, 503, 504}, policy.GetRetryableStatusCodes())
}

func TestUnitMirrorNodeHttpRetryPolicyWithChangesOneField(t *testing.T) {
	t.Parallel()

	original := DefaultMirrorNodeHttpRetryPolicy()
	changed := original.WithMaxAttempts(2)

	assert.Equal(t, uint16(2), changed.GetMaxAttempts())
	assert.Equal(t, uint16(5), original.GetMaxAttempts())
	assert.Equal(t, 30*time.Second, changed.GetPerAttemptTimeout())
	assert.Equal(t, 250*time.Millisecond, changed.GetInitialBackoff())
	assert.Equal(t, 8*time.Second, changed.GetMaxBackoff())
	assert.Equal(t, original.GetRetryableStatusCodes(), changed.GetRetryableStatusCodes())
}

func TestUnitMirrorNodeHttpRetryPolicyKeepsExplicitZero(t *testing.T) {
	t.Parallel()

	policy := DefaultMirrorNodeHttpRetryPolicy().
		WithPerAttemptTimeout(0).
		WithInitialBackoff(0).
		WithMaxBackoff(0).
		WithRetryableStatusCodes(nil)

	assert.Zero(t, policy.GetPerAttemptTimeout())
	assert.Zero(t, policy.GetInitialBackoff())
	assert.Zero(t, policy.GetMaxBackoff())
	assert.Empty(t, policy.GetRetryableStatusCodes())
}

func TestUnitMirrorNodeHttpRetryPolicyStatusCodesAreCopied(t *testing.T) {
	t.Parallel()

	codes := []uint16{http.StatusServiceUnavailable}
	policy := DefaultMirrorNodeHttpRetryPolicy().WithRetryableStatusCodes(codes)
	codes[0] = http.StatusTeapot

	read := policy.GetRetryableStatusCodes()
	read[0] = http.StatusTeapot

	assert.Equal(t, []uint16{http.StatusServiceUnavailable}, policy.GetRetryableStatusCodes())
}

func TestUnitMirrorNodeHttpRetryPolicyRejectsInvalidArguments(t *testing.T) {
	t.Parallel()

	policy := DefaultMirrorNodeHttpRetryPolicy()

	assert.Panics(t, func() { policy.WithMaxAttempts(0) })
	assert.Panics(t, func() { policy.WithPerAttemptTimeout(-time.Second) })
	assert.Panics(t, func() { policy.WithTotalDeadline(-time.Second) })
	assert.Panics(t, func() { policy.WithInitialBackoff(-time.Second) })
	assert.Panics(t, func() { policy.WithMaxBackoff(-time.Second) })
	assert.Panics(t, func() { policy.WithRetryableStatusCodes([]uint16{99}) })
	assert.Panics(t, func() { policy.WithRetryableStatusCodes([]uint16{600}) })
}

func TestUnitMirrorNodeHttpConfigZeroValueIsDefault(t *testing.T) {
	t.Parallel()

	var config MirrorNodeHttpConfig

	assert.Equal(t, DefaultMirrorNodeHttpConfig(), config)
	assert.Nil(t, config.GetTransport())
	assert.Equal(t, DefaultHttpTransportConfiguration(), config.GetTransportConfiguration())
	assert.Equal(t, uint16(5), config.GetRetryPolicy().GetMaxAttempts())
	assert.Empty(t, config.GetRequestHeaders())
}

func TestUnitMirrorNodeHttpConfigWithChangesOneField(t *testing.T) {
	t.Parallel()

	original := DefaultMirrorNodeHttpConfig().WithRequestHeader("Authorization", "Bearer a")
	changed := original.WithRetryPolicy(original.GetRetryPolicy().WithMaxAttempts(2))

	assert.Equal(t, uint16(2), changed.GetRetryPolicy().GetMaxAttempts())
	assert.Equal(t, uint16(5), original.GetRetryPolicy().GetMaxAttempts())
	assert.Equal(t, map[string]string{"authorization": "Bearer a"}, changed.GetRequestHeaders())

	withSecond := changed.WithRequestHeader("x-trace", "1")
	assert.Len(t, changed.GetRequestHeaders(), 1)
	assert.Len(t, withSecond.GetRequestHeaders(), 2)
}

func TestUnitMirrorNodeHttpConfigRejectsUserAgentHeaders(t *testing.T) {
	t.Parallel()

	config := DefaultMirrorNodeHttpConfig()

	assert.Panics(t, func() { config.WithRequestHeader("User-Agent", "my-app") })
	assert.Panics(t, func() { config.WithRequestHeader("X-User-Agent", "my-app") })
	assert.Panics(t, func() { config.WithRequestHeader(" x-user-agent ", "my-app") })
}
