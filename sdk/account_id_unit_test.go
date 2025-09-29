//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitAccountIDChecksumFromString(t *testing.T) {
	t.Parallel()

	id, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	id.ToStringWithChecksum(client)
	id.GetChecksum()
	sol := id.ToSolidityAddress()
	AccountIDFromSolidityAddress(sol)
	err = id.Validate(client)
	require.Error(t, err)
	evmID, err := AccountIDFromEvmAddress(0, 0, "0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	require.NoError(t, err)
	pb := evmID._ToProtobuf()
	_AccountIDFromProtobuf(pb)

	idByte := id.ToBytes()
	AccountIDFromBytes(idByte)

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	alias := key.ToAccountID(0, 0)
	pb = alias._ToProtobuf()
	_AccountIDFromProtobuf(pb)

	require.NoError(t, err)
	assert.Equal(t, id.Account, uint64(123))
}

func TestUnitAccountIDChecksumToString(t *testing.T) {
	t.Parallel()

	id := AccountID{
		Shard:   50,
		Realm:   150,
		Account: 520,
	}
	assert.Equal(t, "50.150.520", id.String())
}

func TestUnitAccountIDFromStringAlias(t *testing.T) {
	t.Parallel()

	key, err := GeneratePrivateKey()
	require.NoError(t, err)
	id, err := AccountIDFromString("0.0." + key.PublicKey().String())
	require.NoError(t, err)
	id2 := key.ToAccountID(0, 0)

	assert.Equal(t, id.String(), id2.String())
}

func TestUnitChecksum(t *testing.T) {
	t.Parallel()

	id, err := LedgerIDFromString("01")
	require.NoError(t, err)
	ad1, err := _ChecksumParseAddress(id, "0.0.3")
	require.NoError(t, err)
	id, err = LedgerIDFromString("10")
	require.NoError(t, err)
	ad2, err := _ChecksumParseAddress(id, "0.0.3")
	require.NoError(t, err)

	require.NotEqual(t, ad1.correctChecksum, ad2.correctChecksum)
}

func TestUnitAccountIDEvm(t *testing.T) {
	t.Parallel()

	id, err := AccountIDFromString("0.0.0011223344556677889900112233445566778899")
	require.NoError(t, err)

	require.Equal(t, id.String(), "0.0.0011223344556677889900112233445566778899")
}

func TestUnitAccountIDPopulateFailForWrongMirrorHost(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateAccount(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateFailWithNoMirror(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	privateKey, err := PrivateKeyGenerateEcdsa()
	require.NoError(t, err)
	publicKey := privateKey.PublicKey()
	evmAddress := publicKey.ToEvmAddress()
	evmAddressAccountID, err := AccountIDFromEvmPublicAddress(evmAddress)
	require.NoError(t, err)
	err = evmAddressAccountID.PopulateAccount(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateEvmFailForWrongMirrorHost(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	id, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	err = id.PopulateEvmAddress(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateEvmFailWithNoMirror(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	id, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	err = id.PopulateEvmAddress(client)
	require.Error(t, err)
}

func TestUnitAccountIDPopulateEvmFailWithNoMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.mirrorNetwork = nil
	client.SetLedgerID(*NewLedgerIDTestnet())
	id, err := AccountIDFromString("0.0.3")
	require.NoError(t, err)
	err = id.PopulateEvmAddress(client)
	require.Error(t, err)
}

func TestUnitAccountIDChecksumError(t *testing.T) {
	t.Parallel()

	id, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)
	client, err := _NewMockClient()
	require.NoError(t, err)
	_, err = id.ToStringWithChecksum(client)
	require.Error(t, err)
}

func TestUnitAccountIDFromEvmAddressIncorrectSize(t *testing.T) {
	t.Parallel()

	// Test with an EVM address that's too short
	_, err := AccountIDFromEvmAddress(0, 0, "abc123")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotCorrectSize)

	// Test with an EVM address that's too long
	_, err = AccountIDFromEvmAddress(0, 0, "0123456789abcdef0123456789abcdef0123456789abcdef")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotCorrectSize)

	// Test with a 0x prefix that gets removed but then is too short
	_, err = AccountIDFromEvmAddress(0, 0, "0xabc123")
	require.Error(t, err)
	require.ErrorIs(t, err, errEvmAddressIsNotCorrectSize)

	// Verify a correct length works
	correctAddress := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	id, err := AccountIDFromEvmAddress(0, 0, correctAddress)
	require.NoError(t, err)
	require.NotNil(t, id.AliasEvmAddress)

	require.Equal(t, strings.ToLower(evmAddress), hex.EncodeToString(*id.AliasEvmAddress))
}

func TestUnitAccountIDFromEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal EVM address
	evmAddress := evmAddress
	bytes, err := hex.DecodeString(evmAddress)
	require.NoError(t, err)
	id, err := AccountIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(0), id.Account)
	require.Equal(t, bytes, *id.AliasEvmAddress)

	// Test with a different shard and realm
	id, err = AccountIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(0), id.Account)
	require.Equal(t, bytes, *id.AliasEvmAddress)

	// Test with a long zero address
	evmAddress = longZeroAddress
	bytes, err = hex.DecodeString(evmAddress)
	require.NoError(t, err)
	id, err = AccountIDFromEvmAddress(0, 0, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(0), id.Shard)
	require.Equal(t, uint64(0), id.Realm)
	require.Equal(t, uint64(0), id.Account)
	require.Equal(t, bytes, *id.AliasEvmAddress)

	// Test with a different shard and realm
	id, err = AccountIDFromEvmAddress(1, 1, evmAddress)
	require.NoError(t, err)
	require.Equal(t, uint64(1), id.Shard)
	require.Equal(t, uint64(1), id.Realm)
	require.Equal(t, uint64(0), id.Account)
	require.Equal(t, bytes, *id.AliasEvmAddress)
}

