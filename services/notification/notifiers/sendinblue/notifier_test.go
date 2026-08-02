package sendinblue_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/models/common"
	"github.com/imcrazytwkr/formdrain/models/form_config/sendinblue"
	sn "github.com/imcrazytwkr/formdrain/services/notification/notifiers/sendinblue"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestSend_Success(t *testing.T) {
	t.Parallel()

	var sawURL string
	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			sawURL = r.URL.String()
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), "body a@b.c") {
				t.Fatalf("body = %s", body)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	tmpl, err := common.NewTemplate("body {{email}}")
	if err != nil {
		t.Fatal(err)
	}

	n := sn.NewSendinblueNotifier("Sender", "from@example.com", client)
	err = n.Send(&sendinblue.SendinblueConfig{
		Recipients: []*sendinblue.EmailContact{{Name: "A", Address: "a@b.c"}},
		Subject:    "subj",
		Template:   tmpl,
	}, map[string]any{"email": "a@b.c"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sawURL, "sendinblue.com") {
		t.Fatalf("url = %q", sawURL)
	}
}

func TestSend_NoRecipients(t *testing.T) {
	t.Parallel()

	n := sn.NewSendinblueNotifier("Sender", "from@example.com", &http.Client{})
	err := n.Send(&sendinblue.SendinblueConfig{
		Recipients: nil,
		Template:   mustTemplate(t, "x"),
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSend_HTTPError(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	n := sn.NewSendinblueNotifier("Sender", "from@example.com", client)
	err := n.Send(&sendinblue.SendinblueConfig{
		Recipients: []*sendinblue.EmailContact{{Address: "a@b.c"}},
		Subject:    "s",
		Template:   mustTemplate(t, "hi"),
	}, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("err = %v", err)
	}
}

func mustTemplate(t *testing.T, raw string) *common.Template {
	t.Helper()
	tmpl, err := common.NewTemplate(raw)
	if err != nil {
		t.Fatal(err)
	}
	return tmpl
}
