//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

func TestIntegrationPrngTransactionCanExecute(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	rang := uint32(100)

	resp, err := NewPrngTransaction().
		SetRange(rang).
		Execute(env.Client)
	require.NoError(t, err)

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)
	require.NotNil(t, record.PrngNumber)
}
