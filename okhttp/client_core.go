package okhttp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/okeefem2/ok-httpclient/core"
	"github.com/okeefem2/ok-httpclient/okhttp_mock"
)

const (
	defaultMaxIdleConnections = 5
	defaultResponseTimeout    = 5 * time.Second
	defaultConnectionTimeout  = 1 * time.Second
)

// Private method on httpClient struct
func (c *httpClient) do(method string, url string, headers http.Header, body interface{}) (*core.Response, error) {

	fullHeaders := c.getRequestHeaders(headers)

	requestBody, err := c.getRequestBody(fullHeaders.Get("Content-Type"), body)

	if err != nil {
		return nil, err
	}

	if mock := okhttp_mock.GetMock(method, url, string(requestBody)); mock != nil {
		return mock.GetResponse()
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, errors.New("unable to create a new request")
	}

	request.Header = fullHeaders

	response, err := c.getHttpClient().Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	finalResponse := core.Response{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Headers:    response.Header,
		Body:       bytes,
	}
	return &finalResponse, nil

}

func (c *httpClient) getHttpClient() *http.Client {

	if c.client != nil {
		return c.client
	}

	// This function will be run 1 time even in concurrent environments
	c.clientOnce.Do(func() {
		// Allow for client override
		if c.builder.client != nil {
			c.client = c.builder.client
			return
		}
		c.client = &http.Client{
			Timeout: c.getConnectionTimeout() + c.getResponseTimeout(),
			Transport: &http.Transport{ // Transport is a RounTripper, which is an interface, so this needs to be a pointer... not sure on that yet but I'll do it
				MaxIdleConnsPerHost:   c.getMaxIdleConnections(), // Should be configured based on the traffic pattern of the app
				ResponseHeaderTimeout: c.getResponseTimeout(),    // How long to wait for a response
				DialContext: (&net.Dialer{
					Timeout: c.getConnectionTimeout(), // Max time to wait for a connection
				}).DialContext, // got error "cannot call pointer method...", fixed by wrapping reference https://stackoverflow.com/questions/44543374/cannot-take-the-address-of-and-cannot-call-pointer-method-on
			},
		}
	})

	return c.client
}

func (c *httpClient) getResponseTimeout() time.Duration {
	if c.builder.disableTimeouts {
		return 0
	}
	if c.builder.responseTimeout > 0 {
		return c.builder.responseTimeout
	}
	return defaultResponseTimeout // Use adefault value of 5
}

func (c *httpClient) getConnectionTimeout() time.Duration {
	if c.builder.disableTimeouts {
		return 0
	}
	if c.builder.connectionTimeout > 0 {
		return c.builder.connectionTimeout
	}
	return defaultConnectionTimeout // Use adefault value of 5
}

func (c *httpClient) getMaxIdleConnections() int {
	if c.builder.maxIdleConnections > 0 {
		return c.builder.maxIdleConnections
	}
	return defaultMaxIdleConnections // Use adefault value of 5
}

func (c *httpClient) getRequestBody(contentType string, body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	switch strings.ToLower(contentType) {
	case "application/json":
		return json.Marshal(body)
	case "application/xml":
		return xml.Marshal(body)
	default:
		return json.Marshal(body)
	}
}

func (c *httpClient) getRequestHeaders(customHeaders http.Header) http.Header {
	result := make(http.Header)
	// Add common headers
	if c.builder != nil {
		for header, value := range c.builder.headers {
			if len(value) > 0 {
				result.Set(header, value[0])
			}
		}
	}

	// Add custom headers
	for header, value := range customHeaders {
		if len(value) > 0 {
			result.Set(header, value[0])
		}
	}

	if c.builder.userAgent != "" && result.Get("User-Agent") != "" {
		result.Set("User-Agent", c.builder.userAgent)
	}

	return result
}
