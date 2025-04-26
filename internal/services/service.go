package services

import (
	"backend-service/internal/storages"
	"github.com/rs/zerolog"
)

type Service struct {
	AuthService     AuthService
	UserRoleService UserRoleService
	ServiceService  ServiceService
}

type ServiceDeps struct {
	Log     zerolog.Logger
	Storage *storages.Storage
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		AuthService:     NewAuthService(deps.Storage.UserRepository),
		UserRoleService: NewUserRoleService(deps.Storage.UserRepository),
		ServiceService:  NewServiceService(deps.Storage.ServiceRepository),
	}
}
