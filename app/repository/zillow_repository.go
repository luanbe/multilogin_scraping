package repository

import "multilogin_scraping/app/models/entity"

type ZillowRepository interface {
	AddZillow(Zillow *entity.ZillowData) error
}

type ZillowRepositoryImpl struct {
	base BaseRepository
}

func NewZillowRepository(br BaseRepository) ZillowRepository {
	return &ZillowRepositoryImpl{br}
}

// func (r *ZillowRepositoryImpl) GetZillow(Zillow *entity.ZillowDatam, id uint64) error {

// }

func (r *ZillowRepositoryImpl) AddZillow(Zillow *entity.ZillowData) error {
	if err := r.base.GetDB().Create(Zillow).Error; err != nil {
		return err
	}
	return nil
}
