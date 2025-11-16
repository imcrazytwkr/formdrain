package form

import (
	"errors"
	"fmt"

	m "github.com/imcrazytwkr/formdrain/models/http"
)

func getErrUnsupportedFormType(contentType m.ContentType) error {
	return fmt.Errorf("unsupported form type: %q", contentType.String())
}

var errFormTooLarge = errors.New("forms data can only be up to 4kb in size")

func getErrMalformedFormData(contentType m.ContentType) error {
	return fmt.Errorf("form data is either not in %q format or is malformed", contentType.String())
}

func getErrFormNotFound(id string) error {
	return fmt.Errorf("could not find form with ID %q", id)
}
