package account

import (
	"errors"
	"strings"
	"testing"
)

func TestHashPassword_RoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(hash, "$argon2id$v=19$m=65536,t=6,p=1$") {
		t.Fatalf("hash = %q", hash)
	}

	err = CheckPassword("correct horse battery staple", hash)
	if err != nil {
		t.Fatalf("CheckPassword err=%v", err)
	}

	err = CheckPassword("wrong password", hash)
	if err != ErrInvalidCredentials {
		t.Fatal("expected mismatch")
	}
}

func TestHashPassword_Empty(t *testing.T) {
	_, err := HashPassword("")
	if !errors.Is(err, ErrEmptyPassword) {
		t.Fatalf("err = %v", err)
	}
}

func TestCheckPassword_Corrupt(t *testing.T) {
	err := CheckPassword("x", "not-a-phc-hash")
	if err == nil {
		t.Fatalf("err=%v", err)
	}
}
