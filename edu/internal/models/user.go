package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID              `json:"id"`
	Email     string                 `json:"email"`
	FirstName *string                `json:"first_name,omitempty"`
	LastName  *string                `json:"last_name,omitempty"`
	Avatar    *string                `json:"avatar,omitempty"`
	Role      string                 `json:"role"`
	Confirmed bool                   `json:"confirmed"`
	Settings  map[string]interface{} `json:"settings,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type UserProfile struct {
	ID        uuid.UUID              `json:"id"`
	Email     string                 `json:"email"`
	FirstName *string                `json:"first_name,omitempty"`
	LastName  *string                `json:"last_name,omitempty"`
	Avatar    *string                `json:"avatar,omitempty"`
	Role      string                 `json:"role"`
	TotalXP   int                    `json:"total_xp"`
	Settings  map[string]interface{} `json:"settings,omitempty"`
}
