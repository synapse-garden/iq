package main_test

import (
	"testing"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) { gc.TestingT(t) }

type MainSuite struct{}

var _ = gc.Suite(&MainSuite{})

func (s *MainSuite) TestMain(c *gc.C) {
	c.Assert(true, gc.Equals, true)
}
