//go:build all || unit
// +build all unit

package hiero

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndpoint_SetAndGetAddress(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{}
	testAddress := []byte{192, 168, 1, 1}

	result := endpoint.SetAddress(testAddress)

	assert.Equal(t, endpoint, result, "SetAddress should return the same endpoint for chaining")
	assert.Equal(t, testAddress, endpoint.GetAddress(), "GetAddress should return the set address")
}

func TestEndpoint_SetAndGetPort(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{}
	testPort := int32(8080)

	result := endpoint.SetPort(testPort)

	assert.Equal(t, endpoint, result, "SetPort should return the same endpoint for chaining")
	assert.Equal(t, testPort, endpoint.GetPort(), "GetPort should return the set port")
}

func TestEndpoint_SetAndGetDomainName(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{}
	testDomain := "example.com"

	result := endpoint.SetDomainName(testDomain)

	assert.Equal(t, endpoint, result, "SetDomainName should return the same endpoint for chaining")
	assert.Equal(t, testDomain, endpoint.GetDomainName(), "GetDomainName should return the set domain name")
}

func TestEndpoint_Validate_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupFunc   func() *Endpoint
		description string
	}{
		{
			name: "ValidWithAddress",
			setupFunc: func() *Endpoint {
				return &Endpoint{
					address: []byte{192, 168, 1, 1},
					port:    8080,
				}
			},
			description: "Should be valid with only address",
		},
		{
			name: "ValidWithDomainName",
			setupFunc: func() *Endpoint {
				return &Endpoint{
					domainName: "example.com",
					port:       8080,
				}
			},
			description: "Should be valid with only domain name",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			endpoint := tc.setupFunc()
			err := endpoint.Validate()

			assert.NoError(t, err, tc.description)
		})
	}
}

func TestEndpoint_Validate_Failure(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupFunc     func() *Endpoint
		expectedError error
	}{
		{
			name: "NeitherAddressNorDomain",
			setupFunc: func() *Endpoint {
				return &Endpoint{
					port: 8080,
				}
			},
			expectedError: errEndpointMustHaveAddressOrDomainName,
		},
		{
			name: "BothAddressAndDomain",
			setupFunc: func() *Endpoint {
				return &Endpoint{
					address:    []byte{192, 168, 1, 1},
					domainName: "example.com",
					port:       8080,
				}
			},
			expectedError: errEndpointCannotHaveBothAddressAndDomainName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			endpoint := tc.setupFunc()
			err := endpoint.Validate()

			require.Error(t, err)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestEndpointFromProtobuf(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		serviceEndpoint *services.ServiceEndpoint
		expectedResult  Endpoint
		description     string
	}{
		{
			name: "NormalPort",
			serviceEndpoint: &services.ServiceEndpoint{
				IpAddressV4: []byte{192, 168, 1, 1},
				Port:        8080,
				DomainName:  "example.com",
			},
			expectedResult: Endpoint{
				address:    []byte{192, 168, 1, 1},
				port:       8080,
				domainName: "example.com",
			},
			description: "Should convert normally when port is not 0 or 50111",
		},
		{
			name: "ZeroPortConvertsTo50211",
			serviceEndpoint: &services.ServiceEndpoint{
				IpAddressV4: []byte{10, 0, 0, 1},
				Port:        0,
				DomainName:  "test.com",
			},
			expectedResult: Endpoint{
				address:    []byte{10, 0, 0, 1},
				port:       50211,
				domainName: "test.com",
			},
			description: "Should convert port 0 to 50211",
		},
		{
			name: "Port50111ConvertsTo50211",
			serviceEndpoint: &services.ServiceEndpoint{
				IpAddressV4: []byte{172, 16, 0, 1},
				Port:        50111,
				DomainName:  "another.com",
			},
			expectedResult: Endpoint{
				address:    []byte{172, 16, 0, 1},
				port:       50211,
				domainName: "another.com",
			},
			description: "Should convert port 50111 to 50211",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := EndpointFromProtobuf(tc.serviceEndpoint)

			assert.Equal(t, tc.expectedResult, result, tc.description)
		})
	}
}

