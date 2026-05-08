//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"net"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----- setters / getters -----

func TestUnitRegisteredNodeAddressBookQuerySettersChain(t *testing.T) {
	t.Parallel()

	q := NewRegisteredNodeAddressBookQuery()
	require.NotNil(t, q)

	assert.Same(t, q, q.SetRegisteredNodeId(42))
	assert.Equal(t, uint64(42), q.GetRegisteredNodeId())

	assert.Same(t, q, q.SetLimit(50))
	assert.Equal(t, int32(50), q.GetLimit())

	assert.Same(t, q, q.SetMaxAttempts(3))
	assert.Equal(t, uint64(3), q.GetMaxAttempts())
}

func TestUnitRegisteredNodeAddressBookQueryDefaults(t *testing.T) {
	t.Parallel()

	q := NewRegisteredNodeAddressBookQuery()
	assert.Equal(t, uint64(0), q.GetRegisteredNodeId())
	assert.Equal(t, int32(0), q.GetLimit())
	assert.Equal(t, uint64(0), q.GetMaxAttempts())
}

func TestUnitRegisteredNodeAddressBookQueryExecuteNilClient(t *testing.T) {
	t.Parallel()

	book, err := NewRegisteredNodeAddressBookQuery().Execute(nil)
	assert.Equal(t, errNoClientProvided, err)
	assert.Equal(t, RegisteredNodeAddressBook{}, book)
}

// ----- buildURL -----

func TestUnitRegisteredNodeAddressBookQueryBuildURL(t *testing.T) {
	t.Parallel()

	id := uint64(7)
	tests := []struct {
		name string
		q    *RegisteredNodeAddressBookQuery
		want string
	}{
		{
			name: "no params",
			q:    &RegisteredNodeAddressBookQuery{},
			want: "https://example/api/v1/network/registered-nodes",
		},
		{
			name: "id only",
			q:    &RegisteredNodeAddressBookQuery{registeredNodeId: &id},
			want: "https://example/api/v1/network/registered-nodes?registerednode.id=7",
		},
		{
			name: "limit only",
			q:    &RegisteredNodeAddressBookQuery{limit: 25},
			want: "https://example/api/v1/network/registered-nodes?limit=25",
		},
		{
			name: "id and limit",
			q:    &RegisteredNodeAddressBookQuery{registeredNodeId: &id, limit: 25},
			want: "https://example/api/v1/network/registered-nodes?limit=25&registerednode.id=7",
		},
		{
			name: "limit zero is omitted",
			q:    &RegisteredNodeAddressBookQuery{registeredNodeId: &id, limit: 0},
			want: "https://example/api/v1/network/registered-nodes?registerednode.id=7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.q.buildURL("https://example/api/v1"))
		})
	}
}

// ----- resolveNextURL -----

func TestUnitResolveNextURL(t *testing.T) {
	t.Parallel()

	base, err := url.Parse("https://mirror.example/api/v1")
	require.NoError(t, err)

	tests := []struct {
		name string
		next string
		want string
	}{
		{
			name: "absolute path replaces base path",
			next: "/api/v1/network/registered-nodes?limit=25&registerednode.id=gt:5",
			want: "https://mirror.example/api/v1/network/registered-nodes?limit=25&registerednode.id=gt:5",
		},
		{
			name: "absolute URL passes through",
			next: "https://other.example/api/v1/network/registered-nodes?limit=25",
			want: "https://other.example/api/v1/network/registered-nodes?limit=25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := resolveNextURL(base, tt.next)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ----- blockNodeApiFromString -----

func TestUnitBlockNodeApiFromStringRecognized(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want BlockNodeApi
	}{
		{"OTHER", BlockNodeApiOther},
		{"STATUS", BlockNodeApiStatus},
		{"PUBLISH", BlockNodeApiPublish},
		{"SUBSCRIBE_STREAM", BlockNodeApiSubscribeStream},
		{"STATE_PROOF", BlockNodeApiStateProof},
		{"status", BlockNodeApiStatus},
	}

	for _, tt := range tests {
		got, err := blockNodeApiFromString(tt.in)
		require.NoError(t, err, "input %q", tt.in)
		assert.Equal(t, tt.want, got, "input %q", tt.in)
	}
}

func TestUnitBlockNodeApiFromStringRejectsUnknown(t *testing.T) {
	t.Parallel()

	cases := []string{
		// "UNKNOWN" only ever appears when BlockNodeApi.String() was called on
		// a value the SDK could not name, so seeing it on the way back in
		// signals a corrupt round-trip.
		"UNKNOWN",
		"UNRECOGNIZED",
		"FUTURE_KIND",
		"",
	}

	for _, in := range cases {
		_, err := blockNodeApiFromString(in)
		assert.Error(t, err, "input %q should error", in)
	}
}

// ----- adminKeyFromJSON -----

func TestUnitAdminKeyFromJSONEd25519(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(mockPrivateKey)
	require.NoError(t, err)
	pub := priv.PublicKey().String()

	key, err := adminKeyFromJSON(adminKeyJSON{Type: "ED25519", Key: pub})
	require.NoError(t, err)
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

func TestUnitAdminKeyFromJSONDefaultBranch(t *testing.T) {
	t.Parallel()

	priv, err := PrivateKeyFromString(mockPrivateKey)
	require.NoError(t, err)
	pub := priv.PublicKey().String()

	key, err := adminKeyFromJSON(adminKeyJSON{Type: "MYSTERY", Key: pub})
	require.NoError(t, err)
	require.NotNil(t, key)
}

func TestUnitAdminKeyFromJSONInvalidKey(t *testing.T) {
	t.Parallel()

	_, err := adminKeyFromJSON(adminKeyJSON{Type: "ED25519", Key: "not-a-key"})
	assert.Error(t, err)
}

// ----- serviceEndpointFromJSON -----

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

	bn, ok := ep.(*BlockNodeServiceEndpoint)
	require.True(t, ok, "expected *BlockNodeServiceEndpoint")
	assert.Equal(t, []byte(net.ParseIP(ip).To4()), bn.GetIPAddress())
	assert.Equal(t, uint32(50211), bn.GetPort())
	assert.True(t, bn.GetRequiresTls())
	assert.Equal(t, []BlockNodeApi{BlockNodeApiPublish, BlockNodeApiStatus}, bn.GetEndpointApis())
}

func TestUnitServiceEndpointFromJSONMirrorNodeIPv4(t *testing.T) {
	t.Parallel()

	ip := "10.0.0.1"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:      "MIRROR_NODE",
		IPAddress: &ip,
		Port:      443,
	})
	require.NoError(t, err)

	mn, ok := ep.(*MirrorNodeServiceEndpoint)
	require.True(t, ok)
	assert.Equal(t, []byte(net.ParseIP(ip).To4()), mn.GetIPAddress())
	assert.Equal(t, uint32(443), mn.GetPort())
	assert.Empty(t, mn.GetDomainName())
}

