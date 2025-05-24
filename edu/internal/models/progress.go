package models

import (
	"time"

	"github.com/google/uuid"
)

type XPEntry struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id,omitempty"`
	LessonID uuid.UUID `json:"lesson_id,omitempty"`
	Type     string    `json:"type"` // lesson_view, test_pass, course_complete
	Amount   int       `json:"amount"`
	EarnedAt time.Time `json:"earned_at"`
}

type CourseProgress struct {
	CourseID         uuid.UUID  `json:"course_id"`
	TotalLessons     int        `json:"total_lessons"`
	CompletedLessons int        `json:"completed_lessons"`
	Percentage       float64    `json:"percentage"`
	XPEarned         int        `json:"xp_earned"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
}
