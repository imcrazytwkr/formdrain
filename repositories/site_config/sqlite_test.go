package site_config_test

import (
	"testing"

	scr "github.com/imcrazytwkr/formdrain/repositories/site_config"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestGetSiteConfigById(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()

	ownerID := testutil.SeedSite(t, db, 7, "forms.example.com")

	repo := scr.NewSqliteSiteConfigRepository(db)

	got, err := repo.GetSiteConfigById(ctx, 7)
	if err != nil {
		t.Fatalf("GetSiteConfigById: %v", err)
	}
	if got == nil || got.SiteId != 7 || got.Hostname != "forms.example.com" || got.OwnerId != ownerID {
		t.Fatalf("got %+v want owner_id=%d", got, ownerID)
	}
	if got.HcaptchaSecret != "" || got.RecaptchaSecret != "" {
		t.Fatalf("unset secrets = %+v", got)
	}

	missing, err := repo.GetSiteConfigById(ctx, 1)
	if err != nil || missing != nil {
		t.Fatalf("missing: %v %v", missing, err)
	}
}

func TestListByOwnerIDAfterID(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	ownerID := testutil.SeedSite(t, db, 1, "c.example")
	testutil.SeedSiteForOwner(t, db, ownerID, 2, "a.example")
	testutil.SeedSiteForOwner(t, db, ownerID, 3, "b.example")

	repo := scr.NewSqliteSiteConfigRepository(db)

	page, err := repo.ListByOwnerIDAfterID(ctx, ownerID, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 2 || page[0].SiteId != 1 || page[1].SiteId != 2 {
		t.Fatalf("first page = %#v", page)
	}

	page, err = repo.ListByOwnerIDAfterID(ctx, ownerID, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 1 || page[0].SiteId != 3 {
		t.Fatalf("second page = %#v", page)
	}

	empty, err := repo.ListByOwnerIDAfterID(ctx, ownerID+99, 0, 10)
	if err != nil || len(empty) != 0 {
		t.Fatalf("empty = %#v err=%v", empty, err)
	}
}

func TestListByOwnerIDAfterHostname(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	ownerID := testutil.SeedSite(t, db, 1, "c.example")
	testutil.SeedSiteForOwner(t, db, ownerID, 2, "a.example")
	testutil.SeedSiteForOwner(t, db, ownerID, 3, "b.example")

	repo := scr.NewSqliteSiteConfigRepository(db)

	page, err := repo.ListByOwnerIDAfterHostname(ctx, ownerID, "", 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 2 || page[0].Hostname != "a.example" || page[1].Hostname != "b.example" {
		t.Fatalf("first page = %#v", page)
	}

	page, err = repo.ListByOwnerIDAfterHostname(ctx, ownerID, "b.example", 3, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 1 || page[0].Hostname != "c.example" {
		t.Fatalf("second page = %#v", page)
	}
}

func TestGetSiteConfigById_CaptchaSecrets(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()

	testutil.SeedSite(t, db, 7, "forms.example.com")
	_, err := db.Exec(
		`UPDATE sites SET hcaptcha_secret = ?, recaptcha_secret = ? WHERE id = ?`,
		"h-secret",
		"r-secret",
		7,
	)
	if err != nil {
		t.Fatalf("update secrets: %v", err)
	}

	repo := scr.NewSqliteSiteConfigRepository(db)

	got, err := repo.GetSiteConfigById(ctx, 7)
	if err != nil {
		t.Fatalf("GetSiteConfigById: %v", err)
	}
	if got == nil || got.HcaptchaSecret != "h-secret" || got.RecaptchaSecret != "r-secret" {
		t.Fatalf("got %+v", got)
	}

	page, err := repo.ListByOwnerIDAfterID(ctx, got.OwnerId, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 1 || page[0].HcaptchaSecret != "h-secret" || page[0].RecaptchaSecret != "r-secret" {
		t.Fatalf("list = %#v", page)
	}
}
