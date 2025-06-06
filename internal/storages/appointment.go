package storages

import (
	"backend-service/internal/entity"
	"backend-service/pkg/database"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *entity.Appointment, serviceIDs []uuid.UUID) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*entity.Appointment, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]*entity.Appointment, error)
	Update(ctx context.Context, appointment *entity.Appointment) error
	UpdateServices(ctx context.Context, appointmentID uuid.UUID, serviceIDs []uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	CheckTimeSlotAvailable(ctx context.Context, appointmentTime string) (bool, error)
}

type appointmentStorage struct {
	pg *database.PostgresDB
}

func NewAppointmentStorage(deps StorageDeps) AppointmentRepository {
	return &appointmentStorage{
		pg: deps.PostgresDB,
	}
}

func (s *appointmentStorage) Create(ctx context.Context, appointment *entity.Appointment, serviceIDs []uuid.UUID) (uuid.UUID, error) {
	tx, err := s.pg.DB.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if appointment.ID == uuid.Nil {
		appointment.ID = uuid.New()
	}

	// Insert appointment
	const appointmentQuery = `
		INSERT INTO appointments (id, user_id, vehicle_id, appointment_time, status, attachments)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	row := tx.QueryRowContext(ctx, appointmentQuery,
		appointment.ID, appointment.UserID, appointment.VehicleID,
		appointment.AppointmentTime, appointment.Status, pq.Array(appointment.Attachments),
	)

	if err := row.Scan(&appointment.ID); err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert appointment: %w", err)
	}

	// Insert appointment services
	const servicesQuery = `
		INSERT INTO appointment_services (id, appointment_id, service_id, price)
		SELECT uuid_generate_v4(), $1, id, price
		FROM services
		WHERE id = ANY($2) AND deleted_at IS NULL;
	`

	if _, err := tx.ExecContext(ctx, servicesQuery, appointment.ID, pq.Array(serviceIDs)); err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert appointment services: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return appointment.ID, nil
}

func (s *appointmentStorage) GetById(ctx context.Context, id uuid.UUID) (*entity.Appointment, error) {
	const query = `
		SELECT 
			a.id, a.user_id, a.vehicle_id, a.appointment_time, a.status, a.attachments,
			COALESCE(json_agg(json_build_object(
				'id', s.id,
				'name', s.name,
				'description', s.description,
				'price', s.price,
				'duration_min', s.duration_min
			)) FILTER (WHERE s.id IS NOT NULL), '[]') as services
		FROM appointments a
		LEFT JOIN appointment_services as_link ON a.id = as_link.appointment_id
		LEFT JOIN services s ON as_link.service_id = s.id AND s.deleted_at IS NULL
		WHERE a.id = $1 AND a.deleted_at IS NULL
		GROUP BY a.id, a.user_id, a.vehicle_id, a.appointment_time, a.status, a.attachments;
	`

	row := s.pg.DB.QueryRowContext(ctx, query, id)

	var appointment entity.Appointment
	var servicesJSON []byte
	if err := row.Scan(
		&appointment.ID, &appointment.UserID, &appointment.VehicleID,
		&appointment.AppointmentTime, &appointment.Status, pq.Array(&appointment.Attachments), &servicesJSON,
	); err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return &appointment, nil
}

func (s *appointmentStorage) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*entity.Appointment, error) {
	const query = `
		SELECT 
			a.id, a.user_id, a.vehicle_id, a.appointment_time, a.status, a.attachments,
			COALESCE(json_agg(json_build_object(
				'id', s.id,
				'name', s.name,
				'description', s.description,
				'price', s.price,
				'duration_min', s.duration_min
			)) FILTER (WHERE s.id IS NOT NULL), '[]') as services
		FROM appointments a
		LEFT JOIN appointment_services as_link ON a.id = as_link.appointment_id
		LEFT JOIN services s ON as_link.service_id = s.id AND s.deleted_at IS NULL
		WHERE a.user_id = $1 AND a.deleted_at IS NULL
		GROUP BY a.id, a.user_id, a.vehicle_id, a.appointment_time, a.status, a.attachments
		ORDER BY a.appointment_time DESC;
	`

	rows, err := s.pg.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query appointments: %w", err)
	}
	defer rows.Close()

	var appointments []*entity.Appointment
	for rows.Next() {
		var appointment entity.Appointment
		var servicesJSON []byte
		if err := rows.Scan(
			&appointment.ID, &appointment.UserID, &appointment.VehicleID,
			&appointment.AppointmentTime, &appointment.Status, pq.Array(&appointment.Attachments), &servicesJSON,
		); err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}

func (s *appointmentStorage) Update(ctx context.Context, appointment *entity.Appointment) error {
	const query = `
		UPDATE appointments
		SET appointment_time = $2, status = $3, attachments = $4, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	result, err := s.pg.DB.ExecContext(ctx, query,
		appointment.ID, appointment.AppointmentTime, appointment.Status, pq.Array(appointment.Attachments),
	)
	if err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("appointment not found")
	}

	return nil
}

func (s *appointmentStorage) UpdateServices(ctx context.Context, appointmentID uuid.UUID, serviceIDs []uuid.UUID) error {
	tx, err := s.pg.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing services
	const deleteQuery = `
		DELETE FROM appointment_services
		WHERE appointment_id = $1;
	`

	if _, err := tx.ExecContext(ctx, deleteQuery, appointmentID); err != nil {
		return fmt.Errorf("failed to delete existing services: %w", err)
	}

	// Insert new services
	const insertQuery = `
		INSERT INTO appointment_services (id, appointment_id, service_id, price)
		SELECT uuid_generate_v4(), $1, id, price
		FROM services
		WHERE id = ANY($2) AND deleted_at IS NULL;
	`

	if _, err := tx.ExecContext(ctx, insertQuery, appointmentID, pq.Array(serviceIDs)); err != nil {
		return fmt.Errorf("failed to insert new services: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *appointmentStorage) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		UPDATE appointments
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	result, err := s.pg.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("appointment not found")
	}

	return nil
}

func (s *appointmentStorage) CheckTimeSlotAvailable(ctx context.Context, appointmentTime string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM appointments
			WHERE appointment_time = $1
			AND deleted_at IS NULL
			AND status != 'cancelled'
		);
	`

	var exists bool
	if err := s.pg.DB.QueryRowContext(ctx, query, appointmentTime).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check time slot: %w", err)
	}

	return !exists, nil
}
