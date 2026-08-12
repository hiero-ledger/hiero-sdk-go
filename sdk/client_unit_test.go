//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitClientFromConfig(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfig([]byte(testClientJSON))
	require.NoError(t, err)

	assert.NotNil(t, client)
	assert.True(t, len(client.network.network) > 0)
	assert.Nil(t, client.operator)
	assert.Equal(t, uint64(3), client.GetShard())
	assert.Equal(t, uint64(5), client.GetRealm())
}

func TestUnitClientFromConfigWithOperator(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfig([]byte(testClientJSONWithOperator))
	require.NoError(t, err)

	assert.NotNil(t, client)

	testOperatorKey, err := PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	require.NoError(t, err)

	assert.True(t, len(client.network.network) > 0)
	assert.NotNil(t, client.operator)
	assert.Equal(t, testOperatorKey.ed25519PrivateKey.keyData, client.operator.privateKey.ed25519PrivateKey.keyData)
	assert.Equal(t, AccountID{Account: 3}.Account, client.operator.accountID.Account)
}

func TestUnitClientFromConfigWithoutMirrorNetwork(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfig([]byte(testClientJSONWithoutMirrorNetwork))
	require.NoError(t, err)
	assert.NotNil(t, client)

	assert.True(t, len(client.network.network) > 0)
	assert.True(t, len(client.GetMirrorNetwork()) == 0)
}

func TestUnitClientFromConfigWrongMirrorNetworkType(t *testing.T) {
	t.Parallel()

	_, err := ClientFromConfig([]byte(testClientJSONWrongTypeMirror))
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "mirrorNetwork is expected to be a string, an array of strings or nil", err.Error())
	}
}

func TestUnitClientFromConfigWrongNetworkType(t *testing.T) {
	t.Parallel()

	_, err := ClientFromConfig([]byte(testClientJSONWrongTypeNetwork))
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network is expected to be map of string to string, or string", err.Error())
	}
}

func TestUnitClientFromConfigWrongAccountIDNetworkType(t *testing.T) {
	_, err := ClientFromConfig([]byte(testClientJSONWrongAccountIDNetwork))
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "expected {shard}.{realm}.{num}", err.Error())
	}
}

func TestUnitClientFromCorrectConfigFile(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfigFile("client-config-with-operator.json")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.operator)
	assert.Equal(t, AccountID{Account: 3}.Account, client.operator.accountID.Account)
	assert.Equal(t, "a608e2130a0a3cb34f86e757303c862bee353d9ab77ba4387ec084f881d420d4", client.operator.privateKey.StringRaw())
}

