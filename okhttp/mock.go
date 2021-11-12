package okhttp

type Mock struct {
	Method string
	Url string
	RequestBody string


	Error error
	ResponseBody string
	ResponseStatusCode int
}

func (m *Mock) GetResponse() (*Response, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	return &Response{
		statusCode: m.ResponseStatusCode,
		body: []byte(m.ResponseBody),
	}, nil
}
