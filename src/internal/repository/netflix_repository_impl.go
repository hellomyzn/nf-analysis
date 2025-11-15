package repository

import (
	"encoding/csv"
	"os"
)

type netflixrepositoryImpl struct {
}

func NewNetflixRepository() NetflixRepository {
	return &netflixrepositoryImpl{}

}

func (r *netflixrepositoryImpl) ReadRawCSV(path string) ([]RawNetflixRecord, error) {
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
