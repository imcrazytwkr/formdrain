package repositories

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/session"
)

type SessionRepository interface {
	Create(ctx context.Context, sess *session.Session) error
	GetByID(ctx context.Context, id string) (*session.Session, error)
	Delete(ctx context.Context, id string) error
}
