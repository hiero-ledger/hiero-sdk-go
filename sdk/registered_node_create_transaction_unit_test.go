//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
)

func TestUnitRegisteredNodeCreateTransactionBuild(t *testing.T) {
	t.Parallel()

	key, _ := PrivateKeyGenerateEd25519()
	endpoint := &BlockNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress:   []byte{10, 0, 0, 1},
			port:        8080,
			requiresTls: true,
		},
		endpointApis: []BlockNodeApi{BlockNodeApiStatus},
	}

	tx := NewRegisteredNodeCreateTransaction().
		SetAdminKey(key.PublicKey()).
		SetDescription("Test Node").
		AddServiceEndpoint(endpoint)

	body := tx.buildProtoBody()
	assert.NotNil(t, body.AdminKey)
	assert.Equal(t, "Test Node", body.Description)
	assert.Len(t, body.ServiceEndpoint, 1)
}

func TestUnitRegisteredNodeCreateTransactionRoundTrip(t *testing.T) {
	t.Parallel()

	key, _ := PrivateKeyGenerateEd25519()
	endpoint := &MirrorNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress:   []byte{192, 168, 1, 1},
			port:        443,
			requiresTls: true,
		},
	}

	tx := NewRegisteredNodeCreateTransaction().
		SetAdminKey(key.PublicKey()).
		SetDescription("My Node").
		AddServiceEndpoint(endpoint)

	body := tx.buildProtoBody()

	pbBody := &services.TransactionBody{
		Data: &services.TransactionBody_RegisteredNodeCreate{
			RegisteredNodeCreate: body,
		},
	}

	restored := _RegisteredNodeCreateTransactionFromProtobuf(*tx.Transaction, pbBody)
	assert.Equal(t, "My Node", restored.GetDescription())
	assert.Len(t, restored.GetServiceEndpoints(), 1)

	_, ok := restored.GetServiceEndpoints()[0].(*MirrorNodeServiceEndpoint)
	assert.True(t, ok, "expected *MirrorNodeServiceEndpoint after round-trip")
}

func TestUnitRegisteredNodeCreateTransactionGetMethod(t *testing.T) {
	t.Parallel()

	tx := NewRegisteredNodeCreateTransaction()
	assert.Equal(t, "RegisteredNodeCreateTransaction", tx.getName())
}

func TestUnitRegisteredNodeCreateTransactionGettersSetters(t *testing.T) {
	t.Parallel()

	key, _ := PrivateKeyGenerateEd25519()
	endpoint := &RpcRelayServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			domainName: "node.example.com",
			port:       443,
		},
	}

	tx := NewRegisteredNodeCreateTransaction().
		SetAdminKey(key.PublicKey()).
		SetDescription("Test").
		SetServiceEndpoints([]RegisteredServiceEndpoint{endpoint})

	assert.Equal(t, key.PublicKey(), tx.GetAdminKey())
	assert.Equal(t, "Test", tx.GetDescription())
	assert.Len(t, tx.GetServiceEndpoints(), 1)
}
