package models

import (
	"time"

	"github.com/google/uuid"
)

type LessonProgress struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	LessonID      uuid.UUID  `json:"lesson_id"`
	ViewedAt      *time.Time `json:"viewed_at"`
	TestScore     *int       `json:"test_score"`
	PassedTest    bool       `json:"passed_test"`
	CompletedAt   *time.Time `json:"completed_at"`
	LastAttemptAt *time.Time `json:"last_attempt_at"`
	IsCompleted   bool       `json:"is_completed"`
	AttemptsCount int        `json:"attempts_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
