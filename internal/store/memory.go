package store

import (
	"auth_test/internal/service"
)

type InMemoryUserStore struct {
	users map[string]*service.User
}

func (s *InMemoryUserStore) Get(username string) (*service.User, error) {
	user, exists := s.users[username]
	if !exists {
		return nil, service.ErrUserNotFound
	}
	return user, nil
}

func (s *InMemoryUserStore) Add(user *service.User) (bool, error) {
	if _, exists := s.users[user.Username]; exists {
		return false, service.ErrUserExists
	}
	s.users[user.Username] = user
	return true, nil
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: map[string]*service.User{
			"Никита": &service.User{Username: "Никита", Password: "password123"},
			"Антон":  &service.User{Username: "Антон", Password: "super123pro"},
			"Адольф": &service.User{Username: "Адольф", Password: "always171wet"},
		},
	}
}
