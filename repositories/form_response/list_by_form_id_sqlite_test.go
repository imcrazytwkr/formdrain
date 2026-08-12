package form_response_test

import (
	"net/netip"
	"testing"
	"time"

	"github.com/imcrazytwkr/formdrain/models/form_response"
	frr "github.com/imcrazytwkr/formdrain/repositories/form_response"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestListFormResponsesByFormID(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()

	testutil.SeedSite(t, db, 1, "example.com")
	testutil.SeedForm(t, db, 10, 1)
	testutil.SeedForm(t, db, 11, 1)

	repo := frr.NewSqliteFormResponseRepository(db)
	ip := netip.MustParseAddr("203.0.113.10")
	t3 := time.Date(2026, 1, 1, 0, 0, 3, 0, time.UTC)
	t2 := time.Date(2026, 1, 1, 0, 0, 2, 0, time.UTC)
	t1 := time.Date(2026, 1, 1, 0, 0, 1, 0, time.UTC)

	idA := "00000000-0000-4000-8000-00000000000a"
	idB := "00000000-0000-4000-8000-00000000000b"
	idC := "00000000-0000-4000-8000-00000000000c"
	idD := "00000000-0000-4000-8000-00000000000d"
	idZ := "ffffffff-ffff-4fff-bfff-ffffffffffff"
	idOther := "00000000-0000-4000-8000-0000000000ff"

	mustSave := func(resp *form_response.FormResponse) {
		t.Helper()
		if _, err := repo.SaveFormResponse(ctx, resp); err != nil {
			t.Fatal(err)
		}
	}

	mustSave(&form_response.FormResponse{
		Id: idA, FormId: 10, CreatedAt: t3, SchemaVersion: 1,
		ClientIP: ip, Payload: map[string]any{"n": "a"},
	})
	mustSave(&form_response.FormResponse{
		Id: idZ, FormId: 10, CreatedAt: t3, SchemaVersion: 1,
		ClientIP: ip, Payload: map[string]any{"n": "z"},
	})
	mustSave(&form_response.FormResponse{
		Id: idB, FormId: 10, CreatedAt: t2, SchemaVersion: 1,
		ClientIP: ip, Payload: map[string]any{"n": "b"},
	})
	mustSave(&form_response.FormResponse{
		Id: idC, FormId: 10, CreatedAt: t1, SchemaVersion: 2,
		ClientIP: ip, Payload: map[string]any{"n": "c"},
	})
	mustSave(&form_response.FormResponse{
		Id: idD, FormId: 10, CreatedAt: t1, SchemaVersion: 2,
		Payload: map[string]any{"n": "d"},
	})
	mustSave(&form_response.FormResponse{
		Id: idOther, FormId: 11, CreatedAt: time.Date(2026, 1, 1, 0, 0, 9, 0, time.UTC),
		SchemaVersion: 1, Payload: map[string]any{"n": "other"},
	})

	got, err := repo.ListFormResponsesByFormID(ctx, 10, time.Time{}, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idA, idZ, idB, idC, idD}; !sameResponseIDs(got, want) {
		t.Fatalf("got %#v want %#v", responseIDs(got), want)
	}
	if got[0].FormId != 10 || got[0].SchemaVersion != 1 {
		t.Fatalf("envelope = %#v", got[0])
	}
	if got[0].Payload["n"] != "a" {
		t.Fatalf("payload = %#v", got[0].Payload)
	}
	if got[0].ClientIP != ip {
		t.Fatalf("ip = %s", got[0].ClientIP)
	}
	if got[4].ClientIP.IsValid() {
		t.Fatalf("expected invalid ip, got %s", got[4].ClientIP)
	}

	page, err := repo.ListFormResponsesByFormID(ctx, 10, t3, idA, 2)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idZ, idB}; !sameResponseIDs(page, want) {
		t.Fatalf("page = %#v want %#v", responseIDs(page), want)
	}

	// idZ is newer than t1 but sorts after idC as TEXT; timestamp must exclude it.
	rest, err := repo.ListFormResponsesByFormID(ctx, 10, t1, idC, 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idD}; !sameResponseIDs(rest, want) {
		t.Fatalf("rest = %#v want %#v", responseIDs(rest), want)
	}

	empty, err := repo.ListFormResponsesByFormID(ctx, 99, time.Time{}, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty) != 0 {
		t.Fatalf("empty = %#v", empty)
	}
}

func responseIDs(rows []*form_response.FormResponse) []string {
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = row.Id
	}
	return out
}

func sameResponseIDs(rows []*form_response.FormResponse, want []string) bool {
	got := responseIDs(rows)
	if len(got) != len(want) {
		return false
	}
	for i := range want {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
