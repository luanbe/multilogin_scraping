package repository

import (
	"multilogin_scraping/app/models/entity"
)

type Maindb3Repository interface {
	ListMaindb3(crawlingStatus string, page, pageSize int) ([]*entity.Maindb3, error)
	UpdateMaindb3(maindb3 *entity.Maindb3, values map[string]interface{}) error
	GetMaindb3(id int) (*entity.Maindb3, error)
	ListMaindb3Interval(day int, crawlingStatus string, page, pageSize int) ([]*entity.Maindb3, error)
}

type Maindb3RepositoryImpl struct {
	base BaseRepository
}

func NewMaindb3Repository(br BaseRepository) Maindb3Repository {
	return &Maindb3RepositoryImpl{br}
}

func (r Maindb3RepositoryImpl) ListMaindb3(crawlingStatus string, page, pageSize int) ([]*entity.Maindb3, error) {
	var items []*entity.Maindb3

	if page == 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	db := r.base.GetDB()
	query := db.Offset(offset).Limit(pageSize)
	if crawlingStatus != "" {
		query = query.Where("crawling_status = ?", crawlingStatus)
	} else {
		query = query.Where("crawling_status is NULL")
	}
	err := query.Find(&items).Error

	if err != nil {
		return nil, err
	}
	return items, nil
}
func (r Maindb3RepositoryImpl) ListMaindb3Interval(day int, crawlingStatus string, page, pageSize int) ([]*entity.Maindb3, error) {
	var items []*entity.Maindb3

	if page == 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	db := r.base.GetDB()
	err := db.Offset(offset).Limit(pageSize).Where("crawling_status = ? AND updated_at < NOW() - INTERVAL ? DAY", crawlingStatus, day).Find(&items).Error
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
