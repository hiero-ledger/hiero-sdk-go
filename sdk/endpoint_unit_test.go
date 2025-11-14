//go:build all || unit

package hiero

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndpointSetAndGetAddress(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{}
	testAddress := []byte{192, 168, 1, 1}

	result := endpoint.SetAddress(testAddress)

	assert.Equal(t, endpoint, result, "SetAddress should return the same endpoint for chaining")
	assert.Equal(t, testAddress, endpoint.GetAddress(), "GetAddress should return the set address")
}

func TestEndpointSetAndGetPort(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{}
	testPort := int32(8080)

	result := endpoint.SetPort(testPort)

	assert.Equal(t, endpoint, result, "SetPort should return the same endpoint for chaining")
	assert.Equal(t, testPort, endpoint.GetPort(), "GetPort should return the set port")
}

func TestEndpointSetAndGetDomainName(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{}
	testDomain := "example.com"

	result := endpoint.SetDomainName(testDomain)

	assert.Equal(t, endpoint, result, "SetDomainName should return the same endpoint for chaining")
	assert.Equal(t, testDomain, endpoint.GetDomainName(), "GetDomainName should return the set domain name")
}

func TestEndpointValidate(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {

			endpoint := tc.setupFunc()
			err := endpoint.Validate()

			assert.NoError(t, err, tc.description)
		})
	}
}

func TestEndpointValidateFailure(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {

			result := EndpointFromProtobuf(tc.serviceEndpoint)

			assert.Equal(t, tc.expectedResult, result, tc.description)
		})
	}
}

func TestEndpointToProtobuf(t *testing.T) {
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

func TestEndpointString(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {

			result := tc.endpoint.String()

			assert.Equal(t, tc.expectedResult, result, tc.description)
		})
	}
}

func TestEndpointEmptyValues(t *testing.T) {
	t.Parallel()

	endpoint := &Endpoint{}

	assert.Nil(t, endpoint.GetAddress())
	assert.Equal(t, int32(0), endpoint.GetPort())
	assert.Equal(t, "", endpoint.GetDomainName())
}

func TestEndpointRoundTripProtobuf(t *testing.T) {
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
