package form_config_test

import (
	"testing"

	fcr "github.com/imcrazytwkr/formdrain/repositories/form_config"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestListFormsByOwnerID(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	repo := fcr.NewSqliteFormConfigRepository(db)

	ownerID := testutil.SeedSite(t, db, 1, "c.example")
	testutil.SeedSiteForOwner(t, db, ownerID, 2, "a.example")
	testutil.SeedSiteForOwner(t, db, ownerID, 3, "b.example")
	testutil.SeedForm(t, db, 11, 1) // c.example
	testutil.SeedForm(t, db, 12, 1)
	testutil.SeedForm(t, db, 21, 2) // a.example
	testutil.SeedForm(t, db, 31, 3) // b.example

	testutil.SeedSite(t, db, 9, "other.example")
	testutil.SeedForm(t, db, 91, 9)

	t.Run("sort_id", func(t *testing.T) {
		got, err := repo.ListFormsByOwnerID(ctx, ownerID, "id", 0, 0, "", 10)
		if err != nil {
			t.Fatal(err)
		}
		if want := []int64{11, 12, 21, 31}; !sameIDs(got, want) {
			t.Fatalf("got %#v want %#v", formIDs(got), want)
		}

		page, err := repo.ListFormsByOwnerID(ctx, ownerID, "id", 12, 0, "", 2)
		if err != nil {
			t.Fatal(err)
		}
		if want := []int64{21, 31}; !sameIDs(page, want) {
			t.Fatalf("page = %#v want %#v", formIDs(page), want)
		}
	})

	t.Run("sort_site_id", func(t *testing.T) {
		got, err := repo.ListFormsByOwnerID(ctx, ownerID, "site_id", 0, 0, "", 10)
		if err != nil {
			t.Fatal(err)
		}
		if want := []int64{11, 12, 21, 31}; !sameIDs(got, want) {
			t.Fatalf("got %#v want %#v", formIDs(got), want)
		}
		if got[0].SiteId != 1 || got[2].SiteId != 2 || got[3].SiteId != 3 {
			t.Fatalf("site order = %#v", got)
		}

		page, err := repo.ListFormsByOwnerID(ctx, ownerID, "site_id", 12, 1, "", 10)
		if err != nil {
			t.Fatal(err)
		}
		if want := []int64{21, 31}; !sameIDs(page, want) {
			t.Fatalf("page = %#v want %#v", formIDs(page), want)
		}
	})

	t.Run("sort_hostname", func(t *testing.T) {
		got, err := repo.ListFormsByOwnerID(ctx, ownerID, "hostname", 0, 0, "", 10)
		if err != nil {
			t.Fatal(err)
		}
		if want := []int64{21, 31, 11, 12}; !sameIDs(got, want) {
			t.Fatalf("got %#v want %#v", formIDs(got), want)
		}
		if got[0].Hostname != "a.example" || got[1].Hostname != "b.example" || got[2].Hostname != "c.example" {
			t.Fatalf("hostnames = %#v", got)
		}

		page, err := repo.ListFormsByOwnerID(ctx, ownerID, "hostname", 21, 0, "a.example", 10)
		if err != nil {
			t.Fatal(err)
		}
		if want := []int64{31, 11, 12}; !sameIDs(page, want) {
			t.Fatalf("page = %#v want %#v", formIDs(page), want)
		}
	})

	t.Run("excludes_other_owners", func(t *testing.T) {
		got, err := repo.ListFormsByOwnerID(ctx, ownerID, "id", 0, 0, "", 100)
		if err != nil {
			t.Fatal(err)
		}
		for _, item := range got {
			if item.Id == 91 {
				t.Fatalf("leaked other owner form: %#v", item)
			}
		}
	})

	t.Run("bad_args", func(t *testing.T) {
		got, err := repo.ListFormsByOwnerID(ctx, 0, "id", 0, 0, "", 10)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("ownerID 0: %#v", got)
		}

		got, err = repo.ListFormsByOwnerID(ctx, ownerID, "id", 0, 0, "", 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("limit 0: %#v", got)
		}
	})
}
