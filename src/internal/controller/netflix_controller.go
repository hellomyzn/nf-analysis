package controller

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/hellomyzn/nf-analysis/internal/model"
)

// Service が満たすべきインターフェース
type NetflixService interface {
	TransformRecords(rawPath string, historyPath string) ([]model.NetflixRecord, error)
	SaveHistory(path string, records []model.NetflixRecord) error
}

type NetflixController struct {
	service NetflixService
}

func NewNetflixController(service NetflixService) *NetflixController {
	return &NetflixController{
		service: service,
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

	outputPath := "src/csv/history.csv"

	// Service で変換
	records, err := c.service.TransformRecords(inputFile, outputPath)
	if err != nil {
		return err
	}

	// 出力 CSV に保存
	return c.service.SaveHistory(outputPath, records)
}
