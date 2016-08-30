package mock

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConstraintBuilder(t *testing.T) {
	builder := builder()

	assert.Len(t, builder.checkers, 0)
}

func TestWithMethod(t *testing.T) {
	builder := builder().WithMethod("PUT")

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("PUT", "/", nil),
		newRequest("PUT", "/path/", nil),
		newRequest("PUT", "/path/", strings.NewReader("something")),
	}, true)

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("GET", "/", nil),
		newRequest("POST", "/", nil),
	}, false)
}

func TestGet(t *testing.T) {
	builder := builder().Get("/some/path")

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("GET", "/some/path", nil),
		newRequest("GET", "/some/path", strings.NewReader("something")),
	}, true)

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("POST", "/some/path", nil),
		newRequest("GET", "/some/path/", nil),
	}, false)
}

func TestPost(t *testing.T) {
	builder := builder().Post("/some/path")

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("POST", "/some/path", nil),
		newRequest("POST", "/some/path", strings.NewReader("something")),
	}, true)

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("GET", "/some/path", nil),
		newRequest("GET", "/some/path/", nil),
	}, false)
}

func TestWithPath(t *testing.T) {
	builder := builder().WithPath("/some/path")

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("POST", "/some/path", nil),
		newRequest("GET", "/some/path", nil),
	}, true)

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("GET", "/some/path/", nil),
		newRequest("GET", "/", nil),
	}, false)
}

func TestWithHeader(t *testing.T) {
	builder := builder().WithHeader("X-Foo", "123")

	r1 := newRequest("POST", "/some/path", nil)
	r1.Header.Add("X-Foo", "123")

	r2 := newRequest("POST", "/some/path", nil)
	r2.Header.Add("X-Foo", "1234")

	r3 := newRequest("POST", "/some/path", nil)

	assertCheckers(t, builder.checkers, []*http.Request{
		r1,
	}, true)

	assertCheckers(t, builder.checkers, []*http.Request{
		r2,
		r3,
	}, false)
}

func TestWithBody(t *testing.T) {
	builder := builder().WithBodyReader(strings.NewReader("some body"))

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("POST", "/some/path", strings.NewReader("some body")),
		newRequest("PUT", "/some/path", strings.NewReader("some body")),
	}, true)

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("POST", "/some/path", nil),
		newRequest("PUT", "/some/path", strings.NewReader("wrong body")),
	}, false)
}

func TestBody(t *testing.T) {
	builder := builder().WithBody("some body")

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("POST", "/some/path", strings.NewReader("some body")),
		newRequest("PUT", "/some/path", strings.NewReader("some body")),
	}, true)

	assertCheckers(t, builder.checkers, []*http.Request{
		newRequest("POST", "/some/path", nil),
		newRequest("PUT", "/some/path", strings.NewReader("wrong body")),
	}, false)
}

func TestReturns(t *testing.T) {
	client := NewClient()
	builder := NewPromiseBuilder(client)
	builder.WillReturn(200, "some response")

	assert.Len(t, client.promises, 1)
	assert.Equal(t, client.promises[0].checkers, builder.checkers)

	buf := &bytes.Buffer{}
	buf.ReadFrom(client.promises[0].response.Body)

	assert.Equal(t, client.promises[0].response.StatusCode, 200)
	assert.Equal(t, buf.String(), "some response")
}

func TestReturnsResponse(t *testing.T) {
	response := &http.Response{StatusCode: 200, Body: &closer{strings.NewReader("some response")}}

	client := NewClient()
	builder := NewPromiseBuilder(client)
	builder.WillReturnResponse(response)

	assert.Len(t, client.promises, 1)
	assert.Equal(t, client.promises[0].checkers, builder.checkers)
	assert.Equal(t, client.promises[0].response, response)
}

func TestCloser(t *testing.T) {
	// Fake test. Just for 100% coverage
	closer := &closer{}
	assert.Nil(t, closer.Close())
}

func assertCheckers(t *testing.T, checkers []CheckerFunc, requests []*http.Request, valid bool) {
	rule := &Promise{checkers: checkers}
	for _, req := range requests {
		_, ok := rule.Check(req)
		assert.Equal(t, ok, valid)
	}
}
