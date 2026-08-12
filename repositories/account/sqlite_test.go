package account_test

import (
	"testing"

	ma "github.com/imcrazytwkr/formdrain/models/account"
	ar "github.com/imcrazytwkr/formdrain/repositories/account"
	"github.com/imcrazytwkr/formdrain/services/account"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestSqliteAccountRepository_GetByEmail(t *testing.T) {
	db := testutil.OpenSqlite(t)
	repo := ar.NewSqliteAccountRepository(db)
	ctx := t.Context()

	hash, err := account.HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}

	acct := &ma.Account{
		Email:        "User@Example.com",
		PasswordHash: hash,
	}
	if err := repo.Create(ctx, acct); err != nil {
		t.Fatal(err)
	}
	if acct.ID < 1 {
		t.Fatalf("ID = %d", acct.ID)
	}

	got, err := repo.GetByEmail(ctx, "user@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected account")
	}
	if got.Email != "User@Example.com" {
		t.Fatalf("Email = %q", got.Email)
	}
	if got.PasswordHash != hash {
		t.Fatal("unexpected hash")
	}

	missing, err := repo.GetByEmail(ctx, "missing@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if missing != nil {
		t.Fatalf("expected nil, got %#v", missing)
	}
}

func TestSqliteAccountRepository_UniqueEmail(t *testing.T) {
	db := testutil.OpenSqlite(t)
	repo := ar.NewSqliteAccountRepository(db)
	ctx := t.Context()

	hash, err := account.HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Create(ctx, &ma.Account{Email: "a@example.com", PasswordHash: hash}); err != nil {
		t.Fatal(err)
	}
	err = repo.Create(ctx, &ma.Account{Email: "A@example.com", PasswordHash: hash})
	if err == nil {
		t.Fatal("expected unique email error")
	}
}

func TestSqliteAccountRepository_GetByID(t *testing.T) {
	db := testutil.OpenSqlite(t)
	repo := ar.NewSqliteAccountRepository(db)
	ctx := t.Context()

	hash, err := account.HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}

	acct := &ma.Account{Email: "id@example.com", PasswordHash: hash}
	if err := repo.Create(ctx, acct); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(ctx, acct.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected account")
	}
	if got.ID != acct.ID || got.Email != acct.Email || got.PasswordHash != hash {
		t.Fatalf("got = %#v", got)
	}

	missing, err := repo.GetByID(ctx, acct.ID+999)
	if err != nil {
		t.Fatal(err)
	}
	if missing != nil {
		t.Fatalf("expected nil, got %#v", missing)
	}

	invalid, err := repo.GetByID(ctx, 0)
	if err != nil || invalid != nil {
		t.Fatalf("id 0: got %#v err %v", invalid, err)
	}
}
