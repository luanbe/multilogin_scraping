package entity

import "multilogin_scraping/app/models/base"

// TODO: Use swagger later
type User struct {
	base.BaseIDModel
	Email     string `gorm:"type:varchar(100); not null; uniqueIndex" json:"email"`
	Password  string `gorm:"not null" json:"password"`
	FirstName string `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string `gorm:"type:varchar(100)" json:"last_name"`
	Phone     string `gorm:"type:varchar(150)" json:"phone"`
	Address   string `gorm:"type:varchar(255)" json:"address"`
}
