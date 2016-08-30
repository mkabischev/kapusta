package kapusta

import (
	"net/http"
)

// Client sends http.Requests and returns http.Responses or errors in case of failure.
type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

// ClientFunc is a function type that implements the Client interface.
type ClientFunc func(r *http.Request) (*http.Response, error)

// Do Client interface support
func (f ClientFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

// MiddlewareFunc wraps a Client with extra behaviour.
type MiddlewareFunc func(Client) Client

// Chain decorates a Client c with all the given middlewares, in order.
func Chain(c Client, mw ...MiddlewareFunc) Client {
	result := c
	for _, middleware := range mw {
		result = middleware(result)
	}
	return result
}
