package testutil

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "modernc.org/sqlite"
)

// OpenSqlite returns a temporary SQLite DB with migrations/001_init.sql applied.
func OpenSqlite(t *testing.T) *sql.DB {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	migrationPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "migrations", "001_init.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		_ = db.Close()
		t.Fatalf("apply migration: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}
