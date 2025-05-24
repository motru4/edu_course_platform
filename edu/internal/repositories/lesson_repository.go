package repositories

import (
	"context"
	"course2/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

type LessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) Create(ctx context.Context, lesson *models.Lesson) error {
	query := `
		INSERT INTO lessons (
			id, course_id, title, content, order_num,
			requires_test, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW(), NOW()
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		lesson.ID, lesson.CourseID, lesson.Title, lesson.Content,
		lesson.OrderNum, lesson.RequiresTest,
	)

	return err
}

func (r *LessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, content, order_num,
			   requires_test, created_at, updated_at
		FROM lessons
		WHERE id = $1
	`

	lesson := &models.Lesson{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Content,
		&lesson.OrderNum, &lesson.RequiresTest, &lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return lesson, nil
}

func (r *LessonRepository) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, content, order_num,
			   requires_test, created_at, updated_at
		FROM lessons
		WHERE course_id = $1
		ORDER BY order_num
	`

	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*models.Lesson
	for rows.Next() {
		lesson := &models.Lesson{}
		err := rows.Scan(
			&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Content,
			&lesson.OrderNum, &lesson.RequiresTest, &lesson.CreatedAt,
			&lesson.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (r *LessonRepository) Update(ctx context.Context, lesson *models.Lesson) error {
	query := `
		UPDATE lessons
		SET title = $1, content = $2, order_num = $3,
			requires_test = $4, updated_at = NOW()
		WHERE id = $5 AND course_id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		lesson.Title, lesson.Content, lesson.OrderNum,
		lesson.RequiresTest, lesson.ID, lesson.CourseID,
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

func (r *LessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM lessons WHERE id = $1`

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

func (r *LessonRepository) GetLesson(ctx context.Context, lessonID uuid.UUID) (*models.Lesson, error) {
	query := `
		SELECT id, course_id, title, content, order_num,
			   requires_test, created_at, updated_at
		FROM lessons
		WHERE id = $1
	`

	lesson := &models.Lesson{}
	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(
		&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Content,
		&lesson.OrderNum, &lesson.RequiresTest, &lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return lesson, nil
}
