package requestcontext

import (
	"github.com/juju/errgo"
)

var (
	NotFoundError = errgo.New("not found")
	maskAny       = errgo.MaskFunc(errgo.Any)
)

func IsNotFound(err error) bool {
	return errgo.Cause(err) == NotFoundError
}
