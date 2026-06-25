//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"math/big"
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

func TestUnitEthereumTransactionDataGetSetDataPerVariant(t *testing.T) {
	t.Parallel()

	variants := []EthereumTransactionBody{
		(&EthereumEIP1559Transaction{}).SetCallData([]byte{0x01}),
		(&EthereumEIP2930Transaction{}).SetCallData([]byte{0x02}),
		(&EthereumEIP7702Transaction{}).SetCallData([]byte{0x03}),
		(&EthereumLegacyTransaction{}).SetCallData([]byte{0x04}),
	}

	for _, v := range variants {
		data := NewEthereumTransactionData(v)
		require.NotNil(t, data)
		require.NotEmpty(t, data.GetData())

		data.SetData([]byte{0xFF})
		assert.Equal(t, []byte{0xFF}, data.GetData())
	}
}

func TestUnitEthereumTransactionDataFromBytesDispatch(t *testing.T) {
	t.Parallel()
	key := newTestEcdsaKey(t)
	to := make([]byte, 20)

	eip1559, err := (&EthereumEIP1559Transaction{}).SetChainId(1).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(0)).Sign(key)
	require.NoError(t, err)
	eip2930, err := (&EthereumEIP2930Transaction{}).SetChainId(1).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(0)).Sign(key)
	require.NoError(t, err)
	eip7702, err := (&EthereumEIP7702Transaction{}).SetChainId(1).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(0)).Sign(key)
	require.NoError(t, err)
	legacy, err := (&EthereumLegacyTransaction{}).SetNonce(0).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(0)).Sign(key)
	require.NoError(t, err)

	cases := []struct {
		name    string
		payload []byte
		isType  func(EthereumTransactionBody) bool
	}{
		{"eip1559", eip1559, func(b EthereumTransactionBody) bool { _, ok := b.(*EthereumEIP1559Transaction); return ok }},
		{"eip2930", eip2930, func(b EthereumTransactionBody) bool { _, ok := b.(*EthereumEIP2930Transaction); return ok }},
		{"eip7702", eip7702, func(b EthereumTransactionBody) bool { _, ok := b.(*EthereumEIP7702Transaction); return ok }},
		{"legacy", legacy, func(b EthereumTransactionBody) bool { _, ok := b.(*EthereumLegacyTransaction); return ok }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := EthereumTransactionDataFromBytes(tc.payload)
			require.NoError(t, err)
			assert.True(t, tc.isType(data.GetTransaction()), "unexpected concrete type for %s", tc.name)

			// The wrapper's ToBytes round-trips back to the original payload.
			out, err := data.ToBytes()
			require.NoError(t, err)
			assert.Equal(t, tc.payload, out)
		})
	}
}

func TestUnitEthereumTransactionDataEmpty(t *testing.T) {
	t.Parallel()
	key := newTestEcdsaKey(t)

	empty := &EthereumTransactionData{}
	assert.Nil(t, empty.GetTransaction())

	_, err := empty.ToBytes()
	assert.Error(t, err)

	_, err = empty.Sign(key)
	assert.Error(t, err)

	_, err = EthereumTransactionDataFromBytes(nil)
	assert.Error(t, err)

	// A recognized prefix (0x01) wrapping a valid but wrong-length RLP list
	// surfaces the variant decode error through the dispatcher.
	list := NewRLPItem(LIST_TYPE)
	list.PushBack(NewRLPItem(VALUE_TYPE).AssignValue([]byte{0x01}))
	body, err := list.Write()
	require.NoError(t, err)
	_, err = EthereumTransactionDataFromBytes(append([]byte{0x01}, body...))
	assert.Error(t, err)
}

func TestUnitEthereumSetEthereumDataFromBodyAllVariants(t *testing.T) {
	t.Parallel()
	key := newTestEcdsaKey(t)
	to := make([]byte, 20)

	bodies := []EthereumTransactionBody{
		(&EthereumEIP1559Transaction{}).SetChainId(1).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(0)),
		(&EthereumEIP2930Transaction{}).SetChainId(1).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(0)),
		(&EthereumEIP7702Transaction{}).SetChainId(1).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(0)),
		(&EthereumLegacyTransaction{}).SetNonce(0).SetGasLimit(21_000).SetTo(to).SetValue(big.NewInt(0)),
	}

	for _, body := range bodies {
		_, err := body.Sign(key)
		require.NoError(t, err)

		tx, err := NewEthereumTransaction().SetEthereumDataFromBody(body)
		require.NoError(t, err)
		assert.NotEmpty(t, tx.GetEthereumData())
	}
}
