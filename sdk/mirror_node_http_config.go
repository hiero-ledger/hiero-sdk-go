package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"
)

const (
	mirrorNodeHttpDefaultMaxAttempts       = 5
	mirrorNodeHttpDefaultPerAttemptTimeout = 30 * time.Second
	mirrorNodeHttpDefaultInitialBackoff    = 250 * time.Millisecond
	mirrorNodeHttpDefaultMaxBackoff        = 8 * time.Second

	// Local mirror nodes may still be starting, so loopback URLs get a longer budget.
	mirrorNodeHttpLoopbackMaxAttempts   = 15
	mirrorNodeHttpLoopbackTotalDeadline = 90 * time.Second
)

// 501, 505, 506, 507, 508, 510 and 511 are absent on purpose: they describe a server that will
// answer the same way next time.
var mirrorNodeHttpDefaultRetryableStatusCodes = []uint16{
	http.StatusRequestTimeout,
	http.StatusTooManyRequests,
	http.StatusInternalServerError,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
}

// MirrorNodeHttpRetryPolicy configures retries for mirror node REST requests. The zero value is
// the default policy. Client.SetMaxAttempts does not apply to mirror node requests.
type MirrorNodeHttpRetryPolicy struct {
	set                  retryPolicyField
	maxAttempts          uint16
	perAttemptTimeout    time.Duration
	totalDeadline        time.Duration
	initialBackoff       time.Duration
	maxBackoff           time.Duration
	retryableStatusCodes []uint16
}

type retryPolicyField uint8

const (
	retryPolicyMaxAttempts retryPolicyField = 1 << iota
	retryPolicyPerAttemptTimeout
	retryPolicyTotalDeadline
	retryPolicyInitialBackoff
	retryPolicyMaxBackoff
	retryPolicyRetryableStatusCodes
)

// DefaultMirrorNodeHttpRetryPolicy returns the default policy, which is also the zero value.
func DefaultMirrorNodeHttpRetryPolicy() MirrorNodeHttpRetryPolicy {
	return MirrorNodeHttpRetryPolicy{}
}

func (p MirrorNodeHttpRetryPolicy) has(field retryPolicyField) bool {
	return p.set&field != 0
}

// GetMaxAttempts returns the maximum number of attempts per request. Each page of a paginated query gets its own.
func (p MirrorNodeHttpRetryPolicy) GetMaxAttempts() uint16 {
	if !p.has(retryPolicyMaxAttempts) {
		return mirrorNodeHttpDefaultMaxAttempts
	}
	return p.maxAttempts
}

// GetPerAttemptTimeout returns the timeout for a single attempt. 0 means none.
func (p MirrorNodeHttpRetryPolicy) GetPerAttemptTimeout() time.Duration {
	if !p.has(retryPolicyPerAttemptTimeout) {
		return mirrorNodeHttpDefaultPerAttemptTimeout
	}
	return p.perAttemptTimeout
}

// GetTotalDeadline returns the deadline for the whole call, retries included. 0 uses Client.GetRequestTimeout.
func (p MirrorNodeHttpRetryPolicy) GetTotalDeadline() time.Duration {
	return p.totalDeadline
}

// GetInitialBackoff returns the backoff before the first retry. Each wait is a random duration up to the backoff.
func (p MirrorNodeHttpRetryPolicy) GetInitialBackoff() time.Duration {
	if !p.has(retryPolicyInitialBackoff) {
		return mirrorNodeHttpDefaultInitialBackoff
	}
	return p.initialBackoff
}

// GetMaxBackoff returns the maximum backoff between attempts.
func (p MirrorNodeHttpRetryPolicy) GetMaxBackoff() time.Duration {
	if !p.has(retryPolicyMaxBackoff) {
		return mirrorNodeHttpDefaultMaxBackoff
	}
	return p.maxBackoff
}

// GetRetryableStatusCodes returns the HTTP statuses that are retried.
func (p MirrorNodeHttpRetryPolicy) GetRetryableStatusCodes() []uint16 {
	return slices.Clone(p.retryableStatuses())
}

func (p MirrorNodeHttpRetryPolicy) retryableStatuses() []uint16 {
	if !p.has(retryPolicyRetryableStatusCodes) {
		return mirrorNodeHttpDefaultRetryableStatusCodes
	}
	return p.retryableStatusCodes
}

// WithMaxAttempts returns a copy with maxAttempts set. It panics if maxAttempts is 0.
func (p MirrorNodeHttpRetryPolicy) WithMaxAttempts(maxAttempts uint16) MirrorNodeHttpRetryPolicy {
	if maxAttempts < 1 {
		panic("maxAttempts must be at least 1")
	}
	p.maxAttempts = maxAttempts
	p.set |= retryPolicyMaxAttempts
	return p
}

// WithPerAttemptTimeout returns a copy with the per-attempt timeout set. It panics if timeout is negative.
func (p MirrorNodeHttpRetryPolicy) WithPerAttemptTimeout(timeout time.Duration) MirrorNodeHttpRetryPolicy {
	requireNonNegativeDuration("perAttemptTimeout", timeout)
	p.perAttemptTimeout = timeout
	p.set |= retryPolicyPerAttemptTimeout
	return p
}

