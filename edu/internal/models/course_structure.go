package models

// CourseStructure представляет структуру курса с уроками и прогрессом
type CourseStructure struct {
	// Основная информация о курсе
	Course *Course `json:"course"`

	// Список уроков курса
	Lessons []*Lesson `json:"lessons"`

	// Статистика прогресса
	TotalLessons     int     `json:"total_lessons"`
	CompletedLessons int     `json:"completed_lessons"`
	Progress         float64 `json:"progress"`
}

// CalculateProgress безопасно вычисляет процент прогресса
func (cs *CourseStructure) CalculateProgress() {
	if cs.TotalLessons > 0 {
		cs.Progress = float64(cs.CompletedLessons) / float64(cs.TotalLessons) * 100
	} else {
		cs.Progress = 0
	}
}
