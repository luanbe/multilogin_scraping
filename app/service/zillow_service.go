package service

import (
	"log"
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/repository"
)

type ZillowService interface {
	AddZillow(zillowData *entity.ZillowData)
}

type ZillowServiceImpl struct {
	//logger      log.Logger
	baseRepo   repository.BaseRepository
	zillowRepo repository.ZillowRepository
}

func NewZillowService(
	//lg log.Logger,
	baseRepo repository.BaseRepository,
	zillowRepo repository.ZillowRepository,
) ZillowService {
	return &ZillowServiceImpl{baseRepo, zillowRepo}
}

func (s *ZillowServiceImpl) AddZillow(zillowData *entity.ZillowData) {
	if err := s.zillowRepo.AddZillow(zillowData); err != nil {
		log.Fatalln(err)
	}
}
