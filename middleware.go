package kapusta

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
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

// BadRequestFunc is used to determine was response successful or not.
// For example you can mark response as successful if http status is 200 or 400.
type GoodRequestFn func(r *http.Response, err error) bool

// RequestStatusIn helper for treat requests with specified statuses valid.
func RequestStatusIn(status ...int) GoodRequestFn {
	return func(r *http.Response, err error) bool {
		if err != nil {
			return false
		}

		for _, s := range status {
			if r.StatusCode == s {
				return true
			}
		}

		return false
	}
}

// BackoffMiddleware implements simple
func BackoffMiddleware(retries int, delay time.Duration, goodRequestFn GoodRequestFn) MiddlewareFunc {
	return func(c Client) Client {
		nextDelay := time.Duration(0)
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			for i := 0; i < retries; i++ {
				res, err := c.Do(r)
				if goodRequestFn(res, err) {
					return res, err
				}

				time.Sleep(nextDelay)
				nextDelay += delay * time.Duration(i)
			}

			return nil, fmt.Errorf("request failed after %v retries.", retries)
		})
	}
}
