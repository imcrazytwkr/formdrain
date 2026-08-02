package errorutil_test

import (
	"errors"
	"testing"

	"github.com/imcrazytwkr/formdrain/utils/errorutil"
)

func TestUnwrapMultiErr(t *testing.T) {
	t.Parallel()

	joined := errors.Join(errors.New("a"), errors.New("b"))
	errs, ok := errorutil.UnwrapMultiErr(joined)
	if !ok || len(errs) != 2 {
		t.Fatalf("got ok=%v errs=%v", ok, errs)
	}

	_, ok = errorutil.UnwrapMultiErr(errors.New("single"))
	if ok {
		t.Fatal("expected false for single error")
	}
}
