package kapusta

import (
	"context"
	"net/http"
)

// IClient sends http.Requests and returns http.Responses or errors in  case of failure.
type Client interface {
	Do(ctx context.Context, r *http.Request) (*http.Response, error)
}

// ClientFunc is a function type that implements the Client interface.
type ClientFunc func(ctx context.Context, r *http.Request) (*http.Response, error)

// Do Client interface support
func (f ClientFunc) Do(ctx context.Context, r *http.Request) (*http.Response, error) {
	return f(ctx, r)
}

// MiddlewareFunc wraps a Client with extra behaviour.
type MiddlewareFunc func(Client) Client

// Decorate decorates a Client c with all the given Decorators, in order.
func Chain(c Client, ds ...MiddlewareFunc) Client {
	result := c
	for _, middleware := range ds {
		result = middleware(result)
	}
	return result
}
