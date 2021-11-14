package okhttp_mock

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"sync"

	"github.com/okeefem2/ok-httpclient/core"
)

// Singleton mock MockServer
var (
	MockServer = mockServer{
		mocks: make(map[string]*Mock),
		httpClient: &httpClientMock{},
	}
)

type mockServer struct {
	enabled bool
	// Can be used to make sure only a single go routine can access this at a given time
	serverMutex sync.Mutex

	mocks map[string]*Mock // TODO prefix trie

	httpClient core.HttpClient
}

func (m *mockServer) Start() {
	m.serverMutex.Lock() // Other go routines will be stopped at this line until unlock is called by the goroutine with the lock
	// I am assuming: this doesn't halt anything, because go will just de schedule the blocked go routines, but those won't get anything done while stuck here
	// so probably a good idea not to have too much computation in locked code
	defer m.serverMutex.Unlock()
	m.enabled = true
}

func (m *mockServer) Stop() {
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()
	m.enabled = false
}

func (m *mockServer) IsEnabled() bool {
	return m.enabled
}

func (m *mockServer) GetClient() core.HttpClient {
	return m.httpClient
}

func (m *mockServer) AddMock(mock Mock) {
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()
	key := m.getMockKey(mock.Method, mock.Url, mock.RequestBody)
	m.mocks[key] = &mock
}

func (m *mockServer) RemoveMocks() {
	m.mocks = make(map[string]*Mock)
}

func (m *mockServer) getMockKey(method string, url string, body string) string {
	// md5 is faster, not as secure, but in this case the security is not as important
	// Curious does map hash the key?
	hasher := md5.New()
	hasher.Write([]byte(method + url + m.cleanBody(body)))
	key := hex.EncodeToString(hasher.Sum(nil))
	return key
}

func (m *mockServer) cleanBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}
	body = strings.ReplaceAll(body, "\t", "")
	body = strings.ReplaceAll(body, "\n", "")
	return body
}

func (m *mockServer) GetMock(method string, url string, body string) *Mock {
	if !m.enabled {
		return nil
	}

	return m.mocks[m.getMockKey(method, url, body)]
}
