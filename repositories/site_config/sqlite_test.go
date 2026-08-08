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

	missing, err := repo.GetSiteConfigById(ctx, 1)
	if err != nil || missing != nil {
		t.Fatalf("missing: %v %v", missing, err)
	}
}
