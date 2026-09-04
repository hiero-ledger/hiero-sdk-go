package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
)

// The mirror REST API version segment, which the base URL already carries.
const mirrorHttpAPIVersionPrefix = "/api/v1"

// mirrorRestPath is a mirror REST path — never an absolute URL. Only newMirrorRestPath can
// produce one, so "no host names above this layer" is a property of the type instead of a
// convention a call site can quietly break. An accidentally absolute "https://…" is rejected
// because it has no leading slash.
//
// This is the structural half of the caller-supplied-cursor problem: a request that could
// address another host cannot be constructed, so nothing downstream has to check for one.
type mirrorRestPath string

func newMirrorRestPath(path string) (mirrorRestPath, error) {
	if path == "" {
		return "", fmt.Errorf("mirror node REST path is empty")
	}
	if !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("mirror node REST path %q must start with %q — absolute URLs are not accepted here", path, "/")
	}
	if strings.ContainsAny(path, " \t\r\n") {
		return "", fmt.Errorf("mirror node REST path %q contains whitespace", path)
	}
	// A "//host" path is harmless under concatenation, but it is never what a caller meant.
	if strings.HasPrefix(path, "//") {
		return "", fmt.Errorf("mirror node REST path %q looks protocol-relative", path)
	}
	if hasDotDotSegment(path) {
		return "", fmt.Errorf("mirror node REST path %q contains a %q segment", path, "..")
	}

	return mirrorRestPath(path), nil
}

func (p mirrorRestPath) String() string {
	return string(p)
}

// hasDotDotSegment reports whether the path could climb out of the API version prefix.
func hasDotDotSegment(path string) bool {
	for segment := range strings.SplitSeq(path, "/") {
		if segment == ".." {
			return true
		}
	}

	return false
}

// resolveMirrorPath joins a mirror base URL and a path by plain concatenation, so both
// spellings of a trailing slash behave identically.
//
// Concatenation rather than URL reference resolution is the point: url.URL.ResolveReference
// lets an absolute or protocol-relative reference replace the host outright, which is exactly
// the failure mode this layer exists to make impossible.
func resolveMirrorPath(baseURL string, path mirrorRestPath) string {
	return strings.TrimSuffix(baseURL, "/") + path.String()
}

// nextPagePath converts a mirror node links.next value into a path. The mirror node returns it
// including the API version prefix, which the base URL also carries, so the prefix is stripped
// to keep resolution plain concatenation.
//
// Anything naming a host is rejected rather than trusted, which is what keeps a caller- or
// server-supplied cursor from redirecting a request carrying mirror credentials.
func nextPagePath(next string) (mirrorRestPath, error) {
	trimmed := strings.TrimSpace(next)
	if trimmed == "" {
		return "", fmt.Errorf("pagination next link is empty")
	}
	if strings.Contains(trimmed, "://") {
		return "", fmt.Errorf("pagination next link %q is absolute; only same-mirror paths are followed", next)
	}

	path, err := newMirrorRestPath(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid pagination next link: %w", err)
	}

	stripped := strings.TrimPrefix(path.String(), mirrorHttpAPIVersionPrefix)
	if stripped == path.String() {
		return path, nil
	}

	return newMirrorRestPath(stripped)
}
