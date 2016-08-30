package kapusta

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummy struct{}

func (d dummy) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{Request: r}, nil
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
