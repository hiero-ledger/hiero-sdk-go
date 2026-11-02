//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
)

func TestUnitRequestTypeFromProtoToString(t *testing.T) {
	t.Parallel()

	for _, code := range services.HederaFunctionality_value {
		requestType := RequestType(code)
		assert.NotPanics(t, func() { _ = requestType.String() })
	}
}

func TestUnitRequestTypeClprStrings(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "LAMBDA_S_STORE", RequestTypeLambdaSStore.String())
	assert.Equal(t, "CLPR_REGISTER_CHANNEL", RequestTypeClprRegisterChannel.String())
	assert.Equal(t, "CLPR_MAINTENANCE", RequestTypeClprMaintenance.String())
	assert.Equal(t, RequestType(services.HederaFunctionality_ClprGetEndpointManifest), RequestTypeClprGetEndpointManifest)
}
