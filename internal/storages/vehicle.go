package storages

import (
	"backend-service/internal/entity"
	"backend-service/pkg/database"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type VehicleRepository interface {
	Create(ctx context.Context, vehicle *entity.Vehicle) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*entity.Vehicle, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]*entity.Vehicle, error)
	GetAll(ctx context.Context) ([]*entity.Vehicle, error)
	Update(ctx context.Context, vehicle *entity.Vehicle) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type vehicleStorage struct {
	pg *database.PostgresDB
}

func NewVehicleStorage(deps StorageDeps) VehicleRepository {
	return &vehicleStorage{
		pg: deps.PostgresDB,
	}
}

func (s *vehicleStorage) Create(ctx context.Context, vehicle *entity.Vehicle) (uuid.UUID, error) {
	if vehicle.ID == uuid.Nil {
		vehicle.ID = uuid.New()
	}

	const query = `
		INSERT INTO vehicles (id, user_id, brand, model, license_plate, year, vin)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	row := s.pg.DB.QueryRowContext(ctx, query,
		vehicle.ID, vehicle.UserID, vehicle.Brand, vehicle.Model,
		vehicle.LicensePlate, vehicle.Year, vehicle.VIN,
	)

	if err := row.Scan(&vehicle.ID); err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert vehicle: %w", err)
	}

	return vehicle.ID, nil
}

func (s *vehicleStorage) GetById(ctx context.Context, id uuid.UUID) (*entity.Vehicle, error) {
	const query = `
		SELECT id, user_id, brand, model, license_plate, year, vin
		FROM vehicles
		WHERE id = $1 AND deleted_at IS NULL;
	`

	row := s.pg.DB.QueryRowContext(ctx, query, id)

	var vehicle entity.Vehicle
	if err := row.Scan(
		&vehicle.ID, &vehicle.UserID, &vehicle.Brand, &vehicle.Model,
		&vehicle.LicensePlate, &vehicle.Year, &vehicle.VIN,
	); err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	return &vehicle, nil
}

func (s *vehicleStorage) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*entity.Vehicle, error) {
	const query = `
		SELECT id, user_id, brand, model, license_plate, year, vin
		FROM vehicles
		WHERE user_id = $1 AND deleted_at IS NULL;
	`

	rows, err := s.pg.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*entity.Vehicle
	for rows.Next() {
		var vehicle entity.Vehicle
		if err := rows.Scan(
			&vehicle.ID, &vehicle.UserID, &vehicle.Brand, &vehicle.Model,
			&vehicle.LicensePlate, &vehicle.Year, &vehicle.VIN,
		); err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}

func (s *vehicleStorage) GetAll(ctx context.Context) ([]*entity.Vehicle, error) {
	const query = `
		SELECT id, user_id, brand, model, license_plate, year, vin
		FROM vehicles
		WHERE deleted_at IS NULL;
	`

	rows, err := s.pg.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*entity.Vehicle
	for rows.Next() {
		var vehicle entity.Vehicle
		if err := rows.Scan(
			&vehicle.ID, &vehicle.UserID, &vehicle.Brand, &vehicle.Model,
			&vehicle.LicensePlate, &vehicle.Year, &vehicle.VIN,
		); err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}

func (s *vehicleStorage) Update(ctx context.Context, vehicle *entity.Vehicle) error {
	const query = `
		UPDATE vehicles
		SET brand = $2, model = $3, license_plate = $4, year = $5, vin = $6
		WHERE id = $1 AND deleted_at IS NULL;
	`

	result, err := s.pg.DB.ExecContext(ctx, query,
		vehicle.ID, vehicle.Brand, vehicle.Model,
		vehicle.LicensePlate, vehicle.Year, vehicle.VIN,
	)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("vehicle not found")
	}

	return nil
}

func (s *vehicleStorage) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE vehicles
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	result, err := s.pg.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("vehicle not found")
	}

	return nil
}
