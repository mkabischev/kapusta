package kapusta

import (
	_ "goblin/test"
	"net/http"
	"testing"

	. "gopkg.in/check.v1"
)

type TestSuite struct {
	dummyClient *dummy
}

func Test(t *testing.T) { TestingT(t) }

type dummy struct{}

func (s *TestSuite) SetUpSuite(c *C) {
	s.dummyClient = &dummy{}
}

func (d *dummy) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{Request: r}, nil
}

func (s *TestSuite) send(r *http.Request, decorator DecoratorFunc) (*http.Response, error) {
	return Decorate(s.dummyClient, decorator).Do(r)
}
