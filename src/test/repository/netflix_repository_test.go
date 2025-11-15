package repository_test

import (
	"os"
	"testing"

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
