//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitContractIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := ContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	sol := id.ToSolidityAddress()
	ContractIDFromSolidityAddress(sol)
	err = id.Validate(client)
	require.Error(t, err)
	evmID, err := ContractIDFromEvmAddress(0, 0, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4")
	require.NoError(t, err)
	pb := evmID._ToProtobuf()
	_ContractIDFromProtobuf(pb)

	idByte := id.ToBytes()
	ContractIDFromBytes(idByte)

	id._ToProtoKey()

	assert.Equal(t, id.Contract, uint64(123))
}

func TestUnitContractIDFromEvmAddressIncorrectSize(t *testing.T) {
	t.Parallel()

	// Test with an EVM address that's too short
	_, err := ContractIDFromEvmAddress(0, 0, "abc123")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotCorrectSize)

	// Test with an EVM address that's too long
	_, err = ContractIDFromEvmAddress(0, 0, "0123456789abcdef0123456789abcdef0123456789abcdef")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotCorrectSize)

	// Test with a 0x prefix that gets removed but then is too short
	_, err = ContractIDFromEvmAddress(0, 0, "0xabc123")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotCorrectSize)

	// Verify a correct length works
	correctAddress := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	id, err := ContractIDFromEvmAddress(0, 0, correctAddress)
	require.NoError(t, err)
	require.NotNil(t, id.EvmAddress)

	require.Equal(t, strings.ToLower(evmAddress), hex.EncodeToString(id.EvmAddress))
}

func TestUnitContractIDChecksumToString(t *testing.T) {
	t.Parallel()

	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitContractIDFromStringEVM(t *testing.T) {
	t.Parallel()

	id, err := ContractIDFromString("0.0.0011223344556677889900112233445577889900")
	require.NoError(t, err)

	require.Equal(t, "0.0.0011223344556677889900112233445577889900", id.String())
}

func TestUnitContractIDProtobuf(t *testing.T) {
	t.Parallel()

	id, err := ContractIDFromString("0.0.0011223344556677889900112233445577889900")
	require.NoError(t, err)

	pb := id._ToProtobuf()

	decoded, err := hex.DecodeString("0011223344556677889900112233445577889900")
	require.NoError(t, err)

	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_EvmAddress{EvmAddress: decoded},
	})

	pbFrom := _ContractIDFromProtobuf(pb)

	require.Equal(t, id, *pbFrom)
}

func TestUnitContractIDEvm(t *testing.T) {
	t.Parallel()

	hexString, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	id, err := ContractIDFromString(fmt.Sprintf("0.0.%s", hexString.PublicKey().String()))
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(id.EvmAddress), hexString.PublicKey().String())

	pb := id._ToProtobuf()
	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_EvmAddress{EvmAddress: id.EvmAddress},
	})

	id, err = ContractIDFromString("0.0.123")
	require.NoError(t, err)
	require.Equal(t, id.Contract, uint64(123))
	require.Nil(t, id.EvmAddress)

	pb = id._ToProtobuf()
	require.Equal(t, pb, &services.ContractID{
		ShardNum: 0,
		RealmNum: 0,
		Contract: &services.ContractID_ContractNum{ContractNum: 123},
	})
}

func TestUnitContractIDPopulateFailForWrongMirrorHost(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := ContractIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateContract(client)
	require.Error(t, err)
}

func TestUnitContractIDPopulateFailWithNoMirror(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := ContractIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateContract(client)
	require.Error(t, err)
}

func TestUnitContractIDChecksumError(t *testing.T) {
	t.Parallel()

	id, err := ContractIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	client, err := _NewMockClient()

	_, err = id.ToStringWithChecksum(*client)
	require.ErrorContains(t, err, "can't derive checksum for ID")

	id.EvmAddress = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	_, err = id.ToStringWithChecksum(*client)
	require.ErrorContains(t, err, "EvmAddress doesn't support checksums")
}

func TestUnitContractIDFromEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal EVM address
	evmAddress := evmAddress
	bytes, err := hex.DecodeString(evmAddress)
	require.NoError(t, err)
	id, err := ContractIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(0), id.Contract)
	require.Equal(t, bytes, id.EvmAddress)

	// Test with a different shard and realm
	id, err = ContractIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(0), id.Contract)
	require.Equal(t, bytes, id.EvmAddress)

	// Test with a long zero address
	evmAddress = longZeroAddress
	bytes, err = hex.DecodeString(evmAddress)
	require.NoError(t, err)
	id, err = ContractIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(0), id.Contract)
	require.Equal(t, bytes, id.EvmAddress)

	// Test with a different shard and realm
	id, err = ContractIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(0), id.Contract)
	require.Equal(t, bytes, id.EvmAddress)
}

func TestUnitContractIDToEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal contract ID
	id := ContractID{Shard: 0, Realm: 0, Contract: 123}
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a different shard and realm
	id = ContractID{Shard: 1, Realm: 1, Contract: 123}
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a long zero address
	longZeroAddress := longZeroAddress
	bytes, err := hex.DecodeString(longZeroAddress)
	id = ContractID{Shard: 1, Realm: 1, EvmAddress: bytes}
	require.NoError(t, err)
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a normal EVM address
	evmAddress := evmAddress
	bytes, err = hex.DecodeString(evmAddress)
	id = ContractID{Shard: 0, Realm: 0, EvmAddress: bytes}
	expected := strings.ToLower(evmAddress)
	require.NoError(t, err)
	require.Equal(t, expected, id.ToEvmAddress())

	// Test with different shard and realm
	id = ContractID{Shard: 1, Realm: 1, EvmAddress: bytes}
	require.Equal(t, expected, id.ToEvmAddress())
}

func TestUnitContractIDPopulateContractWithDifferentPorts(t *testing.T) {
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
			// Test PopulateContract
			t.Run("PopulateContract", func(t *testing.T) {
				// Create a mock server that responds with contract data
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verify the request path contains contracts
					assert.Contains(t, r.URL.Path, "contracts")

					response := map[string]interface{}{
						"contract_id": "0.0.12345",
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

				// Create a contract ID with EVM address
				evmAddressBytes, err := hex.DecodeString(evmAddress)
				require.NoError(t, err)
				contractID := ContractID{
					Shard:      0,
					Realm:      0,
					Contract:   0,
					EvmAddress: evmAddressBytes,
				}

				// Test PopulateContract
				err = contractID.PopulateContract(client)
				require.NoError(t, err, "PopulateContract should succeed for %s", test.description)
				assert.Equal(t, uint64(12345), contractID.Contract)
			})
		})
	}
}
