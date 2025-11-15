package repository

import "github.com/hellomyzn/nf-analysis/internal/model"

type RawNetflixRecord struct {
	Title string
	Date  string
}

type NetflixRepository interface {
        ReadRawCSV(path string) ([]RawNetflixRecord, error)
        ReadHistory(path string) ([]model.NetflixRecord, error)
        SaveCSV(path string, records []model.NetflixRecord) error
}
