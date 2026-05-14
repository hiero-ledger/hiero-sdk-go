//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
)

func TestUnitContractLogInfo_ProtobufRoundTrip(t *testing.T) {
	t.Parallel()

	logInfo := ContractLogInfo{
		ContractID: ContractID{Contract: 7},
		Bloom:      []byte{1, 2, 3},
		Topics:     [][]byte{{4, 5}, {6, 7, 8}},
		Data:       []byte{9, 10, 11},
	}

	actual := _ContractLogInfoFromProtobuf(logInfo._ToProtobuf())

	assert.Equal(t, logInfo, actual)
}

func TestUnitContractLogInfo_FromProtobufNil(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ContractLogInfo{}, _ContractLogInfoFromProtobuf(nil))
}

func TestUnitContractLogInfo_FromProtobufNilContractID(t *testing.T) {
	t.Parallel()

	pb := &services.ContractLoginfo{
		ContractID: nil,
		Bloom:      []byte{1, 2, 3},
		Topic:      [][]byte{{4, 5}, {6}},
		Data:       []byte{7, 8},
	}

	actual := _ContractLogInfoFromProtobuf(pb)

	assert.Equal(t, ContractID{}, actual.ContractID)
	assert.Equal(t, []byte{1, 2, 3}, actual.Bloom)
	assert.Equal(t, [][]byte{{4, 5}, {6}}, actual.Topics)
	assert.Equal(t, []byte{7, 8}, actual.Data)
}

func TestUnitContractLogInfo_EmptyTopicsRoundTrip(t *testing.T) {
	t.Parallel()

	logInfo := ContractLogInfo{
		ContractID: ContractID{Contract: 9},
		Bloom:      []byte{1},
		Topics:     [][]byte{},
		Data:       []byte{2},
	}

	actual := _ContractLogInfoFromProtobuf(logInfo._ToProtobuf())

	assert.NotNil(t, actual.Topics)
	assert.Equal(t, 0, len(actual.Topics))
	assert.Equal(t, logInfo, actual)
}
