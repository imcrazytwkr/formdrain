package form_config_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/models/form_config"
	fcr "github.com/imcrazytwkr/formdrain/repositories/form_config"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestGetFormConfigById(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()

	testutil.SeedSite(t, db, 1, "example.com")
	_, err := db.Exec(`
		INSERT INTO forms (id, site_id, captcha_type, redirect_to, field_schema, schema_version, notifiers)
		VALUES (
			10,
			1,
			'hcaptcha',
			'https://example.com/thanks',
			'{"version":1,"fields":[{"name":"email","type":"string","required":true}]}',
			1,
			'{"discord":null,"brevo":null}'
		);
	`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := fcr.NewSqliteFormConfigRepository(db)

	got, err := repo.GetFormConfigById(ctx, 10)
	if err != nil {
		t.Fatalf("GetFormConfigById: %v", err)
	}
	if got == nil {
		t.Fatal("expected form config")
	}
	if got.FormId != 10 || got.SiteId != 1 {
		t.Fatalf("ids: got form=%d site=%d", got.FormId, got.SiteId)
	}
	if got.CaptchaType != form_config.CaptchaTypeHcaptcha {
		t.Fatalf("captcha: got %v", got.CaptchaType)
	}
	if got.CaptchaField != "" {
		t.Fatalf("captcha_field: %q", got.CaptchaField)
	}
	if got.RedirectTo != "https://example.com/thanks" {
		t.Fatalf("redirect: %q", got.RedirectTo)
	}
	if got.SchemaVersion != 1 || len(got.FieldSchema.Fields) != 1 || got.FieldSchema.Fields[0].Name != "email" {
		t.Fatalf("schema: %+v", got.FieldSchema)
	}

	missing, err := repo.GetFormConfigById(ctx, 999)
	if err != nil {
		t.Fatalf("missing: %v", err)
	}
	if missing != nil {
		t.Fatal("expected nil for missing form")
	}

	invalid, err := repo.GetFormConfigById(ctx, 0)
	if err != nil || invalid != nil {
		t.Fatalf("id 0: got %v err %v", invalid, err)
	}
}

func TestGetFormConfigById_CorruptFieldSchema(t *testing.T) {
	db := testutil.OpenSqlite(t)

	testutil.SeedSite(t, db, 1, "example.com")
	_, err := db.Exec(`
		INSERT INTO forms (id, site_id, captcha_type, field_schema, schema_version, notifiers)
		VALUES (20, 1, 'hcaptcha', '"not-an-object"', 1, '{}');
	`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := fcr.NewSqliteFormConfigRepository(db)
	_, err = repo.GetFormConfigById(t.Context(), 20)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}

func TestGetFormConfigById_InvalidCaptchaType(t *testing.T) {
	db := testutil.OpenSqlite(t)

	testutil.SeedSite(t, db, 1, "example.com")
	_, err := db.Exec(`
		INSERT INTO forms (id, site_id, captcha_type, field_schema, schema_version, notifiers)
		VALUES (21, 1, 'bogus', '{"version":1,"fields":[]}', 1, '{}');
	`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := fcr.NewSqliteFormConfigRepository(db)
	_, err = repo.GetFormConfigById(t.Context(), 21)
	if err == nil {
		t.Fatal("expected captcha type error")
	}
}

func TestGetFormConfigById_CorruptNotifiers(t *testing.T) {
	db := testutil.OpenSqlite(t)

	testutil.SeedSite(t, db, 1, "example.com")
	_, err := db.Exec(`
		INSERT INTO forms (id, site_id, captcha_type, field_schema, schema_version, notifiers)
		VALUES (22, 1, 'hcaptcha', '{"version":1,"fields":[]}', 1, '"nope"');
	`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := fcr.NewSqliteFormConfigRepository(db)
	_, err = repo.GetFormConfigById(t.Context(), 22)
	if err == nil {
		t.Fatal("expected notifiers unmarshal error")
	}
}

func TestGetFormConfigById_CaptchaField(t *testing.T) {
	db := testutil.OpenSqlite(t)

	testutil.SeedSite(t, db, 1, "example.com")
	_, err := db.Exec(`
		INSERT INTO forms (id, site_id, captcha_type, captcha_field, field_schema, schema_version, notifiers)
		VALUES (30, 1, 'hcaptcha', 'cf-turnstile-response', '{"version":1,"fields":[]}', 1, '{}');
	`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := fcr.NewSqliteFormConfigRepository(db)
	got, err := repo.GetFormConfigById(t.Context(), 30)
	if err != nil {
		t.Fatalf("GetFormConfigById: %v", err)
	}
	if got.CaptchaField != "cf-turnstile-response" {
		t.Fatalf("captcha_field = %q", got.CaptchaField)
	}
	if got.CaptchaTokenField() != "cf-turnstile-response" {
		t.Fatalf("CaptchaTokenField = %q", got.CaptchaTokenField())
	}
}

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

	otherOwner := testutil.SeedSite(t, db, 9, "other.example")
	_ = otherOwner
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
