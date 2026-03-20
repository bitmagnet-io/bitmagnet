package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExtensionFromPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		wantExt   string
		wantValid bool
	}{
		{
			name:      "simple extension",
			path:      "short film.mkv",
			wantExt:   "mkv",
			wantValid: true,
		},
		{
			name:      "multi part extension",
			path:      "short film.mkv.mk2",
			wantExt:   "mk2",
			wantValid: true,
		},
		{
			name:      "double dot before extension",
			path:      "short film..mkv",
			wantExt:   "mkv",
			wantValid: true,
		},
		{
			name:      "triple dot before extension",
			path:      "short film...mkv",
			wantExt:   "mkv",
			wantValid: true,
		},
		{
			name:      "hidden file is not treated as an extension",
			path:      ".mkv",
			wantValid: false,
		},
		{
			name:      "hidden file in nested path is not treated as an extension",
			path:      "Movies/.mkv",
			wantValid: false,
		},
		{
			name:      "trailing dot is invalid",
			path:      "short film.",
			wantValid: false,
		},
		{
			name:      "extension with punctuation is invalid",
			path:      "short film.m-kv",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FileExtensionFromPath(tt.path)
			assert.Equal(t, tt.wantValid, got.Valid)
			assert.Equal(t, tt.wantExt, got.String)
		})
	}
}
