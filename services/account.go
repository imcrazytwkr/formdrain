package services

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/account"
)

type AccountService interface {
	Login(ctx context.Context, email, password string) (*account.Account, error)
}
