package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
)

// getDBURL reads DBURL (e.g. sqlite:/path/to/file.db).
func getDBURL() (*url.URL, error) {
	raw := os.Getenv("DBURL")
	if len(raw) < 1 {
		return nil, errors.New("DBURL is not set")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("DBURL %q is invalid", raw)
	}

	if len(parsed.Scheme) < 1 {
		return nil, fmt.Errorf("DBURL %q has no scheme", raw)
	}

	return parsed, nil
}

func sqliteFilePath(u *url.URL) (string, error) {
	if u.Scheme != "sqlite" {
		return "", fmt.Errorf("DBURL scheme %q is not sqlite", u.Scheme)
	}

	path := u.Path
	if len(path) < 1 {
		path = u.Opaque
	}
	if len(path) < 1 {
		return "", fmt.Errorf("DBURL %q has an empty sqlite path", u.Redacted())
	}

	return path, nil
}

func getBrevoAPIKey() (string, error) {
	key := os.Getenv("BREVO_API_KEY")
	if len(key) < 1 {
		return "", errors.New("BREVO_API_KEY is not set")
	}
	return key, nil
}
