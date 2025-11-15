package repository

import (
	"encoding/csv"
	"os"

	"github.com/hellomyzn/nf-analysis/internal/model"
)

type netflixRepositoryImpl struct {
}

func NewNetflixRepository() NetflixRepository {
	return &netflixRepositoryImpl{}

}

func (r *netflixRepositoryImpl) ReadRawCSV(path string) ([]RawNetflixRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var records []RawNetflixRecord

	// スキップ：ヘッダー

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}

		records = append(records, RawNetflixRecord{
			Title: row[0],
			Date:  row[1],
		})
	}

	return records, nil

}

func (r *netflixRepositoryImpl) SaveCSV(path string, records []model.NetflixRecord) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// header
	writer.Write([]string{"date", "title", "season", "episode"})

	for _, rec := range records {
		writer.Write([]string{
			rec.Date,
			rec.Title,
			rec.Season,
			rec.Episode,
		})
	}

	return nil
}
