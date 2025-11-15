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

func Test_SaveCSV_WritesRecordsInGivenOrder(t *testing.T) {
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

	expectedLines := []string{
		"id,date,title,season,episode",
		"vid-2,2025-11-13,\"Later Episode\",\"Season 1\",\"Episode 2\"",
		"vid-1,2025-11-14,\"Earlier Episode\",\"Season 1\",\"Episode 1\"",
	}

	for i, line := range lines {
		if line != expectedLines[i] {
			t.Fatalf("line %d = %q, want %q", i, line, expectedLines[i])
		}
	}
}

func Test_ReadHistory_ReturnsRecords(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "history-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	contents := strings.Join([]string{
		"id,date,title,season,episode",
		"vid-1,2025-11-13,\"Later Episode\",\"Season 1\",\"Episode 2\"",
		"vid-2,2025-11-14,\"Earlier Episode\",\"Season 1\",\"Episode 1\"",
	}, "\n")

	if err := os.WriteFile(tmpFile.Name(), []byte(contents), 0644); err != nil {
		t.Fatalf("failed to seed temp file: %v", err)
	}

	repo := repository.NewNetflixRepository()
	records, err := repo.ReadHistory(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadHistory returned error: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	if records[0].ID != "vid-1" || records[1].Title != "Earlier Episode" {
		t.Fatalf("unexpected records: %+v", records)
	}
}
