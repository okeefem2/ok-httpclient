package example

import (
	"net/http"
	"time"

	"github.com/okeefem2/ok-httpclient/okhttp"
)

var (
	httpClient = getHttpClient()
)

func getHttpClient() okhttp.Client {
	clientBuilder :=  okhttp.NewBuilder()

	// Using make because http Header is an alias to map
	commonHeaders := make(http.Header)
	// commonHeaders.Set("Authorization", "Bearer ABC-123")
	return clientBuilder.SetHeaders(commonHeaders).
			SetConnectionTimeout(1 * time.Second).
			SetResponseTimeout(4 * time.Second).
			Build()
}
