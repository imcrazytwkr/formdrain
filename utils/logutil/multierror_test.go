package logutil_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/imcrazytwkr/formdrain/utils/logutil"
	"github.com/rs/zerolog"
)

func TestUnwrapErr_Joined(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := zerolog.New(&buf)
	err := errors.Join(errors.New("a"), errors.New("b"))

	logutil.UnwrapErr(log.Error(), err).Msg("fail")

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	errs, ok := payload["errors"].([]any)
	if !ok || len(errs) != 2 {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestUnwrapErr_Single(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := zerolog.New(&buf)
	logutil.UnwrapErr(log.Error(), errors.New("solo")).Msg("fail")

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["error"] != "solo" {
		t.Fatalf("payload = %#v", payload)
	}
}
