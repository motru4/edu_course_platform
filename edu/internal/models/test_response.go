package models

// TestResponse представляет ответ с тестом и информацией о прогрессе
type TestResponse struct {
	// Основная информация о тесте
	Test *Test `json:"test"`

	// Вопросы теста (без правильных ответов)
	Questions []*Question `json:"questions"`

	// Проходной балл
	PassingScore int `json:"passing_score"`

	// Информация о прогрессе
	LastScore     *int `json:"last_score,omitempty"`
	Passed        bool `json:"passed"`
	AttemptsCount int  `json:"attempts_count"`
}
