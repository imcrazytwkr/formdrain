package templates

import "io"

type Template interface {
	Execute(w io.Writer, data map[string]any) error
	ExecuteString(data map[string]any) (string, error)
}
