package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/prateek041/user-management-service/internal/data"
	"github.com/prateek041/user-management-service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUser(id string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
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
	return s.repo.Create(user)
}

func (s *DefaultUserService) GetUser(id string) (*model.User, error) {
	return s.repo.Get(id)
}

func (s *DefaultUserService) GetUserByUsername(username string) (*model.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *DefaultUserService) AuthenticateUser(username, password string) (*model.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	if err != nil {
		return nil, fmt.Errorf("Invalid credentials")
	}

	return user, nil
}
