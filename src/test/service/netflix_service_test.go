package service_test

import (
	"testing"

	"github.com/hellomyzn/nf-analysis/internal/repository"
	"github.com/hellomyzn/nf-analysis/internal/service"
)

// モック Repository
type mockRepo struct{}

func (m *mockRepo) ReadRawCSV(path string) ([]repository.RawNetflixRecord, error) {
	return []repository.RawNetflixRecord{
		{
			Title: "The Walking Dead: Season 5: Four Walls and a Roof",
			Date:  "11/14/25",
		},
		{
			Title: "One Piece: Egghead Arc ②: Episode 1149",
			Date:  "11/13/25",
		},
	}, nil
}

func Test_TransformRecords(t *testing.T) {
	srv := service.NewNetflixService(&mockRepo{})

	records, err := srv.TransformRecords("dummy.csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	// 1件目の検証
	r1 := records[0]
	if r1.Title != "The Walking Dead" {
		t.Errorf("Title mismatch: %v", r1.Title)
	}
	if r1.Season != "Season 5" {
		t.Errorf("Season mismatch: %v", r1.Season)
	}
	if r1.Episode != "Four Walls and a Roof" {
		t.Errorf("Episode mismatch: %v", r1.Episode)
	}
	if r1.Date != "2025-11-14" {
		t.Errorf("Date mismatch: %v", r1.Date)
	}
}
