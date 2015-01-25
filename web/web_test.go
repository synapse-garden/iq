package web_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/synapse-garden/iq/web"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) { gc.TestingT(t) }

type WebSuite struct {
	iqr *web.IqRunner
}

var _ = gc.Suite(&WebSuite{})

func (s *WebSuite) SetUpTest(c *gc.C) {
	s.iqr = web.CreateRunner(25000, nil)
	s.iqr.StartRun()

	go func() {
		err := <-s.iqr.Errors()
		c.Assert(err, gc.IsNil)
	}()
}

func (s *WebSuite) TearDownTest(c *gc.C) {
	s.iqr.Kill()

	err, ok := <-s.iqr.Errors()
	if !ok {
		// Channel was already closed
		c.Log("channel was already closed")
	}
	c.Assert(err, gc.IsNil)
}

func (s *WebSuite) TestStartRun(c *gc.C) {
	baseUri := "http://localhost:25000"
	for i, t := range []struct {
		should       string
		path         string
		params       map[string]string
		expectedCode int
		expectedResp string
		expectedErr  string
	}{{
		should:       "pass",
		path:         "/default",
		expectedCode: 200,
		expectedResp: "Hello default!",
	}, {
		should: "fail",
		path:   "/failing",
		params: map[string]string{
			"unknown": "param",
		},
		expectedCode: 404,
		expectedErr:  "page not found",
	}} {
		c.Logf("test %d: should %s\n", i, t.should)

		resp := httptest.NewRecorder()

		uri := baseUri + t.path

		params := make(url.Values)
		for k, v := range t.params {
			params.Add(k, v)
		}

		req, err := http.NewRequest("GET", uri+params.Encode(), nil)
		if err != nil {
			c.Check(err, gc.ErrorMatches, t.expectedErr)
			continue
		}

		http.DefaultServeMux.ServeHTTP(resp, req)

		c.Logf("%#v\n", resp)

		srBody := resp.Body.String()
		c.Logf("%#v\n", srBody)

		c.Check(resp.Code, gc.Equals, t.expectedCode)

		if t.expectedErr != "" {
			c.Check(srBody, gc.Matches, strconv.Itoa(t.expectedCode)+" "+t.expectedErr+"[\\s]*")
			continue
		}

		c.Check(srBody, gc.Matches, ".*"+t.expectedResp+".*")
		c.Check(srBody, gc.Not(gc.Matches), ".*Error.*")
	}
}
