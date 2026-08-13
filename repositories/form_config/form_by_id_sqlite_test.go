package form_config_test

import (
	"testing"

	c "github.com/imcrazytwkr/formdrain/models/common"
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
	if got.CaptchaType != c.CaptchaTypeHcaptcha {
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
