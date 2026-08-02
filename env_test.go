package main

import (
	"net/url"
	"testing"
)

func TestGetHost(t *testing.T) {
	t.Setenv("HOST", "")
	got, err := getHost()
	if err != nil || got != "" {
		t.Fatalf("empty: got %q err %v", got, err)
	}

	t.Setenv("HOST", "127.0.0.1")
	got, err = getHost()
	if err != nil || got != "127.0.0.1" {
		t.Fatalf("valid: got %q err %v", got, err)
	}

	t.Setenv("HOST", "not-an-ip")
	_, err = getHost()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetPort(t *testing.T) {
	t.Setenv("PORT", "")
	got, err := getPort()
	if err != nil || got != "8080" {
		t.Fatalf("default: got %q err %v", got, err)
	}

	t.Setenv("PORT", "3000")
	got, err = getPort()
	if err != nil || got != "3000" {
		t.Fatalf("valid: got %q err %v", got, err)
	}

	t.Setenv("PORT", "99999")
	_, err = getPort()
	if err == nil {
		t.Fatal("expected error")
	}
}

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
