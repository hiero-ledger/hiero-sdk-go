//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireClosedSoon(t *testing.T, ch <-chan struct{}, msg string) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		require.FailNow(t, msg)
	}
}

func TestUnitHttpCancellationZeroValueIsNone(t *testing.T) {
	t.Parallel()

	var zero Cancellation

	assert.False(t, zero.IsCancelled())
	assert.False(t, CancellationNone().IsCancelled())
	assert.NotNil(t, zero.GetContext())
}

func TestUnitHttpCancellationIsCancelledBeforeCallbacksRun(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancellation := CancellationFromContext(ctx)

	seen := make(chan bool, 1)
	release := cancellation.OnCancel(func() { seen <- cancellation.IsCancelled() })
	defer release()

	cancel()
	cancel()

	select {
	case cancelled := <-seen:
		assert.True(t, cancelled)
	case <-time.After(time.Second):
		require.FailNow(t, "the callback never ran")
	}
}

func TestUnitHttpCancellationOnCancelAfterCancelRunsImmediately(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ran := make(chan struct{})
	release := CancellationFromContext(ctx).OnCancel(func() { close(ran) })
	defer release()

	requireClosedSoon(t, ran, "a registration on an already-cancelled instance must run rather than never")
}

func TestUnitHttpCancellationReleasedRegistrationNeverRuns(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	ran := make(chan struct{}, 1)

	release := CancellationFromContext(ctx).OnCancel(func() { ran <- struct{}{} })
	release()
	cancel()

	select {
	case <-ran:
		assert.Fail(t, "a released registration must not run")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestUnitHttpCancellationPanickingCallbackDoesNotStopOthers(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancellation := CancellationFromContext(ctx)

	ran := make(chan struct{})
	defer cancellation.OnCancel(func() { panic("a callback that throws") })()
	defer cancellation.OnCancel(func() { close(ran) })()

	cancel()

	requireClosedSoon(t, ran, "one callback panicking must not stop another, nor crash the process")
}

// countingContext is a cancellable context that counts live AfterFunc registrations, which is
// how context.WithCancel and WithTimeout attach to a parent that is not a stdlib context.
type countingContext struct {
	done   chan struct{}
	mu     sync.Mutex
	active int
	funcs  map[int]func()
	next   int
}

func newCountingContext() *countingContext {
	return &countingContext{done: make(chan struct{}), funcs: map[int]func(){}}
}

func (c *countingContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c *countingContext) Done() <-chan struct{}       { return c.done }
func (c *countingContext) Value(any) any               { return nil }

func (c *countingContext) Err() error {
	select {
	case <-c.done:
		return context.Canceled
	default:
		return nil
	}
}

func (c *countingContext) AfterFunc(f func()) func() bool {
	c.mu.Lock()
	id := c.next
	c.next++
	c.funcs[id] = f
	c.active++
	c.mu.Unlock()

	return func() bool {
		c.mu.Lock()
		defer c.mu.Unlock()
		if _, ok := c.funcs[id]; !ok {
			return false
		}
		delete(c.funcs, id)
		c.active--
		return true
	}
}

func (c *countingContext) liveRegistrations() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.active
}
