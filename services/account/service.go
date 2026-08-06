package account

import (
	"context"

	ma "github.com/imcrazytwkr/formdrain/models/account"
	"github.com/imcrazytwkr/formdrain/repositories"
	"github.com/imcrazytwkr/formdrain/services"
)

type service struct {
	accounts repositories.AccountRepository
}

func NewService(accounts repositories.AccountRepository) services.AccountService {
	return &service{accounts}
}

func (s *service) Login(ctx context.Context, email, password string) (*ma.Account, error) {
	account, err := s.accounts.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, ErrInvalidCredentials
	}

	err = CheckPassword(password, account.PasswordHash)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return account, nil
}
