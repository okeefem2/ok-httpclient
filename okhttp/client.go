package okhttp

import (
	"net/http"
	"sync"
)

type Client interface {
	Get(url string, headers http.Header) (*Response, error)
	Put(url string, headers http.Header, body interface{}) (*Response, error)
	Post(url string, headers http.Header, body interface{}) (*Response, error)
	Patch(url string, headers http.Header, body interface{}) (*Response, error)
	Delete(url string, headers http.Header) (*Response, error)
}

// Via duck typing, this struct implementes the HttpClient interface
// This is private because it is camelcase
type httpClient struct{
	builder *clientBuilder

	client *http.Client
	clientOnce sync.Once
}

// Add receiver type to assign the function as a method to that given receiver type
func (c *httpClient)Get(url string, headers http.Header)  (*Response, error) {
	return c.do(http.MethodGet, url, headers, nil)
}

func (c *httpClient)Post(url string, headers http.Header, body interface{})  (*Response, error) {
	return c.do(http.MethodGet, url, headers, body)
}

func (c *httpClient)Put(url string, headers http.Header, body interface{})  (*Response, error) {
	return c.do(http.MethodGet, url, headers, body)
}

func (c *httpClient)Patch(url string, headers http.Header, body interface{})  (*Response, error) {
	return c.do(http.MethodGet, url, headers, body)
}

func (c *httpClient)Delete(url string, headers http.Header)  (*Response, error) {
	return c.do(http.MethodGet, url, headers, nil)
}
