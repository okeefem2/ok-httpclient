package okhttp_mock

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"sync"

	"github.com/okeefem2/ok-httpclient/core"
)

// Singleton mock mocks
var (
	mocks = mockServer{
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

func IsMockServerEnabled() bool {
	return mocks.enabled
}

func GetMockedClient() core.HttpClient {
	return mocks.httpClient
}

func AddMock(mock Mock) {
	mocks.serverMutex.Lock()
	defer mocks.serverMutex.Unlock()
	key := mocks.getMockKey(mock.Method, mock.Url, mock.RequestBody)
	mocks.mocks[key] = &mock
}

func RemoveMocks() {
	mocks.mocks = make(map[string]*Mock)
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

func GetMock(method string, url string, body string) *Mock {
	if !mocks.enabled {
		return nil
	}

	return mocks.mocks[mocks.getMockKey(method, url, body)]
}
