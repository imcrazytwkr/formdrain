package mappers_test

import (
	"testing"

	models "github.com/imcrazytwkr/formdrain/models/account"
	"github.com/imcrazytwkr/formdrain/routes/auth/mappers"
)

func TestUser(t *testing.T) {
	t.Parallel()

	account := &models.Account{
		ID:           42,
		Email:        "user@example.com",
		PasswordHash: "secret-hash",
	}

	got := mappers.User(account)
	if got == nil {
		t.Fatal("expected user")
	}

	if got.ID != 42 || got.Email != "user@example.com" {
		t.Fatalf("got = %#v", got)
	}
}
