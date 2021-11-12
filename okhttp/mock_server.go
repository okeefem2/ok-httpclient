package okhttp

import (
	"fmt"
	"sync"
)

// Singleton mock mocks
var (
	mocks = mockServer{
		mocks: make(map[string]*Mock),
	}
)

type mockServer struct {
	enabled bool
	// Can be used to make sure only a single go routine can access this at a given time
	serverMutex sync.Mutex

	mocks map[string]*Mock // TODO prefix trie
}

func StartMockServer() {
	mocks.serverMutex.Lock() // Other go routines will be stopped at this line until unlock is called by the goroutine with the lock
	// I am assuming: this doesn't halt anything, because go will just de schedule the blocked go routines, but those won't get anything done while stuck here
	// so probably a good idea not to have too much computation in locked code
	defer mocks.serverMutex.Unlock()
	mocks.enabled = true

}

func StopMockServer() {
	mocks.serverMutex.Lock()
	defer mocks.serverMutex.Unlock()
	mocks.enabled = false
}

func AddMock(mock Mock) {
	mocks.serverMutex.Lock()
	defer mocks.serverMutex.Unlock()
	key := mocks.getMockKey(mock.Method, mock.Url,  mock.RequestBody)
	mocks.mocks[key] = &mock
}

func (m *mockServer) getMockKey(method string, url string, body string) string {
	return method + url + body
}

func (m *mockServer) getMock(method string, url string, body string) *Mock {
	if !m.enabled {
		return nil
	}
	mock := m.mocks[m.getMockKey(method, url, body)]

	if mock != nil {
		return mock
	}

	return &Mock{
		Error: fmt.Errorf("no mock matching %s from %s with the given body", method, url),
	}
}
