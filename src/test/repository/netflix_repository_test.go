package repository_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hellomyzn/nf-analysis/internal/model"
	"github.com/hellomyzn/nf-analysis/internal/repository"
)

const sampleCSV = `Title,Date
The Walking Dead: Season 5: Four Walls and a Roof,11/14/25
One Piece: Egghead Arc ②: Episode 1149,11/13/25
`

func Test_ReadRawCSV(t *testing.T) {
	testFile := "test_raw.csv"

	// テスト用ファイルを作成
	err := os.WriteFile(testFile, []byte(sampleCSV), 0644)
	if err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}
	defer os.Remove(testFile)

	repo := repository.NewNetflixRepository()

	rows, err := repo.ReadRawCSV(testFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}

	if rows[0].Title != "The Walking Dead: Season 5: Four Walls and a Roof" {
		t.Errorf("unexpected Title: %v", rows[0].Title)
	}
	if rows[0].Date != "11/14/25" {
		t.Errorf("unexpected Date: %v", rows[0].Date)
	}
}

func Test_SaveCSV_SortsByDateDesc(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "history-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	records := []model.NetflixRecord{
		{
			ID:      "vid-2",
			Title:   "Later Episode",
			Season:  "Season 1",
			Episode: "Episode 2",
			Date:    "2025-11-13",
		},
		{
			ID:      "vid-1",
			Title:   "Earlier Episode",
			Season:  "Season 1",
			Episode: "Episode 1",
			Date:    "2025-11-14",
		},
	}

	repo := repository.NewNetflixRepository()
	if err := repo.SaveCSV(tmpFile.Name(), records); err != nil {
		t.Fatalf("SaveCSV returned error: %v", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read temp file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("unexpected number of lines: %d", len(lines))
	}

	expectedOrder := []string{
		"id,date,title,season,episode",
		"vid-1,2025-11-14,Earlier Episode,Season 1,Episode 1",
		"vid-2,2025-11-13,Later Episode,Season 1,Episode 2",
	}

	for i, line := range lines {
		if line != expectedOrder[i] {
			t.Fatalf("line %d = %q, want %q", i, line, expectedOrder[i])
		}
	}
}

func Test_SaveCSV_AssignsIDsWhenMissing(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "history-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	records := []model.NetflixRecord{
		{
			Title:   "Newest Episode",
			Season:  "Season 1",
			Episode: "Episode 2",
			Date:    "2025-11-14",
		},
		{
			Title:   "Older Episode",
			Season:  "Season 1",
			Episode: "Episode 1",
			Date:    "2025-11-13",
		},
	}

	repo := repository.NewNetflixRepository()
	if err := repo.SaveCSV(tmpFile.Name(), records); err != nil {
		t.Fatalf("SaveCSV returned error: %v", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read temp file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("unexpected number of lines: %d", len(lines))
	}

	if !strings.HasPrefix(lines[1], "1,") {
		t.Fatalf("expected first record to start with ID 1, got %q", lines[1])
	}
	if !strings.HasPrefix(lines[2], "2,") {
		t.Fatalf("expected second record to start with ID 2, got %q", lines[2])
	}
}

func Test_SaveCSV_AssignsIDsUsingExistingHistory(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "history-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	existing := strings.Join([]string{
		"id,date,title,season,episode",
		"vid-0041,2025-11-15,Existing Episode,Season 1,Episode 1",
		"vid-0042,2025-11-14,Existing Episode 2,Season 1,Episode 2",
	}, "\n")

	if err := os.WriteFile(tmpFile.Name(), []byte(existing), 0644); err != nil {
		t.Fatalf("failed to seed temp file: %v", err)
	}

	records := []model.NetflixRecord{
		{
			Title:   "Newest Episode",
			Season:  "Season 2",
			Episode: "Episode 1",
			Date:    "2025-11-16",
		},
		{
			Title:   "Older Episode",
			Season:  "Season 2",
			Episode: "Episode 2",
			Date:    "2025-11-10",
		},
	}

	repo := repository.NewNetflixRepository()
	if err := repo.SaveCSV(tmpFile.Name(), records); err != nil {
		t.Fatalf("SaveCSV returned error: %v", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read temp file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("unexpected number of lines: %d", len(lines))
	}

	if !strings.HasPrefix(lines[1], "vid-0043,") {
		t.Fatalf("expected first record to have ID vid-0043, got %q", lines[1])
	}
	if !strings.HasPrefix(lines[2], "vid-0044,") {
		t.Fatalf("expected second record to have ID vid-0044, got %q", lines[2])
	}
}