func TestUnitClientFromMissingConfigFile(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfigFile("missing.json")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestUnitClientSetNetworkExtensive(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	nodes := make(map[string]AccountID, 2)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	network := client.GetNetwork()
	assert.Equal(t, 2, len(network))
	assert.Equal(t, network["0.testnet.hedera.com:50211"], AccountID{0, 0, 3, nil, nil, nil})
	assert.Equal(t, network["1.testnet.hedera.com:50211"], AccountID{0, 0, 4, nil, nil, nil})

	nodes = make(map[string]AccountID, 2)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["2.testnet.hedera.com:50211"] = AccountID{0, 0, 5, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	network = client.GetNetwork()
	assert.Equal(t, 3, len(network))
	assert.Equal(t, network["0.testnet.hedera.com:50211"], AccountID{0, 0, 3, nil, nil, nil})
	assert.Equal(t, network["1.testnet.hedera.com:50211"], AccountID{0, 0, 4, nil, nil, nil})
	assert.Equal(t, network["2.testnet.hedera.com:50211"], AccountID{0, 0, 5, nil, nil, nil})

	nodes = make(map[string]AccountID, 1)
	nodes["2.testnet.hedera.com:50211"] = AccountID{0, 0, 5, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	network = client.GetNetwork()
	networkMirror := client.GetMirrorNetwork()
	assert.Equal(t, 1, len(network))
	assert.Equal(t, network["2.testnet.hedera.com:50211"], AccountID{0, 0, 5, nil, nil, nil})
	// There is only one mirror address, no matter the transport security
	assert.Equal(t, "nonexistent-mirror-testnet:443", networkMirror[0])

	client.SetTransportSecurity(true)
	client.SetCertificateVerification(true)
	network = client.GetNetwork()
	networkMirror = client.GetMirrorNetwork()
	assert.Equal(t, network["2.testnet.hedera.com:50212"], AccountID{0, 0, 5, nil, nil, nil})
	assert.Equal(t, "nonexistent-mirror-testnet:443", networkMirror[0])

	err = client.Close()
	require.NoError(t, err)
}

func TestUnitClientSetMirrorNetwork(t *testing.T) {
	t.Parallel()

	mirrorNetworkString := "testnet.mirrornode.hedera.com:443"
	mirrorNetwork1String := "testnet1.mirrornode.hedera.com:443"
	defaultNetwork := make([]string, 0)
	defaultNetwork = append(defaultNetwork, mirrorNetworkString)
	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetMirrorNetwork(defaultNetwork)

	mirrorNetwork := client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, mirrorNetworkString, mirrorNetwork[0])

	defaultNetworkWithExtraNode := make([]string, 0)
	defaultNetworkWithExtraNode = append(defaultNetworkWithExtraNode, mirrorNetworkString)
	defaultNetworkWithExtraNode = append(defaultNetworkWithExtraNode, mirrorNetwork1String)

	client.SetMirrorNetwork(defaultNetworkWithExtraNode)
	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 2, len(mirrorNetwork))
	require.True(t, contains(mirrorNetwork, mirrorNetworkString))
	require.True(t, contains(mirrorNetwork, mirrorNetwork1String))

	defaultNetwork = make([]string, 0)
	defaultNetwork = append(defaultNetwork, mirrorNetwork1String)

	client.SetMirrorNetwork(defaultNetwork)
	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, mirrorNetwork1String, mirrorNetwork[0])

	defaultNetwork = make([]string, 0)
	defaultNetwork = append(defaultNetwork, mirrorNetworkString)

	client.SetMirrorNetwork(defaultNetwork)
	mirrorNetwork = client.GetMirrorNetwork()
	assert.Equal(t, 1, len(mirrorNetwork))
	assert.Equal(t, mirrorNetworkString, mirrorNetwork[0])

	client.SetTransportSecurity(true)
	mirrorNetwork = client.GetMirrorNetwork()
	// SetTransportSecurity is deprecated, so the mirror node should not be updated
	assert.Equal(t, mirrorNetworkString, mirrorNetwork[0])

	err = client.Close()
	require.NoError(t, err)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestUnitClientSetMultipleNetwork(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	nodes := make(map[string]AccountID, 8)
	nodes["0.testnet.hedera.com:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["34.94.106.61:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["50.18.132.211:50211"] = AccountID{0, 0, 3, nil, nil, nil}
	nodes["138.91.142.219:50211"] = AccountID{0, 0, 3, nil, nil, nil}

	nodes["1.testnet.hedera.com:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["35.237.119.55:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["3.212.6.13:50211"] = AccountID{0, 0, 4, nil, nil, nil}
	nodes["52.168.76.241:50211"] = AccountID{0, 0, 4, nil, nil, nil}

	err = client.SetNetwork(nodes)
	require.NoError(t, err)
	net := client.GetNetwork()

	if val, ok := net["0.testnet.hedera.com:50211"]; ok {
		require.Equal(t, val.String(), "0.0.3")
	}

	if val, ok := net["1.testnet.hedera.com:50211"]; ok {
		require.Equal(t, val.String(), "0.0.4")
	}

	if val, ok := net["50.18.132.211:50211"]; ok {
		require.Equal(t, val.String(), "0.0.3")
	}

	if val, ok := net["3.212.6.13:50211"]; ok {
		require.Equal(t, val.String(), "0.0.4")
	}

}

func TestUnitClientLogger(t *testing.T) {
	client := ClientForTestnet()

	var buf bytes.Buffer
	writer := zerolog.ConsoleWriter{Out: &buf, TimeFormat: time.RFC3339}

	hederaLoger := NewLogger("test", LoggerLevelTrace)

	l := zerolog.New(&writer)
	hederaLoger.logger = &l

	client.SetLogger(hederaLoger)
	client.SetLogLevel(LoggerLevelInfo)

	client.logger.Trace("trace message", "traceKey", "traceValue")
	client.logger.Debug("debug message", "debugKey", "debugValue")
	client.logger.Info("info message", "infoKey", "infoValue")
	client.logger.Warn("warn message", "warnKey", "warnValue")
	client.logger.Error("error message", "errorKey", "errorValue")

	assert.NotContains(t, buf.String(), "trace message")
	assert.NotContains(t, buf.String(), "debug message")
	assert.Contains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "error message")

	buf.Reset()
	client.SetLogLevel(LoggerLevelWarn)
	client.logger.Trace("trace message", "traceKey", "traceValue")
	client.logger.Debug("debug message", "debugKey", "debugValue")
	client.logger.Info("info message", "infoKey", "infoValue")
	client.logger.Warn("warn message", "warnKey", "warnValue")
	client.logger.Error("error message", "errorKey", "errorValue")

	assert.NotContains(t, buf.String(), "trace message")
	assert.NotContains(t, buf.String(), "debug message")
	assert.NotContains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "error message")

	buf.Reset()
	client.SetLogLevel(LoggerLevelTrace)
	client.logger.Trace("trace message", "traceKey", "traceValue")
	client.logger.Debug("debug message", "debugKey", "debugValue")
	client.logger.Info("info message", "infoKey", "infoValue")
	client.logger.Warn("warn message", "warnKey", "warnValue")
	client.logger.Error("error message", "errorKey", "errorValue")

	assert.Contains(t, buf.String(), "trace message")
	assert.Contains(t, buf.String(), "debug message")
	assert.Contains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "error message")

	hl := client.GetLogger()
	assert.Equal(t, hl, hederaLoger)
}

