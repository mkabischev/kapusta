package mock

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	client := NewClient()
	client.Get("/path").WillReturn(200, "response")

	res, err := client.Do(newRequest("GET", "/path", nil))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, readerToString(res.Body), "response")
}

func TestManyConstraintsOnSamePath(t *testing.T) {
	client := NewClient()
	client.Get("/path").WillReturn(200, "response")
	client.Get("/path").WillReturn(200, "response_2")

	res, err := client.Do(newRequest("GET", "/path", nil))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, readerToString(res.Body), "response")
}

func TestManyMethodsOnSamePath(t *testing.T) {
	client := NewClient()
	client.Get("/path").WillReturn(200, "response")
	client.Post("/path").WillReturn(200, "response_2")

	res, err := client.Do(newRequest("POST", "/path", nil))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, readerToString(res.Body), "response_2")
}

func TestManyBodyOnSamePath(t *testing.T) {
	client := NewClient()
	client.Post("/path").WithBody("request_1").WillReturn(200, "response")
	client.Post("/path").WithBody("request_2").WillReturn(200, "response_2")

	res, err := client.Do(newRequest("POST", "/path", strings.NewReader("request_2")))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, readerToString(res.Body), "response_2")
}

func TestNoSuitableRequestWrongPath(t *testing.T) {
	client := NewClient()
	client.Get("/").WillReturn(200, "response")

	res, err := client.Do(newRequest("GET", "/path", nil))

	assert.Nil(t, res)
	assert.Error(t, err, "No suitable request found")
}

func TestNoSuitableRequestWrongBody(t *testing.T) {
	client := NewClient()
	client.Post("/").WithBody("request").WillReturn(200, "response")

	res, err := client.Do(newRequest("POST", "/", strings.NewReader("request_")))

	assert.Nil(t, res)
	assert.Error(t, err, "No suitable request found")
}

func TestGetWithQueryParams(t *testing.T) {
	client := NewClient()
	client.Get("/path?param=value").WillReturn(200, "response")

	res, err := client.Do(newRequest("GET", "/path?param=value", nil))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, readerToString(res.Body), "response")
}

func TestRequest(t *testing.T) {
	client := NewClient()
	client.Request().WillReturn(200, "response")

	res, err := client.Do(newRequest("GET", "/path", nil))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, readerToString(res.Body), "response")
}
