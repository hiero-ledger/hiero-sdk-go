package hiero

// SPDX-License-Identifier: Apache-2.0

import "context"

// Cancellation lets a caller end an HTTP exchange early. It wraps a context.Context; the zero value never fires.
type Cancellation struct {
	ctx context.Context
}

// CancellationNone returns a Cancellation that never fires.
func CancellationNone() Cancellation {
	return Cancellation{ctx: context.Background()}
}

// CancellationFromContext returns a Cancellation that fires when ctx is done.
func CancellationFromContext(ctx context.Context) Cancellation {
	return Cancellation{ctx: ctx}
}

// GetContext returns the context.Context the cancellation wraps.
func (c Cancellation) GetContext() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

// IsCancelled reports whether the cancellation has fired.
func (c Cancellation) IsCancelled() bool {
	return c.GetContext().Err() != nil
}

// OnCancel runs callback on its own goroutine once the cancellation fires, or at once if it has.
// Call the returned release func when the exchange ends. A panic in callback is recovered.
func (c Cancellation) OnCancel(callback func()) (release func()) {
	stop := context.AfterFunc(c.GetContext(), func() {
		defer func() { _ = recover() }()
		callback()
	})

	return func() { stop() }
}
