package services

import (
	"backend-service/internal/entity"
	"backend-service/internal/storages"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, input entity.UserRegister) (uuid.UUID, error)
	Login(ctx context.Context, input entity.UserLogin) (*entity.User, error)
}

type authService struct {
	userRepo storages.UserRepository
}

func NewAuthService(userRepo storages.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, input entity.UserRegister) (uuid.UUID, error) {
	user := input.UserRegisterToUser()
	user.ID = uuid.New()

	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *authService) Login(ctx context.Context, input entity.UserLogin) (*entity.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if !user.CheckPasswordHash(input.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}
