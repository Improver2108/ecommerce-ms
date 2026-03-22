package user

import (
	"errors"
)

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

func (s *service) UpdateUser(userId string, req UpdateUserRequest) (*User, error) {
	user, ok := s.users[userId]

	if !ok {
		return nil, errors.New("No user found")
	}

	user.Name = req.Name
	user.Email = req.Name

	return &user, nil
}
