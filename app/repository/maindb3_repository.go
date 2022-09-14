package repository

import (
	"multilogin_scraping/app/models/entity"
)

type Maindb3Repository interface {
	ListMaindb3(crawlingStatus string, limit int) ([]*entity.Maindb3, error)
	UpdateMaindb3(maindb3 *entity.Maindb3, values map[string]interface{}) error
	GetMaindb3(id int) (*entity.Maindb3, error)
	ListMaindb3Interval(day int, crawlingStatus string, limit int) ([]*entity.Maindb3, error)
}

type Maindb3RepositoryImpl struct {
	base BaseRepository
}

func NewMaindb3Repository(br BaseRepository) Maindb3Repository {
	return &Maindb3RepositoryImpl{br}
}

func (r Maindb3RepositoryImpl) ListMaindb3(crawlingStatus string, limit int) ([]*entity.Maindb3, error) {
	var items []*entity.Maindb3
	query := r.base.GetDB().Model(&entity.Maindb3{})
	err := query.Limit(limit).Where("crawling_status != ?", crawlingStatus).Or("crawling_status is NULL").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
func (r Maindb3RepositoryImpl) ListMaindb3Interval(day int, crawlingStatus string, limit int) ([]*entity.Maindb3, error) {
	var items []*entity.Maindb3
	query := r.base.GetDB().Model(&entity.Maindb3{})
	err := query.Limit(limit).Where("crawling_status = ? AND updated_at < NOW() - INTERVAL ? DAY", crawlingStatus, day).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
func (r Maindb3RepositoryImpl) GetMaindb3(id int) (*entity.Maindb3, error) {
	maindb3 := &entity.Maindb3{}
	query := r.base.GetDB().Model(&entity.Maindb3{})
	err := query.First(maindb3, id).Error
	if err != nil {
		return nil, err
	}
	return maindb3, nil
}

func (r Maindb3RepositoryImpl) UpdateMaindb3(maindb3 *entity.Maindb3, values map[string]interface{}) error {
	if err := r.base.GetDB().Model(maindb3).Updates(values).Error; err != nil {
		return err
	}
	return nil
}
