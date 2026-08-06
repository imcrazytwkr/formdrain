package bodyparser_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/httpserver/contenttype"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/bodyparser"
)

const maxBodySize = 4096

func jsonRequest(t *testing.T, body string, contentLength int64) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.ContentLength = contentLength
	return req.WithContext(contenttype.WithContentType(req.Context(), m.ContentTypeJSON))
}

func formRequest(t *testing.T, body string, contentLength int64) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.ContentLength = contentLength
	return req.WithContext(contenttype.WithContentType(req.Context(), m.ContentTypeHTML))
}

func TestParse_JSONSuccess(t *testing.T) {
	body := `{"email":"a@b.c"}`
	req := jsonRequest(t, body, int64(len(body)))

	got, err := bodyparser.Parse(req)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if got["email"] != "a@b.c" {
		t.Fatalf("email = %#v", got["email"])
	}
}

func TestParse_FormSuccess(t *testing.T) {
	body := "email=a%40b.c&name=Ada"
	req := formRequest(t, body, int64(len(body)))

	got, err := bodyparser.Parse(req)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if got["email"] != "a@b.c" {
		t.Fatalf("email = %#v", got["email"])
	}
	if got["name"] != "Ada" {
		t.Fatalf("name = %#v", got["name"])
	}
}

func TestParse_UnsupportedContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))

	_, err := bodyparser.Parse(req)
	if !errors.Is(err, bodyparser.ErrUnsupportedContentType) {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_MalformedJSON_UnexpectedEOF(t *testing.T) {
	body := `{`
	req := jsonRequest(t, body, int64(len(body)))

	_, err := bodyparser.Parse(req)
	if !errors.Is(err, bodyparser.ErrMalformedBody) {
		t.Fatalf("err = %v", err)
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestParse_BodyReadError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", errReader{})
	req.ContentLength = 1
	req = req.WithContext(contenttype.WithContentType(req.Context(), m.ContentTypeJSON))

	_, err := bodyparser.Parse(req)
	if !errors.Is(err, bodyparser.ErrMalformedBody) {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_ExactLimitOK(t *testing.T) {
	// {"email":"…"} is 12 bytes of framing; pad fills the rest to maxBodySize.
	pad := strings.Repeat("a", maxBodySize-12)
	body := `{"email":"` + pad + `"}`
	if len(body) != maxBodySize {
		t.Fatalf("fixture length = %d, want %d", len(body), maxBodySize)
	}
	req := jsonRequest(t, body, -1)

	got, err := bodyparser.Parse(req)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if got["email"] != pad {
		t.Fatalf("email length = %d", len(got["email"].(string)))
	}
}

func TestParse_ExactOverRead(t *testing.T) {
	pad := strings.Repeat("a", maxBodySize+1-12)
	body := `{"email":"` + pad + `"}`
	if len(body) != maxBodySize+1 {
		t.Fatalf("fixture length = %d, want %d", len(body), maxBodySize+1)
	}
	req := jsonRequest(t, body, -1)

	_, err := bodyparser.Parse(req)
	if !errors.Is(err, bodyparser.ErrBodyTooLarge) {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_BodyTooLarge_UnderstatedContentLength(t *testing.T) {
	body := `{"email":"` + strings.Repeat("a", 5000) + `"}`
	req := jsonRequest(t, body, 100)

	_, err := bodyparser.Parse(req)
	if !errors.Is(err, bodyparser.ErrBodyTooLarge) {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_BodyTooLarge_ContentLength(t *testing.T) {
	body := `{"email":"` + strings.Repeat("a", 5000) + `"}`
	req := jsonRequest(t, body, int64(len(body)))

	_, err := bodyparser.Parse(req)
	if !errors.Is(err, bodyparser.ErrBodyTooLarge) {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_BodyTooLarge_Truncation(t *testing.T) {
	body := `{"email":"` + strings.Repeat("a", 5000) + `"}`
	req := jsonRequest(t, body, -1)

	_, err := bodyparser.Parse(req)
	if !errors.Is(err, bodyparser.ErrBodyTooLarge) {
		t.Fatalf("err = %v", err)
	}
}
