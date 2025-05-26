package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/prateek041/user-management-service/internal/data"
	"github.com/prateek041/user-management-service/internal/model"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUser(id string) (*model.User, error)
	AuthenticateUser(username, password string) (*model.User, error)
}

type DefaultUserService struct {
	repo data.UserRepository
}

func NewDefaultUserService(repo data.UserRepository) *DefaultUserService {
	return &DefaultUserService{repo: repo}
}

func (s *DefaultUserService) CreateUser(user *model.User) error {
	user.ID = uuid.New().String() // Generate a unique ID
	// In a real application, you would hash the password before saving
	return s.repo.Create(user)
}

func (s *DefaultUserService) GetUser(id string) (*model.User, error) {
	return s.repo.Get(id)
}

func (s *DefaultUserService) AuthenticateUser(username, password string) (*model.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	// In a real application, you would compare the provided password with the hashed password
	if user.Password != password {
		return nil, fmt.Errorf("invalid credentials")
	}
	return user, nil
}
