package repositories

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/account"
)

type AccountRepository interface {
	Create(ctx context.Context, acct *account.Account) error
	GetByEmail(ctx context.Context, email string) (*account.Account, error)
	GetByID(ctx context.Context, id int64) (*account.Account, error)
}
