package entity

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
)

type Appointment struct {
	ID              uuid.UUID         `json:"id"`
	UserID          uuid.UUID         `json:"user_id"`
	VehicleID       uuid.UUID         `json:"vehicle_id"`
	AppointmentTime time.Time         `json:"appointment_time"`
	Status          AppointmentStatus `json:"status"`
	Services        []*Service        `json:"services,omitempty"`
	Attachments     []string          `json:"attachments"`
	CreatedAt       *time.Time        `json:"created_at,omitempty"`
	UpdatedAt       *time.Time        `json:"updated_at,omitempty"`
	DeletedAt       *time.Time        `json:"deleted_at,omitempty"`
}

type AppointmentCreate struct {
	VehicleID       uuid.UUID   `json:"vehicle_id"`
	AppointmentTime time.Time   `json:"appointment_time"`
	ServiceIDs      []uuid.UUID `json:"service_ids"`
	Attachments     []string    `json:"attachments"`
}

func (a *AppointmentCreate) Validate() error {
	if a.VehicleID == uuid.Nil {
		return fmt.Errorf("vehicle_id is required")
	}

	if a.AppointmentTime.IsZero() {
		return fmt.Errorf("appointment_time is required")
	}

	// Check if appointment time is in the future
	if a.AppointmentTime.Before(time.Now()) {
		return fmt.Errorf("appointment time must be in the future")
	}

	// Check if appointment is within business hours (assuming 9 AM to 6 PM)
	hour := a.AppointmentTime.Hour()
	if hour < 9 || hour >= 18 {
		return fmt.Errorf("appointments must be scheduled between 9 AM and 6 PM")
	}

	// Check if it's a weekday
	if a.AppointmentTime.Weekday() == time.Saturday || a.AppointmentTime.Weekday() == time.Sunday {
		return fmt.Errorf("appointments can only be scheduled on weekdays")
	}

	if len(a.ServiceIDs) == 0 {
		return fmt.Errorf("at least one service must be selected")
	}

	return nil
}

func (a *AppointmentCreate) ToAppointment(userID uuid.UUID) *Appointment {
	return &Appointment{
		UserID:          userID,
		VehicleID:       a.VehicleID,
		AppointmentTime: a.AppointmentTime,
		Status:          AppointmentStatusScheduled,
	}
}

type AppointmentUpdate struct {
	AppointmentTime *time.Time         `json:"appointment_time,omitempty"`
	Status          *AppointmentStatus `json:"status,omitempty"`
	ServiceIDs      []uuid.UUID        `json:"service_ids,omitempty"`
	Attachments     []string           `json:"attachments,omitempty"`
}

func (a *AppointmentUpdate) Validate() error {
	if a.AppointmentTime != nil {
		// Check if appointment time is in the future
		if a.AppointmentTime.Before(time.Now()) {
			return fmt.Errorf("appointment time must be in the future")
		}

		// Check if appointment is within business hours (assuming 9 AM to 6 PM)
		hour := a.AppointmentTime.Hour()
		if hour < 9 || hour >= 18 {
			return fmt.Errorf("appointments must be scheduled between 9 AM and 6 PM")
		}

		// Check if it's a weekday
		if a.AppointmentTime.Weekday() == time.Saturday || a.AppointmentTime.Weekday() == time.Sunday {
			return fmt.Errorf("appointments can only be scheduled on weekdays")
		}
	}

	if a.Status != nil {
		switch *a.Status {
		case AppointmentStatusScheduled, AppointmentStatusCompleted, AppointmentStatusCancelled:
			// Valid status
		default:
			return fmt.Errorf("invalid status: must be one of scheduled, completed, or cancelled")
		}
	}

	if len(a.ServiceIDs) > 0 {
		// If services are being updated, ensure at least one is selected
		if len(a.ServiceIDs) == 0 {
			return fmt.Errorf("at least one service must be selected")
		}
	}

	return nil
}

type AppointmentResponse struct {
	ID              uuid.UUID         `json:"id"`
	Vehicle         *Vehicle          `json:"vehicle"`
	AppointmentTime time.Time         `json:"appointment_time"`
	Status          AppointmentStatus `json:"status"`
	Services        []*Service        `json:"services"`
	TotalPrice      float64           `json:"total_price"`
	TotalDuration   int               `json:"total_duration_min"`
	Attachments     []string          `json:"attachments"`
	CreatedAt       *time.Time        `json:"created_at,omitempty"`
}
