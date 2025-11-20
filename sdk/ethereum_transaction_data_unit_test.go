//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitEthereumEIP2930TransactionData(t *testing.T) {
	t.Parallel()

	validRlpHexes := []string{
		"01f86e82012a0585a54f4c3c00832dc6c094000000000000000000000000000000000000041d8502540be40080c001a017138522c841d5864818ab16f5d5dc7e353009c7a9ecb951e84a37991d14e47aa04bb17a3123f800f63e881dd0fa753bbd8b41af0093d85c3a35e66d9ba66b0bee",
		"01f86c800685a54f4c3c00832dc6c094b19f7459342cf8c4fbd5053f5ef6c5981afa52d68502540be40080c080a07c66e585eda1a789ff33cfbd5d4b01c42b86ef64ef82097f20f1e700b7c57685a065a34302820466199cce5ed4e66e2bd59c93032ca5c4478fc9925c10972223b3",
		"01f86e82012a0485a54f4c3c00832dc6c094d19e85e39e96020bb5eb6b5806784b5a09b9254f8502540be40080c001a086d729022621ac1e4411d4ea18b43763970a2f75e3a69ac7d641b129f57c2623a006930b9d645c7d035c7e76c7ad2410911f6a9c515e926c54dbf18d94b9bcaa86",
	}

	var validRlpBytes [][]byte
	for _, rlpHex := range validRlpHexes {
		rlpBytes, err := hex.DecodeString(rlpHex)
		assert.NoError(t, err)
		validRlpBytes = append(validRlpBytes, rlpBytes)
	}

	var eip2930TsData []*EthereumTransactionData
	for _, rlpBytes := range validRlpBytes {
		eip2930TData, err := EthereumTransactionDataFromBytes(rlpBytes)
		assert.NoError(t, err)
		eip2930TsData = append(eip2930TsData, eip2930TData)
	}

	var validRlpHexesOut []string
	for _, eip2930TData := range eip2930TsData {
		eip2930DataBytes, err := eip2930TData.ToBytes()
		assert.NoError(t, err)
		validRlpHexesOut = append(validRlpHexesOut, hex.EncodeToString(eip2930DataBytes))
	}

	for i, validRlpHex := range validRlpHexes {
		require.Equal(t, validRlpHex, validRlpHexesOut[i])
	}
}
