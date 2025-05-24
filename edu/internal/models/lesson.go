package models

import (
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	ID           uuid.UUID `json:"id"`
	CourseID     uuid.UUID `json:"course_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	OrderNum     int       `json:"order_num"`
	RequiresTest bool      `json:"requires_test"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Поля для отображения прогресса
	Completed  bool       `json:"completed"`
	TestScore  *int       `json:"test_score,omitempty"`
	PassedTest bool       `json:"passed_test,omitempty"`
	ViewedAt   *time.Time `json:"viewed_at,omitempty"`
	HasTest    bool       `json:"has_test,omitempty"`
}
