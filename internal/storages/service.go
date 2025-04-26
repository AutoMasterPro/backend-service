package storages

import (
	"backend-service/internal/entity"
	"backend-service/pkg/database"
	"context"
	"github.com/google/uuid"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *entity.Service) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	GetAll(ctx context.Context) ([]*entity.Service, error)
	Update(ctx context.Context, service *entity.Service) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type serviceStorage struct {
	pg *database.PostgresDB
}

func NewServiceStorage(deps StorageDeps) ServiceRepository {
	return &serviceStorage{
		pg: deps.PostgresDB,
	}
}

func (s *serviceStorage) Create(ctx context.Context, service *entity.Service) (uuid.UUID, error) {
	if service.ID == uuid.Nil {
		service.ID = uuid.New()
	}

	const query = `
		INSERT INTO services (id, name, description, price, duration_min)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	row := s.pg.DB.QueryRowContext(ctx, query,
		service.ID, service.Name, service.Description, service.Price, service.DurationMin,
	)

	if err := row.Scan(&service.ID); err != nil {
		return uuid.Nil, err
	}

	return service.ID, nil
}

func (s *serviceStorage) GetById(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	const query = `
		SELECT id, name, description, price, duration_min
		FROM services
		WHERE id = $1 AND deleted_at IS NULL;
	`

	row := s.pg.DB.QueryRowContext(ctx, query, id)

	var service entity.Service
	if err := row.Scan(&service.ID, &service.Name, &service.Description, &service.Price, &service.DurationMin); err != nil {
		return nil, err
	}

	return &service, nil
}

func (s *serviceStorage) GetAll(ctx context.Context) ([]*entity.Service, error) {
	const query = `
		SELECT id, name, description, price, duration_min
		FROM services
		WHERE deleted_at IS NULL;
	`

	rows, err := s.pg.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*entity.Service
	for rows.Next() {
		var service entity.Service
		if err := rows.Scan(&service.ID, &service.Name, &service.Description, &service.Price, &service.DurationMin); err != nil {
			return nil, err
		}
		services = append(services, &service)
	}

	return services, nil
}

func (s *serviceStorage) Update(ctx context.Context, service *entity.Service) (uuid.UUID, error) {
	const query = `
		UPDATE services
		SET name = $2, description = $3, price = $4, duration_min = $5
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id;
	`

	row := s.pg.DB.QueryRowContext(ctx, query,
		service.ID, service.Name, service.Description, service.Price, service.DurationMin,
	)

	if err := row.Scan(&service.ID); err != nil {
		return uuid.Nil, err
	}

	return service.ID, nil
}

func (s *serviceStorage) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE services
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	_, err := s.pg.DB.ExecContext(ctx, query, id)
	return err
}
