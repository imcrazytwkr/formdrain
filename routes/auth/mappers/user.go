package mappers

import m "github.com/imcrazytwkr/formdrain/models/account"

func User(account *m.Account) *m.User {
	return &m.User{
		ID:    account.ID,
		Email: account.Email,
	}
}
