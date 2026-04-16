//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
)

func TestUnitBlockNodeApiValues(t *testing.T) {
	t.Parallel()

	assert.Equal(t, BlockNodeApi(0), BlockNodeApiOther)
	assert.Equal(t, BlockNodeApi(1), BlockNodeApiStatus)
	assert.Equal(t, BlockNodeApi(2), BlockNodeApiPublish)
	assert.Equal(t, BlockNodeApi(3), BlockNodeApiSubscribeStream)
	assert.Equal(t, BlockNodeApi(4), BlockNodeApiStateProof)
}

func TestUnitBlockNodeApiString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "OTHER", BlockNodeApiOther.String())
	assert.Equal(t, "STATUS", BlockNodeApiStatus.String())
	assert.Equal(t, "PUBLISH", BlockNodeApiPublish.String())
	assert.Equal(t, "SUBSCRIBE_STREAM", BlockNodeApiSubscribeStream.String())
	assert.Equal(t, "STATE_PROOF", BlockNodeApiStateProof.String())
}

func TestUnitBlockNodeServiceEndpointRoundTripIP(t *testing.T) {
	t.Parallel()

	ip := []byte{192, 168, 1, 1}
	endpoint := &BlockNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress:   ip,
			port:        8080,
			requiresTls: true,
		},
		endpointApis: []BlockNodeApi{BlockNodeApiStatus},
	}

	pb := endpoint._ToProtobuf()
	restored := _RegisteredServiceEndpointFromProtobuf(pb)

	block, ok := restored.(*BlockNodeServiceEndpoint)
	assert.True(t, ok, "expected *BlockNodeServiceEndpoint")
	assert.Equal(t, ip, block.GetIPAddress())
	assert.Equal(t, uint32(8080), block.GetPort())
	assert.True(t, block.GetRequiresTls())
	assert.Equal(t, []BlockNodeApi{BlockNodeApiStatus}, block.GetEndpointApis())
	assert.Equal(t, "", block.GetDomainName())
}

func TestUnitMirrorNodeServiceEndpointRoundTripDomain(t *testing.T) {
	t.Parallel()

	endpoint := &MirrorNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			domainName:  "mirror.example.com",
			port:        443,
			requiresTls: true,
		},
	}

	pb := endpoint._ToProtobuf()
	restored := _RegisteredServiceEndpointFromProtobuf(pb)

	mirror, ok := restored.(*MirrorNodeServiceEndpoint)
	assert.True(t, ok, "expected *MirrorNodeServiceEndpoint")
	assert.Equal(t, "mirror.example.com", mirror.GetDomainName())
	assert.Equal(t, uint32(443), mirror.GetPort())
	assert.True(t, mirror.GetRequiresTls())
	assert.Nil(t, mirror.GetIPAddress())
}

func TestUnitBlockNodeServiceEndpointValidateBothAddresses(t *testing.T) {
	t.Parallel()

	endpoint := &BlockNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress:  []byte{10, 0, 0, 1},
			domainName: "example.com",
		},
	}

	err := endpoint.Validate()
	assert.Equal(t, errEndpointCannotHaveBothAddressAndDomainName, err)
}

func TestUnitMirrorNodeServiceEndpointValidateNoAddress(t *testing.T) {
	t.Parallel()

	endpoint := &MirrorNodeServiceEndpoint{}

	err := endpoint.Validate()
	assert.Equal(t, errEndpointMustHaveAddressOrDomainName, err)
}

func TestUnitRpcRelayServiceEndpointProtobufOneof(t *testing.T) {
	t.Parallel()

	endpoint := &RpcRelayServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			domainName: "rpc.example.com",
			port:       8545,
		},
	}

	pb := endpoint._ToProtobuf()

	_, ok := pb.EndpointType.(*services.RegisteredServiceEndpoint_RpcRelay)
	assert.True(t, ok, "expected EndpointType to be *RegisteredServiceEndpoint_RpcRelay")
}

func TestUnitBlockNodeServiceEndpointValidateValid(t *testing.T) {
	t.Parallel()

	endpoint := &BlockNodeServiceEndpoint{
		registeredEndpointBase: registeredEndpointBase{
			ipAddress: []byte{10, 0, 0, 1},
			port:      8080,
		},
	}

	err := endpoint.Validate()
	assert.NoError(t, err)
}

