package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"errors"
	"fmt"
	"strings"
)

// The mirror REST API version segment, which the base URL already carries.
const mirrorHttpAPIVersionPrefix = "/api/v1"

// mirrorNodeRestPath is a path relative to the mirror node REST base URL. Create one with newMirrorNodeRestPath.
type mirrorNodeRestPath string

func newMirrorNodeRestPath(path string) (mirrorNodeRestPath, error) {
	if path == "" {
		return "", errors.New("mirror node REST path is empty")
	}
	if !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("mirror node REST path %q must start with %q", path, "/")
	}
	if strings.ContainsAny(path, " \t\r\n") {
		return "", fmt.Errorf("mirror node REST path %q contains whitespace", path)
	}
	if strings.HasPrefix(path, "//") {
		return "", fmt.Errorf("mirror node REST path %q looks protocol-relative", path)
	}
	if hasDotDotSegment(path) {
		return "", fmt.Errorf("mirror node REST path %q contains a %q segment", path, "..")
	}

	return mirrorNodeRestPath(path), nil
}

func (p mirrorNodeRestPath) String() string {
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

// resolveMirrorPath appends path to baseURL. It does not use url.ResolveReference, which would let path replace the host.
func resolveMirrorPath(baseURL string, path mirrorNodeRestPath) string {
	return strings.TrimSuffix(baseURL, "/") + path.String()
}

// nextPagePath converts a links.next value into a path, stripping the /api/v1 prefix. Absolute URLs are rejected.
func nextPagePath(next string) (mirrorNodeRestPath, error) {
	trimmed := strings.TrimSpace(next)
	if trimmed == "" {
		return "", errors.New("pagination next link is empty")
	}
	if strings.Contains(trimmed, "://") {
		return "", fmt.Errorf("pagination next link %q is absolute; only same-mirror paths are followed", next)
	}

	path, err := newMirrorNodeRestPath(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid pagination next link: %w", err)
	}

	stripped := strings.TrimPrefix(path.String(), mirrorHttpAPIVersionPrefix)
	if stripped == path.String() {
		return path, nil
	}

	return newMirrorNodeRestPath(stripped)
}
