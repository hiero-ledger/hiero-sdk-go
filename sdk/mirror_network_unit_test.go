//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitMirrorNetworkSelectsNodesRoundRobin(t *testing.T) {
	t.Parallel()

	network := _NewMirrorNetwork()
	require.NoError(t, network._SetNetwork([]string{"a.example:443", "b.example:443", "c.example:443"}))

	first := make([]string, 0, 3)
	for range 3 {
		node, err := network._GetNextMirrorNode()
		require.NoError(t, err)
		first = append(first, node._GetAddress())
	}
	assert.ElementsMatch(t, []string{"a.example:443", "b.example:443", "c.example:443"}, first)

	for _, expected := range first {
		node, err := network._GetNextMirrorNode()
		require.NoError(t, err)
		assert.Equal(t, expected, node._GetAddress())
	}
}

func TestUnitMirrorNetworkAdvancesCursorConcurrently(t *testing.T) {
	t.Parallel()

	network := _NewMirrorNetwork()
	require.NoError(t, network._SetNetwork([]string{"a.example:443", "b.example:443"}))

	const callers = 64
	var wg sync.WaitGroup
	wg.Add(callers)
	for range callers {
		go func() {
			defer wg.Done()
			_, err := network._GetNextMirrorNode()
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	assert.Equal(t, uint64(callers), network.nextNode.Load())
}

func TestUnitMirrorNetworkGetNextMirrorNodeErrorsWithNoNodes(t *testing.T) {
	t.Parallel()

	network := _NewMirrorNetwork()

	_, err := network._GetNextMirrorNode()
	assert.ErrorContains(t, err, "no healthy nodes")
}
