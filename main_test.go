package main

import (
	"os"
	"testing"
)

func TestIsGameReadyText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "game ready notification",
			text: "sender=:1.42 body=[\"Dota 2\", \"Your game is ready\"]",
			want: true,
		},
		{
			name: "game unpausing notification",
			text: "sender=:1.42 body=[\"Dota 2\", \"The game is unpausing...\"]",
			want: true,
		},
		{
			name: "unrelated notification",
			text: "sender=:1.42 body=[\"Dota 2\", \"Friend is now online\"]",
			want: false,
		},
		{
			name: "case sensitive mismatch",
			text: "sender=:1.42 body=[\"Dota 2\", \"your game is ready\"]",
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isGameReadyText(tt.text); got != tt.want {
				t.Fatalf("isGameReadyText(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestDetectDisplayServer(t *testing.T) {
	tests := []struct {
		name    string
		wayland string
		session string
		want    string
	}{
		{name: "wayland via WAYLAND_DISPLAY", wayland: "wayland-0", session: "", want: "wayland"},
		{name: "wayland via XDG_SESSION_TYPE", wayland: "", session: "wayland", want: "wayland"},
		{name: "x11 default", wayland: "", session: "", want: "x11"},
		{name: "x11 with session type x11", wayland: "", session: "x11", want: "x11"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("WAYLAND_DISPLAY")
			os.Unsetenv("XDG_SESSION_TYPE")

			if tt.wayland != "" {
				t.Setenv("WAYLAND_DISPLAY", tt.wayland)
			}
			if tt.session != "" {
				t.Setenv("XDG_SESSION_TYPE", tt.session)
			}

			if got := detectDisplayServer(); got != tt.want {
				t.Fatalf("detectDisplayServer() = %q, want %q", got, tt.want)
			}
		})
	}
}
