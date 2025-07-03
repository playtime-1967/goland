package services

import (
	"grain/entities"
	"grain/repositories"
)

type UserService interface {
	Create(name, email string) (entities.User, error)
}

type userService struct {
	repo repositories.UserRepo
}

func NewUserService(repo repositories.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(name, email string) (entities.User, error) {
	user := entities.NewUser(name, email)
	err := s.repo.Create(user)
	return user, err
}
