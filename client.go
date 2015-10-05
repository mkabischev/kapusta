package kapusta

import "net/http"

// IClient sends http.Requests and returns http.Responses or errors in  case of failure.
type IClient interface {
	Do(*http.Request) (*http.Response, error)
}

// ClientFunc is a function type that implements the Client interface.
type ClientFunc func(*http.Request) (*http.Response, error)

// Do Client interface support
func (f ClientFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}
