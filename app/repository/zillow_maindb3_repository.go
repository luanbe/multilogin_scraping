package repository

import (
	"multilogin_scraping/app/models/entity"
)

type Maindb3Repository interface {
	GetMaindb3(id int) (*entity.ZillowMaindb3Address, error)
	CheckMaindb3SearchExist(maindb3 *entity.ZillowMaindb3Address, url, addressStreet, addressCity string) bool
	ListMaindb3(crawlingStatus string, page, pageSize int) ([]*entity.ZillowMaindb3Address, error)
	UpdateMaindb3(maindb3 *entity.ZillowMaindb3Address, values map[string]interface{}) error
	AddMaindb3(maindb3 *entity.ZillowMaindb3Address) error
	ListMaindb3Interval(day int, crawlingStatus string, page, pageSize int) ([]*entity.ZillowMaindb3Address, error)
}

type Maindb3RepositoryImpl struct {
	base BaseRepository
}

func NewMaindb3Repository(br BaseRepository) Maindb3Repository {
	return &Maindb3RepositoryImpl{br}
}

func (r Maindb3RepositoryImpl) ListMaindb3(crawlingStatus string, page, pageSize int) ([]*entity.ZillowMaindb3Address, error) {
	var items []*entity.ZillowMaindb3Address

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
func (r Maindb3RepositoryImpl) ListMaindb3Interval(day int, crawlingStatus string, page, pageSize int) ([]*entity.ZillowMaindb3Address, error) {
	var items []*entity.ZillowMaindb3Address

	if page == 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	db := r.base.GetDB()
	if err := db.Offset(offset).Limit(pageSize).
		Where("crawling_status = ? AND updated_at < NOW() - INTERVAL ? DAY", crawlingStatus, day).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
func (r Maindb3RepositoryImpl) GetMaindb3(id int) (*entity.ZillowMaindb3Address, error) {
	maindb3 := &entity.ZillowMaindb3Address{}

	query := r.base.GetDB().Model(&entity.ZillowMaindb3Address{})
	err := query.First(maindb3, id).Error
	if err != nil {
		return nil, err
	}
	return maindb3, nil
}

func (r Maindb3RepositoryImpl) UpdateMaindb3(maindb3 *entity.ZillowMaindb3Address, values map[string]interface{}) error {
	if err := r.base.GetDB().Model(maindb3).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

func (r Maindb3RepositoryImpl) AddMaindb3(maindb3 *entity.ZillowMaindb3Address) error {
	if err := r.base.GetDB().Create(maindb3).Error; err != nil {
		return err
	}
	return nil
}

func (r Maindb3RepositoryImpl) CheckMaindb3SearchExist(maindb3 *entity.ZillowMaindb3Address, url, addressStreet, addressCity string) bool {
	query := r.base.GetDB().
		Where("LOWER(Address_Street) = LOWER(?) AND LOWER(Address_City) = LOWER(?)", addressStreet, addressCity).
		Or("url = ?", url).
		First(maindb3)

	if query.RowsAffected > 0 {
		return true
	}
	return false

}
