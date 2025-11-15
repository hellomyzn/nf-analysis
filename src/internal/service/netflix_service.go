package service

import (
	"sort"
	"strconv"
	"strings"

	"github.com/hellomyzn/nf-analysis/internal/model"
	"github.com/hellomyzn/nf-analysis/internal/repository"
	"github.com/hellomyzn/nf-analysis/internal/util"
)

type NetflixService struct {
	repo repository.NetflixRepository
}

func NewNetflixService(repo repository.NetflixRepository) *NetflixService {
	return &NetflixService{
		repo: repo,
	}
}

func (s *NetflixService) TransformRecords(path string) ([]model.NetflixRecord, error) {
	rawRecords, err := s.repo.ReadRawCSV(path)
	if err != nil {
		return nil, err
	}

	var results []model.NetflixRecord

	for _, r := range rawRecords {
		date, err := util.ConvertDate(r.Date)
		if err != nil {
			return nil, err
		}

		title, season, episode := util.SplitTitle(r.Title)
		results = append(results, model.NetflixRecord{
			Title:   title,
			Season:  season,
			Episode: episode,
			Date:    date,
		})
	}

	return results, nil
}

func (s *NetflixService) SaveHistory(path string, incoming []model.NetflixRecord) error {
	existing, err := s.repo.ReadHistory(path)
	if err != nil {
		return err
	}

	existingByID := make(map[string]struct{}, len(existing))
	existingBySignature := make(map[string]struct{}, len(existing))

	for i := range existing {
		existing[i] = normalizeRecord(existing[i])
		if existing[i].ID != "" {
			existingByID[existing[i].ID] = struct{}{}
		}
		sig := recordSignature(existing[i])
		if sig != "" {
			existingBySignature[sig] = struct{}{}
		}
	}

	var toAppend []model.NetflixRecord
	for _, rec := range incoming {
		normalized := normalizeRecord(rec)

		if normalized.ID != "" {
			if _, ok := existingByID[normalized.ID]; ok {
				continue
			}
		}

		sig := recordSignature(normalized)
		if sig != "" {
			if _, ok := existingBySignature[sig]; ok {
				continue
			}
		}

		toAppend = append(toAppend, normalized)
		if normalized.ID != "" {
			existingByID[normalized.ID] = struct{}{}
		}
		if sig != "" {
			existingBySignature[sig] = struct{}{}
		}
	}

	combined := append(existing, toAppend...)

	sort.SliceStable(combined, func(i, j int) bool {
		if combined[i].Date == combined[j].Date {
			return combined[i].Title < combined[j].Title
		}
		return combined[i].Date < combined[j].Date
	})

	nextID := newIDGenerator(existing)
	for i := range combined {
		if combined[i].ID == "" {
			combined[i].ID = nextID()
		}
	}

	return s.repo.SaveCSV(path, combined)
}

func normalizeRecord(rec model.NetflixRecord) model.NetflixRecord {
	rec.ID = strings.TrimSpace(rec.ID)
	rec.Date = strings.TrimSpace(rec.Date)
	rec.Title = strings.TrimSpace(rec.Title)
	rec.Season = strings.TrimSpace(rec.Season)
	rec.Episode = strings.TrimSpace(rec.Episode)
	return rec
}

func recordSignature(rec model.NetflixRecord) string {
	parts := []string{rec.Date, rec.Title, rec.Season, rec.Episode}

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	joined := strings.Join(parts, "\x1f")
	if strings.Trim(joined, "\x1f") == "" {
		return ""
	}
	return joined
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

func newIDGenerator(existing []model.NetflixRecord) func() string {
	var (
		next   = 1
		prefix string
		width  int
	)

	for _, rec := range existing {
		id := strings.TrimSpace(rec.ID)
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
