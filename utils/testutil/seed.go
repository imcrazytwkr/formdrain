package testutil

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

// SeedSite inserts an account and a site owned by that account. Returns the account id.
func SeedSite(t *testing.T, db *sql.DB, siteID int64, hostname string) int64 {
	t.Helper()

	res, err := db.Exec(
		`INSERT INTO accounts (email, password_hash, created_at) VALUES (?, ?, ?)`,
		fmt.Sprintf("owner-%d@example.com", siteID),
		"unused-hash",
		time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	ownerID, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("account id: %v", err)
	}

	_, err = db.Exec(
		`INSERT INTO sites (id, hostname, owner_id) VALUES (?, ?, ?)`,
		siteID,
		hostname,
		ownerID,
	)
	if err != nil {
		t.Fatalf("seed site: %v", err)
	}

	return ownerID
}

// SeedSiteForOwner inserts a site owned by an existing account.
func SeedSiteForOwner(t *testing.T, db *sql.DB, ownerID, siteID int64, hostname string) {
	t.Helper()

	_, err := db.Exec(
		`INSERT INTO sites (id, hostname, owner_id) VALUES (?, ?, ?)`,
		siteID,
		hostname,
		ownerID,
	)
	if err != nil {
		t.Fatalf("seed site: %v", err)
	}
}

// SeedForm inserts a minimal form row for list/get repository tests.
func SeedForm(t *testing.T, db *sql.DB, formID, siteID int64) {
	t.Helper()

	_, err := db.Exec(`
		INSERT INTO forms (id, site_id, captcha_type, field_schema, schema_version, notifiers)
		VALUES (?, ?, 'hcaptcha', '{"version":1,"fields":[]}', 1, '{}');
	`, formID, siteID)
	if err != nil {
		t.Fatalf("seed form: %v", err)
	}
}
