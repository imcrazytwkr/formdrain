package site_config_test

import (
	"testing"

	scr "github.com/imcrazytwkr/formdrain/repositories/site_config"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestGetSiteConfigById(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()

	_, err := db.Exec(`INSERT INTO sites (id, hostname) VALUES (7, 'forms.example.com')`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := scr.NewSqliteSiteConfigRepository(db)

	got, err := repo.GetSiteConfigById(ctx, 7)
	if err != nil {
		t.Fatalf("GetSiteConfigById: %v", err)
	}
	if got == nil || got.SiteId != 7 || got.Hostname != "forms.example.com" {
		t.Fatalf("got %+v", got)
	}

	missing, err := repo.GetSiteConfigById(ctx, 1)
	if err != nil || missing != nil {
		t.Fatalf("missing: %v %v", missing, err)
	}
}
