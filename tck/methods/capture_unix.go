//go:build unix

package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"io"
	"os"
	"sync"

	"golang.org/x/sys/unix"
)

// stdoutMu serializes captureStdout, because file descriptor 1 is shared by the whole process
var stdoutMu sync.Mutex

// captureStdout calls fn while file descriptor 1 is redirected into a pipe and returns what was written
// to it. The SDK's package level logger holds os.Stdout itself and exposes no writer to swap, so the
// descriptor is the only place its output can be intercepted. The output is passed on to stdout afterwards.
func captureStdout[T any](fn func() T) (T, []byte, error) {
	stdoutMu.Lock()
	defer stdoutMu.Unlock()

	var result T
	r, w, err := os.Pipe()
	if err != nil {
		return result, nil, err
	}
	defer r.Close()
	defer w.Close()

	saved, err := unix.Dup(1)
	if err != nil {
		return result, nil, err
	}
	defer unix.Close(saved)

	if err := unix.Dup2(int(w.Fd()), 1); err != nil {
		return result, nil, err
	}
	// Restores stdout on every path, including a panic in fn, before the pipe is closed
	defer func() { _ = unix.Dup2(saved, 1) }()

	captured := make(chan []byte, 1)
	go func() {
		out, _ := io.ReadAll(r)
		captured <- out
	}()

	result = fn()

	if err := unix.Dup2(saved, 1); err != nil {
		return result, nil, err
	}
	// Both write ends are closed now, so the reader sees EOF
	w.Close()
	out := <-captured

	_, _ = os.Stdout.Write(out)
	return result, out, nil
}
