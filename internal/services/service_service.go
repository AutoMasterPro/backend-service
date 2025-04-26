package services

import (
	"backend-service/internal/entity"
	"backend-service/internal/storages"
	"context"
	"github.com/google/uuid"
)

type ServiceService interface {
	Create(ctx context.Context, service *entity.Service) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	GetAll(ctx context.Context) ([]*entity.Service, error)
	Update(ctx context.Context, service *entity.Service) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type serviceService struct {
	repo storages.ServiceRepository
}

func NewServiceService(repo storages.ServiceRepository) ServiceService {
	return &serviceService{
		repo: repo,
	}
}

func (s *serviceService) Create(ctx context.Context, service *entity.Service) (uuid.UUID, error) {
	service.ID = uuid.New()
	_, err := s.repo.Create(ctx, service)
	if err != nil {
		return uuid.Nil, err
	}

	return service.ID, nil
}

func (s *serviceService) GetById(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	return s.repo.GetById(ctx, id)
}

func (s *serviceService) GetAll(ctx context.Context) ([]*entity.Service, error) {
	return s.repo.GetAll(ctx)
}

func (s *serviceService) Update(ctx context.Context, service *entity.Service) (uuid.UUID, error) {
	return s.repo.Update(ctx, service)
}

func (s *serviceService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
