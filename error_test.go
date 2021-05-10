package errors

import (
	"testing"

	util "github.com/geektime007/util/testutil"
)

func TestErrorNil(t *testing.T) {
	util.MustT(t, Nil == nil)
	util.MustT(t, Nil.IsNil())
	util.MustT(t, Nil.GetCode() == 0)
	util.MustT(t, Nil.Error() == "ok")
}

func TestCloneWithOriginError(t *testing.T) {
	util.MustT(t, Nil.GetCode() == Nil.CloneWithOriginError(nil).GetCode())
}
