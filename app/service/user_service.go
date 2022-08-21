package service

import (
	"github.com/luanbe/golang-web-app-structure/app/models/entity"
	"github.com/luanbe/golang-web-app-structure/app/repository"
)

type UserService interface {
	AddUser(username, email string) (*entity.User, error)
	GetUser(email string) (*entity.User, error)
}

type UserServiceImpl struct {
	//logger      log.Logger
	baseRepo repository.BaseRepository
	userRepo repository.UserRepository
}

func NewUserService(
	//lg log.Logger,
	baseRepo repository.BaseRepository,
	userRepo repository.UserRepository,
) UserService {
	return &UserServiceImpl{baseRepo, userRepo}
}

func (s *UserServiceImpl) AddUser(email, password string) (*entity.User, error) {
	s.baseRepo.BeginTx()
	User := &entity.User{
		Email:    email,
		Password: password,
	}
	err := s.userRepo.AddUser(User)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (s *UserServiceImpl) GetUser(email string) (*entity.User, error) {
	s.baseRepo.BeginTx()
	User := &entity.User{}
	err := s.userRepo.GetUser(User, email)
	if err != nil {
		return nil, err
	}
	return User, nil
}
