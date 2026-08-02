package notification

import (
	"errors"
	"testing"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
	"github.com/imcrazytwkr/formdrain/models/form_config/sendinblue"
)

type stubDiscord struct {
	calls int
	err   error
}

func (s *stubDiscord) Send(_ *discord.DiscordConfig, _ map[string]any) error {
	s.calls++
	return s.err
}

type stubSendinblue struct {
	calls int
	err   error
}

func (s *stubSendinblue) Send(_ *sendinblue.SendinblueConfig, _ map[string]any) error {
	s.calls++
	return s.err
}

func TestSend_NilConfigs(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{}
	s := &stubSendinblue{}
	svc := &httpNotificationService{discordNotifier: d, sendinblueNotifier: s}

	err := svc.Send(fc.NotifiersConfig{}, map[string]any{"email": "a@b.c"})
	if err != nil {
		t.Fatal(err)
	}
	if d.calls != 0 || s.calls != 0 {
		t.Fatalf("calls discord=%d sendinblue=%d", d.calls, s.calls)
	}
}

func TestSend_BothNotifiers(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{}
	s := &stubSendinblue{}
	svc := &httpNotificationService{discordNotifier: d, sendinblueNotifier: s}

	err := svc.Send(fc.NotifiersConfig{
		Discord:    &discord.DiscordConfig{},
		Sendinblue: &sendinblue.SendinblueConfig{},
	}, map[string]any{"email": "a@b.c"})
	if err != nil {
		t.Fatal(err)
	}
	if d.calls != 1 || s.calls != 1 {
		t.Fatalf("calls discord=%d sendinblue=%d", d.calls, s.calls)
	}
}

func TestSend_JoinsErrors(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{err: errors.New("discord failed")}
	s := &stubSendinblue{err: errors.New("sendinblue failed")}
	svc := &httpNotificationService{discordNotifier: d, sendinblueNotifier: s}

	err := svc.Send(fc.NotifiersConfig{
		Discord:    &discord.DiscordConfig{},
		Sendinblue: &sendinblue.SendinblueConfig{},
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, d.err) || !errors.Is(err, s.err) {
		t.Fatalf("joined = %v", err)
	}
}

func TestSend_OneNotifierOnly(t *testing.T) {
	t.Parallel()

	d := &stubDiscord{}
	s := &stubSendinblue{}
	svc := &httpNotificationService{discordNotifier: d, sendinblueNotifier: s}

	err := svc.Send(fc.NotifiersConfig{Discord: &discord.DiscordConfig{}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if d.calls != 1 || s.calls != 0 {
		t.Fatalf("calls discord=%d sendinblue=%d", d.calls, s.calls)
	}
}
