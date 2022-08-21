package base

import (
	util "github.com/luanbe/golang-web-app-structure/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type BaseIDModel struct {
	//gorm.Model
	// Id of as primary key
	// in: int64
	ID uint64 `gorm:"primary_key:auto_increment" json:"id"`

	IsDeleted bool      `json:"is_deleted"`
	CreatedBy string    `json:"-" gorm:"type:varchar"`
	CreatedAt time.Time `json:"-"`
	UpdatedBy string    `json:"-" gorm:"type:varchar"`
	UpdatedAt time.Time `json:"-"`
}

func (base *BaseIDModel) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("IsDeleted", false)
	tx.Statement.SetColumn("CreatedAt", time.Now().In(util.Loc).Format("2006-01-02 15:04:05"))
	tx.Statement.SetColumn("UpdatedAt", time.Now().In(util.Loc).Format("2006-01-02 15:04:05"))
	return nil
}

func (base *BaseIDModel) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now().In(util.Loc).Format("2006-01-02 15:04:05"))
	return nil
}
