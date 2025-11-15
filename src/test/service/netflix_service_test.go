package service_test

import (
	"testing"

	"github.com/hellomyzn/nf-analysis/internal/model"
	"github.com/hellomyzn/nf-analysis/internal/repository"
	"github.com/hellomyzn/nf-analysis/internal/service"
)

// モック Repository
type mockRepo struct {
	rawRecords     []repository.RawNetflixRecord
	historyRecords []model.NetflixRecord
	savedPath      string
	savedRecords   []model.NetflixRecord
}

func (m *mockRepo) ReadRawCSV(path string) ([]repository.RawNetflixRecord, error) {
	return m.rawRecords, nil
}

func (m *mockRepo) SaveCSV(path string, records []model.NetflixRecord) error {
	m.savedPath = path
	m.savedRecords = records
	return nil
}

func (m *mockRepo) ReadHistory(path string) ([]model.NetflixRecord, error) {
	return m.historyRecords, nil
}

func Test_TransformRecords(t *testing.T) {
	mock := &mockRepo{
		rawRecords: []repository.RawNetflixRecord{
			{
				Title: "The Walking Dead: Season 5: Four Walls and a Roof",
				Date:  "11/14/25",
			},
			{
				Title: "One Piece: Egghead Arc ②: Episode 1149",
				Date:  "11/13/25",
			},
			{
				Title: "Breaking Bad: Season 1: Episode 1",
				Date:  "11/15/25",
			},
			{
				Title: "One Piece: Egghead Arc ②: Episode 1149",
				Date:  "11/13/25",
			},
		},
		historyRecords: []model.NetflixRecord{
			{
				ID:      "vid-0041",
				Title:   "The Walking Dead",
				Season:  "Season 5",
				Episode: "Four Walls and a Roof",
				Date:    "2025-11-14",
			},
			{
				ID:      "vid-0042",
				Title:   "Existing Show",
				Season:  "Season 1",
				Episode: "Episode 1",
				Date:    "2025-11-12",
			},
		},
	}
	srv := service.NewNetflixService(mock)

	records, err := srv.TransformRecords("dummy.csv", "history.csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	if records[0].Date != "2025-11-13" {
		t.Fatalf("expected first record date 2025-11-13, got %s", records[0].Date)
	}

	if records[0].ID != "vid-0043" {
		t.Fatalf("expected first record ID vid-0043, got %s", records[0].ID)
	}

	if records[0].Title != "One Piece" || records[0].Season != "Egghead Arc ②" || records[0].Episode != "Episode 1149" {
		t.Fatalf("unexpected first record content: %+v", records[0])
	}

	if records[1].Date != "2025-11-15" {
		t.Fatalf("expected second record date 2025-11-15, got %s", records[1].Date)
	}

	if records[1].ID != "vid-0044" {
		t.Fatalf("expected second record ID vid-0044, got %s", records[1].ID)
	}

}

func Test_SaveHistory_SimplyMergesExistingAndIncoming(t *testing.T) {
	mock := &mockRepo{
		historyRecords: []model.NetflixRecord{
			{
				ID:      "vid-0041",
				Title:   "Existing Episode",
				Season:  "Season 1",
				Episode: "Episode 1",
				Date:    "2025-11-15",
			},
		},
	}
	srv := service.NewNetflixService(mock)

	incoming := []model.NetflixRecord{
		{
			ID:      "vid-0042",
			Title:   "New Episode",
			Season:  "Season 1",
			Episode: "Episode 2",
			Date:    "2025-11-16",
		},
	}

	if err := srv.SaveHistory("history.csv", incoming); err != nil {
		t.Fatalf("SaveHistory returned error: %v", err)
	}

	if mock.savedPath != "history.csv" {
		t.Fatalf("expected SaveCSV to be invoked with history.csv, got %q", mock.savedPath)
	}

	if len(mock.savedRecords) != 2 {
		t.Fatalf("expected 2 records saved, got %d", len(mock.savedRecords))
	}

	if mock.savedRecords[0].ID != "vid-0041" {
		t.Fatalf("expected first record to remain existing entry, got %q", mock.savedRecords[0].ID)
	}

	if mock.savedRecords[1] != incoming[0] {
		t.Fatalf("expected incoming record to be appended without modification")
	}
}

func Test_SaveHistory_PersistsExistingWhenNoIncoming(t *testing.T) {
	mock := &mockRepo{
		historyRecords: []model.NetflixRecord{
			{
				ID:      "vid-0041",
				Title:   "Existing Episode",
				Season:  "Season 1",
				Episode: "Episode 1",
				Date:    "2025-11-15",
			},
		},
	}
	srv := service.NewNetflixService(mock)

	if err := srv.SaveHistory("history.csv", nil); err != nil {
		t.Fatalf("SaveHistory returned error: %v", err)
	}

	if len(mock.savedRecords) != 1 {
		t.Fatalf("expected existing record to be saved, got %d", len(mock.savedRecords))
	}

	if mock.savedRecords[0].ID != "vid-0041" {
		t.Fatalf("expected existing record to remain unchanged, got %q", mock.savedRecords[0].ID)
	}
}
