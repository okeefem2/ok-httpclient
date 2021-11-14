package okhttp_mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpClientMock struct {

}

func (c *httpClientMock) Do(request *http.Request) (*http.Response, error) {
		requestBody, err := request.GetBody()
		if err != nil {
			return nil, err
		}

		defer requestBody.Close()

		body, err := ioutil.ReadAll(requestBody)

		if err != nil {
			return nil, err
		}

		if mock := GetMock(request.Method, request.URL.String(), string(body)); mock != nil {

			if mock.Error != nil {
				return nil, mock.Error
			}
			response := http.Response{
				StatusCode: mock.ResponseStatusCode,
				Body: ioutil.NopCloser(strings.NewReader(mock.ResponseBody)),
				ContentLength: int64(len(mock.ResponseBody)),
			}
			return &response, nil
		}

		return nil, fmt.Errorf("no mock matching %s from %s with the given body", request.Method, request.URL)
}
