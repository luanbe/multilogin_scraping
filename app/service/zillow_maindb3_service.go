package service

import (
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/repository"
)

type Maindb3Service interface {
	ListMaindb3Data(crawlingStatus string, page, limit int) ([]*entity.ZillowMaindb3Address, error)
	UpdateStatus(maindb3 *entity.ZillowMaindb3Address, status string) error
	GetMaindb3(id int) (*entity.ZillowMaindb3Address, error)
	ListMaindb3IntervalData(day int, crawlingStatus string, page, limit int) ([]*entity.ZillowMaindb3Address, error)
	AddMaindb3SearchData(url, addressStreet, addressCity string) (bool, *entity.ZillowMaindb3Address, error)
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

func (s *Maindb3ServiceImpl) ListMaindb3Data(crawlingStatus string, page, limit int) ([]*entity.ZillowMaindb3Address, error) {
	//s.baseRepo.BeginTx()
	result, err := s.maindb3Repo.ListMaindb3(crawlingStatus, page, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *Maindb3ServiceImpl) ListMaindb3IntervalData(day int, crawlingStatus string, page, limit int) ([]*entity.ZillowMaindb3Address, error) {
	//s.baseRepo.BeginTx()
	result, err := s.maindb3Repo.ListMaindb3Interval(day, crawlingStatus, page, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *Maindb3ServiceImpl) GetMaindb3(id int) (*entity.ZillowMaindb3Address, error) {
	result, err := s.maindb3Repo.GetMaindb3(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Maindb3ServiceImpl) UpdateStatus(maindb3 *entity.ZillowMaindb3Address, status string) error {
	if err := s.maindb3Repo.UpdateMaindb3(maindb3, map[string]interface{}{"crawling_status": status}); err != nil {
		return err
	}
	return nil
}
func (s *Maindb3ServiceImpl) AddMaindb3SearchData(url, addressStreet, addressCity string) (bool, *entity.ZillowMaindb3Address, error) {
	maindb3 := &entity.ZillowMaindb3Address{}
	status := s.maindb3Repo.CheckMaindb3SearchExist(maindb3, url, addressStreet, addressCity)
	if status == true {
		return false, maindb3, nil
	}

	maindb3.URL = url
	maindb3.AddressStreet = addressStreet
	maindb3.AddressCity = addressCity
	if err := s.maindb3Repo.AddMaindb3(maindb3); err != nil {
		return false, nil, err
	}

	return true, maindb3, nil
}
