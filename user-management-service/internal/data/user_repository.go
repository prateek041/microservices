package data

import (
	"fmt"
	"sync"

	"github.com/prateek041/user-management-service/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	Get(id string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
}

type InMemoryUserRepository struct {
	users map[string]*model.User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*model.User),
	}
}

func (r *InMemoryUserRepository) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user with ID '%s' already exists", user.ID)
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Get(id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID '%s' not found", id)
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByUsername(username string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with username '%s' not found", username)
}
