//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMirrorNodeContractQuerySetAndGetContractID(t *testing.T) {
	contractID := ContractID{Shard: 0, Realm: 0, Contract: 1234}

	query1 := NewMirrorNodeContractEstimateGasQuery()
	query1.SetContractID(contractID)
	assert.Equal(t, contractID, query1.GetContractID())

	query2 := NewMirrorNodeContractCallQuery()
	query2.SetContractID(contractID)
	assert.Equal(t, contractID, query2.GetContractID())
}

func TestMirrorNodeContractQuerySetAndGetSenderEvmAddress(t *testing.T) {
	evmAddress := "0x1234567890abcdef1234567890abcdef12345678"

	query1 := NewMirrorNodeContractEstimateGasQuery()
	query1.SetSenderEvmAddress(evmAddress)
	assert.Equal(t, evmAddress, query1.GetSenderEvmAddress())
	assert.Equal(t, AccountID{}, query1.GetSender())

	query2 := NewMirrorNodeContractCallQuery()
	query2.SetSenderEvmAddress(evmAddress)
	assert.Equal(t, evmAddress, query2.GetSenderEvmAddress())
	assert.Equal(t, AccountID{}, query2.GetSender())
}

func TestMirrorNodeContractQuerySetAndGetContractEvmAddress(t *testing.T) {
	evmAddress := "0x1234567890abcdef1234567890abcdef12345678"

	query1 := NewMirrorNodeContractEstimateGasQuery()
	query1.SetContractEvmAddress(evmAddress)
	assert.Equal(t, evmAddress, query1.GetContractEvmAddress())
	assert.Equal(t, ContractID{}, query1.GetContractID())

	query2 := NewMirrorNodeContractCallQuery()
	query2.SetContractEvmAddress(evmAddress)
	assert.Equal(t, evmAddress, query2.GetContractEvmAddress())
	assert.Equal(t, ContractID{}, query2.GetContractID())
}

func TestMirrorNodeContractQuerySetAndGetCallData(t *testing.T) {
	params := []byte("test")

	query1 := NewMirrorNodeContractEstimateGasQuery()
	query1.SetFunctionParameters(params)
	assert.Equal(t, params, query1.GetCallData())

	query2 := NewMirrorNodeContractCallQuery()
	query2.SetFunctionParameters(params)
	assert.Equal(t, params, query2.GetCallData())
}

func TestMirrorNodeContractQuerySetFunctionWithoutParameters(t *testing.T) {
	query1 := &mirrorNodeContractQuery{}
	query1.setFunction("myFunction", nil)
	assert.NotNil(t, query1.GetCallData())

	query2 := NewMirrorNodeContractEstimateGasQuery()
	query2.SetFunction("myFunction", nil)
	assert.NotNil(t, query1.GetCallData())

	query3 := NewMirrorNodeContractCallQuery()
	query3.SetFunction("myFunction", nil)
	assert.NotNil(t, query2.GetCallData())
}

func TestMirrorNodeContractQuerySetAndGetBlockNumber(t *testing.T) {
	blockNumber := int64(123456)

	query1 := NewMirrorNodeContractCallQuery()
	query1.SetBlockNumber(blockNumber)
	assert.Equal(t, blockNumber, query1.GetBlockNumber())
}

func TestMirrorNodeContractQuerySetAndGetValue(t *testing.T) {
	value := int64(1000)

	query1 := NewMirrorNodeContractEstimateGasQuery()
	query1.SetValue(value)
	assert.Equal(t, value, query1.GetValue())

	query2 := NewMirrorNodeContractCallQuery()
	query2.SetValue(value)
	assert.Equal(t, value, query2.GetValue())
}

func TestMirrorNodeContractQuerySetAndGetGasLimit(t *testing.T) {
	gas := int64(50000)

	query1 := NewMirrorNodeContractEstimateGasQuery()
	query1.SetGasLimit(gas)
	assert.Equal(t, gas, query1.GetGasLimit())

	query2 := NewMirrorNodeContractCallQuery()
	query2.SetGasLimit(gas)
	assert.Equal(t, gas, query2.GetGasLimit())
}

