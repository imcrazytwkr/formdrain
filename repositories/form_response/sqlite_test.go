package form_response_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/netip"
	"testing"

	"github.com/imcrazytwkr/formdrain/models/form_response"
	frr "github.com/imcrazytwkr/formdrain/repositories/form_response"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestSaveFormResponse(t *testing.T) {
	db := testutil.OpenSqlite(t)

	_, err := db.Exec(`
		INSERT INTO sites (id, hostname) VALUES (1, 'example.com');
		INSERT INTO forms (id, site_id, captcha_type, field_schema, schema_version, notifiers)
		VALUES (10, 1, 'hcaptcha', '{"version":1,"fields":[]}', 3, '{}');
	`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := frr.NewSqliteFormResponseRepository(db)
	ip := netip.MustParseAddr("203.0.113.10")
	id, err := repo.SaveFormResponse(t.Context(), &form_response.FormResponse{
		FormId:        10,
		SchemaVersion: 3,
		ClientIP:      ip,
		Payload:       map[string]any{"email": "a@b.c"},
	})
	if err != nil {
		t.Fatalf("SaveFormResponse: %v", err)
	}
	if len(id) < 1 {
		t.Fatal("empty id")
	}

	var formId int64
	var schemaVersion int
	var clientIP string
	var payload string
	var createdAt string
	err = db.QueryRow(`
		SELECT form_id, schema_version, client_ip, payload, created_at
		FROM form_responses WHERE id = ?`, id,
	).Scan(&formId, &schemaVersion, &clientIP, &payload, &createdAt)
	if err != nil {
		t.Fatalf("select: %v", err)
	}
	if formId != 10 || schemaVersion != 3 || clientIP != "203.0.113.10" {
		t.Fatalf("envelope form=%d ver=%d ip=%q", formId, schemaVersion, clientIP)
	}
	if len(createdAt) < 1 {
		t.Fatal("missing created_at")
	}

	var parsed map[string]any
	err = json.Unmarshal([]byte(payload), &parsed)
	if err != nil {
		t.Fatalf("payload json: %v", err)
	}
	if parsed["email"] != "a@b.c" {
		t.Fatalf("payload: %#v", parsed)
	}
}

func TestSaveFormResponse_Nil(t *testing.T) {
	db := testutil.OpenSqlite(t)
	repo := frr.NewSqliteFormResponseRepository(db)
	_, err := repo.SaveFormResponse(t.Context(), nil)
	if !errors.Is(err, frr.ErrEmptyFormResponse) {
		t.Fatalf("err = %v", err)
	}
}

func TestSaveFormResponse_NoClientIP(t *testing.T) {
	db := testutil.OpenSqlite(t)

	_, err := db.Exec(`
		INSERT INTO sites (id, hostname) VALUES (1, 'example.com');
		INSERT INTO forms (id, site_id, captcha_type, field_schema, schema_version, notifiers)
		VALUES (10, 1, 'hcaptcha', '{"version":1,"fields":[]}', 1, '{}');
	`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := frr.NewSqliteFormResponseRepository(db)
	id, err := repo.SaveFormResponse(t.Context(), &form_response.FormResponse{
		FormId:        10,
		SchemaVersion: 1,
		Payload:       map[string]any{},
	})
	if err != nil {
		t.Fatal(err)
	}

	var clientIP sql.NullString
	err = db.QueryRow(`SELECT client_ip FROM form_responses WHERE id = ?`, id).Scan(&clientIP)
	if err != nil {
		t.Fatal(err)
	}
	if clientIP.Valid {
		t.Fatalf("expected NULL client_ip, got %q", clientIP.String)
	}
}
