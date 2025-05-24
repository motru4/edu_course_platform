package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleStudent Role = "student"
	RoleAuthor  Role = "author"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID                uuid.UUID `json:"id" db:"id"`
	Email             string    `json:"email" db:"email"`
	PasswordHash      string    `json:"-" db:"password_hash"`
	Role              Role      `json:"role" db:"role"`
	Confirmed         bool      `json:"confirmed" db:"confirmed"`
	GoogleID          *string   `json:"google_id,omitempty" db:"google_id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at" db:"password_changed_at"`
}

type UserCreate struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	//Role     Role   `json:"role" binding:"required,oneof=student author admin"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
