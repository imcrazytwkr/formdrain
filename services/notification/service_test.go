package notification

import (
	"errors"
	"testing"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/models/form_config/brevo"
	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
)

type stubDiscord struct {
	calls int
	err   error
}

func (s *stubDiscord) Send(_ *discord.DiscordConfig, _ map[string]any) error {
	s.calls++
	return s.err
}

type stubBrevo struct {
	calls int
	err   error
}

func (s *stubBrevo) Send(_ *brevo.BrevoConfig, _ map[string]any) error {
	s.calls++
	return s.err
}

func TestSend_NilConfigs(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{}
	b := &stubBrevo{}
	svc := &httpNotificationService{discordNotifier: d, brevoNotifier: b}

	err := svc.Send(fc.NotifiersConfig{}, map[string]any{"email": "a@b.c"})
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

	err := svc.Send(fc.NotifiersConfig{
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

	err := svc.Send(fc.NotifiersConfig{
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

	err := svc.Send(fc.NotifiersConfig{Discord: &discord.DiscordConfig{}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if d.calls != 1 || b.calls != 0 {
		t.Fatalf("calls discord=%d brevo=%d", d.calls, b.calls)
	}
}
