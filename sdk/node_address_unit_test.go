//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitNodeAddressStringTest(t *testing.T) {
	t.Parallel()

	// Prepare an AccountID for the test
	id, err := AccountIDFromString("0.0.123")
	require.NoError(t, err)

	// Prepare a NodeAddress struct for testing
	nodeAddress := NodeAddress{
		PublicKey: "sample-public-key",
		AccountID: &id,
		NodeID:    1234,
		CertHash:  []byte("sample-cert-hash"),
		Addresses: []Endpoint{
			{
				address:    []byte("192.168.1.1"),
				port:       8080,
				domainName: "example.com",
			},
		},
		Description: "Sample Node",
	}

	// Generate the string representation
	result := nodeAddress.String()

	// Check if fields are present in the result
	require.Contains(t, result, "NodeAccountId: 0.0.123")
	require.Contains(t, result, "CertHash: sample-cert-hash")
	require.Contains(t, result, "NodeId: 1234")
	require.Contains(t, result, "PubKey: sample-public-key")
	require.Contains(t, result, "example.com")
}
