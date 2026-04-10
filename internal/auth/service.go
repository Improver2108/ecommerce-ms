package auth

import (
	"context"
	"errors"
	"fmt"

	"example.com/ecommerce/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *service {
	return &service{
		queries: queries,
	}
}

func (s *service) Register(req RegisterRequest) (*AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash passoword")
	}

	user, err := s.queries.CreateUser(context.Background(), db.CreateUserParams{
		Name:        req.Name,
		Email:       req.Email,
		PasswordHsh: string(hash),
		Phone:       pgtype.Text{Valid: false},
	})

	if err != nil {
		return nil, errors.New("email already registered or db error")
	}

	token, err := GenerateTokens(fmt.Sprintf("%d", user.ID), user.Email)

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

	user, err := s.queries.GetUserByEmail(context.Background(), req.Email)

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHsh), []byte(req.Password))

	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	token, err := GenerateTokens(fmt.Sprintf("%d", user.ID), user.Email)

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  UserData{ID: user.ID, Name: user.Name, Email: user.Email},
	}, nil
}
