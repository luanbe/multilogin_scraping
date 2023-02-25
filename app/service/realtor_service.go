package service

import (
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/repository"
)

type RealtorService interface {
	AddRealtor(realtorData *entity.Realtor) error
	UpdateRealtor(realtorData *entity.Realtor, id uint64) error
	GetRealtorByID(id uint64) (*entity.Realtor, error)
	GetRealtorByURL(url string) (*entity.Realtor, error)
}

type RealtorServiceImpl struct {
	//logger      log.Logger
	baseRepo    repository.BaseRepository
	realtorRepo repository.RealtorRepository
}

func NewRealtorService(
	//lg log.Logger,
	baseRepo repository.BaseRepository,
	realtorRepo repository.RealtorRepository,
) RealtorService {
	return &RealtorServiceImpl{baseRepo, realtorRepo}
}

func (s *RealtorServiceImpl) AddRealtor(realtorData *entity.Realtor) error {
	oldRealtorData, _ := s.realtorRepo.GetRealtorFirst(map[string]interface{}{"url": realtorData.URL})
	if oldRealtorData == nil {
		if err := s.realtorRepo.AddRealtor(realtorData); err != nil {
			return err
		}
	}
	return nil
}

func (s *RealtorServiceImpl) UpdateRealtor(realtorData *entity.Realtor, id uint64) error {
	if err := s.realtorRepo.UpdateRealtor(realtorData, map[string]interface{}{"id": id}); err != nil {
		return err
	}
	return nil
}

func (s *RealtorServiceImpl) GetRealtorByID(id uint64) (*entity.Realtor, error) {
	realtorData, err := s.realtorRepo.GetRealtorFirst(map[string]interface{}{"maindb3_id": id})
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	return realtorData, nil
}

func (s *RealtorServiceImpl) GetRealtorByURL(url string) (*entity.Realtor, error) {
	realtorData, err := s.realtorRepo.GetRealtorFirst(map[string]interface{}{"url": url})
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	return realtorData, nil
}
