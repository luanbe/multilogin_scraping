package repository

import "multilogin_scraping/app/models/entity"

type MovotoRepository interface {
	AddMovoto(movoto *entity.Movoto) error
	GetMovotoFirst(query map[string]interface{}) (*entity.Movoto, error)
	UpdateMovoto(movoto *entity.Movoto, query map[string]interface{}) error
}

type MovotoRepositoryImpl struct {
	base BaseRepository
}

func NewMovotoRepository(br BaseRepository) MovotoRepository {
	return &MovotoRepositoryImpl{br}
}

func (r *MovotoRepositoryImpl) AddMovoto(movoto *entity.Movoto) error {
	if err := r.base.GetDB().Create(movoto).Error; err != nil {
		return err
	}
	return nil
}

func (r *MovotoRepositoryImpl) UpdateMovoto(movoto *entity.Movoto, query map[string]interface{}) error {
	if err := r.base.GetDB().Where(query).Updates(movoto).Error; err != nil {
		return err
	}
	return nil
}

func (r *MovotoRepositoryImpl) GetMovotoFirst(query map[string]interface{}) (*entity.Movoto, error) {
	var movotoData = &entity.Movoto{}
	if err := r.base.GetDB().Where(query).First(movotoData).Error; err != nil {
		return nil, err
	}
	return movotoData, nil
}
