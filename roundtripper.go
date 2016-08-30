package kapusta

import "net/http"

// RoundTripper is adapter for using Client as transport for http.Client
type RoundTripper struct {
	client Client
}

// RoundTrip implementation of RoundTripper interface
func (t RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.client.Do(r)
}