func TestUnitAccountIDToEvmAddress(t *testing.T) {
	t.Parallel()

	// Test with a normal account ID
	id := AccountID{Shard: 0, Realm: 0, Account: 123}
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a different shard and realm
	id = AccountID{Shard: 1, Realm: 1, Account: 123}
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a long zero address
	bytes, err := hex.DecodeString(longZeroAddress)
	id = AccountID{Shard: 1, Realm: 1, AliasEvmAddress: &bytes}
	require.NoError(t, err)
	require.Equal(t, longZeroAddress, id.ToEvmAddress())

	// Test with a normal EVM address
	emvAddress := evmAddress
	bytes, err = hex.DecodeString(emvAddress)
	id = AccountID{Shard: 0, Realm: 0, AliasEvmAddress: &bytes}
	expected := strings.ToLower(evmAddress)
	require.NoError(t, err)
	require.Equal(t, expected, id.ToEvmAddress())

	// Test with different shard and realm
	id = AccountID{Shard: 1, Realm: 1, AliasEvmAddress: &bytes}
	require.Equal(t, expected, id.ToEvmAddress())
}

func TestUnitAccountIDPopulateWithDifferentPorts(t *testing.T) {
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
			t.Run("PopulateAccount", func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Contains(t, r.URL.Path, "accounts")

					response := map[string]interface{}{
						"account": "0.0.12345",
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

				// Create an account ID with EVM address
				evmAddressBytes, err := hex.DecodeString(evmAddress)
				require.NoError(t, err)
				accountID := AccountID{
					Shard:           0,
					Realm:           0,
					Account:         0,
					AliasEvmAddress: &evmAddressBytes,
				}

				// Test PopulateAccount
				err = accountID.PopulateAccount(client)
				require.NoError(t, err, "PopulateAccount should succeed for %s", test.description)
				assert.Equal(t, uint64(12345), accountID.Account)
			})

			// Test PopulateEvmAddress
			t.Run("PopulateEvmAddress", func(t *testing.T) {
				// Create a mock server that responds with EVM address data
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verify the request path
					assert.Contains(t, r.URL.Path, "accounts/0.0.789")

					response := map[string]interface{}{
						"evm_address": "0x" + evmAddress,
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

				// Create an account ID with account number
				accountID := AccountID{
					Shard:   0,
					Realm:   0,
					Account: 789,
				}

				// Test PopulateEvmAddress
				err = accountID.PopulateEvmAddress(client)
				require.NoError(t, err, "PopulateEvmAddress should succeed for %s", test.description)
				require.NotNil(t, accountID.AliasEvmAddress)

				expectedBytes, err := hex.DecodeString(evmAddress)
				require.NoError(t, err)
				assert.Equal(t, expectedBytes, *accountID.AliasEvmAddress)
			})
		})
	}
}
