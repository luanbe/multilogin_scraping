package repository

import "github.com/luanbe/golang-web-app-structure/app/models/entity"

type UserRepository interface {
	AddUser(user *entity.User) error
	GetUser(User *entity.User, email string) error
}

type UserRepositoryImpl struct {
	base BaseRepository
}

func NewUserRepository(br BaseRepository) UserRepository {
	return &UserRepositoryImpl{br}
}

func (r *UserRepositoryImpl) AddUser(User *entity.User) error {
	if err := r.base.GetDB().Create(User).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) GetUser(User *entity.User, email string) error {
	if err := r.base.GetDB().Where("email = ?", email).First(User).Error; err != nil {
		return err
	}
	return nil
}
