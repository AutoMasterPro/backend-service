package services

import (
	"backend-service/internal/entity"
	"backend-service/internal/storages"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type AppointmentService interface {
	Create(ctx context.Context, userID uuid.UUID, input *entity.AppointmentCreate) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*entity.Appointment, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]*entity.Appointment, error)
	Update(ctx context.Context, id uuid.UUID, input *entity.AppointmentUpdate) error
	Cancel(ctx context.Context, id uuid.UUID) error
}

type appointmentService struct {
	appointmentRepo storages.AppointmentRepository
	vehicleRepo     storages.VehicleRepository
	serviceRepo     storages.ServiceRepository
}

func NewAppointmentService(
	appointmentRepo storages.AppointmentRepository,
	vehicleRepo storages.VehicleRepository,
	serviceRepo storages.ServiceRepository,
) AppointmentService {
	return &appointmentService{
		appointmentRepo: appointmentRepo,
		vehicleRepo:     vehicleRepo,
		serviceRepo:     serviceRepo,
	}
}

func (s *appointmentService) Create(ctx context.Context, userID uuid.UUID, input *entity.AppointmentCreate) (uuid.UUID, error) {
	if err := input.Validate(); err != nil {
		return uuid.Nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if the vehicle belongs to the user
	vehicle, err := s.vehicleRepo.GetById(ctx, input.VehicleID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get vehicle: %w", err)
	}
	if vehicle.UserID != userID {
		return uuid.Nil, fmt.Errorf("vehicle does not belong to the user")
	}

	// Check if the time slot is available
	available, err := s.appointmentRepo.CheckTimeSlotAvailable(ctx, input.AppointmentTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to check time slot: %w", err)
	}
	if !available {
		return uuid.Nil, fmt.Errorf("time slot is not available")
	}

	// Create appointment
	appointment := input.ToAppointment(userID)
	appointment.Attachments = input.Attachments
	return s.appointmentRepo.Create(ctx, appointment, input.ServiceIDs)
}

func (s *appointmentService) GetById(ctx context.Context, id uuid.UUID) (*entity.Appointment, error) {
	return s.appointmentRepo.GetById(ctx, id)
}

func (s *appointmentService) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*entity.Appointment, error) {
	return s.appointmentRepo.GetByUserId(ctx, userId)
}

func (s *appointmentService) Update(ctx context.Context, id uuid.UUID, input *entity.AppointmentUpdate) error {
	if err := input.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	appointment, err := s.appointmentRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get appointment: %w", err)
	}

	//if input.AppointmentTime != nil {
	//	// Check if the new time slot is available
	//	available, err := s.appointmentRepo.CheckTimeSlotAvailable(ctx, input.AppointmentTime.Format("2006-01-02 15:04:05"))
	//	if err != nil {
	//		return fmt.Errorf("failed to check time slot: %w", err)
	//	}
	//	if !available {
	//		return fmt.Errorf("time slot is not available")
	//	}
	//	appointment.AppointmentTime = *input.AppointmentTime
	//}

	if input.Status != nil {
		appointment.Status = *input.Status
	}

	if input.Attachments != nil {
		appointment.Attachments = input.Attachments
	}

	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	if len(input.ServiceIDs) > 0 {
		if err := s.appointmentRepo.UpdateServices(ctx, id, input.ServiceIDs); err != nil {
			return fmt.Errorf("failed to update services: %w", err)
		}
	}

	return nil
}

func (s *appointmentService) Cancel(ctx context.Context, id uuid.UUID) error {
	appointment, err := s.appointmentRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get appointment: %w", err)
	}

	status := entity.AppointmentStatusCancelled
	appointment.Status = status

	return s.appointmentRepo.Update(ctx, appointment)
}
