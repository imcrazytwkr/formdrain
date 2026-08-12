package mappers_test

import (
	"net/netip"
	"testing"
	"time"

	fr "github.com/imcrazytwkr/formdrain/models/form_response"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/mappers"
)

func TestFormResponse(t *testing.T) {
	t.Parallel()

	id := "00000000-0000-4000-8000-000000000001"
	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	ip := netip.MustParseAddr("203.0.113.10")
	src := &fr.FormResponse{
		Id:            id,
		FormId:        10,
		CreatedAt:     created,
		SchemaVersion: 3,
		ClientIP:      ip,
		Payload:       map[string]any{"email": "a@b.c"},
	}

	got, err := mappers.FormResponse(src)
	if err != nil {
		t.Fatal(err)
	}
	if got.Id.String() != id {
		t.Fatalf("Id = %s", got.Id)
	}
	if got.FormId != 10 || got.SchemaVersion != 3 {
		t.Fatalf("got = %#v", got)
	}
	if !got.CreatedAt.Equal(created) {
		t.Fatalf("CreatedAt = %s", got.CreatedAt)
	}
	if got.ClientIp == nil || *got.ClientIp != "203.0.113.10" {
		t.Fatalf("ClientIp = %#v", got.ClientIp)
	}
	if got.Payload["email"] != "a@b.c" {
		t.Fatalf("Payload = %#v", got.Payload)
	}
}

func TestFormResponse_NoClientIP(t *testing.T) {
	t.Parallel()

	got, err := mappers.FormResponse(&fr.FormResponse{
		Id:      "00000000-0000-4000-8000-000000000002",
		FormId:  10,
		Payload: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ClientIp != nil {
		t.Fatalf("ClientIp = %#v", got.ClientIp)
	}
	if got.Payload == nil {
		t.Fatal("expected empty payload map")
	}
}