func TestMirrorNodeContractQuerySetAndGetGasPrice(t *testing.T) {
	gasPrice := int64(200)

	query1 := NewMirrorNodeContractEstimateGasQuery()
	query1.SetGasPrice(gasPrice)
	assert.Equal(t, gasPrice, query1.GetGasPrice())

	query2 := NewMirrorNodeContractCallQuery()
	query2.SetGasPrice(gasPrice)
	assert.Equal(t, gasPrice, query2.GetGasPrice())
}

func TestMirrorNodeContractQueryEstimateGasWithMissingContractIDOrEvmAddressThrowsException(t *testing.T) {
	query1 := &mirrorNodeContractQuery{}
	query1.setFunction("testFunction", NewContractFunctionParameters().AddString("params"))
	_, err1 := query1.estimateGas(nil)
	require.Error(t, err1)

	query2 := NewMirrorNodeContractEstimateGasQuery()
	query2.setFunction("testFunction", NewContractFunctionParameters().AddString("params"))
	_, err2 := query2.Execute(nil)
	require.Error(t, err2)
}

func TestMirrorNodeContractQueryCreateJSONPayloadAllFieldsSet(t *testing.T) {
	query := &mirrorNodeContractQuery{
		callData:           []byte("testData"),
		senderEvmAddress:   stringPtr("0x1234567890abcdef1234567890abcdef12345678"),
		contractEvmAddress: stringPtr("0xabcdefabcdefabcdefabcdefabcdefabcdef"),
		gasLimit:           int64Ptr(50000),
		gasPrice:           int64Ptr(2000),
		value:              int64Ptr(1000),
		blockNumber:        int64Ptr(123456),
	}

	jsonPayload, err := query.createJSONPayload(true, "latest")
	assert.NoError(t, err)

	expectedJson := `{"data":"7465737444617461","to":"0xabcdefabcdefabcdefabcdefabcdefabcdef","estimate":true,"blockNumber":"latest","from":"0x1234567890abcdef1234567890abcdef12345678","gas":50000,"gasPrice":2000,"value":1000}`
	assert.JSONEq(t, expectedJson, jsonPayload)
}

func TestMirrorNodeContractQueryCreateJSONPayloadOnlyRequiredFieldsSet(t *testing.T) {
	query := &mirrorNodeContractQuery{
		callData:           []byte("testData"),
		contractEvmAddress: stringPtr("0xabcdefabcdefabcdefabcdefabcdefabcdef"),
	}

	jsonPayload, err := query.createJSONPayload(true, "latest")
	assert.NoError(t, err)

	expectedJson := `{"data":"7465737444617461","to":"0xabcdefabcdefabcdefabcdefabcdefabcdef","estimate":true,"blockNumber":"latest"}`
	assert.JSONEq(t, expectedJson, jsonPayload)
}

func TestMirrorNodeContractQueryCreateJSONPayloadSomeOptionalFieldsSet(t *testing.T) {
	query := &mirrorNodeContractQuery{
		callData:           []byte("testData"),
		senderEvmAddress:   stringPtr("0x1234567890abcdef1234567890abcdef12345678"),
		contractEvmAddress: stringPtr("0xabcdefabcdefabcdefabcdefabcdefabcdef"),
		gasLimit:           int64Ptr(50000),
		value:              int64Ptr(1000),
	}

	jsonPayload, err := query.createJSONPayload(false, "latest")
	assert.NoError(t, err)

	expectedJson := `{"data":"7465737444617461","to":"0xabcdefabcdefabcdefabcdefabcdefabcdef","estimate":false,"blockNumber":"latest","from":"0x1234567890abcdef1234567890abcdef12345678","gas":50000,"value":1000}`
	assert.JSONEq(t, expectedJson, jsonPayload)
}

