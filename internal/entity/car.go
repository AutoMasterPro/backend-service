package entity

import (
	"github.com/google/uuid"
	"time"
)

type Car struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	StateNumber string     `json:"state_number"`
	Model       string     `json:"model"`
	VIN         string     `json:"vin"`
	YearRelease string     `json:"year_release"`
	User        User       `json:"user"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
