package controller_test

import (
	"os"
	"testing"

	"github.com/hellomyzn/nf-analysis/internal/controller"
	"github.com/hellomyzn/nf-analysis/internal/model"
)

// モック service とモック repository を用意する

type mockService struct {
	lastRawPath     string
	lastHistoryPath string
	savedPath       string
	savedRecords    []model.NetflixRecord
}

func (m *mockService) TransformRecords(rawPath string, historyPath string) ([]model.NetflixRecord, error) {
	m.lastRawPath = rawPath
	m.lastHistoryPath = historyPath
	return []model.NetflixRecord{
		{
			ID:    "mock-id",
			Title: "The Walking Dead: Season 5: Four Walls and a Roof",
			Date:  "2025-11-14",
		},
	}, nil
}

func (m *mockService) SaveHistory(path string, records []model.NetflixRecord) error {
	m.savedPath = path
	m.savedRecords = records
	return nil
}

func Test_Controller_Run(t *testing.T) {
	// テスト用の入力 CSV を置くディレクトリ
	os.MkdirAll("src/csv/netflix", 0755)
	os.WriteFile("src/csv/netflix/input.csv", []byte("dummy"), 0644)
	defer os.RemoveAll("src")

	mockService := &mockService{}

	c := controller.NewNetflixController(mockService)

	err := c.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mockService.lastRawPath != "src/csv/netflix/input.csv" {
		t.Errorf("unexpected raw path: %v", mockService.lastRawPath)
	}

	if mockService.lastHistoryPath != "src/csv/history.csv" {
		t.Errorf("unexpected history path: %v", mockService.lastHistoryPath)
	}

	// 出力先の確認
	if mockService.savedPath != "src/csv/history.csv" {
		t.Errorf("unexpected output path: %v", mockService.savedPath)
	}

	// 保存されたレコードの検証
	if len(mockService.savedRecords) != 1 {
		t.Fatalf("expected 1 record, got %d", len(mockService.savedRecords))
	}
}
