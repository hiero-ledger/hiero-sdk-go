//go:build !unix

package methods

// SPDX-License-Identifier: Apache-2.0

// captureStdout only calls fn: stdout redirection is Unix-only, so this platform reports no warning.
func captureStdout[T any](fn func() T) (T, []byte, error) {
	return fn(), nil, nil
}
