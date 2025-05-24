package services

import (
	"backend-service/internal/storages"
	"github.com/rs/zerolog"
)

type Service struct {
	AuthService        AuthService
	UserRoleService    UserRoleService
	ServiceService     ServiceService
	VehicleService     VehicleService
	AppointmentService AppointmentService
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
		VehicleService:  NewVehicleService(deps.Storage.VehicleRepository),
		AppointmentService: NewAppointmentService(
			deps.Storage.AppointmentRepository,
			deps.Storage.VehicleRepository,
			deps.Storage.ServiceRepository,
		),
	}
}
