package account

import "time"

type Account struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
