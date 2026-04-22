//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitRegisteredNodeAddressBookQuerySetterChaining(t *testing.T) {
	t.Parallel()

	q := NewRegisteredNodeAddressBookQuery()
	require.NotNil(t, q)

	result := q.SetRegisteredNodeId(42)
	assert.Same(t, q, result, "SetRegisteredNodeId should return the same query for chaining")
	assert.Equal(t, uint64(42), q.GetRegisteredNodeId())
}

func TestUnitRegisteredNodeAddressBookQueryGetRegisteredNodeIdDefault(t *testing.T) {
	t.Parallel()

	q := NewRegisteredNodeAddressBookQuery()
	assert.Equal(t, uint64(0), q.GetRegisteredNodeId(), "unset ID should default to 0")
}

func TestUnitRegisteredNodeAddressBookQueryExecuteNilClient(t *testing.T) {
	t.Parallel()

	book, err := NewRegisteredNodeAddressBookQuery().Execute(nil)
	assert.Equal(t, errNoClientProvided, err)
	assert.Equal(t, RegisteredNodeAddressBook{}, book)
}

func TestUnitBlockNodeApiFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want BlockNodeApi
	}{
		{"STATUS", BlockNodeApiStatus},
		{"PUBLISH", BlockNodeApiPublish},
		{"SUBSCRIBE_STREAM", BlockNodeApiSubscribeStream},
		{"STATE_PROOF", BlockNodeApiStateProof},
		{"OTHER", BlockNodeApiOther},
		{"unrecognised", BlockNodeApiOther},
		{"", BlockNodeApiOther},
		{"status", BlockNodeApiStatus}, // lower-case should still match
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, blockNodeApiFromString(tt.in), "input %q", tt.in)
	}
}

func TestUnitAdminKeyFromJSONEd25519(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(mockPrivateKey)
	require.NoError(t, err)
	pub := priv.PublicKey().String()

	key, err := adminKeyFromJSON(adminKeyJSON{Type: "ED25519", Key: pub})
	require.NoError(t, err)
	require.NotNil(t, key)
	assert.Equal(t, pub, key.(PublicKey).String())
}

func TestUnitAdminKeyFromJSONEcdsa(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	pub := priv.PublicKey().StringRaw()

	key, err := adminKeyFromJSON(adminKeyJSON{Type: "ECDSA_SECP256K1", Key: pub})
	require.NoError(t, err)
	require.NotNil(t, key)
}

func TestUnitAdminKeyFromJSONInvalid(t *testing.T) {
	t.Parallel()

	_, err := adminKeyFromJSON(adminKeyJSON{Type: "ED25519", Key: "not-a-key"})
	assert.Error(t, err)
}

func TestUnitServiceEndpointFromJSONBlockNode(t *testing.T) {
	t.Parallel()

	ip := "192.168.1.1"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:        "BLOCK_NODE",
		IPAddress:   &ip,
		Port:        50211,
		RequiresTls: true,
		BlockNode:   &blockNodeJSON{EndpointApis: []string{"PUBLISH", "STATUS"}},
	})
	require.NoError(t, err)

	block, ok := ep.(*BlockNodeServiceEndpoint)
	require.True(t, ok, "expected *BlockNodeServiceEndpoint")
	assert.Equal(t, net.IPv4(192, 168, 1, 1).To4(), net.IP(block.GetIPAddress()))
	assert.Equal(t, uint32(50211), block.GetPort())
	assert.True(t, block.GetRequiresTls())
	assert.Equal(t, []BlockNodeApi{BlockNodeApiPublish, BlockNodeApiStatus}, block.GetEndpointApis())
}

func TestUnitServiceEndpointFromJSONMirrorNode(t *testing.T) {
	t.Parallel()

	domain := "mirror.example.com"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:       "MIRROR_NODE",
		DomainName: &domain,
		Port:       443,
	})
	require.NoError(t, err)

	mirror, ok := ep.(*MirrorNodeServiceEndpoint)
	require.True(t, ok, "expected *MirrorNodeServiceEndpoint")
	assert.Equal(t, "mirror.example.com", mirror.GetDomainName())
	assert.Equal(t, uint32(443), mirror.GetPort())
}

