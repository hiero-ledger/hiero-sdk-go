//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitMirrorRestPathAcceptsRelativePaths(t *testing.T) {
	t.Parallel()

	accepted := []string{
		testMirrorAccountPath,
		"/accounts/0.0.1?limit=25",
		"/accounts/0.0.1/tokens?token.id=gt:0.0.5&order=asc",
		"/",
		"/accounts/0x00000000000000000000000000000000000004d2",
	}

	for _, raw := range accepted {
		path, err := newMirrorRestPath(raw)
		require.NoError(t, err, "%q should be accepted", raw)
		assert.Equal(t, raw, path.String())
	}
}

func TestUnitMirrorRestPathRejectsAnythingNamingAHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
	}{
		{name: "empty", raw: ""},
		{name: "no leading slash", raw: "accounts/0.0.1"},
		{name: "absolute https", raw: "https://evil.example.com/steal"},
		{name: "absolute http", raw: "http://evil.example.com/steal"},
		{name: "protocol relative", raw: "//evil.example.com/steal"},
		{name: "climbs out of the api prefix", raw: "/accounts/../../etc/passwd"},
		{name: "embedded space", raw: "/accounts/0.0.1 ?limit=1"},
		{name: "embedded newline", raw: "/accounts\n/0.0.1"},
		{name: "embedded tab", raw: "/accounts\t0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := newMirrorRestPath(tt.raw)
			require.Error(t, err, "%q must not be constructible", tt.raw)
		})
	}
}

func TestUnitResolveMirrorPathConcatenates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		base string
		path string
		want string
	}{
		{
			base: "https://mainnet.mirrornode.hedera.com/api/v1",
			path: testMirrorAccountPath,
			want: "https://mainnet.mirrornode.hedera.com/api/v1/accounts/0.0.1",
		},
		{
			base: "https://mainnet.mirrornode.hedera.com/api/v1/",
			path: testMirrorAccountPath,
			want: "https://mainnet.mirrornode.hedera.com/api/v1/accounts/0.0.1",
		},
		{
			base: "http://localhost:8084/api/v1",
			path: "/accounts/0.0.1?limit=25",
			want: "http://localhost:8084/api/v1/accounts/0.0.1?limit=25",
		},
	}

	for _, tt := range tests {
		path, err := newMirrorRestPath(tt.path)
		require.NoError(t, err)
		assert.Equal(t, tt.want, resolveMirrorPath(tt.base, path))
	}
}

// Concatenation is chosen over url.ResolveReference specifically because reference resolution
// lets a hostile value replace the host. This pins the difference so nobody "simplifies" it
// back to ResolveReference later.
func TestUnitResolveMirrorPathCannotBeRedirectedOffHost(t *testing.T) {
	t.Parallel()

	const base = "https://mainnet.mirrornode.hedera.com/api/v1"

	baseURL, err := url.Parse(base)
	require.NoError(t, err)

	hostile := "//evil.example.com/steal"

	// What url.ResolveReference does with it — the behaviour we are avoiding.
	reference, err := url.Parse(hostile)
	require.NoError(t, err)
	assert.Equal(t, "evil.example.com", baseURL.ResolveReference(reference).Host,
		"reference resolution hands the host to the attacker")

	// What this layer does with it: refuses to build the request at all.
	_, err = newMirrorRestPath(hostile)
	require.Error(t, err)

	// And even if the guard were bypassed, concatenation keeps it on the configured host.
	assert.Equal(t,
		"https://mainnet.mirrornode.hedera.com/api/v1//evil.example.com/steal",
		resolveMirrorPath(base, mirrorRestPath(hostile)))
}

func TestUnitNextPagePathStripsAPIVersionPrefix(t *testing.T) {
	t.Parallel()

	// The mirror node returns links.next including /api/v1, which the base URL also carries.
	path, err := nextPagePath("/api/v1/accounts/0.0.1/tokens?token.id=gt:0.0.5")
	require.NoError(t, err)
	assert.Equal(t, "/accounts/0.0.1/tokens?token.id=gt:0.0.5", path.String())

	resolved := resolveMirrorPath("https://mainnet.mirrornode.hedera.com/api/v1", path)
	assert.Equal(t, "https://mainnet.mirrornode.hedera.com/api/v1/accounts/0.0.1/tokens?token.id=gt:0.0.5", resolved)
}

func TestUnitNextPagePathKeepsPrefixlessPaths(t *testing.T) {
	t.Parallel()

	path, err := nextPagePath("/accounts/0.0.1/tokens?limit=2")
	require.NoError(t, err)
	assert.Equal(t, "/accounts/0.0.1/tokens?limit=2", path.String())
}

func TestUnitNextPagePathRejectsHostileCursors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		next string
	}{
		{name: "empty", next: ""},
		{name: "whitespace only", next: "   "},
		{name: "absolute https", next: "https://evil.example.com/steal"},
		{name: "absolute http", next: "http://evil.example.com/api/v1/accounts"},
		{name: "protocol relative", next: "//evil.example.com/api/v1/accounts"},
		{name: "relative without slash", next: "accounts?limit=1"},
		{name: "path traversal", next: "/api/v1/../../secrets"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := nextPagePath(tt.next)
			require.Error(t, err, "%q must be rejected", tt.next)
		})
	}
}
