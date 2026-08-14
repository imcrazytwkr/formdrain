package config

import "testing"

func TestLogMode_RoundTrip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		mode LogMode
		text string
	}{
		{LogModeTrace, "trace"},
		{LogModeRelease, "release"},
	}

	for _, tc := range cases {
		if tc.mode.String() != tc.text {
			t.Fatalf("String(%v) = %q", tc.mode, tc.mode.String())
		}

		text, err := tc.mode.MarshalText()
		if err != nil || string(text) != tc.text {
			t.Fatalf("%s: MarshalText = %q err=%v", tc.text, text, err)
		}

		var got LogMode
		err = got.UnmarshalText([]byte(tc.text))
		if err != nil || got != tc.mode {
			t.Fatalf("%s: UnmarshalText got %v err=%v", tc.text, got, err)
		}
	}
}

func TestLogMode_UnmarshalEmpty(t *testing.T) {
	t.Parallel()

	var m LogMode = LogModeRelease
	err := m.UnmarshalText([]byte(""))
	if err != nil || m != LogModeTrace {
		t.Fatalf("got %v err=%v", m, err)
	}
}

func TestLogMode_UnmarshalInvalid(t *testing.T) {
	t.Parallel()

	var m LogMode
	err := m.UnmarshalText([]byte("debug"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLogMode_MarshalInvalid(t *testing.T) {
	t.Parallel()

	_, err := LogMode(99).MarshalText()
	if err == nil {
		t.Fatal("expected error")
	}
}
