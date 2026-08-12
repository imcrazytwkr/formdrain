package form_config_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/models/form_config"
	fcr "github.com/imcrazytwkr/formdrain/repositories/form_config"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestListFormsBySiteID(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	repo := fcr.NewSqliteFormConfigRepository(db)

	testutil.SeedSite(t, db, 1, "example.com")
	testutil.SeedForm(t, db, 10, 1)
	testutil.SeedForm(t, db, 20, 1)
	testutil.SeedForm(t, db, 30, 1)

	got, err := repo.ListFormsBySiteID(ctx, 1, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].Id != 10 || got[1].Id != 20 || got[2].Id != 30 {
		t.Fatalf("order = %#v", formIDs(got))
	}
	if got[0].SiteId != 1 || got[0].Hostname != "" {
		t.Fatalf("item = %#v", got[0])
	}

	page, err := repo.ListFormsBySiteID(ctx, 1, 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 2 || page[0].Id != 20 || page[1].Id != 30 {
		t.Fatalf("page after 10 = %#v", formIDs(page))
	}

	empty, err := repo.ListFormsBySiteID(ctx, 1, 30, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected empty after last, got %#v", formIDs(empty))
	}

	badSite, err := repo.ListFormsBySiteID(ctx, 0, 0, 10)
	if err != nil || badSite != nil {
		t.Fatalf("siteID 0: got %#v err %v", badSite, err)
	}

	badLimit, err := repo.ListFormsBySiteID(ctx, 1, 0, 0)
	if err != nil || badLimit != nil {
		t.Fatalf("limit 0: got %#v err %v", badLimit, err)
	}
}

func formIDs(items []*form_config.FormListItem) []int64 {
	out := make([]int64, len(items))
	for i, item := range items {
		out[i] = item.Id
	}
	return out
}

func sameIDs(items []*form_config.FormListItem, want []int64) bool {
	if len(items) != len(want) {
		return false
	}
	for i := range want {
		if items[i].Id != want[i] {
			return false
		}
	}
	return true
}
