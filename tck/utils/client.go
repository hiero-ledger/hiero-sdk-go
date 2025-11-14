// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"sync"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type SafeClientMap struct {
	clients map[string]*hiero.Client
	mu      sync.Mutex
}

func NewSafeClientMap() *SafeClientMap {
	return &SafeClientMap{
		clients: make(map[string]*hiero.Client),
	}
}

func (m *SafeClientMap) Get(sessionId string) *hiero.Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients[sessionId]
}

func (m *SafeClientMap) Set(sessionId string, client *hiero.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[sessionId] = client
}

func (m *SafeClientMap) Delete(sessionId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, sessionId)
}
