package store

import (
	"auth_test/internal/service"
	"fmt"
)

type InMemoryUserStore struct {
	Users map[string]*service.User
}

func (s *InMemoryUserStore) Get(username string) (*service.User, error) {
	user, exists := s.Users[username]
	if !exists {
		return nil, fmt.Errorf("пользователь не существует")
	}
	return user, nil
}

func (s *InMemoryUserStore) Add(user *service.User) (bool, error) {
	if _, exists := s.Users[user.Username]; exists {
		return false, fmt.Errorf("пользователь уже существует")
	}
	s.Users[user.Username] = user
	return true, nil
}