// WithTotalDeadline returns a copy with the total deadline set. It panics if deadline is negative.
func (p MirrorNodeHttpRetryPolicy) WithTotalDeadline(deadline time.Duration) MirrorNodeHttpRetryPolicy {
	requireNonNegativeDuration("totalDeadline", deadline)
	p.totalDeadline = deadline
	p.set |= retryPolicyTotalDeadline
	return p
}

// WithInitialBackoff returns a copy with the initial backoff set. It panics if backoff is negative.
func (p MirrorNodeHttpRetryPolicy) WithInitialBackoff(backoff time.Duration) MirrorNodeHttpRetryPolicy {
	requireNonNegativeDuration("initialBackoff", backoff)
	p.initialBackoff = backoff
	p.set |= retryPolicyInitialBackoff
	return p
}

// WithMaxBackoff returns a copy with the maximum backoff set. It panics if backoff is negative.
func (p MirrorNodeHttpRetryPolicy) WithMaxBackoff(backoff time.Duration) MirrorNodeHttpRetryPolicy {
	requireNonNegativeDuration("maxBackoff", backoff)
	p.maxBackoff = backoff
	p.set |= retryPolicyMaxBackoff
	return p
}

// WithRetryableStatusCodes returns a copy with the retried statuses set. It panics if a code is not 100–599.
func (p MirrorNodeHttpRetryPolicy) WithRetryableStatusCodes(codes []uint16) MirrorNodeHttpRetryPolicy {
	for _, code := range codes {
		if code < 100 || code > 599 {
			panic(fmt.Sprintf("%d is not an HTTP status code", code))
		}
	}
	p.retryableStatusCodes = slices.Clone(codes)
	p.set |= retryPolicyRetryableStatusCodes
	return p
}

// withLoopbackDefaults applies the local-network max attempts and total deadline where they are unset.
func (p MirrorNodeHttpRetryPolicy) withLoopbackDefaults() MirrorNodeHttpRetryPolicy {
	if !p.has(retryPolicyMaxAttempts) {
		p = p.WithMaxAttempts(mirrorNodeHttpLoopbackMaxAttempts)
	}
	if !p.has(retryPolicyTotalDeadline) {
		p = p.WithTotalDeadline(mirrorNodeHttpLoopbackTotalDeadline)
	}
	return p
}

// MirrorNodeHttpConfig configures how a Client sends mirror node REST requests. The zero value is the default.
type MirrorNodeHttpConfig struct {
	transport              HttpTransport
	transportConfiguration HttpTransportConfiguration
	retryPolicy            MirrorNodeHttpRetryPolicy
	requestHeaders         map[string]string
}

// DefaultMirrorNodeHttpConfig returns the default configuration, which is also the zero value.
func DefaultMirrorNodeHttpConfig() MirrorNodeHttpConfig {
	return MirrorNodeHttpConfig{}
}

// GetTransport returns the injected transport, or nil if the SDK's default is used.
func (c MirrorNodeHttpConfig) GetTransport() HttpTransport {
	return c.transport
}

// GetTransportConfiguration returns the configuration of the transport the SDK builds.
func (c MirrorNodeHttpConfig) GetTransportConfiguration() HttpTransportConfiguration {
	return c.transportConfiguration
}

// GetRetryPolicy returns the retry policy.
func (c MirrorNodeHttpConfig) GetRetryPolicy() MirrorNodeHttpRetryPolicy {
	return c.retryPolicy
}

// GetRequestHeaders returns the headers sent with every mirror node request.
func (c MirrorNodeHttpConfig) GetRequestHeaders() map[string]string {
	return maps.Clone(c.requestHeaders)
}

// WithTransport returns a copy that uses transport instead of the SDK's own; nil restores the default.
// The caller owns transport: Client.Close does not close it, and the transport configuration is ignored.
func (c MirrorNodeHttpConfig) WithTransport(transport HttpTransport) MirrorNodeHttpConfig {
	c.transport = transport
	return c
}

// WithTransportConfiguration returns a copy with the transport configuration set.
func (c MirrorNodeHttpConfig) WithTransportConfiguration(config HttpTransportConfiguration) MirrorNodeHttpConfig {
	c.transportConfiguration = config
	return c
}

// WithRetryPolicy returns a copy with the retry policy set.
func (c MirrorNodeHttpConfig) WithRetryPolicy(policy MirrorNodeHttpRetryPolicy) MirrorNodeHttpConfig {
	c.retryPolicy = policy
	return c
}

// WithRequestHeader returns a copy with a header sent on every mirror node request. The name is
// lowercased. It panics for User-Agent and x-user-agent.
func (c MirrorNodeHttpConfig) WithRequestHeader(name, value string) MirrorNodeHttpConfig {
	if lowered := strings.ToLower(strings.TrimSpace(name)); lowered == "user-agent" || lowered == httpHeaderXUserAgent {
		panic(name + " is reserved for the SDK")
	}
	c.requestHeaders = withHeader(c.requestHeaders, name, value)
	return c
}
