package entity

import "github.com/luanbe/golang-web-app-structure/app/models/base"

// TODO: Use swagger later
type User struct {
	base.BaseIDModel
	Email     string `gorm:"not null, uniqueIndex" json:"email"`
	Password  string `gorm:"not null" json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Gender    string `json:"gender"`
}
