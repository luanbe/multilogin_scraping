package service

import (
	"log"
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/repository"
)

type Maindb3Service interface {
	ListMaindb3Data(crawlingStatus string, limit int) []*entity.Maindb3
	UpdateStatus(maindb3 *entity.Maindb3, status string)
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

func (s *Maindb3ServiceImpl) ListMaindb3Data(crawlingStatus string, limit int) []*entity.Maindb3 {
	//s.baseRepo.BeginTx()
	result, err := s.maindb3Repo.ListMaindb3(crawlingStatus, limit)
	if err != nil {
		log.Fatalln(err)
	}
	return result
}

func (s *Maindb3ServiceImpl) UpdateStatus(maindb3 *entity.Maindb3, status string) {
	if err := s.maindb3Repo.UpdateMaindb3(maindb3, map[string]interface{}{"crawling_status": status}); err != nil {
		log.Fatalln(err)
	}
}
