package account

import "time"

type Account struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func (a *Account) AsUser() *User {
	return &User{
		ID:    a.ID,
		Email: a.Email,
	}
}
