package entity

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	Price       float64    `json:"price" db:"price"`
	DurationMin int        `json:"duration_min" db:"duration_min"`
	CreatedAt   *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (s *Service) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}
	if s.Price < 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if s.DurationMin < 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	return nil
}
