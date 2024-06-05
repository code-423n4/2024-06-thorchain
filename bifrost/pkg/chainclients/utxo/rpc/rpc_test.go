package rpc

import (
	"errors"
	"testing"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type RPCSuite struct{}

var _ = Suite(&RPCSuite{})

func (s *RPCSuite) TestRetry(c *C) {
	cl := Client{maxRetries: 3}
	called := 0
	err := cl.retry(func() error {
		called++
		return nil
	})
	c.Assert(err, IsNil)
	c.Assert(called, Equals, 1)

	called = 0
	err = cl.retry(func() error {
		called++
		return errors.New("error")
	})
	c.Assert(err, NotNil)
	c.Assert(called, Equals, 1)

	called = 0
	err = cl.retry(func() error {
		called++
		return errors.New("500 Internal Server Error: work queue depth exceeded")
	})
	c.Assert(err, NotNil)
	c.Assert(called, Equals, 4)

	called = 0
	err = cl.retry(func() error {
		called++
		if called < 2 {
			return errors.New("500 Internal Server Error: work queue depth exceeded")
		}
		return nil
	})
	c.Assert(err, IsNil)
	c.Assert(called, Equals, 2)
}
