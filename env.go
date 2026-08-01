package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

func getHost() (string, error) {
	listenHost := os.Getenv("HOST")
	if len(listenHost) == 0 {
		return "", nil
	}

	// ParseIP returns nil on invalid IP
	if net.ParseIP(listenHost) == nil {
		return "", fmt.Errorf("listen host %q is not a valid IP address", listenHost)
	}

	return listenHost, nil
}

func getPort() (string, error) {
	listenPort := os.Getenv("PORT")
	if len(listenPort) == 0 {
		log.Debug().Msg("Listen port is unset or empty, falling back to default")
		listenPort = "8080"
	}

	// Checking is port number fits in Uint16
	_, err := strconv.ParseUint(listenPort, 10, 16)
	if err != nil {
		return "", fmt.Errorf("listen port number %q is invalid; Valid range is 0-65535", listenPort)
	}

	return listenPort, nil
}

func getConnString() (*url.URL, error) {
	mongoConnstring := os.Getenv("MONGO_DBURL")
	if len(mongoConnstring) < 1 {
		return nil, errors.New("mongoDB connstring is not set")
	}

	url, err := url.Parse(mongoConnstring)
	if err != nil {
		return nil, fmt.Errorf("mongoDB connstring %q is invalid", mongoConnstring)
	}

	if url.Scheme != "mongo" {
		return nil, fmt.Errorf("mongoDB connstring %q has invalid scheme", url.Scheme)
	}

	return url, nil
}

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