func TestUnitClientClientFromConfigWithoutScheduleNetworkUpdate(t *testing.T) {
	t.Parallel()

	client, err := ClientFromConfigWithoutScheduleNetworkUpdate([]byte(testClientJSON))
	require.NoError(t, err)
	assert.NotNil(t, client)

	assert.True(t, len(client.network.network) > 0)
	assert.Equal(t, time.Duration(0), client.GetNetworkUpdatePeriod())
}

func TestUnitClientPersistsShardAndRealm(t *testing.T) {
	t.Parallel()

	network := _NewNetwork()
	client := _NewClient(network, []string{}, NewLedgerIDTestnet(), true, 1, 2)
	assert.Equal(t, uint64(1), client.GetShard())
	assert.Equal(t, uint64(2), client.GetRealm())
}

func TestUnitClientForNameLocalhost(t *testing.T) {
	t.Parallel()

	client, err := ClientForName("localhost")
	require.NoError(t, err)
	defer client.Close()

	// The localhost preset targets the Solo local consensus node and mirror gRPC endpoint.
	network := client.GetNetwork()
	assert.Equal(t, AccountID{Account: 3}, network["127.0.0.1:35211"])
	assert.Contains(t, client.GetMirrorNetwork(), "127.0.0.1:5600")
}

func TestUnitClientForNetworkV2(t *testing.T) {
	t.Parallel()
	network := map[string]AccountID{
		"127.0.0.1:50211": {Account: 3, Shard: 1, Realm: 2},
		"127.0.0.1:50212": {Account: 4, Shard: 1, Realm: 2},
		"127.0.0.1:50213": {Account: 5, Shard: 1, Realm: 2},
	}
	client, err := ClientForNetworkV2(network)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), client.GetShard())
	assert.Equal(t, uint64(2), client.GetRealm())

	network = map[string]AccountID{
		"127.0.0.1:50211": {Account: 3, Shard: 2, Realm: 2},
		"127.0.0.1:50212": {Account: 4, Shard: 1, Realm: 2},
		"127.0.0.1:50213": {Account: 5, Shard: 1, Realm: 2},
	}

	client, err = ClientForNetworkV2(network)
	require.Error(t, err)
	assert.Equal(t, err.Error(), "network is not valid, all nodes must be in the same shard and realm")

	network = map[string]AccountID{
		"127.0.0.1:50211": {Account: 3, Shard: 1, Realm: 1},
		"127.0.0.1:50212": {Account: 4, Shard: 1, Realm: 2},
		"127.0.0.1:50213": {Account: 5, Shard: 1, Realm: 2},
	}

	client, err = ClientForNetworkV2(network)
	require.Error(t, err)
	assert.Equal(t, err.Error(), "network is not valid, all nodes must be in the same shard and realm")

	network = make(map[string]AccountID)
	client, err = ClientForNetworkV2(network)
	require.Error(t, err)
	assert.Equal(t, err.Error(), "network is empty")
}

