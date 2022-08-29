package service

import (
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/repository"
)

type Maindb3Service interface {
	ListMaindb3Data(crawlingStatus string, limit int) ([]*entity.Maindb3, error)
	UpdateStatus(maindb3 *entity.Maindb3, status string) error
}

type Maindb3ServiceImpl struct {
	//logger      log.Logger
	baseRepo    repository.BaseRepository
	maindb3Repo repository.Maindb3Repository
}

func NewMaindb3Service(
	//lg log.Logger,
	baseRepo repository.BaseRepository,
	maindb3Repo repository.Maindb3Repository,
) Maindb3Service {
	return &Maindb3ServiceImpl{baseRepo, maindb3Repo}
}

func (s *Maindb3ServiceImpl) ListMaindb3Data(crawlingStatus string, limit int) ([]*entity.Maindb3, error) {
	//s.baseRepo.BeginTx()
	result, err := s.maindb3Repo.ListMaindb3(crawlingStatus, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Maindb3ServiceImpl) UpdateStatus(maindb3 *entity.Maindb3, status string) error {
	if err := s.maindb3Repo.UpdateMaindb3(maindb3, map[string]interface{}{"crawling_status": status}); err != nil {
		return err
	}
	return nil
}
