package controller_test

import (
	"os"
	"testing"

	"github.com/hellomyzn/nf-analysis/internal/controller"
	"github.com/hellomyzn/nf-analysis/internal/model"
	"github.com/hellomyzn/nf-analysis/internal/repository"
)

// モック service とモック repository を用意する

type mockService struct{}

func (m *mockService) TransformRecords(path string) ([]model.NetflixRecord, error) {
	return []model.NetflixRecord{
		{
			ID:      "mock-id",
			Title:   "The Walking Dead",
			Season:  "Season 5",
			Episode: "Four Walls and a Roof",
			Date:    "2025-11-14",
		},
	}, nil
}

type mockRepo struct {
	savedPath    string
	savedRecords []model.NetflixRecord
}

func (m *mockRepo) SaveCSV(path string, records []model.NetflixRecord) error {
	m.savedPath = path
	m.savedRecords = records
	return nil
}

// このインターフェースは既存の ReadRawCSV 用。
// Controller では SaveCSV のみ必要。
func (m *mockRepo) ReadRawCSV(path string) ([]repository.RawNetflixRecord, error) {
	return nil, nil
}

func Test_Controller_Run(t *testing.T) {
	// テスト用の入力 CSV を置くディレクトリ
	os.MkdirAll("src/csv/netflix", 0755)
	os.WriteFile("src/csv/netflix/input.csv", []byte("dummy"), 0644)
	defer os.RemoveAll("src")

	mockRepo := &mockRepo{}
	mockService := &mockService{}

	c := controller.NewNetflixController(mockService, mockRepo)

	err := c.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 出力先の確認
	if mockRepo.savedPath != "src/csv/history.csv" {
		t.Errorf("unexpected output path: %v", mockRepo.savedPath)
	}

	// 保存されたレコードの検証
	if len(mockRepo.savedRecords) != 1 {
		t.Fatalf("expected 1 record, got %d", len(mockRepo.savedRecords))
	}
}
