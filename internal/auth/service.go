package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type storedUser struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type service struct {
	users map[string]storedUser
}

func NewService() *service {
	return &service{
		users: make(map[string]storedUser),
	}
}

func (s *service) Register(req RegisterRequest) (*AuthResponse, error) {
	for _, u := range s.users {
		if u.Email == req.Email {
			return nil, errors.New("Email already registered")
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash passoword")
	}

	user := storedUser{
		ID:           uuid.NewString(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	s.users[user.ID] = user

	token, err := GenerateTokens(user.ID, user.Email)

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User: UserData{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *service) Login(req LoginRequest) (*AuthResponse, error) {
	var found *storedUser
	for _, user := range s.users {
		if user.Email == req.Email {
			found = &user
			break
		}
	}
	if found == nil {
		return nil, errors.New("invalid username or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(found.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	token, err := GenerateTokens(found.ID, found.Email)

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  UserData{ID: found.ID, Name: found.Name, Email: found.Email},
	}, nil
}
