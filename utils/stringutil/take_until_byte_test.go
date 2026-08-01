package stringutil_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/utils/stringutil"
)

func TestTakeUntilByte(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		c    byte
		want string
	}{
		{in: "a,b,c", c: ',', want: "a"},
		{in: "abc", c: ',', want: "abc"},
		{in: ",abc", c: ',', want: ""},
		{in: "", c: ',', want: ""},
		{in: "application/json; charset=utf-8", c: ';', want: "application/json"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			if got := stringutil.TakeUntilByte(tt.in, tt.c); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
