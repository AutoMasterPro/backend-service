package entity

import (
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"time"
)

type Vehicle struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	Brand        string     `json:"brand"`
	Model        string     `json:"model"`
	LicensePlate string     `json:"license_plate"`
	Year         int        `json:"year"`
	VIN          string     `json:"vin,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type VehicleCreate struct {
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	LicensePlate string `json:"license_plate"`
	Year         int    `json:"year"`
	VIN          string `json:"vin,omitempty"`
}

func (v *VehicleCreate) Validate() error {
	if v.Brand == "" {
		return fmt.Errorf("brand is required")
	}
	if v.Model == "" {
		return fmt.Errorf("model is required")
	}
	if v.LicensePlate == "" {
		return fmt.Errorf("license plate is required")
	}

	// Basic license plate format validation (adjust regex based on your country's format)
	licensePlateRegex := regexp.MustCompile(`^[A-Z0-9]{3,10}$`)
	if !licensePlateRegex.MatchString(v.LicensePlate) {
		return fmt.Errorf("invalid license plate format")
	}

	currentYear := time.Now().Year()
	if v.Year < 1900 || v.Year > currentYear+1 {
		return fmt.Errorf("invalid year: must be between 1900 and %d", currentYear+1)
	}

	// VIN validation if provided
	if v.VIN != "" {
		vinRegex := regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)
		if !vinRegex.MatchString(v.VIN) {
			return fmt.Errorf("invalid VIN format")
		}
	}

	return nil
}

func (v *VehicleCreate) ToVehicle(userID uuid.UUID) *Vehicle {
	return &Vehicle{
		UserID:       userID,
		Brand:        v.Brand,
		Model:        v.Model,
		LicensePlate: v.LicensePlate,
		Year:         v.Year,
		VIN:          v.VIN,
	}
}

type VehicleUpdate struct {
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	LicensePlate string `json:"license_plate"`
	Year         int    `json:"year"`
	VIN          string `json:"vin,omitempty"`
}

func (v *VehicleUpdate) Validate() error {
	if v.Brand == "" {
		return fmt.Errorf("brand is required")
	}
	if v.Model == "" {
		return fmt.Errorf("model is required")
	}
	if v.LicensePlate == "" {
		return fmt.Errorf("license plate is required")
	}

	// Basic license plate format validation
	licensePlateRegex := regexp.MustCompile(`^[A-Z0-9]{3,10}$`)
	if !licensePlateRegex.MatchString(v.LicensePlate) {
		return fmt.Errorf("invalid license plate format")
	}

	currentYear := time.Now().Year()
	if v.Year < 1900 || v.Year > currentYear+1 {
		return fmt.Errorf("invalid year: must be between 1900 and %d", currentYear+1)
	}

	// VIN validation if provided
	if v.VIN != "" {
		vinRegex := regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)
		if !vinRegex.MatchString(v.VIN) {
			return fmt.Errorf("invalid VIN format")
		}
	}

	return nil
}

func (v *VehicleUpdate) ToVehicle(id uuid.UUID) *Vehicle {
	return &Vehicle{
		ID:           id,
		Brand:        v.Brand,
		Model:        v.Model,
		LicensePlate: v.LicensePlate,
		Year:         v.Year,
		VIN:          v.VIN,
	}
}
