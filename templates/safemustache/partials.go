package safemustache

import (
	"fmt"

	"github.com/cbroglie/mustache"
)

var noPartials mustache.PartialProvider = &denyPartials{}

type denyPartials struct{}

func (denyPartials) Get(name string) (string, error) {
	return "", fmt.Errorf("%w: %q", ErrPartial, name)
}
