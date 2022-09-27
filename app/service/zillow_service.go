package service

import (
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/repository"
)

type ZillowService interface {
	AddZillow(zillowData *entity.ZillowDetail) error
	UpdateZillow(zillowData *entity.ZillowDetail, id uint64) error
	GetZillowByID(id uint64) (*entity.ZillowDetail, error)
	GetZillowByURL(url string) (*entity.ZillowDetail, error)
	UpdateZillowPriceHistory(zillowPriceHistories []*entity.ZillowPriceHistory) error
	UpdateZillowPublicTaxHistory(zillowPublicTaxHistories []*entity.ZillowPublicTaxHistory) error
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

func (s *ZillowServiceImpl) AddZillow(zillowData *entity.ZillowDetail) error {
	if err := s.zillowRepo.AddZillow(zillowData); err != nil {
		return err
	}
	return nil
}

func (s *ZillowServiceImpl) UpdateZillow(zillowData *entity.ZillowDetail, id uint64) error {
	if err := s.zillowRepo.UpdateZillow(zillowData, map[string]interface{}{"id": id}); err != nil {
		return err
	}
	return nil
}

func (s *ZillowServiceImpl) GetZillowByID(id uint64) (*entity.ZillowDetail, error) {
	zillowData, err := s.zillowRepo.GetZillowFirst(map[string]interface{}{"maindb3_id": id})
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	return zillowData, nil
}

func (s *ZillowServiceImpl) GetZillowByURL(url string) (*entity.ZillowDetail, error) {
	zillowData, err := s.zillowRepo.GetZillowFirst(map[string]interface{}{"url": url})
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	return zillowData, nil
}

func (s *ZillowServiceImpl) UpdateZillowPriceHistory(zillowPriceHistories []*entity.ZillowPriceHistory) error {
	for _, history := range zillowPriceHistories {
		historyData, err := s.zillowRepo.GetZillowPriceHistoryFirst(map[string]interface{}{"maindb3_id": history.Maindb3ID, "date": history.Date})
		if err != nil && err.Error() != "record not found" {
			return err
		}
		if historyData == nil {
			if err := s.zillowRepo.AddZillowPriceHistory(history); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *ZillowServiceImpl) UpdateZillowPublicTaxHistory(zillowPublicTaxHistories []*entity.ZillowPublicTaxHistory) error {
	for _, history := range zillowPublicTaxHistories {
		historyData, err := s.zillowRepo.GetZillowPublicTaxHistoryFirst(map[string]interface{}{"maindb3_id": history.Maindb3ID, "year": history.Year})
		if err != nil && err.Error() != "record not found" {
			return err
		}
		if historyData == nil {
			if err := s.zillowRepo.AddZillowPublicTaxHistory(history); err != nil {
				return err
			}
		}
	}
	return nil
}
