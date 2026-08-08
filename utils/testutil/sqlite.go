package testutil

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	_ "modernc.org/sqlite"
)

// OpenSqlite returns a temporary SQLite DB with all migrations/*.sql applied in name order.
func OpenSqlite(t *testing.T) *sql.DB {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	migrationsDir := filepath.Join(filepath.Dir(thisFile), "..", "..", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("read migrations: %v", err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		names = append(names, e.Name())
	}

	if len(names) < 1 {
		t.Fatal("no migration SQL files found")
	}

	slices.Sort(names)

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		_ = db.Close()
		t.Fatalf("enable foreign_keys: %v", err)
	}

	for _, name := range names {
		migrationSQL, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			_ = db.Close()
			t.Fatalf("read migration %s: %v", name, err)
		}
		_, err = db.Exec(string(migrationSQL))
		if err != nil {
			_ = db.Close()
			t.Fatalf("apply migration %s: %v", name, err)
		}
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}
