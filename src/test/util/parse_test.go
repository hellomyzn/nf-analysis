package util_test

import (
	"testing"

	"github.com/hellomyzn/nf-analysis/internal/util"
)

func Test_ConvertDate(t *testing.T) {
	input := "11/14/25"
	want := "2025-11-14"

	got, err := util.ConvertDate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_SplitTitle(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		wantTitle   string
		wantSeason  string
		wantEpisode string
	}{
		{
			name:        "The Walking Dead S5",
			input:       "The Walking Dead: Season 5: Four Walls and a Roof",
			wantTitle:   "The Walking Dead",
			wantSeason:  "Season 5",
			wantEpisode: "Four Walls and a Roof",
		},
		{
			name:        "One Piece Egghead",
			input:       "One Piece: Egghead Arc ②: Episode 1149",
			wantTitle:   "One Piece",
			wantSeason:  "Egghead Arc ②",
			wantEpisode: "Episode 1149",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			title, season, episode := util.SplitTitle(tt.input)

			if title != tt.wantTitle {
				t.Errorf("title = %q, want %q", title, tt.wantTitle)
			}
			if season != tt.wantSeason {
				t.Errorf("season = %q, want %q", season, tt.wantSeason)
			}
			if episode != tt.wantEpisode {
				t.Errorf("episode = %q, want %q", episode, tt.wantEpisode)
			}
		})
	}
}
