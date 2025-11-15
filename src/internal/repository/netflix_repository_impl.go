package repository

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

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

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	if _, err := writer.WriteString("id,date,title,season,episode\n"); err != nil {
		return err
	}

	for _, rec := range records {
		line := strings.Join([]string{
			formatCSVField(rec.ID, false),
			formatCSVField(rec.Date, false),
			formatCSVField(rec.Title, true),
			formatCSVField(rec.Season, true),
			formatCSVField(rec.Episode, true),
		}, ",") + "\n"

		if _, err := writer.WriteString(line); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func formatCSVField(value string, forceQuote bool) string {
	trimmed := strings.TrimSpace(value)
	escaped := strings.ReplaceAll(trimmed, "\"", "\"\"")

	needsQuote := forceQuote
	if strings.ContainsAny(escaped, ",\n\r") {
		needsQuote = true
	}

	if needsQuote {
		return fmt.Sprintf("\"%s\"", escaped)
	}

	return escaped
}

func (r *netflixRepositoryImpl) ReadHistory(path string) ([]model.NetflixRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var records []model.NetflixRecord
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		rec := model.NetflixRecord{}
		rec.ID = strings.TrimSpace(row[0])
		if len(row) > 1 {
			rec.Date = strings.TrimSpace(row[1])
		}
		if len(row) > 2 {
			rec.Title = strings.TrimSpace(row[2])
		}
		if len(row) > 3 {
			rec.Season = strings.TrimSpace(row[3])
		}
		if len(row) > 4 {
			rec.Episode = strings.TrimSpace(row[4])
		}

		records = append(records, rec)
	}

	return records, nil
}
