package services

import (
	"backend-service/internal/storages"
	"context"
	"github.com/google/uuid"
)

type UserRoleService interface {
	IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error)
}

type userRoleService struct {
	userService storages.UserRepository
}

func NewUserRoleService(userService storages.UserRepository) UserRoleService {
	return &userRoleService{
		userService: userService,
	}
}

func (u *userRoleService) IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error) {
	user, err := u.userService.GetById(ctx, userId)
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}
