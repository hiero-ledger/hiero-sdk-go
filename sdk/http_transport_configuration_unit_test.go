//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnitHttpTransportConfigurationZeroValueIsDefault(t *testing.T) {
	t.Parallel()

	var config HttpTransportConfiguration

	assert.Equal(t, DefaultHttpTransportConfiguration(), config)
	assert.Zero(t, config.GetConnectTimeout())
	assert.Equal(t, uint16(5), config.GetMaxRedirects())
	assert.Equal(t, int64(32<<20), config.GetMaxResponseBytes())
	assert.Empty(t, config.GetDefaultHeaders())
}

func TestUnitHttpTransportConfigurationWithChangesOneField(t *testing.T) {
	t.Parallel()

	original := DefaultHttpTransportConfiguration()
	changed := original.WithMaxRedirects(0)

	assert.Zero(t, changed.GetMaxRedirects())
	assert.Equal(t, uint16(5), original.GetMaxRedirects())
	assert.Equal(t, int64(32<<20), changed.GetMaxResponseBytes())
	assert.Equal(t, 2*time.Second, original.WithConnectTimeout(2*time.Second).GetConnectTimeout())
}

func TestUnitHttpTransportConfigurationDefaultHeadersAreLowercasedCopies(t *testing.T) {
	t.Parallel()

	first := DefaultHttpTransportConfiguration().WithDefaultHeader("X-Trace", "a")
	second := first.WithDefaultHeader("x-trace", "b")

	assert.Equal(t, map[string]string{"x-trace": "a"}, first.GetDefaultHeaders())
	assert.Equal(t, map[string]string{"x-trace": "b"}, second.GetDefaultHeaders())

	read := second.GetDefaultHeaders()
	read["x-trace"] = "mutated"
	assert.Equal(t, "b", second.GetDefaultHeaders()["x-trace"])
}

func TestUnitHttpTransportConfigurationRejectsInvalidArguments(t *testing.T) {
	t.Parallel()

	config := DefaultHttpTransportConfiguration()

	assert.Panics(t, func() { config.WithConnectTimeout(-time.Second) })
	assert.Panics(t, func() { config.WithMaxResponseBytes(0) })
	assert.Panics(t, func() { config.WithDefaultHeader(" ", "value") })
}
