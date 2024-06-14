package provider

import "net/http"

type authedTransport struct {
	endpoint string
	wrapped  http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	return t.wrapped.RoundTrip(req)
}