func TestUnitServiceEndpointFromJSONMirrorNodeDomain(t *testing.T) {
	t.Parallel()

	domain := "mirror.example.com"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:       "MIRROR_NODE",
		DomainName: &domain,
		Port:       443,
	})
	require.NoError(t, err)

	mn, ok := ep.(*MirrorNodeServiceEndpoint)
	require.True(t, ok)
	assert.Empty(t, mn.GetIPAddress())
	assert.Equal(t, domain, mn.GetDomainName())
}

func TestUnitServiceEndpointFromJSONRpcRelay(t *testing.T) {
	t.Parallel()

	domain := "relay.example.com"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:       "RPC_RELAY",
		DomainName: &domain,
		Port:       7546,
	})
	require.NoError(t, err)

	rr, ok := ep.(*RpcRelayServiceEndpoint)
	require.True(t, ok)
	assert.Equal(t, domain, rr.GetDomainName())
	assert.Equal(t, uint32(7546), rr.GetPort())
}

func TestUnitServiceEndpointFromJSONGeneralService(t *testing.T) {
	t.Parallel()

	domain := "general.example.com"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:           "GENERAL_SERVICE",
		DomainName:     &domain,
		Port:           9000,
		GeneralService: &generalServiceJSON{Description: "custom-service"},
	})
	require.NoError(t, err)

	gs, ok := ep.(*GeneralServiceEndpoint)
	require.True(t, ok)
	assert.Equal(t, domain, gs.GetDomainName())
	assert.Equal(t, "custom-service", gs.GetDescription())
}

func TestUnitServiceEndpointFromJSONGeneralServiceMissingNested(t *testing.T) {
	t.Parallel()

	domain := "general.example.com"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:       "GENERAL_SERVICE",
		DomainName: &domain,
		Port:       9000,
	})
	require.NoError(t, err)

	gs, ok := ep.(*GeneralServiceEndpoint)
	require.True(t, ok)
	assert.Empty(t, gs.GetDescription())
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

	got := ep.GetIPAddress()
	require.Len(t, got, 16, "IPv6 should be 16 bytes")
	assert.Equal(t, []byte(net.ParseIP(ip).To16()), got)
}

func TestUnitServiceEndpointFromJSONInvalidIP(t *testing.T) {
	t.Parallel()

	bad := "not-an-ip"
	_, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:      "MIRROR_NODE",
		IPAddress: &bad,
		Port:      443,
	})
	assert.Error(t, err)
}

func TestUnitServiceEndpointFromJSONUnknownType(t *testing.T) {
	t.Parallel()

	_, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type: "MYSTERY",
		Port: 443,
	})
	assert.Error(t, err)
}

func TestUnitServiceEndpointFromJSONBlockNodeMissingNested(t *testing.T) {
	t.Parallel()

	ip := "10.0.0.1"
	ep, err := serviceEndpointFromJSON(serviceEndpointJSON{
		Type:      "BLOCK_NODE",
		IPAddress: &ip,
		Port:      50211,
	})
	require.NoError(t, err)

	bn, ok := ep.(*BlockNodeServiceEndpoint)
	require.True(t, ok)
	assert.Empty(t, bn.GetEndpointApis())
}

// ----- registeredNodeFromJSON -----

func TestUnitRegisteredNodeFromJSONBadAdminKey(t *testing.T) {
	t.Parallel()

	_, err := registeredNodeFromJSON(registeredNodeJSON{
		AdminKey:         &adminKeyJSON{Type: "ED25519", Key: "not-a-key"},
		RegisteredNodeID: 1,
	})
	assert.Error(t, err)
}

// ----- resolveAttempts / Execute precondition guards -----

func TestUnitResolveAttempts(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)

	q := NewRegisteredNodeAddressBookQuery()
	assert.Equal(t, uint64(1), q.resolveAttempts(client), "fallback to 1 when neither is set")

	client.SetMaxAttempts(7)
	assert.Equal(t, uint64(7), q.resolveAttempts(client), "uses client default when query unset")

	q.SetMaxAttempts(3)
	assert.Equal(t, uint64(3), q.resolveAttempts(client), "query setting wins over client")
}

func TestUnitRegisteredNodeAddressBookQueryExecuteNoMirror(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork(nil)

	_, err = NewRegisteredNodeAddressBookQuery().Execute(client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mirror node is not set")
}
