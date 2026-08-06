package account

import (
	"errors"
	"testing"

	"github.com/imcrazytwkr/formdrain/models/account"
	ar "github.com/imcrazytwkr/formdrain/repositories/account"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestService_Login(t *testing.T) {
	db := testutil.OpenSqlite(t)
	accounts := ar.NewSqliteAccountRepository(db)
	svc := NewService(accounts)

	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}
	acct := &account.Account{Email: "a@example.com", PasswordHash: hash}
	if err := accounts.Create(t.Context(), acct); err != nil {
		t.Fatal(err)
	}

	gotAcct, err := svc.Login(t.Context(), "a@example.com", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if gotAcct.ID != acct.ID || gotAcct.Email != acct.Email {
		t.Fatalf("acct=%#v", gotAcct)
	}

	_, err = svc.Login(t.Context(), "a@example.com", "wrong")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("err = %v", err)
	}

	_, err = svc.Login(t.Context(), "missing@example.com", "secret")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("err = %v", err)
	}
}
