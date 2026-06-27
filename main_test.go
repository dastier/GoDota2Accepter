package main

import "testing"

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
