package account

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/imcrazytwkr/formdrain/models/account"
	"github.com/imcrazytwkr/formdrain/repositories"
)

type sqliteAccountRepository struct {
	db *sql.DB
}

func NewSqliteAccountRepository(db *sql.DB) repositories.AccountRepository {
	return &sqliteAccountRepository{db: db}
}

const insertAccount = `
INSERT INTO accounts (
	email,
	password_hash,
	created_at
) VALUES (?, ?, ?)
`

func (r *sqliteAccountRepository) Create(ctx context.Context, acct *account.Account) error {
	if acct == nil {
		return ErrNilAccount
	}

	email := strings.TrimSpace(acct.Email)
	if len(email) < 1 || len(acct.PasswordHash) < 1 {
		return ErrNoRequiredParams
	}

	createdAt := acct.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	res, err := r.db.ExecContext(
		ctx,
		insertAccount,
		email,
		acct.PasswordHash,
		createdAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	acct.ID = id
	acct.Email = email
	acct.CreatedAt = createdAt.UTC()
	return nil
}

const selectAccountByEmail = `
SELECT id, email, password_hash, created_at
FROM accounts
WHERE email = ?
`

func (r *sqliteAccountRepository) GetByEmail(ctx context.Context, email string) (*account.Account, error) {
	email = strings.TrimSpace(email)
	if len(email) < 1 {
		return nil, nil
	}

	return r.scanAccount(r.db.QueryRowContext(ctx, selectAccountByEmail, email))
}

const selectAccountByID = `
SELECT id, email, password_hash, created_at
FROM accounts
WHERE id = ?
`

func (r *sqliteAccountRepository) GetByID(ctx context.Context, id int64) (*account.Account, error) {
	if id < 1 {
		return nil, nil
	}

	return r.scanAccount(r.db.QueryRowContext(ctx, selectAccountByID, id))
}

func (r *sqliteAccountRepository) scanAccount(row *sql.Row) (*account.Account, error) {
	var acct account.Account
	var createdAt string

	err := row.Scan(&acct.ID, &acct.Email, &acct.PasswordHash, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	parsed, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, err
	}

	acct.CreatedAt = parsed.UTC()
	return &acct, nil
}
