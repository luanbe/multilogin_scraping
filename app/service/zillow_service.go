package service

import (
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/repository"
)

type ZillowService interface {
	AddZillow(zillowData *entity.ZillowData) error
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

func (s *ZillowServiceImpl) AddZillow(zillowData *entity.ZillowData) error {
	if err := s.zillowRepo.AddZillow(zillowData); err != nil {
		return err
	}
	return nil
}
