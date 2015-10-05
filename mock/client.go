package mock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client implementation of kapusta.Client for testing purpose
type Client struct {
	promises []Promise
}

// NewClient creates new client
func NewClient() *Client {
	return &Client{}
}

// Do comparing request with all constrains and returns first matching
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	requestBody := readerToString(r.Body)
	for _, promise := range c.promises {
		// Request copy is required because body can be read only once
		requestCopy := new(http.Request)
		*requestCopy = *r
		requestCopy.Body = closer{strings.NewReader(requestBody)}
		if response, ok := promise.Check(requestCopy); ok {
			return response, nil
		}
	}
	return nil, fmt.Errorf("No suitable request found")
}

// Get helper method for using RuleBuilder
func (c *Client) Get(path string) *PromiseBuilder {
	return NewPromiseBuilder(c).Get(path)
}

// Post helpep method for using RuleBuilder
func (c *Client) Post(path string) *PromiseBuilder {
	return NewPromiseBuilder(c).Post(path)
}

func readerToString(r io.Reader) string {
	if r == nil {
		return ""
	}
	buf := &bytes.Buffer{}
	buf.ReadFrom(r)

	return buf.String()
}
