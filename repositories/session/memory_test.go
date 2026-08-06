package session_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/imcrazytwkr/formdrain/models/session"
	sr "github.com/imcrazytwkr/formdrain/repositories/session"
)

const testSessionId = "sess-1"
const expiredSessionId = "expired"

func TestMemorySessionRepository(t *testing.T) {
	ctx := t.Context()
	repo := sr.NewMemorySessionRepository()

	s := &session.Session{
		// ID:        testSessionId,
		AccountID: 42,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}

	err := repo.Create(ctx, s)
	if err != nil {
		t.Fatal(err)
	}

	id := s.ID
	err = uuid.Validate(id)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.AccountID != 42 {
		t.Fatalf("got = %#v", got)
	}

	err = repo.Delete(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	got, err = repo.GetByID(ctx, testSessionId)
	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Fatalf("expected nil after delete, got %#v", got)
	}
}

func TestMemorySessionRepository_Expired(t *testing.T) {
	ctx := t.Context()
	repo := sr.NewMemorySessionRepository()
	now := time.Now().UTC()

	s := &session.Session{
		AccountID: 1,
		ExpiresAt: now.Add(-time.Hour),
	}

	err := repo.Create(ctx, s)
	if err != nil {
		t.Fatal(err)
	}

	id := s.ID
	err = uuid.Validate(id)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Fatalf("expected nil for expired, got %#v", got)
	}
}