func TestUnitServiceEndpointFromJSONRpcRelay(t *testing.T) {
	t.Parallel()

	domain := "rpc.example.com"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:       "RPC_RELAY",
		DomainName: &domain,
		Port:       8545,
	})
	require.NoError(t, err)

	rpc, ok := ep.(*RpcRelayServiceEndpoint)
	require.True(t, ok, "expected *RpcRelayServiceEndpoint")
	assert.Equal(t, "rpc.example.com", rpc.GetDomainName())
	assert.Equal(t, uint32(8545), rpc.GetPort())
}

func TestUnitServiceEndpointFromJSONUnknownType(t *testing.T) {
	t.Parallel()

	domain := "x.example.com"
	_, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:       "SOMETHING_ELSE",
		DomainName: &domain,
		Port:       80,
	})
	assert.ErrorContains(t, err, "unknown endpoint type")
}

func TestUnitServiceEndpointFromJSONIPv6(t *testing.T) {
	t.Parallel()

	ip := "2001:db8::1"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:      "MIRROR_NODE",
		IPAddress: &ip,
		Port:      443,
	})
	require.NoError(t, err)

	mirror, ok := ep.(*MirrorNodeServiceEndpoint)
	require.True(t, ok)
	assert.Equal(t, net.ParseIP(ip).To16(), net.IP(mirror.GetIPAddress()))
}

func TestUnitServiceEndpointFromJSONInvalidIP(t *testing.T) {
	t.Parallel()

	ip := "not-an-ip"
	_, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:      "MIRROR_NODE",
		IPAddress: &ip,
		Port:      443,
	})
	assert.ErrorContains(t, err, "invalid IP address")
}

func TestUnitRegisteredNodeFromJSON(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(mockPrivateKey)
	require.NoError(t, err)
	pub := priv.PublicKey().String()

	ip := "10.0.0.1"
	domain := "rpc.example.com"
	raw := registeredNodeJSON{
		AdminKey:         &adminKeyJSON{Type: "ED25519", Key: pub},
		CreatedTimestamp: "1700000000.000000000",
		Description:      "node description",
		RegisteredNodeID: 7,
		ServiceEndpoints: []serviceEndpointJSON{
			{
				Type:      "BLOCK_NODE",
				IPAddress: &ip,
				Port:      50211,
				BlockNode: &blockNodeJSON{EndpointApis: []string{"PUBLISH"}},
			},
			{
				Type:       "RPC_RELAY",
				DomainName: &domain,
				Port:       8545,
			},
		},
	}

	node, err := registeredNodeFromJSON(raw)
	require.NoError(t, err)

	assert.Equal(t, "node description", node.Description)
	assert.Equal(t, uint64(7), node.RegisteredNodeID)
	assert.Equal(t, "1700000000.000000000", node.CreatedTimestamp)
	require.NotNil(t, node.AdminKey)
	assert.Equal(t, pub, node.AdminKey.(PublicKey).String())
	require.Len(t, node.ServiceEndpoints, 2)

	_, isBlock := node.ServiceEndpoints[0].(*BlockNodeServiceEndpoint)
	assert.True(t, isBlock)
	_, isRpc := node.ServiceEndpoints[1].(*RpcRelayServiceEndpoint)
	assert.True(t, isRpc)
}

func TestUnitRegisteredNodeFromJSONPropagatesEndpointError(t *testing.T) {
	t.Parallel()

	ip := "bad-ip"
	raw := registeredNodeJSON{
		RegisteredNodeID: 1,
		ServiceEndpoints: []serviceEndpointJSON{
			{Type: "MIRROR_NODE", IPAddress: &ip, Port: 443},
		},
	}

	_, err := registeredNodeFromJSON(raw)
	assert.ErrorContains(t, err, "failed to parse service endpoint")
}
