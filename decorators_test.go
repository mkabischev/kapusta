package kapusta

import (
	"net/http"

	. "gopkg.in/check.v1"
)

type TestDecoratorSuite struct {
	TestSuite
}

type callOrder []string

var (
	_ = Suite(&TestDecoratorSuite{})
)

func createDecorator(name string, order *callOrder) DecoratorFunc {
	return func(c IClient) IClient {
		return ClientFunc(func(r *http.Request) (*http.Response, error) {
			*order = append(*order, name)
			return c.Do(r)
		})
	}
}

type SimpleDecorator struct {
	name      string
	callOrder []string
}

func (s *TestDecoratorSuite) TestDecorate(c *C) {
	var order callOrder
	client := Decorate(s.dummyClient, createDecorator("A", &order), createDecorator("B", &order), createDecorator("C", &order))

	client.Do(&http.Request{})

	c.Assert(callOrder{"C", "B", "A"}, DeepEquals, order)
}

func (s *TestDecoratorSuite) TestHeaderDecorator(c *C) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := s.send(r, HeaderDecorator("X-Foo", "123"))

	c.Assert(res.Request.Header.Get("X-Foo"), Equals, "123")
}

func (s *TestDecoratorSuite) TestHeadersDecorator(c *C) {
	r, _ := http.NewRequest("GET", "/", nil)
	res, _ := s.send(r, HeadersDecorator(map[string]string{"X-Foo": "123", "X-Bar": "456"}))

	c.Assert(res.Request.Header.Get("X-Foo"), Equals, "123")
	c.Assert(res.Request.Header.Get("X-Bar"), Equals, "456")
}

func (s *TestDecoratorSuite) TestPanicDecorator(c *C) {
	panicTriggerDecorator := func(c IClient) IClient {
		return ClientFunc(func(r *http.Request) (res *http.Response, err error) {
			panic("oops")
		})
	}
	r, _ := http.NewRequest("GET", "/", nil)
	client := Decorate(s.dummyClient, panicTriggerDecorator, RecoverDecorator())
	res, err := client.Do(r)

	c.Assert(res, IsNil)
	c.Assert(err, ErrorMatches, "*oops")
}
