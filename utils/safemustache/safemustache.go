// Package safemustache wraps github.com/cbroglie/mustache for owner-supplied
// templates. It fails closed on partials, unescaped interpolation, and delimiter
// changes. Parse returns a validated *mustache.Template; application types (e.g.
// models/common.Template) own JSON and Execute helpers.
//
// Callers must pass sealed data only: map[string]any, map[string]string, or small
// structs with exported fields and no methods. Do not pass values with methods;
// the underlying library may invoke them via reflection.
//
// Importing this package sets mustache.AllowMissingVariables to false for the
// whole process.
package safemustache

import (
	"fmt"
	"strings"

	"github.com/cbroglie/mustache"
)

const MaxSourceBytes = 64 << 10

func init() {
	mustache.AllowMissingVariables = false
}

// Parse validates and compiles src. Safe to call at config-save time.
func Parse(src string) (*mustache.Template, error) {
	if len(src) < 1 {
		return nil, ErrEmpty
	}

	if len(src) > MaxSourceBytes {
		return nil, ErrTooLarge
	}

	if err := scanForbidden(src); err != nil {
		return nil, err
	}

	tmpl, err := mustache.ParseStringPartials(src, noPartials)
	if err != nil {
		return nil, err
	}

	err = rejectPartials(tmpl.Tags())
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func scanForbidden(src string) error {
	if strings.Contains(src, "{{{") {
		return ErrRawInterpolation
	}

	if strings.Contains(src, "{{&") {
		return ErrRawInterpolation
	}

	if strings.Contains(src, "{{=") {
		return ErrDelimiterChange
	}

	return nil
}

func rejectPartials(tags []mustache.Tag) error {
	for _, tag := range tags {
		switch tag.Type() {
		case mustache.Partial:
			return fmt.Errorf("%w: %q", ErrPartial, tag.Name())
		case mustache.Section, mustache.InvertedSection:
			if err := rejectPartials(tag.Tags()); err != nil {
				return err
			}
		}
	}
	return nil
}
