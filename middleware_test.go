package kapusta

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dummy struct{}

func (d dummy) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{Request: r}, nil
}

type failableClient struct {
	fails          int // number of first requests which will fail
	calls          int // number of requests processed
	responseOnFail *http.Response
	errorOnFail    error
}

func newFailableClient(fails int, r *http.Response, err error) Client {
	return &failableClient{
		fails:          fails,
		responseOnFail: r,
		errorOnFail:    err,
	}
}

func (c *failableClient) Do(r *http.Request) (*http.Response, error) {
	c.calls++

	if c.calls < c.fails {
		return c.responseOnFail, c.errorOnFail
	}

	return &http.Response{Request: r, StatusCode: 200}, nil
}

type callOrder []string

func send(r *http.Request, middleware MiddlewareFunc) (*http.Response, error) {
	return Chain(&dummy{}, middleware).Do(r)
}

func createMiddleware(name string, order *callOrder) MiddlewareFunc {
	return func(c Client) Client {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			*order = append(*order, name)
			return c.Do(r)
		})
	}
}

func TestChain(t *testing.T) {
	var order callOrder
	client := Chain(&dummy{}, createMiddleware("A", &order), createMiddleware("B", &order), createMiddleware("C", &order))

	client.Do(&http.Request{})

	assert.Equal(t, callOrder{"C", "B", "A"}, order)
}

func TestHeaderMiddleware(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := send(r, HeaderMiddleware("X-Foo", "123"))

	assert.Equal(t, res.Request.Header.Get("X-Foo"), "123")
}

func TestHeadersMiddleware(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := send(r, HeadersMiddleware(map[string]string{"X-Foo": "123", "X-Bar": "456"}))

	assert.Equal(t, res.Request.Header.Get("X-Foo"), "123")
	assert.Equal(t, res.Request.Header.Get("X-Bar"), "456")
}

func TestBaseURLMiddleware(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := send(r, BaseURLMiddleware("http://example.com"))

	assert.Equal(t, res.Request.URL.String(), "http://example.com/")
}

func TestBaseURLMiddlewareFail(t *testing.T) {
	assert.Panics(t, func() {
		BaseURLMiddleware("://example.com")
	})
}

func TestPanicMiddleware(t *testing.T) {
	panicTriggerDecorator := func(c Client) Client {
		return ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			panic("oops")
		})
	}
	r, _ := http.NewRequest("GET", "/", nil)
	client := Chain(&dummy{}, panicTriggerDecorator, RecoverMiddleware())
	res, err := client.Do(r)

	assert.Nil(t, res)
	assert.Error(t, err, "oops")
}

func TestBackoffMiddleware(t *testing.T) {
	c := newFailableClient(5, &http.Response{StatusCode: 500}, nil)
	c = Chain(c, BackoffMiddleware(10, time.Second, RequestStatusIn(200, 400)))

	r, _ := http.NewRequest("GET", "/", nil)
	_, err := c.Do(r)

	assert.Nil(t, err)

}
