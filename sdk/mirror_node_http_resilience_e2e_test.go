//go:build all || e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mirrorRestLoopbackAddress = "127.0.0.1:38081"

func newLocalMirrorClient(t *testing.T, mirror ...string) *Client {
	t.Helper()

	if os.Getenv("HEDERA_NETWORK") != "localhost" {
		t.Skip("needs HEDERA_NETWORK=localhost and a local Solo network")
	}

	client, err := ClientForNetworkV2(map[string]AccountID{"127.0.0.1:50211": {Account: 3}})
	require.NoError(t, err)
	client.SetMirrorNetwork(mirror)
	t.Cleanup(func() { _ = client.Close() })

	return client
}

// MIRROR_REST_UP_CMD must serve mirror REST on 127.0.0.1:38081, for Solo
// "kubectl --context kind-solo-cluster port-forward svc/mirror-1-rest -n solo 38081:80".
// The test runs it 20 seconds after the queries start.
func TestIntegrationMirrorNodeHttpWaitsForStartingMirrorNode(t *testing.T) {
	upCmd := os.Getenv("MIRROR_REST_UP_CMD")
	if upCmd == "" {
		t.Skip("set MIRROR_REST_UP_CMD to the command that serves mirror REST on " + mirrorRestLoopbackAddress)
	}
	if conn, err := net.DialTimeout("tcp", mirrorRestLoopbackAddress, time.Second); err == nil {
		_ = conn.Close()
		t.Skip("mirror REST is already up on " + mirrorRestLoopbackAddress + "; stop it so the test can bring it up late")
	}

	client := newLocalMirrorClient(t, "127.0.0.1:5600")

	const upAfter = 20 * time.Second
	up := exec.Command("sh", "-c", "exec "+upCmd)
	timer := time.AfterFunc(upAfter, func() {
		if err := up.Start(); err != nil {
			t.Errorf("starting %q: %v", upCmd, err)
		}
	})
	t.Cleanup(func() {
		timer.Stop()
		if up.Process != nil {
			_ = up.Process.Kill()
			_ = up.Wait()
		}
	})

	calls := map[string]func() error{
		"MirrorNodeAccountBalanceQuery": func() error {
			_, err := NewMirrorNodeAccountBalanceQuery().SetAccountID(AccountID{Account: 2}).Execute(client)
			return err
		},
		"AccountID.PopulateEvmAddress": func() error {
			id := AccountID{Account: 2}
			return id.PopulateEvmAddress(client)
		},
	}

	start := time.Now()
	var wg sync.WaitGroup
	for name, call := range calls {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := call()
			elapsed := time.Since(start).Round(100 * time.Millisecond)
			t.Logf("%-30s returned after %5s: %v", name, elapsed, err)

			// Any answer from the mirror node counts, even "not found" on a mirror still
			// importing: the question is whether the SDK reached it or gave up first.
			assert.False(t, errors.Is(err, syscall.ECONNREFUSED),
				"%s gave up after %s, while the mirror node was still starting (it came up at %s)", name, elapsed, upAfter)
		}()
	}
	wg.Wait()
}

func TestIntegrationMirrorNodeHttpSelectsMirrorNodesRoundRobin(t *testing.T) {
	client := newLocalMirrorClient(t, "localhost:5600", "127.0.0.1:5600")

	previous := ""
	for call := range 8 {
		baseURL, err := client.GetMirrorRestApiBaseUrl()
		require.NoError(t, err)
		if call > 0 {
			assert.NotEqual(t, previous, baseURL, "call %d reused the node call %d used", call, call-1)
		}
		previous = baseURL
	}
}

func TestIntegrationMirrorNodeHttpRemoteHostIsNotLocal(t *testing.T) {
	client := newLocalMirrorClient(t, "mylocalhost.example.com:443")

	_, err := NewRegisteredNodeAddressBookQuery().SetMaxAttempts(1).Execute(client)
	require.Error(t, err)
}
