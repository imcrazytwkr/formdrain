package session

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	m "github.com/imcrazytwkr/formdrain/models/session"
	"github.com/imcrazytwkr/formdrain/repositories"
)

type memorySessionRepository struct {
	mu   sync.RWMutex
	byID map[string]*m.Session
}

func NewMemorySessionRepository() repositories.SessionRepository {
	return &memorySessionRepository{
		byID: make(map[string]*m.Session),
	}
}

func (r *memorySessionRepository) Create(_ context.Context, session *m.Session) error {
	if session == nil || session.AccountID < 1 {
		return ErrInvalidSession
	}

	// @NOTE: do not modify source object until insertion has succeeded
	s := *session
	if uuid.Validate(s.ID) != nil {
		s.ID = uuid.New().String()
	}

	r.mu.Lock()
	r.byID[s.ID] = &s
	r.mu.Unlock()

	session.ID = s.ID
	return nil
}

func (r *memorySessionRepository) GetByID(ctx context.Context, id string) (*m.Session, error) {
	if uuid.Validate(id) != nil {
		return nil, nil
	}

	r.mu.RLock()
	session, ok := r.byID[id]
	if !ok {
		r.mu.RUnlock()
		return nil, nil
	}

	if !session.ExpiresAt.IsZero() && !session.ExpiresAt.After(time.Now().UTC()) {
		r.mu.RUnlock()
		return nil, r.Delete(ctx, id)
	}

	s := *session
	r.mu.RUnlock()
	return &s, nil
}

func (r *memorySessionRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	delete(r.byID, id)
	r.mu.Unlock()

	return nil
}
