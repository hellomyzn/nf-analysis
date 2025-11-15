package service

import (
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
