//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
)

func TestUnitClprTransactionFromBytesNamesUnsupportedBodies(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		body *services.TransactionBody
	}{
		{"clpr_endpoint_publication", &services.TransactionBody{Data: &services.TransactionBody_ClprEndpointPublication{ClprEndpointPublication: &services.ClprEndpointPublicationTransactionBody{}}}},
		{"clpr_maintenance", &services.TransactionBody{Data: &services.TransactionBody_ClprMaintenance{ClprMaintenance: &services.ClprMaintenanceTransactionBody{}}}},
		{"clprEvictPeerEndpoint", &services.TransactionBody{Data: &services.TransactionBody_ClprEvictPeerEndpoint{ClprEvictPeerEndpoint: &services.ClprEvictPeerEndpointTransactionBody{}}}},
	} {
		_, err := TransactionFromBytes(clprUnsignedBytes(t, tt.body))
		assert.ErrorIs(t, err, errClprBodyUnsupported)
		assert.ErrorContains(t, err, "does not support "+tt.name)
	}
}
