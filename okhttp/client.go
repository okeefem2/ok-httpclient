package okhttp

import (
	"net/http"
	"sync"

	"github.com/okeefem2/ok-httpclient/core"
)

type Client interface {
	// ... is a variatic and means there can be 0+ instances of the hearder type. Makes the argument optional
	// Get(url string, headers ...http.Header) (*core.Response, error) but then you just choose only the first one... seems odd haha
	Get(url string, headers ...http.Header) (*core.Response, error)
	Put(url string, body interface{}, headers ...http.Header) (*core.Response, error)
	Post(url string, body interface{}, headers ...http.Header) (*core.Response, error)
	Patch(url string, body interface{}, headers ...http.Header) (*core.Response, error)
	Delete(url string, headers ...http.Header) (*core.Response, error)
	Options(url string, headers ...http.Header) (*core.Response, error)
	Head(url string, headers ...http.Header) (*core.Response, error)
}

// Via duck typing, this struct implementes the HttpClient interface
// This is private because it is camelcase
type httpClient struct {
	builder *clientBuilder

	client     *http.Client
	clientOnce sync.Once
}

// Add receiver type to assign the function as a method to that given receiver type
func (c *httpClient) Get(url string, headers ...http.Header) (*core.Response, error) {
	return c.do(http.MethodGet, url, getHeaders(headers...), nil)
}

func (c *httpClient) Post(url string, body interface{}, headers ...http.Header) (*core.Response, error) {
	return c.do(http.MethodPost, url, getHeaders(headers...), body)
}

func (c *httpClient) Put(url string, body interface{}, headers ...http.Header) (*core.Response, error) {
	return c.do(http.MethodPut, url, getHeaders(headers...), body)
}

func (c *httpClient) Patch(url string, body interface{}, headers ...http.Header) (*core.Response, error) {
	return c.do(http.MethodPatch, url, getHeaders(headers...), body)
}

func (c *httpClient) Delete(url string, headers ...http.Header) (*core.Response, error) {
	return c.do(http.MethodDelete, url, getHeaders(headers...), nil)
}

func (c *httpClient) Options(url string, headers ...http.Header) (*core.Response, error) {
	return c.do(http.MethodOptions, url, getHeaders(headers...), nil)
}

func (c *httpClient) Head(url string, headers ...http.Header) (*core.Response, error) {
	return c.do(http.MethodHead, url, getHeaders(headers...), nil)
}
