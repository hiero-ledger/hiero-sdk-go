//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitEvmHookCallNew(t *testing.T) {
	t.Parallel()

	evmHookCall := NewEvmHookCall()
	require.NotNil(t, evmHookCall)
	require.Empty(t, evmHookCall.GetData())
	require.Equal(t, uint64(0), evmHookCall.GetGasLimit())
}

func TestUnitEvmHookCallSetData(t *testing.T) {
	t.Parallel()

	evmHookCall := NewEvmHookCall()
	testData := []byte{0x01, 0x02, 0x03, 0x04}

	evmHookCall.SetData(testData)
	require.Equal(t, testData, evmHookCall.GetData())
}

func TestUnitEvmHookCallSetGasLimit(t *testing.T) {
	t.Parallel()

	evmHookCall := NewEvmHookCall()
	gasLimit := uint64(25000)

	evmHookCall.SetGasLimit(gasLimit)
	require.Equal(t, gasLimit, evmHookCall.GetGasLimit())
}

func TestUnitEvmHookCallChaining(t *testing.T) {
	t.Parallel()

	testData := []byte{0x01, 0x02, 0x03}
	gasLimit := uint64(30000)

	evmHookCall := NewEvmHookCall().
		SetData(testData).
		SetGasLimit(gasLimit)

	require.Equal(t, testData, evmHookCall.GetData())
	require.Equal(t, gasLimit, evmHookCall.GetGasLimit())
}

func TestUnitEvmHookCallToProtobuf(t *testing.T) {
	t.Parallel()

	testData := []byte{0x11, 0x22, 0x33, 0x44, 0x55}
	gasLimit := uint64(50000)

	evmHookCall := NewEvmHookCall().
		SetData(testData).
		SetGasLimit(gasLimit)

	proto := evmHookCall.toProtobuf()
	require.NotNil(t, proto)
	require.Equal(t, testData, proto.Data)
	require.Equal(t, gasLimit, proto.GasLimit)
}

func TestUnitEvmHookCallFromProtobuf(t *testing.T) {
	t.Parallel()

	testData := []byte{0xaa, 0xbb, 0xcc}
	gasLimit := uint64(60000)

	proto := &services.EvmHookCall{
		Data:     testData,
		GasLimit: gasLimit,
	}

	evmHookCall := evmHookCallFromProtobuf(proto)
	require.Equal(t, testData, evmHookCall.GetData())
	require.Equal(t, gasLimit, evmHookCall.GetGasLimit())
}

func TestUnitEvmHookCallRoundTripProtobuf(t *testing.T) {
	t.Parallel()

	testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	gasLimit := uint64(75000)

	// Create original
	original := NewEvmHookCall().
		SetData(testData).
		SetGasLimit(gasLimit)

	// Convert to protobuf
	proto := original.toProtobuf()
	require.NotNil(t, proto)

	// Convert back from protobuf
	reconstructed := evmHookCallFromProtobuf(proto)

	// Verify round-trip
	require.Equal(t, original.GetData(), reconstructed.GetData())
	require.Equal(t, original.GetGasLimit(), reconstructed.GetGasLimit())
}

func TestUnitEvmHookCallEmptyData(t *testing.T) {
	t.Parallel()

	evmHookCall := NewEvmHookCall().
		SetData([]byte{}).
		SetGasLimit(10000)

	proto := evmHookCall.toProtobuf()
	require.NotNil(t, proto)
	require.Empty(t, proto.Data)
	require.Equal(t, uint64(10000), proto.GasLimit)

	reconstructed := evmHookCallFromProtobuf(proto)
	require.Empty(t, reconstructed.GetData())
	require.Equal(t, uint64(10000), reconstructed.GetGasLimit())
}

func TestUnitEvmHookCallNilProtobuf(t *testing.T) {
	t.Parallel()

	// Test that fromProtobuf handles nil gracefully
	reconstructed := evmHookCallFromProtobuf(nil)
	require.Empty(t, reconstructed.GetData())
	require.Equal(t, uint64(0), reconstructed.GetGasLimit())
}

func TestUnitEvmHookCallMultipleRoundTrips(t *testing.T) {
	t.Parallel()

	testData := []byte{0xde, 0xad, 0xbe, 0xef}
	gasLimit := uint64(42000)

	original := NewEvmHookCall().
		SetData(testData).
		SetGasLimit(gasLimit)

	// First round trip
	proto1 := original.toProtobuf()
	reconstructed1 := evmHookCallFromProtobuf(proto1)

	// Second round trip
	proto2 := reconstructed1.toProtobuf()
	reconstructed2 := evmHookCallFromProtobuf(proto2)

	// Third round trip
	proto3 := reconstructed2.toProtobuf()
	reconstructed3 := evmHookCallFromProtobuf(proto3)

	// All should be equal
	assert.Equal(t, original.GetData(), reconstructed1.GetData())
	assert.Equal(t, original.GetData(), reconstructed2.GetData())
	assert.Equal(t, original.GetData(), reconstructed3.GetData())
	assert.Equal(t, original.GetGasLimit(), reconstructed1.GetGasLimit())
	assert.Equal(t, original.GetGasLimit(), reconstructed2.GetGasLimit())
	assert.Equal(t, original.GetGasLimit(), reconstructed3.GetGasLimit())
}

func TestUnitEvmHookCallUpdateData(t *testing.T) {
	t.Parallel()

	evmHookCall := NewEvmHookCall()

	data1 := []byte{0x01, 0x02}
	evmHookCall.SetData(data1)
	require.Equal(t, data1, evmHookCall.GetData())

	data2 := []byte{0x03, 0x04, 0x05}
	evmHookCall.SetData(data2)
	require.Equal(t, data2, evmHookCall.GetData())
	require.NotEqual(t, data1, evmHookCall.GetData())
}

func TestUnitEvmHookCallUpdateGasLimit(t *testing.T) {
	t.Parallel()

	evmHookCall := NewEvmHookCall()

	evmHookCall.SetGasLimit(10000)
	require.Equal(t, uint64(10000), evmHookCall.GetGasLimit())

	evmHookCall.SetGasLimit(50000)
	require.Equal(t, uint64(50000), evmHookCall.GetGasLimit())
}
