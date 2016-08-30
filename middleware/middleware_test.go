package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/mkabischev/kapusta"
	. "gopkg.in/check.v1"
)

type TestSuite struct {
	dummyClient *dummy
}

type dummy struct{}

func (d *dummy) Do(ctx context.Context, r *http.Request) (*http.Response, error) {
	return &http.Response{Request: r}, nil
}

type callOrder []string

var (
	_ = Suite(&TestSuite{})
)

func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C) {
	s.dummyClient = &dummy{}
}

func (s *TestSuite) send(ctx context.Context, r *http.Request, middleware kapusta.MiddlewareFunc) (*http.Response, error) {
	return kapusta.Chain(s.dummyClient, middleware).Do(ctx, r)
}

func createMiddleware(name string, order *callOrder) kapusta.MiddlewareFunc {
	return func(c kapusta.Client) kapusta.Client {
		return kapusta.ClientFunc(func(ctx context.Context, r *http.Request) (*http.Response, error) {
			*order = append(*order, name)
			return c.Do(ctx, r)
		})
	}
}

func (s *TestSuite) TestChain(c *C) {
	var order callOrder
	client := kapusta.Chain(s.dummyClient, createMiddleware("A", &order), createMiddleware("B", &order), createMiddleware("C", &order))

	client.Do(context.Background(), &http.Request{})

	c.Assert(callOrder{"C", "B", "A"}, DeepEquals, order)
}

func (s *TestSuite) TestHeaderMiddleware(c *C) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := s.send(context.Background(), r, HeaderMiddleware("X-Foo", "123"))

	c.Assert(res.Request.Header.Get("X-Foo"), Equals, "123")
}

func (s *TestSuite) TestHeadersMiddleware(c *C) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := s.send(context.Background(), r, HeadersMiddleware(map[string]string{"X-Foo": "123", "X-Bar": "456"}))

	c.Assert(res.Request.Header.Get("X-Foo"), Equals, "123")
	c.Assert(res.Request.Header.Get("X-Bar"), Equals, "456")
}

func (s *TestSuite) TestPanicMiddleware(c *C) {
	panicTriggerDecorator := func(c kapusta.Client) kapusta.Client {
		return kapusta.ClientFunc(func(ctx context.Context, r *http.Request) (res *http.Response, err error) {
			panic("oops")
		})
	}
	r, _ := http.NewRequest("GET", "/", nil)
	client := kapusta.Chain(s.dummyClient, panicTriggerDecorator, RecoverMiddleware())
	res, err := client.Do(context.Background(), r)

	c.Assert(res, IsNil)
	c.Assert(err, ErrorMatches, "*oops")
}