func TestUnitBlockNodeServiceEndpointApiRoundTrip(t *testing.T) {
	t.Parallel()

	apis := []BlockNodeApi{
		BlockNodeApiOther,
		BlockNodeApiStatus,
		BlockNodeApiPublish,
		BlockNodeApiSubscribeStream,
		BlockNodeApiStateProof,
	}

	for _, api := range apis {
		endpoint := &BlockNodeServiceEndpoint{
			registeredEndpointBase: registeredEndpointBase{
				ipAddress: []byte{10, 0, 0, 1},
				port:      8080,
			},
			endpointApis: []BlockNodeApi{api},
		}

		pb := endpoint._ToProtobuf()
		restored := _RegisteredServiceEndpointFromProtobuf(pb)

		block, ok := restored.(*BlockNodeServiceEndpoint)
		assert.True(t, ok)
		assert.Equal(t, []BlockNodeApi{api}, block.GetEndpointApis(), "BlockNodeApi %s did not round-trip", api.String())
	}
}

func TestUnitBlockNodeServiceEndpointSetters(t *testing.T) {
	t.Parallel()

	endpoint := &BlockNodeServiceEndpoint{}
	endpoint.SetIPAddress([]byte{192, 168, 0, 1}).
		SetPort(443).
		SetRequiresTls(true).
		AddEndpointApi(BlockNodeApiSubscribeStream)

	assert.Equal(t, []byte{192, 168, 0, 1}, endpoint.GetIPAddress())
	assert.Equal(t, uint32(443), endpoint.GetPort())
	assert.True(t, endpoint.GetRequiresTls())
	assert.Equal(t, []BlockNodeApi{BlockNodeApiSubscribeStream}, endpoint.GetEndpointApis())
}

func TestUnitMirrorNodeServiceEndpointSetters(t *testing.T) {
	t.Parallel()

	endpoint := &MirrorNodeServiceEndpoint{}
	endpoint.SetDomainName("mirror.example.com").
		SetPort(443).
		SetRequiresTls(true)

	assert.Equal(t, "mirror.example.com", endpoint.GetDomainName())
	assert.Equal(t, uint32(443), endpoint.GetPort())
	assert.True(t, endpoint.GetRequiresTls())
}

func TestUnitRpcRelayServiceEndpointSetters(t *testing.T) {
	t.Parallel()

	endpoint := &RpcRelayServiceEndpoint{}
	endpoint.SetDomainName("rpc.example.com").
		SetPort(8545)

	assert.Equal(t, "rpc.example.com", endpoint.GetDomainName())
	assert.Equal(t, uint32(8545), endpoint.GetPort())
	assert.False(t, endpoint.GetRequiresTls())
}

func TestUnitRegisteredServiceEndpointFromProtobufDispatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint RegisteredServiceEndpoint
	}{
		{"BlockNode", &BlockNodeServiceEndpoint{registeredEndpointBase: registeredEndpointBase{ipAddress: []byte{1, 2, 3, 4}, port: 1}}},
		{"MirrorNode", &MirrorNodeServiceEndpoint{registeredEndpointBase: registeredEndpointBase{ipAddress: []byte{1, 2, 3, 4}, port: 2}}},
		{"RpcRelay", &RpcRelayServiceEndpoint{registeredEndpointBase: registeredEndpointBase{ipAddress: []byte{1, 2, 3, 4}, port: 3}}},
	}

	for _, tt := range tests {
		pb := tt.endpoint._ToProtobuf()
		restored := _RegisteredServiceEndpointFromProtobuf(pb)

		switch tt.name {
		case "BlockNode":
			_, ok := restored.(*BlockNodeServiceEndpoint)
			assert.True(t, ok, "expected *BlockNodeServiceEndpoint")
		case "MirrorNode":
			_, ok := restored.(*MirrorNodeServiceEndpoint)
			assert.True(t, ok, "expected *MirrorNodeServiceEndpoint")
		case "RpcRelay":
			_, ok := restored.(*RpcRelayServiceEndpoint)
			assert.True(t, ok, "expected *RpcRelayServiceEndpoint")
		}
	}
}
