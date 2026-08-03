package notification

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/models/form_config/brevo"
	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
)

type stubDiscord struct {
	mu     sync.Mutex
	calls  int
	err    error
	ctxs   []context.Context
	done   chan struct{}
	onSend func(context.Context)
}

func (s *stubDiscord) Send(ctx context.Context, _ *discord.DiscordConfig, _ map[string]any) error {
	if s.onSend != nil {
		s.onSend(ctx)
	}

	s.mu.Lock()
	s.calls++
	s.ctxs = append(s.ctxs, ctx)
	err := s.err
	s.mu.Unlock()

	if s.done != nil {
		select {
		case s.done <- struct{}{}:
		default:
		}
	}
	return err
}

type stubBrevo struct {
	mu    sync.Mutex
	calls int
	err   error
}

func (s *stubBrevo) Send(_ context.Context, _ *brevo.BrevoConfig, _ map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	return s.err
}

func TestSend_NilConfigs(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{}
	b := &stubBrevo{}
	svc := &httpNotificationService{discordNotifier: d, brevoNotifier: b}

	err := svc.Send(t.Context(), fc.NotifiersConfig{}, map[string]any{"email": "a@b.c"})
	if err != nil {
		t.Fatal(err)
	}
	if d.calls != 0 || b.calls != 0 {
		t.Fatalf("calls discord=%d brevo=%d", d.calls, b.calls)
	}
}

func TestSend_BothNotifiers(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{}
	b := &stubBrevo{}
	svc := &httpNotificationService{discordNotifier: d, brevoNotifier: b}

	err := svc.Send(t.Context(), fc.NotifiersConfig{
		Discord: &discord.DiscordConfig{},
		Brevo:   &brevo.BrevoConfig{},
	}, map[string]any{"email": "a@b.c"})
	if err != nil {
		t.Fatal(err)
	}
	if d.calls != 1 || b.calls != 1 {
		t.Fatalf("calls discord=%d brevo=%d", d.calls, b.calls)
	}
}

func TestSend_JoinsErrors(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{err: errors.New("discord failed")}
	b := &stubBrevo{err: errors.New("brevo failed")}
	svc := &httpNotificationService{discordNotifier: d, brevoNotifier: b}

	err := svc.Send(t.Context(), fc.NotifiersConfig{
		Discord: &discord.DiscordConfig{},
		Brevo:   &brevo.BrevoConfig{},
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, d.err) || !errors.Is(err, b.err) {
		t.Fatalf("joined = %v", err)
	}
}

func TestSend_OneNotifierOnly(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{}
	b := &stubBrevo{}
	svc := &httpNotificationService{discordNotifier: d, brevoNotifier: b}

	err := svc.Send(t.Context(), fc.NotifiersConfig{Discord: &discord.DiscordConfig{}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if d.calls != 1 || b.calls != 0 {
		t.Fatalf("calls discord=%d brevo=%d", d.calls, b.calls)
	}
}

func TestSendAsync_IgnoresParentCancel(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{
		done: make(chan struct{}, 1),
		onSend: func(ctx context.Context) {
			if err := ctx.Err(); err != nil {
				t.Errorf("notify ctx err during Send = %v", err)
			}
			deadline, ok := ctx.Deadline()
			if !ok {
				t.Error("expected notify ctx deadline")
				return
			}
			remaining := time.Until(deadline)
			if remaining < notificationTimeout-time.Second || remaining > notificationTimeout {
				t.Errorf("deadline remaining = %v, want ~%v", remaining, notificationTimeout)
			}
		},
	}
	b := &stubBrevo{}
	svc := &httpNotificationService{discordNotifier: d, brevoNotifier: b}

	parent, cancel := context.WithCancel(t.Context())
	cancel()

	svc.SendAsync(parent, fc.NotifiersConfig{Discord: &discord.DiscordConfig{}}, nil)

	select {
	case <-d.done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for async send")
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	if d.calls != 1 {
		t.Fatalf("calls = %d", d.calls)
	}
}
