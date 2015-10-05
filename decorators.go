package kapusta

import (
	"fmt"
	"net/http"
	"net/url"
)

// DecoratorFunc wraps a Client with extra behaviour.
type DecoratorFunc func(IClient) IClient

// Decorate decorates a Client c with all the given Decorators, in order.
func Decorate(c IClient, ds ...DecoratorFunc) IClient {
	result := c
	for _, decorate := range ds {
		result = decorate(result)
	}
	return result
}

// HeaderDecorator returns a DecoratorFunc that adds the given HTTP header to every request done by a Client.
func HeaderDecorator(name, value string) DecoratorFunc {
	return HeadersDecorator(map[string]string{name: value})
}

// HeadersDecorator returns a DecoratorFunc that adds the given HTTP headers to every request done by a Client.
func HeadersDecorator(values map[string]string) DecoratorFunc {
	return func(c IClient) IClient {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			for key, value := range values {
				r.Header.Add(key, value)
			}
			return c.Do(r)
		})
	}
}

// RecoverDecorator returns a DecoratorFunc that recovers panic and convert it to error
func RecoverDecorator() DecoratorFunc {
	return func(c IClient) IClient {
		return ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("recovered panic: %v", r)
				}
			}()
			return c.Do(r)
		})
	}
}

// BaseURLDecorator returns a DecoratorFunc that replaces scheme and host in request from specified baseURL
func BaseURLDecorator(baseURL string) DecoratorFunc {
	parsed, err := url.Parse(baseURL)

	if err != nil {
		panic(err)
	}

	return func(c IClient) IClient {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			r.URL.Scheme = parsed.Scheme
			r.URL.Host = parsed.Host

			return c.Do(r)
		})
	}
}
