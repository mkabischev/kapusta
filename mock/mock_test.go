package mock

import (
	"bytes"
	"io"
	"net/http"
)

type dummyClient struct{}

func (c *dummyClient) Do(r *http.Request) (*http.Response, error) {
	return nil, nil
}

func writeHeaders(headers http.Header) string {
	buf := &bytes.Buffer{}
	headers.Write(buf)

	return buf.String()
}

func builder() *PromiseBuilder {
	return NewPromiseBuilder(NewClient())
}

func newRequest(method, path string, body io.Reader) *http.Request {
	r, _ := http.NewRequest(method, path, body)

	return r
}
