package controller

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/hellomyzn/nf-analysis/internal/model"
)

// Service が満たすべきインターフェース
type NetflixService interface {
	TransformRecords(path string) ([]model.NetflixRecord, error)
}

// Repository が満たすべきインターフェース
type NetflixRepository interface {
	SaveCSV(path string, records []model.NetflixRecord) error
}

type NetflixController struct {
	service NetflixService
	repo    NetflixRepository
}

func NewNetflixController(service NetflixService, repo NetflixRepository) *NetflixController {
	return &NetflixController{
		service: service,
		repo:    repo,
	}
}
func (c *NetflixController) Run() error {
	// 入力 CSV を探す
	inputDir := "src/csv/netflix"
	var inputFile string

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// ディレクトリではなく、CSV らしいファイルを探す
		if !info.IsDir() && filepath.Ext(info.Name()) == ".csv" {
			inputFile = path
		}
		return nil
	})
	if err != nil {
		return err
	}

	if inputFile == "" {
		return errors.New("no CSV file found in src/csv/netflix")
	}

	// Service で変換
	records, err := c.service.TransformRecords(inputFile)
	if err != nil {
		return err
	}

	// 出力 CSV に保存
	outputPath := "src/csv/netflix.csv"
	return c.repo.SaveCSV(outputPath, records)
}
