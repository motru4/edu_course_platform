package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CategoryID    uuid.UUID `json:"category_id"`
	Level         string    `json:"level"`
	Duration      int       `json:"duration"`
	Rating        float64   `json:"rating"`
	StudentsCount int       `json:"students_count"`
	Thumbnail     string    `json:"thumbnail"`
	Price         float64   `json:"price"`
	Status        string    `json:"status"`
	CreatedBy     uuid.UUID `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Test struct {
	ID           uuid.UUID `json:"id"`
	LessonID     uuid.UUID `json:"lesson_id"`
	PassingScore int       `json:"passing_score"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Question struct {
	ID            uuid.UUID `json:"id"`
	TestID        uuid.UUID `json:"test_id"`
	QuestionText  string    `json:"question_text"`
	Options       []string  `json:"options"`
	CorrectAnswer int       `json:"correct_answer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
