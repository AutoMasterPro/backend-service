package services

import (
	"backend-service/internal/entity"
	"backend-service/internal/storages"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type VehicleService interface {
	Create(ctx context.Context, userID uuid.UUID, input *entity.VehicleCreate) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*entity.Vehicle, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]*entity.Vehicle, error)
	Update(ctx context.Context, id uuid.UUID, input *entity.VehicleUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type vehicleService struct {
	repo storages.VehicleRepository
}

func NewVehicleService(repo storages.VehicleRepository) VehicleService {
	return &vehicleService{
		repo: repo,
	}
}

func (s *vehicleService) Create(ctx context.Context, userID uuid.UUID, input *entity.VehicleCreate) (uuid.UUID, error) {
	if err := input.Validate(); err != nil {
		return uuid.Nil, fmt.Errorf("validation error: %w", err)
	}

	vehicle := input.ToVehicle(userID)
	return s.repo.Create(ctx, vehicle)
}

func (s *vehicleService) GetById(ctx context.Context, id uuid.UUID) (*entity.Vehicle, error) {
	return s.repo.GetById(ctx, id)
}

func (s *vehicleService) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*entity.Vehicle, error) {
	return s.repo.GetByUserId(ctx, userId)
}

func (s *vehicleService) Update(ctx context.Context, id uuid.UUID, input *entity.VehicleUpdate) error {
	if err := input.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	vehicle := input.ToVehicle(id)
	return s.repo.Update(ctx, vehicle)
}

func (s *vehicleService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
