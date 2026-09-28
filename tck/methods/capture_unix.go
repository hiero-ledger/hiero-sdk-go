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

	captured := make(chan []byte, 1)
	go func() {
		out, _ := io.ReadAll(r)
		captured <- out
	}()

	if err := runRedirected(w, func() { result = fn() }); err != nil {
		return result, nil, err
	}
	out := <-captured

	_, _ = os.Stdout.Write(out)
	return result, out, nil
}

// runRedirected runs fn with file descriptor 1 pointing at w. It restores the descriptor and closes w on
// every path, including a panic in fn, so the reader always sees EOF.
func runRedirected(w *os.File, fn func()) (err error) {
	defer w.Close()

	saved, err := unix.Dup(1)
	if err != nil {
		return err
	}
	defer unix.Close(saved)

	if err := unix.Dup2(int(w.Fd()), 1); err != nil {
		return err
	}
	defer func() {
		if restoreErr := unix.Dup2(saved, 1); err == nil {
			err = restoreErr
		}
	}()

	fn()
	return nil
}
