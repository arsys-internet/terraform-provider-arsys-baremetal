package client

import "net/http"

type TokenTransport struct {
	Token string
	Base  http.RoundTripper
}

func (t *TokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-TOKEN", t.Token)
	return t.Base.RoundTrip(req)
}
