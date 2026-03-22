package user

import (
	"errors"
)

type Service interface {
	GetById(id string) (*User, error)
}

type service struct {
	users map[string]User
}

func NewService() *service {
	return &service{
		users: make(map[string]User),
	}
}

func (s *service) GetById(id string) (*User, error) {
	u, userFound := s.users[id]

	if !userFound {
		return nil, errors.New("No user found")
	}
	return &u, nil
}