func TestEndpoint_ToProtobuf(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{
		address:    []byte{192, 168, 1, 1},
		port:       8080,
		domainName: "example.com",
	}

	result := endpoint._ToProtobuf()

	expected := &services.ServiceEndpoint{
		IpAddressV4: []byte{192, 168, 1, 1},
		Port:        8080,
		DomainName:  "example.com",
	}

	assert.Equal(t, expected, result)
}

func TestEndpoint_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		endpoint       *Endpoint
		expectedResult string
		description    string
	}{
		{
			name: "WithDomainName",
			endpoint: &Endpoint{
				domainName: "example.com",
				port:       8080,
			},
			expectedResult: "example.com:8080",
			description:    "Should use domain name format when domain name is set",
		},
		{
			name: "WithIPAddress",
			endpoint: &Endpoint{
				address: []byte{192, 168, 1, 1},
				port:    8080,
			},
			expectedResult: "192.168.1.1:8080",
			description:    "Should use IP address format when only address is set",
		},
		{
			name: "WithBothPrefersDomainName",
			endpoint: &Endpoint{
				address:    []byte{10, 0, 0, 1},
				domainName: "test.com",
				port:       9090,
			},
			expectedResult: "test.com:9090",
			description:    "Should prefer domain name when both are set",
		},
		{
			name: "DifferentIPAddress",
			endpoint: &Endpoint{
				address: []byte{172, 16, 254, 100},
				port:    443,
			},
			expectedResult: "172.16.254.100:443",
			description:    "Should correctly format different IP addresses",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := tc.endpoint.String()

			assert.Equal(t, tc.expectedResult, result, tc.description)
		})
	}
}

func TestEndpoint_MethodChaining(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{}
	testAddress := []byte{192, 168, 1, 1}
	testPort := int32(8080)
	testDomain := "example.com"

	result := endpoint.
		SetAddress(testAddress).
		SetPort(testPort).
		SetDomainName(testDomain)

	assert.Equal(t, endpoint, result)
	assert.Equal(t, testAddress, endpoint.GetAddress())
	assert.Equal(t, testPort, endpoint.GetPort())
	assert.Equal(t, testDomain, endpoint.GetDomainName())
}

func TestEndpoint_EmptyValues(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{}

	assert.Nil(t, endpoint.GetAddress())
	assert.Equal(t, int32(0), endpoint.GetPort())
	assert.Equal(t, "", endpoint.GetDomainName())
}

func TestEndpoint_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("EmptyAddress", func(t *testing.T) {
		t.Parallel()

		endpoint := &Endpoint{}
		endpoint.SetAddress([]byte{})

		assert.Equal(t, []byte{}, endpoint.GetAddress())
	})

	t.Run("NilAddressAfterSet", func(t *testing.T) {
		t.Parallel()

		endpoint := &Endpoint{}
		endpoint.SetAddress(nil)

		assert.Nil(t, endpoint.GetAddress())
	})

	t.Run("EmptyDomainName", func(t *testing.T) {
		t.Parallel()

		endpoint := &Endpoint{}
		endpoint.SetDomainName("")

		assert.Equal(t, "", endpoint.GetDomainName())
	})

	t.Run("NegativePort", func(t *testing.T) {
		t.Parallel()

		endpoint := &Endpoint{}
		endpoint.SetPort(-1)

		assert.Equal(t, int32(-1), endpoint.GetPort())
	})
}

func TestEndpoint_RoundTripProtobuf(t *testing.T) {
	t.Parallel()

	originalEndpoint := &Endpoint{
		address:    []byte{203, 0, 113, 45},
		port:       9000,
		domainName: "roundtrip.test",
	}

	// Convert to protobuf and back
	proto := originalEndpoint._ToProtobuf()
	reconstructed := EndpointFromProtobuf(proto)

	assert.Equal(t, originalEndpoint.address, reconstructed.address)
	assert.Equal(t, originalEndpoint.port, reconstructed.port)
	assert.Equal(t, originalEndpoint.domainName, reconstructed.domainName)
}
