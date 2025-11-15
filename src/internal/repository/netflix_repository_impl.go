package repository

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
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
	sorted := make([]model.NetflixRecord, len(records))
	copy(sorted, records)

	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Date == sorted[j].Date {
			return sorted[i].Title < sorted[j].Title
		}
		return sorted[i].Date < sorted[j].Date
	})

	nextID := r.newIDGenerator(path)
	for i := range sorted {
		if sorted[i].ID == "" {
			sorted[i].ID = nextID()
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// header
	writer.Write([]string{"id", "date", "title", "season", "episode"})

	for _, rec := range sorted {
		writer.Write([]string{
			rec.ID,
			rec.Date,
			rec.Title,
			rec.Season,
			rec.Episode,
		})
	}

	return nil
}

func (r *netflixRepositoryImpl) newIDGenerator(path string) func() string {
	var (
		next   = 1
		prefix string
		width  int
	)

	if f, err := os.Open(path); err == nil {
		defer f.Close()

		reader := csv.NewReader(f)
		if rows, err := reader.ReadAll(); err == nil {
			for i := 1; i < len(rows); i++ {
				row := rows[i]
				if len(row) == 0 {
					continue
				}

				id := strings.TrimSpace(row[0])
				if id == "" {
					continue
				}

				if p, w, value, ok := parseIDComponents(id); ok {
					if value >= next {
						prefix = p
						width = w
						next = value + 1
					}
					continue
				}

				if value, err := strconv.Atoi(id); err == nil {
					if value >= next {
						prefix = ""
						width = 0
						next = value + 1
					}
				}
			}
		}
	}

	return func() string {
		value := next
		next++

		digits := strconv.Itoa(value)
		if width > 0 {
			if len(digits) < width {
				digits = strings.Repeat("0", width-len(digits)) + digits
			}
			return prefix + digits
		}

		return prefix + digits
	}
}

func parseIDComponents(id string) (prefix string, width int, value int, ok bool) {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return "", 0, 0, false
	}

	idx := len(trimmed) - 1
	for idx >= 0 && trimmed[idx] >= '0' && trimmed[idx] <= '9' {
		idx--
	}
	idx++

	if idx >= len(trimmed) {
		return "", 0, 0, false
	}

	digits := trimmed[idx:]
	value, err := strconv.Atoi(digits)
	if err != nil {
		return "", 0, 0, false
	}

	return trimmed[:idx], len(digits), value, true
}
