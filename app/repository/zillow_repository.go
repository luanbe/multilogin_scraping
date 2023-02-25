package repository

import "multilogin_scraping/app/models/entity"

type ZillowRepository interface {
	AddZillow(Zillow *entity.Zillow) error
	GetZillowFirst(query map[string]interface{}) (*entity.Zillow, error)
	UpdateZillow(Zillow *entity.Zillow, query map[string]interface{}) error
	GetZillowPriceHistoryFirst(query map[string]interface{}) (*entity.ZillowPriceHistory, error)
	AddZillowPriceHistory(zillowHistory *entity.ZillowPriceHistory) error
	GetZillowPublicTaxHistoryFirst(query map[string]interface{}) (*entity.ZillowPublicTaxHistory, error)
	AddZillowPublicTaxHistory(zillowPubicTaxHistory *entity.ZillowPublicTaxHistory) error
}

type ZillowRepositoryImpl struct {
	base BaseRepository
}

func NewZillowRepository(br BaseRepository) ZillowRepository {
	return &ZillowRepositoryImpl{br}
}

// func (r *ZillowRepositoryImpl) GetZillowByID(Zillow *entity.ZillowDatam, id uint64) error {

// }

func (r *ZillowRepositoryImpl) AddZillow(zillow *entity.Zillow) error {
	if err := r.base.GetDB().Create(zillow).Error; err != nil {
		return err
	}
	return nil
}

func (r *ZillowRepositoryImpl) UpdateZillow(zillow *entity.Zillow, query map[string]interface{}) error {
	if err := r.base.GetDB().Where(query).Updates(zillow).Error; err != nil {
		return err
	}
	return nil
}

func (r *ZillowRepositoryImpl) GetZillowFirst(query map[string]interface{}) (*entity.Zillow, error) {
	var zillowData = &entity.Zillow{}
	if err := r.base.GetDB().Where(query).First(zillowData).Error; err != nil {
		return nil, err
	}
	return zillowData, nil
}

func (r *ZillowRepositoryImpl) GetZillowPriceHistoryFirst(query map[string]interface{}) (*entity.ZillowPriceHistory, error) {
	var zillowPriceHistory = &entity.ZillowPriceHistory{}
	if err := r.base.GetDB().Where(query).First(zillowPriceHistory).Error; err != nil {
		return nil, err
	}
	return zillowPriceHistory, nil
}
func (r *ZillowRepositoryImpl) AddZillowPriceHistory(zillowPriceHistory *entity.ZillowPriceHistory) error {
	if err := r.base.GetDB().Create(zillowPriceHistory).Error; err != nil {
		return err
	}
	return nil
}

func (r *ZillowRepositoryImpl) GetZillowPublicTaxHistoryFirst(query map[string]interface{}) (*entity.ZillowPublicTaxHistory, error) {
	var zillowPublicTaxHistory = &entity.ZillowPublicTaxHistory{}
	if err := r.base.GetDB().Where(query).First(zillowPublicTaxHistory).Error; err != nil {
		return nil, err
	}
	return zillowPublicTaxHistory, nil
}

func (r *ZillowRepositoryImpl) AddZillowPublicTaxHistory(zillowPubicTaxHistory *entity.ZillowPublicTaxHistory) error {
	if err := r.base.GetDB().Create(zillowPubicTaxHistory).Error; err != nil {
		return err
	}
	return nil
}
