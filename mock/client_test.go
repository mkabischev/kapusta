package mock

import (
	"context"
	"strings"

	. "gopkg.in/check.v1"
)

type TestClientSuite struct {
	TestSuite
}

var _ = Suite(&TestClientSuite{})

func (s *TestClientSuite) TestDo(c *C) {
	client := NewClient()
	client.Get("/path").WillReturn(200, "response")

	res, err := client.Do(context.Background(), s.newRequest("GET", "/path", nil))

	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(readerToString(res.Body), Equals, "response")
}

func (s *TestClientSuite) TestManyConstraintsOnSamePath(c *C) {
	client := NewClient()
	client.Get("/path").WillReturn(200, "response")
	client.Get("/path").WillReturn(200, "response_2")

	res, err := client.Do(context.Background(), s.newRequest("GET", "/path", nil))

	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(readerToString(res.Body), Equals, "response")
}

func (s *TestClientSuite) TestManyMethodsOnSamePath(c *C) {
	client := NewClient()
	client.Get("/path").WillReturn(200, "response")
	client.Post("/path").WillReturn(200, "response_2")

	res, err := client.Do(context.Background(), s.newRequest("POST", "/path", nil))

	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(readerToString(res.Body), Equals, "response_2")
}

func (s *TestClientSuite) TestManyBodyOnSamePath(c *C) {
	client := NewClient()
	client.Post("/path").WithBody("request_1").WillReturn(200, "response")
	client.Post("/path").WithBody("request_2").WillReturn(200, "response_2")

	res, err := client.Do(context.Background(), s.newRequest("POST", "/path", strings.NewReader("request_2")))

	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(readerToString(res.Body), Equals, "response_2")
}

func (s *TestClientSuite) TestNoSuitableRequestWrongPath(c *C) {
	client := NewClient()
	client.Get("/").WillReturn(200, "response")

	res, err := client.Do(context.Background(), s.newRequest("GET", "/path", nil))
	c.Assert(res, IsNil)
	c.Assert(err, ErrorMatches, "No suitable request found")
}

func (s *TestClientSuite) TestNoSuitableRequestWrongBody(c *C) {
	client := NewClient()
	client.Post("/").WithBody("request").WillReturn(200, "response")

	res, err := client.Do(context.Background(), s.newRequest("POST", "/", strings.NewReader("request_")))
	c.Assert(res, IsNil)
	c.Assert(err, ErrorMatches, "No suitable request found")
}

func (s *TestClientSuite) TestGetWithQueryParams(c *C) {
	client := NewClient()
	client.Get("/path?param=value").WillReturn(200, "response")

	res, err := client.Do(context.Background(), s.newRequest("GET", "/path?param=value", nil))

	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(readerToString(res.Body), Equals, "response")
}

func (s *TestClientSuite) TestRequest(c *C) {
	client := NewClient()
	client.Request().WillReturn(200, "response")

	res, err := client.Do(context.Background(), s.newRequest("GET", "/path", nil))

	c.Assert(err, IsNil)
	c.Assert(res.StatusCode, Equals, 200)
	c.Assert(readerToString(res.Body), Equals, "response")
}
