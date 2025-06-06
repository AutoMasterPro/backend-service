package entity

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           uuid.UUID  `json:"id,omitempty"`
	FullName     string     `json:"full_name,omitempty"`
	Phone        string     `json:"phone,omitempty"`
	Email        string     `json:"email,omitempty"`
	PasswordHash string     `json:"password_hash,omitempty"`
	IsAdmin      bool       `json:"is_admin,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func (e *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(e.PasswordHash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