func TestUnitClientGetMirrorRestApiBaseUrl(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		domain         string
		expectedScheme string
	}{
		{
			name:           "HTTP for port 80",
			domain:         "mirror.example.com:80",
			expectedScheme: "http",
		},
		{
			name:           "HTTPS for port 443",
			domain:         "mirror.example.com:443",
			expectedScheme: "https",
		},
		{
			name:           "HTTPS for custom port",
			domain:         "mirror.example.com:8080",
			expectedScheme: "https",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := _NewMockClient()
			require.NoError(t, err)
			client.SetLedgerID(*NewLedgerIDTestnet())
			client.SetMirrorNetwork([]string{test.domain})

			baseURL, err := client.GetMirrorRestApiBaseUrl()
			require.NoError(t, err)

			parsedURL, err := url.Parse(baseURL)
			require.NoError(t, err)
			assert.Equal(t, test.expectedScheme, parsedURL.Scheme)

			assert.Equal(t, test.domain, parsedURL.Host)
			assert.Equal(t, "/api/v1", parsedURL.Path)
		})
	}
}

func TestUnitClientGetMirrorRestApiBaseUrlLocalHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		domain         string
		expectedScheme string
		expectedPort   string
	}{
		{
			name:           "HTTP local host",
			domain:         "localhost:80",
			expectedScheme: "http",
		},
		{
			name:           "HTTP 127.0.0.1",
			domain:         "127.0.0.1:8080",
			expectedScheme: "http",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := _NewMockClient()
			require.NoError(t, err)
			client.SetLedgerID(*NewLedgerIDTestnet())
			client.SetMirrorNetwork([]string{test.domain})

			baseURL, err := client.GetMirrorRestApiBaseUrl()
			require.NoError(t, err)

			parsedURL, err := url.Parse(baseURL)
			require.NoError(t, err)
			assert.Equal(t, test.expectedScheme, parsedURL.Scheme)

			if test.domain == "localhost:80" {
				assert.Equal(t, "localhost:38081", parsedURL.Host)
			} else {
				assert.Equal(t, "127.0.0.1:38081", parsedURL.Host)
			}

			assert.Equal(t, "/api/v1", parsedURL.Path)
		})
	}
}

// _CostAnswer returns a COST_ANSWER handler reporting the given precheck status, recording every
// request it received.
func _CostAnswer(status services.ResponseCodeEnum, captured *[]*services.Query) func(*services.Query) *services.Response {
	var mu sync.Mutex
	return func(request *services.Query) *services.Response {
		mu.Lock()
		*captured = append(*captured, request)
		mu.Unlock()

		return &services.Response{
			Response: &services.Response_CryptoGetInfo{
				CryptoGetInfo: &services.CryptoGetInfoResponse{
					Header: &services.ResponseHeader{
						NodeTransactionPrecheckCode: status,
						ResponseType:                services.ResponseType_COST_ANSWER,
						Cost:                        25,
					},
				},
			},
		}
	}
}

// _AssertIsCostAnswerProbe asserts the request is the liveness probe: a COST_ANSWER getAccountInfo
// for the treasury account, carrying no payment.
func _AssertIsCostAnswerProbe(t *testing.T, request *services.Query) {
	t.Helper()

	info, ok := request.Query.(*services.Query_CryptoGetInfo)
	require.True(t, ok, "probe must be a getAccountInfo query, was: %v", request.Query)
	require.Equal(t, services.ResponseType_COST_ANSWER, info.CryptoGetInfo.Header.ResponseType)
	require.Nil(t, info.CryptoGetInfo.Header.Payment, "the probe must not attach a payment transaction")
	require.Equal(t,
		AccountID{Account: treasuryAccountNum}._ToProtobuf().String(),
		info.CryptoGetInfo.AccountID.String())
}

// Ping probes with an unpaid COST_ANSWER getAccountInfo query, and needs no operator.
func TestUnitClientPingReachableNode(t *testing.T) {
	t.Parallel()

	var captured []*services.Query
	client, server := NewMockClientAndServer([][]interface{}{{_CostAnswer(services.ResponseCodeEnum_OK, &captured)}})
	defer server.Close()

	client.operator = nil

	require.NoError(t, client.Ping(AccountID{Account: 3}))
	require.Len(t, captured, 1, "Ping must send exactly one probe query")
	_AssertIsCostAnswerProbe(t, captured[0])
}

