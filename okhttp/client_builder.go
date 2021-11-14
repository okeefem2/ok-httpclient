package okhttp

import (
	"net/http"
	"time"

	"github.com/okeefem2/ok-httpclient/core"
)

type ClientBuilder interface {
	SetHeaders(headers http.Header) ClientBuilder
	SetConnectionTimeout(timeout time.Duration) ClientBuilder
	SetResponseTimeout(timeout time.Duration) ClientBuilder
	SetMaxIdleConnections(max int) ClientBuilder
	DisableTimeouts(disable bool) ClientBuilder
	SetHttpClient(c core.HttpClient) ClientBuilder
	SetUserAgent(agent string) ClientBuilder

	Build() Client
}

type clientBuilder struct {
	maxIdleConnections int
	connectionTimeout  time.Duration
	responseTimeout    time.Duration
	disableTimeouts    bool
	headers            http.Header
	userAgent          string
	// baseUrl string
	client core.HttpClient
}

func (c *clientBuilder) SetHeaders(headers http.Header) ClientBuilder {
	c.headers = headers
	return c
}

func (c *clientBuilder) SetConnectionTimeout(timeout time.Duration) ClientBuilder {
	c.connectionTimeout = timeout
	return c
}

func (c *clientBuilder) SetResponseTimeout(timeout time.Duration) ClientBuilder {
	c.responseTimeout = timeout
	return c
}

func (c *clientBuilder) SetMaxIdleConnections(max int) ClientBuilder {
	c.maxIdleConnections = max
	return c
}

func (c *clientBuilder) DisableTimeouts(disable bool) ClientBuilder {
	c.disableTimeouts = disable
	return c
}

func (c *clientBuilder) SetHttpClient(client core.HttpClient) ClientBuilder {
	c.client = client
	return c
}

func (c *clientBuilder) SetUserAgent(agent string) ClientBuilder {
	c.userAgent = agent
	return c
}

// Titlecase = exported and public, this is our public API
func NewBuilder() ClientBuilder {
	builder := &clientBuilder{}
	return builder
}

func (c *clientBuilder) Build() Client {
	client := httpClient{
		builder: c,
	}

	return &client
}
