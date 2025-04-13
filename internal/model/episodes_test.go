package model

import "testing"

func TestEpisodesString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		episodes Episodes
		want     string
	}{
		{
			name:     "empty",
			episodes: make(Episodes),
			want:     "",
		},
		{
			name: "single whole season",
			episodes: Episodes{
				1: {},
			},
			want: "S01",
		},
		{
			name: "range of whole seasons",
			episodes: Episodes{
				1: {},
				2: {},
				3: {},
			},
			want: "S01-03",
		},
		{
			name: "single season with episodes",
			episodes: Episodes{
				1: {
					1: {},
					2: {},
				},
			},
			want: "S01E01-02",
		},
		{
			name: "multiple seasons",
			episodes: Episodes{
				1: {
					1: {},
					2: {},
				},
				2: {},
			},
			want: "S01E01-02, S02",
		},
		{
			name: "mixed bag",
			episodes: Episodes{
				1: {},
				2: {},
				3: {},
				5: {
					1: {},
					2: {},
					4: {},
				},
				6: {},
				7: {},
				9: {},
			},
			want: "S01-03, S05E01-02,E04, S06-07, S09",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.episodes.String(); got != tt.want {
				t.Errorf("Episodes.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