func TestMirrorNodeContractQueryCreateJSONPayloadAllOptionalFieldsDefault(t *testing.T) {
	query := &mirrorNodeContractQuery{
		callData:           []byte("testData"),
		contractEvmAddress: stringPtr("0xabcdefabcdefabcdefabcdefabcdefabcdef"),
	}

	jsonPayload, err := query.createJSONPayload(false, "latest")
	assert.NoError(t, err)

	expectedJson := `{"data":"7465737444617461","to":"0xabcdefabcdefabcdefabcdefabcdefabcdef","estimate":false,"blockNumber":"latest"}`
	assert.JSONEq(t, expectedJson, jsonPayload)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func TestUnitMirrorNodeContractQueryWithDifferentPorts(t *testing.T) {
	// Note: Not running in parallel since we modify global http.DefaultTransport

	tests := []struct {
		name           string
		domain         string
		expectedScheme string
		description    string
	}{
		{
			name:           "port 80 uses HTTP",
			domain:         "mirror80.example.com:80",
			expectedScheme: "http",
			description:    "Port 80 should use HTTP scheme",
		},
		{
			name:           "port 443 uses HTTPS",
			domain:         "mirror443.example.com:443",
			expectedScheme: "https",
			description:    "Port 443 should use HTTPS scheme",
		},
		{
			name:           "port 8443 uses HTTPS",
			domain:         "mirror8443.example.com:8443",
			expectedScheme: "https",
			description:    "Other ports should use HTTPS scheme for security",
		},
		{
			name:           "port 9999 uses HTTPS",
			domain:         "mirror9999.example.com:9999",
			expectedScheme: "https",
			description:    "Any non-standard port should use HTTPS scheme",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Test EstimateGas
			t.Run("EstimateGas", func(t *testing.T) {
				// Create a mock server that responds with gas estimation
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verify the request path contains contracts/call
					assert.Contains(t, r.URL.Path, "contracts/call")

					response := map[string]interface{}{
						"result": "0x5208", // 21000 gas in hex
					}
					w.Header().Set("Content-Type", "application/json")
					err := json.NewEncoder(w).Encode(response)
					require.NoError(t, err)
				}))
				defer server.Close()

				// Setup mock transport
				cleanup := SetupMockTransportForDomain(test.domain, server.URL)
				defer cleanup()

				// Setup client with the test domain as the mirror network
				client, err := _NewMockClient()
				require.NoError(t, err)
				client.SetLedgerID(*NewLedgerIDTestnet())
				client.SetMirrorNetwork([]string{test.domain})

				// Create a contract query
				query := NewMirrorNodeContractEstimateGasQuery()
				query.SetContractEvmAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
				query.SetFunction("testFunction", NewContractFunctionParameters().AddString("test"))

				// Test EstimateGas
				gasEstimate, err := query.Execute(client)
				require.NoError(t, err, "EstimateGas should succeed for %s", test.description)
				assert.Equal(t, uint64(21000), gasEstimate)
			})

			// Test ContractCall
			t.Run("ContractCall", func(t *testing.T) {
				// Create a mock server that responds with contract call result
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verify the request path contains contracts/call
					assert.Contains(t, r.URL.Path, "contracts/call")

					response := map[string]interface{}{
						"result": "0x0000000000000000000000000000000000000000000000000000000000000001",
					}
					w.Header().Set("Content-Type", "application/json")
					err := json.NewEncoder(w).Encode(response)
					require.NoError(t, err)
				}))
				defer server.Close()

				// Setup mock transport
				cleanup := SetupMockTransportForDomain(test.domain, server.URL)
				defer cleanup()

				// Setup client with the test domain as the mirror network
				client, err := _NewMockClient()
				require.NoError(t, err)
				client.SetLedgerID(*NewLedgerIDTestnet())
				client.SetMirrorNetwork([]string{test.domain})

				// Create a contract call query
				query := NewMirrorNodeContractCallQuery()
				query.SetContractEvmAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
				query.SetFunction("testFunction", NewContractFunctionParameters().AddString("test"))

				// Test ContractCall
				result, err := query.Execute(client)
				require.NoError(t, err, "ContractCall should succeed for %s", test.description)
				assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", result)
			})
		})
	}
}
