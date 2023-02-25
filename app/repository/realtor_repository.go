package repository

import "multilogin_scraping/app/models/entity"

type RealtorRepository interface {
	AddRealtor(realtor *entity.Realtor) error
	GetRealtorFirst(query map[string]interface{}) (*entity.Realtor, error)
	UpdateRealtor(realtor *entity.Realtor, query map[string]interface{}) error
}

type RealtorRepositoryImpl struct {
	base BaseRepository
}

func NewRealtorRepository(br BaseRepository) RealtorRepository {
	return &RealtorRepositoryImpl{br}
}

func (r *RealtorRepositoryImpl) AddRealtor(realtor *entity.Realtor) error {
	if err := r.base.GetDB().Create(realtor).Error; err != nil {
		return err
	}
	return nil
}

func (r *RealtorRepositoryImpl) UpdateRealtor(realtor *entity.Realtor, query map[string]interface{}) error {
	if err := r.base.GetDB().Where(query).Updates(realtor).Error; err != nil {
		return err
	}
	return nil
}

func (r *RealtorRepositoryImpl) GetRealtorFirst(query map[string]interface{}) (*entity.Realtor, error) {
	var realtorData = &entity.Realtor{}
	if err := r.base.GetDB().Where(query).First(realtorData).Error; err != nil {
		return nil, err
	}
	return realtorData, nil
}