// PingAll probes every node in the network.
func TestUnitClientPingAllProbesEveryNode(t *testing.T) {
	t.Parallel()

	var node3, node4, node5 []*services.Query
	client, server := NewMockClientAndServer([][]interface{}{
		{_CostAnswer(services.ResponseCodeEnum_OK, &node3)},
		{_CostAnswer(services.ResponseCodeEnum_OK, &node4)},
		{_CostAnswer(services.ResponseCodeEnum_OK, &node5)},
	})
	defer server.Close()

	require.Len(t, client.GetNetwork(), 3)
	client.PingAll()

	for _, captured := range [][]*services.Query{node3, node4, node5} {
		require.Len(t, captured, 1)
		_AssertIsCostAnswerProbe(t, captured[0])
	}
}

// A successful probe halves the node's backoff, so a recovered node is retried sooner if it fails
// again.
func TestUnitClientPingSuccessDecreasesNodeBackoff(t *testing.T) {
	t.Parallel()

	var captured []*services.Query
	client, server := NewMockClientAndServer([][]interface{}{{_CostAnswer(services.ResponseCodeEnum_OK, &captured)}})
	defer server.Close()

	node, ok := client.network._GetNodeForAccountID(AccountID{Account: 3})
	require.True(t, ok)
	node.minBackoff = 250 * time.Millisecond
	node.currentBackoff = 4 * time.Second

	require.NoError(t, client.Ping(AccountID{Account: 3}))
	assert.Equal(t, 2*time.Second, node.currentBackoff)
}

// Ping cannot readmit a node: a node in backoff is never contacted, so the probe fails without
// reaching the network even if the node has recovered. Readmission happens when the backoff elapses.
func TestUnitClientPingSkipsNodeInBackoff(t *testing.T) {
	t.Parallel()

	var captured []*services.Query
	client, server := NewMockClientAndServer([][]interface{}{{_CostAnswer(services.ResponseCodeEnum_OK, &captured)}})
	defer server.Close()

	node, ok := client.network._GetNodeForAccountID(AccountID{Account: 3})
	require.True(t, ok)
	readmitTime := time.Now().Add(time.Hour)
	node.readmitTime = &readmitTime
	require.False(t, node._IsHealthy())

	client.SetMaxAttempts(1)

	require.Error(t, client.Ping(AccountID{Account: 3}))
	assert.Empty(t, captured, "a node in backoff must not be probed")
}

// The other half of that contract: a successful probe lowers the backoff but leaves the readmit time
// and failure count untouched, so a node stays unhealthy until its backoff elapses.
func TestUnitClientPingSuccessDoesNotClearBackoff(t *testing.T) {
	t.Parallel()

	var captured []*services.Query
	client, server := NewMockClientAndServer([][]interface{}{{_CostAnswer(services.ResponseCodeEnum_OK, &captured)}})
	defer server.Close()

	node, ok := client.network._GetNodeForAccountID(AccountID{Account: 3})
	require.True(t, ok)
	// One millisecond past the readmit time: the probe reaches the network while the bookkeeping from
	// the earlier failure is still in place.
	readmitTime := time.Now().Add(-time.Millisecond)
	node.readmitTime = &readmitTime
	node.badGrpcStatusCount = 3

	require.NoError(t, client.Ping(AccountID{Account: 3}))
	require.Len(t, captured, 1)
	assert.Equal(t, &readmitTime, node._GetReadmitTime(), "a successful probe must not clear the readmit time")
	assert.Equal(t, int64(3), node._GetAttempts(), "a successful probe must not reset the failure count")
}

// A non-retryable precheck failure surfaces to the caller without a retry.
func TestUnitClientPingPrecheckFailure(t *testing.T) {
	t.Parallel()

	var captured []*services.Query
	invalid := _CostAnswer(services.ResponseCodeEnum_INVALID_ACCOUNT_ID, &captured)

	// Queue extra responses so a wrongly retrying Ping would consume more than one.
	client, server := NewMockClientAndServer([][]interface{}{{invalid, invalid, invalid}})
	defer server.Close()

	err := client.Ping(AccountID{Account: 3})
	require.ErrorContains(t, err, "INVALID_ACCOUNT_ID")
	require.Len(t, captured, 1, "a non-retryable precheck must not be retried")
}

// _BusyCostAnswer returns a COST_ANSWER handler that always reports BUSY, a retryable precheck
// status, and counts its calls.
func _BusyCostAnswer(calls *int) func(*services.Query) *services.Response {
	return func(request *services.Query) *services.Response {
		*calls++
		return &services.Response{
			Response: &services.Response_CryptoGetInfo{
				CryptoGetInfo: &services.CryptoGetInfoResponse{
					Header: &services.ResponseHeader{
						NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
						ResponseType:                services.ResponseType_COST_ANSWER,
					},
				},
			},
		}
	}
}

