package entity

import (
	"multilogin_scraping/app/models/base"
)

// TODO: Use swagger later
type CrawlerLog struct {
	base.BaseIDModel
	CrawlerName string `gorm:"type:varchar(255)" json:"crawler_name"`
	Level       string `gorm:"type:varchar(100)" json:"type"`
	Status      bool   `json:"status"`
	Message     string `json:"message"`
	DataID      int    `json:"data_id"`
}
