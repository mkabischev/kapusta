package kapusta

import (
	"fmt"
	"net/http"
	"net/url"
)

// HeaderMiddleware returns a DecoratorFunc that adds the given HTTP header to every request done by a Client.
func HeaderMiddleware(name, value string) MiddlewareFunc {
	return HeadersMiddleware(map[string]string{name: value})
}

// HeadersMiddleware returns a DecoratorFunc that adds the given HTTP headers to every request done by a Client.
func HeadersMiddleware(values map[string]string) MiddlewareFunc {
	return func(c Client) Client {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			for key, value := range values {
				r.Header.Add(key, value)
			}
			return c.Do(r)
		})
	}
}

// RecoverMiddleware returns a DecoratorFunc that recovers panic and convert it to error
func RecoverMiddleware() MiddlewareFunc {
	return func(c Client) Client {
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

// BaseURLMiddleware returns a DecoratorFunc that replaces scheme and host in request from specified baseURL
func BaseURLMiddleware(baseURL string) MiddlewareFunc {
	parsed, err := url.Parse(baseURL)

	if err != nil {
		panic(err)
	}

	return func(c Client) Client {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			r.URL.Scheme = parsed.Scheme
			r.URL.Host = parsed.Host

			return c.Do(r)
		})
	}
}
