package storages

import (
	"backend-service/pkg/database"
	"github.com/rs/zerolog"
)

type Storage struct {
	UserRepository        UserRepository
	ServiceRepository     ServiceRepository
	VehicleRepository     VehicleRepository
	AppointmentRepository AppointmentRepository
}

type StorageDeps struct {
	PostgresDB *database.PostgresDB
	Log        zerolog.Logger
}

func NewStorage(deps StorageDeps) *Storage {
	return &Storage{
		UserRepository:        NewUserStorage(deps),
		ServiceRepository:     NewServiceStorage(deps),
		VehicleRepository:     NewVehicleStorage(deps),
		AppointmentRepository: NewAppointmentStorage(deps),
	}
}
