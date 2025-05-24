package models

import (
	"time"

	"github.com/google/uuid"
)

type VerificationType string

const (
	VerificationTypeRegistration VerificationType = "registration"
	VerificationTypeLogin        VerificationType = "login"
	VerificationTypePassword     VerificationType = "password"
)

type VerificationCode struct {
	ID        uuid.UUID        `db:"id"`
	UserID    uuid.UUID        `db:"user_id"`
	Email     string           `db:"email"`
	Code      string           `db:"code"`
	Type      VerificationType `db:"type"`
	Used      bool             `db:"used"`
	ExpiresAt time.Time        `db:"expires_at"`
	CreatedAt time.Time        `db:"created_at"`
}

type VerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}
