package thorchain

import (
	"errors"
	"fmt"
	"strings"

	se "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/hashicorp/go-multierror"
	. "gopkg.in/check.v1"
)

type ErrorsTestSuite struct{}

var _ = Suite(&ErrorsTestSuite{})

func (ErrorsTestSuite) TestErrInternal(c *C) {
	codeErr := errBadVersion
	_, code, log := se.ABCIInfo(codeErr, false)
	c.Check(code, Equals, uint32(101))
	c.Check(strings.Contains(log, "bad version"), Equals, true)

	codelessErr := fmt.Errorf("codeless error")
	_, code, log = se.ABCIInfo(codelessErr, false)
	c.Check(int(code), Equals, 1)
	c.Check(strings.Contains(log, "codeless error"), Equals, false) // Redacted error.
	c.Check(log, Equals, "internal")

	internalErr := ErrInternal(codeErr, codelessErr.Error())
	_, code, log = se.ABCIInfo(internalErr, false)
	c.Check(int(code), Equals, 1)
	c.Check(strings.Contains(log, "codeless error"), Equals, false) // Redacted error.
	c.Check(log, Equals, "internal")

	appendedError := multierror.Append(codeErr, codelessErr)
	_, code, log = se.ABCIInfo(appendedError, false)
	c.Check(int(code), Equals, 1)
	c.Check(strings.Contains(log, "codeless error"), Equals, false) // Redacted error.
	c.Check(log, Equals, "internal")

	joinedError := errors.Join(codeErr, codelessErr)
	_, code, log = se.ABCIInfo(joinedError, false)
	c.Check(int(code), Equals, 1)
	c.Check(strings.Contains(log, "codeless error"), Equals, false) // Redacted error.
	c.Check(log, Equals, "internal")

	wrappedErr := se.Wrap(codeErr, codelessErr.Error())
	_, code, log = se.ABCIInfo(wrappedErr, false)
	c.Assert(int(code), Equals, 101) // codeErr's code preserved.
	c.Check(strings.Contains(log, "bad version"), Equals, true)
	c.Check(strings.Contains(log, "codeless error"), Equals, true)
	// Both errors' contents are preserved.
}
