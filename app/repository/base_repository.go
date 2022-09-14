package repository

import "gorm.io/gorm"

type BaseRepository interface {
	GetDB() *gorm.DB
	BeginTx()
	CommitTx()
	RollbackTx()
}

type BaseRepositoryImpl struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return &BaseRepositoryImpl{db}
}

func (br *BaseRepositoryImpl) GetDB() *gorm.DB {
	return br.db
}

func (br *BaseRepositoryImpl) BeginTx() {
	br.db = br.GetDB().Begin()
}

func (br *BaseRepositoryImpl) CommitTx() {
	br.GetDB().Commit()
}

func (br *BaseRepositoryImpl) RollbackTx() {
	br.GetDB().Rollback()
}
