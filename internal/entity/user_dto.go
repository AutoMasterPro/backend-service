package entity

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type UserRegister struct {
	FullName string `json:"full_name,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (e *UserRegister) Validate() error {
	if e.FullName == "" {
		return fmt.Errorf("full name is required")
	}

	if e.Phone == "" {
		return fmt.Errorf("phone is required")
	}

	if e.Email == "" {
		return fmt.Errorf("email is required")
	}

	// Регулярное выражение для проверки формата email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(e.Email) {
		return fmt.Errorf("invalid email format")
	}

	if e.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

func (e *UserRegister) UserRegisterToUser() *User {
	bytes, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return &User{
		FullName:     e.FullName,
		Phone:        e.Phone,
		Email:        e.Email,
		PasswordHash: string(bytes),
	}
}

type UserLogin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (e *UserLogin) Validate() error {
	if e.Email == "" {
		return fmt.Errorf("email is required")
	}
	if e.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (e *UserLogin) UserLoginToUser() *User {
	bytes, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return &User{
		Email:        e.Email,
		PasswordHash: string(bytes),
	}
}
