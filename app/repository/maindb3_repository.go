package repository

import "multilogin_scraping/app/models/entity"

type Maindb3Repository interface {
	ListMaindb3(crawlingStatus string, limit int) ([]*entity.Maindb3, error)
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
