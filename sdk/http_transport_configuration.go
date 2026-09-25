package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"maps"
	"strings"
	"time"
)

const (
	httpDefaultMaxRedirects = 5
	// Bodies are read into memory, so cap what a mirror node can make us allocate.
	httpDefaultMaxResponseBytes = 32 << 20
)

// HttpTransportConfiguration configures the transport the SDK builds. The zero value is the default.
type HttpTransportConfiguration struct {
	set              transportConfigField
	connectTimeout   time.Duration
	maxRedirects     uint16
	maxResponseBytes int64
	defaultHeaders   map[string]string
}

type transportConfigField uint8

const (
	transportConfigMaxRedirects transportConfigField = 1 << iota
	transportConfigMaxResponseBytes
)

// DefaultHttpTransportConfiguration returns the default configuration, which is also the zero
// value.
func DefaultHttpTransportConfiguration() HttpTransportConfiguration {
	return HttpTransportConfiguration{}
}

// GetConnectTimeout returns the timeout for the TCP connect and TLS handshake. 0 keeps the net/http defaults.
func (c HttpTransportConfiguration) GetConnectTimeout() time.Duration {
	return c.connectTimeout
}

// GetMaxRedirects returns the maximum number of redirects followed. The default is 5.
func (c HttpTransportConfiguration) GetMaxRedirects() uint16 {
	if c.set&transportConfigMaxRedirects == 0 {
		return httpDefaultMaxRedirects
	}
	return c.maxRedirects
}

// GetMaxResponseBytes returns the maximum response body size. The default is 32 MiB.
func (c HttpTransportConfiguration) GetMaxResponseBytes() int64 {
	if c.set&transportConfigMaxResponseBytes == 0 {
		return httpDefaultMaxResponseBytes
	}
	return c.maxResponseBytes
}

// GetDefaultHeaders returns the headers sent with every request. Request headers override them.
func (c HttpTransportConfiguration) GetDefaultHeaders() map[string]string {
	return maps.Clone(c.defaultHeaders)
}

// WithConnectTimeout returns a copy with the connect timeout set. It panics if timeout is negative.
func (c HttpTransportConfiguration) WithConnectTimeout(timeout time.Duration) HttpTransportConfiguration {
	requireNonNegativeDuration("connectTimeout", timeout)
	c.connectTimeout = timeout
	return c
}

// WithMaxRedirects returns a copy with the redirect limit set. 0 follows no redirects.
func (c HttpTransportConfiguration) WithMaxRedirects(maxRedirects uint16) HttpTransportConfiguration {
	c.maxRedirects = maxRedirects
	c.set |= transportConfigMaxRedirects
	return c
}

// WithMaxResponseBytes returns a copy with the response size limit set. It panics if maxResponseBytes is less than 1.
func (c HttpTransportConfiguration) WithMaxResponseBytes(maxResponseBytes int64) HttpTransportConfiguration {
	if maxResponseBytes < 1 {
		panic("maxResponseBytes must be at least 1")
	}
	c.maxResponseBytes = maxResponseBytes
	c.set |= transportConfigMaxResponseBytes
	return c
}

// WithDefaultHeader returns a copy with a default header set. The name is lowercased.
func (c HttpTransportConfiguration) WithDefaultHeader(name, value string) HttpTransportConfiguration {
	c.defaultHeaders = withHeader(c.defaultHeaders, name, value)
	return c
}

// withHeader returns a copy of headers with the lowercased name set to value.
func withHeader(headers map[string]string, name, value string) map[string]string {
	if strings.TrimSpace(name) == "" {
		panic("a header name must not be empty")
	}

	copied := make(map[string]string, len(headers)+1)
	maps.Copy(copied, headers)
	copied[strings.ToLower(name)] = value

	return copied
}

func requireNonNegativeDuration(name string, value time.Duration) {
	if value < 0 {
		panic(name + " must not be negative")
	}
}
