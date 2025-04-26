package entity

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           uuid.UUID
	FullName     string
	Phone        string
	Email        string
	PasswordHash string
	IsAdmin      bool
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
}

func (e *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(e.PasswordHash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
