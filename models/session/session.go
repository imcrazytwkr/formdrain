package session

import (
	"time"
)

type Session struct {
	ID        string
	AccountID int64
	ExpiresAt time.Time
}
