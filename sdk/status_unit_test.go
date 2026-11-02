//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
)

// TestUnitStatusFromProtoToString tests pulling all codes from the proto generated code,
// converting it to the sdk enum, and calling String()
//
// Ideally this will catch any changes to _Response codes when the protobufs get updated
func TestUnitStatusFromProtoToString(t *testing.T) {
	t.Parallel()

	for _, code := range services.ResponseCodeEnum_value {
		status := Status(code)
		assert.NotPanics(t, func() { _ = status.String() })
	}
}

func TestUnitStatusClprStrings(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "CONFIG_FILE_PART_UPLOADED", CONFIG_FILE_PART_UPLOADED.String())
	assert.Equal(t, "CLPR_NOT_ENABLED", CLPR_NOT_ENABLED.String())
	assert.Equal(t, "CLPR_NO_PROGRESS", CLPR_NO_PROGRESS.String())
	assert.Equal(t, Status(services.ResponseCodeEnum_CLPR_EMPTY_BUNDLE), CLPR_EMPTY_BUNDLE)
}
