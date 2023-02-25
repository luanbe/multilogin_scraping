package service

import (
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/repository"
)

type MovotoService interface {
	AddMovoto(movotoData *entity.Movoto) error
	UpdateMovoto(movotoData *entity.Movoto, id uint64) error
	GetMovotoByID(id uint64) (*entity.Movoto, error)
	GetMovotoByURL(url string) (*entity.Movoto, error)
}

type MovotoServiceImpl struct {
	//logger      log.Logger
	baseRepo   repository.BaseRepository
	movotoRepo repository.MovotoRepository
}

func NewMovotoService(
	//lg log.Logger,
	baseRepo repository.BaseRepository,
	movotoRepo repository.MovotoRepository,
) MovotoService {
	return &MovotoServiceImpl{baseRepo, movotoRepo}
}

func (s *MovotoServiceImpl) AddMovoto(movotoData *entity.Movoto) error {
	oldMovotoData, _ := s.movotoRepo.GetMovotoFirst(map[string]interface{}{"url": movotoData.URL})
	if oldMovotoData == nil {
		if err := s.movotoRepo.AddMovoto(movotoData); err != nil {
			return err
		}
	}
	return nil
}

func (s *MovotoServiceImpl) UpdateMovoto(movotoData *entity.Movoto, id uint64) error {
	if err := s.movotoRepo.UpdateMovoto(movotoData, map[string]interface{}{"id": id}); err != nil {
		return err
	}
	return nil
}

func (s *MovotoServiceImpl) GetMovotoByID(id uint64) (*entity.Movoto, error) {
	movotoData, err := s.movotoRepo.GetMovotoFirst(map[string]interface{}{"maindb3_id": id})
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	return movotoData, nil
}

func (s *MovotoServiceImpl) GetMovotoByURL(url string) (*entity.Movoto, error) {
	movotoData, err := s.movotoRepo.GetMovotoFirst(map[string]interface{}{"url": url})
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	return movotoData, nil
}
