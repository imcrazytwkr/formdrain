package main

import (
	"net/url"
	"testing"
)

func TestGetDBURL(t *testing.T) {
	t.Setenv("DBURL", "")
	_, err := getDBURL()
	if err == nil {
		t.Fatal("expected unset error")
	}

	t.Setenv("DBURL", "://bad")
	_, err = getDBURL()
	if err == nil {
		t.Fatal("expected parse error")
	}

	t.Setenv("DBURL", "sqlite:/tmp/x.db")
	u, err := getDBURL()
	if err != nil || u.Scheme != "sqlite" {
		t.Fatalf("got %#v err %v", u, err)
	}
}

func TestGetBrevoAPIKey(t *testing.T) {
	t.Setenv("BREVO_API_KEY", "")
	_, err := getBrevoAPIKey()
	if err == nil {
		t.Fatal("expected unset error")
	}

	t.Setenv("BREVO_API_KEY", "secret-key")
	got, err := getBrevoAPIKey()
	if err != nil || got != "secret-key" {
		t.Fatalf("got %q err %v", got, err)
	}
}

func TestSqliteFilePath(t *testing.T) {
	u, err := url.Parse("sqlite:/tmp/form.db")
	if err != nil {
		t.Fatal(err)
	}
	path, err := sqliteFilePath(u)
	if err != nil || path != "/tmp/form.db" {
		t.Fatalf("path=%q err=%v", path, err)
	}

	u, _ = url.Parse("postgres://x")
	_, err = sqliteFilePath(u)
	if err == nil {
		t.Fatal("expected scheme error")
	}

	u, _ = url.Parse("sqlite:")
	_, err = sqliteFilePath(u)
	if err == nil {
		t.Fatal("expected empty path error")
	}
}