// With MaxAttempts(1) — the configuration for pruning dead nodes with PingAll — a retryable BUSY
// precheck fails Ping without a retry.
func TestUnitClientPingBusyNodeSingleAttempt(t *testing.T) {
	t.Parallel()

	probeCalls := 0
	busy := _BusyCostAnswer(&probeCalls)

	// Queue several responses so a wrongly retrying Ping would consume more than one.
	client, server := NewMockClientAndServer([][]interface{}{{busy, busy, busy}})
	defer server.Close()

	client.SetMaxAttempts(1)

	err := client.Ping(AccountID{Account: 3})
	require.Error(t, err)
	require.ErrorContains(t, err, "BUSY")
	require.Equal(t, 1, probeCalls, "Ping must not retry when the client allows a single attempt")
}

// Ping imposes no retry limit of its own: BUSY is retried up to the client's MaxAttempts, matching
// the Java SDK, which leaves the fail-fast choice to the caller.
func TestUnitClientPingBusyNodeFollowsClientMaxAttempts(t *testing.T) {
	t.Parallel()

	probeCalls := 0
	busy := _BusyCostAnswer(&probeCalls)

	client, server := NewMockClientAndServer([][]interface{}{{busy, busy, busy, busy, busy}})
	defer server.Close()

	client.SetMaxAttempts(3)

	err := client.Ping(AccountID{Account: 3})
	require.Error(t, err)
	require.ErrorContains(t, err, "BUSY")
	require.Equal(t, 3, probeCalls, "Ping must retry up to the client's MaxAttempts")
}

// The retry behaviour Ping relies on: a query with no max retry of its own follows the client's
// MaxAttempts.
func TestUnitQueryClientMaxAttemptsAppliesWithoutMaxRetry(t *testing.T) {
	t.Parallel()

	queryCalls := 0
	busy := _BusyCostAnswer(&queryCalls)

	client, server := NewMockClientAndServer([][]interface{}{{busy, busy, busy, busy, busy}})
	defer server.Close()

	client.SetMaxAttempts(3)

	_, err := NewAccountInfoQuery().
		SetAccountID(AccountID{Account: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		GetCost(client)
	require.Error(t, err)
	require.Equal(t, 3, queryCalls, "the client's MaxAttempts must still apply when the request sets no max retry")
}

// Ping rejects a node that is not part of the client's network.
func TestUnitClientPingUnknownNode(t *testing.T) {
	t.Parallel()

	client, server := NewMockClientAndServer([][]interface{}{{}})
	defer server.Close()

	unknown := AccountID{Account: 99}
	var invalidNode ErrInvalidNodeAccountIDSet
	require.ErrorAs(t, client.Ping(unknown), &invalidNode)
	assert.Equal(t, unknown, invalidNode.NodeAccountID)
}

// Ping surfaces a gRPC transport failure and records it against the node, which is what lets PingAll
// prune dead nodes.
func TestUnitClientPingUnreachableNode(t *testing.T) {
	t.Parallel()

	client, server := NewMockClientAndServer([][]interface{}{{}})
	// Close the server so the probe fails at the gRPC layer.
	server.Close()
	client.SetGrpcDeadline(500 * time.Millisecond)
	client.SetMaxAttempts(1)

	node, ok := client.network._GetNodeForAccountID(AccountID{Account: 3})
	require.True(t, ok)

	err := client.Ping(AccountID{Account: 3})
	require.ErrorContains(t, err, "Unavailable")
	assert.Equal(t, int64(1), node._GetAttempts(), "a failed probe must be recorded against the node")
}

// Ping fails fast against an address nothing is listening on, the case PingAll with
// SetMaxAttempts(1) is meant to prune.
func TestUnitClientPingUnreachableAddress(t *testing.T) {
	t.Parallel()

	// Port 1 is not bound, so the node resolves within the network but the connection is refused.
	unreachable := AccountID{Account: 99}
	client, err := ClientForNetworkV2(map[string]AccountID{"127.0.0.1:1": unreachable})
	require.NoError(t, err)
	client.SetMaxAttempts(1)

	start := time.Now()
	err = client.Ping(unreachable)
	elapsed := time.Since(start)

	require.ErrorContains(t, err, "Unavailable")
	assert.Less(t, elapsed, 30*time.Second, "a single-attempt Ping must fail fast")
}
