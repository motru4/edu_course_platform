package models

import (
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	UserID            uuid.UUID `json:"user_id"`
	Role              Role      `json:"role"`
	PasswordChangedAt int64     `json:"pwd_changed"`
}

type RefreshSession struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    int64     `db:"expires_at"`
	CreatedAt    int64     `db:"created_at"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
