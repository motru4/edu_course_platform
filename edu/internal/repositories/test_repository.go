package repositories

import (
	"context"
	"course2/internal/models"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type TestRepository struct {
	db *sql.DB
}

func NewTestRepository(db *sql.DB) *TestRepository {
	return &TestRepository{db: db}
}

func (r *TestRepository) Create(ctx context.Context, test *models.Test) error {
	query := `
		INSERT INTO tests (
			id, lesson_id, passing_score, created_at, updated_at
		) VALUES (
			$1, $2, $3, NOW(), NOW()
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		test.ID, test.LessonID, test.PassingScore,
	)

	return err
}

func (r *TestRepository) GetByLessonID(ctx context.Context, lessonID uuid.UUID) (*models.Test, error) {
	query := `
		SELECT id, lesson_id, passing_score, created_at, updated_at
		FROM tests
		WHERE lesson_id = $1
	`

	test := &models.Test{}
	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(
		&test.ID, &test.LessonID, &test.PassingScore,
		&test.CreatedAt, &test.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (r *TestRepository) Update(ctx context.Context, test *models.Test) error {
	query := `
		UPDATE tests
		SET passing_score = $1, updated_at = NOW()
		WHERE id = $2 AND lesson_id = $3
	`

	result, err := r.db.ExecContext(ctx, query,
		test.PassingScore, test.ID, test.LessonID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tests WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Методы для работы с вопросами
func (r *TestRepository) CreateQuestion(ctx context.Context, question *models.Question) error {
	query := `
		INSERT INTO questions (
			id, test_id, question_text, options, correct_answer,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, NOW(), NOW()
		)
	`

	optionsJSON, err := json.Marshal(question.Options)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query,
		question.ID, question.TestID, question.QuestionText,
		optionsJSON, question.CorrectAnswer,
	)

	return err
}

func (r *TestRepository) GetQuestions(ctx context.Context, testID uuid.UUID) ([]*models.Question, error) {
	query := `
		SELECT id, test_id, question_text, options, correct_answer,
			   created_at, updated_at
		FROM questions
		WHERE test_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(ctx, query, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		question := &models.Question{}
		var optionsJSON []byte
		err := rows.Scan(
			&question.ID, &question.TestID, &question.QuestionText,
			&optionsJSON, &question.CorrectAnswer, &question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(optionsJSON, &question.Options)
		if err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func (r *TestRepository) UpdateQuestion(ctx context.Context, question *models.Question) error {
	query := `
		UPDATE questions
		SET question_text = $1, options = $2, correct_answer = $3,
			updated_at = NOW()
		WHERE id = $4 AND test_id = $5
	`

	optionsJSON, err := json.Marshal(question.Options)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query,
		question.QuestionText, optionsJSON, question.CorrectAnswer,
		question.ID, question.TestID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TestRepository) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM questions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
