package repositories

import (
	"context"
	"course2/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(ctx context.Context, course *models.Course) error {
	query := `
		INSERT INTO courses (
			id, title, description, category_id, level, duration,
			rating, students_count, thumbnail, price, status, created_by,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW()
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		course.ID, course.Title, course.Description, course.CategoryID,
		course.Level, course.Duration, course.Rating, course.StudentsCount,
		course.Thumbnail, course.Price, course.Status, course.CreatedBy,
	)

	return err
}

func (r *CourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	query := `
		SELECT id, title, description, category_id, level, duration,
			   rating, students_count, thumbnail, price, status, created_by,
			   created_at, updated_at
		FROM courses
		WHERE id = $1
	`

	course := &models.Course{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID, &course.Title, &course.Description, &course.CategoryID,
		&course.Level, &course.Duration, &course.Rating, &course.StudentsCount,
		&course.Thumbnail, &course.Price, &course.Status, &course.CreatedBy,
		&course.CreatedAt, &course.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) List(ctx context.Context, offset, limit int) ([]*models.Course, error) {
	query := `
		SELECT id, title, description, category_id, level, duration,
			   rating, students_count, thumbnail, price, status, created_by,
			   created_at, updated_at
		FROM courses
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := &models.Course{}
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.CategoryID,
			&course.Level, &course.Duration, &course.Rating, &course.StudentsCount,
			&course.Thumbnail, &course.Price, &course.Status, &course.CreatedBy,
			&course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (r *CourseRepository) Update(ctx context.Context, course *models.Course) error {
	query := `
		UPDATE courses
		SET title = $1, description = $2, category_id = $3, level = $4,
			duration = $5, rating = $6, students_count = $7, thumbnail = $8,
			price = $9, status = $10, updated_at = NOW()
		WHERE id = $11
	`

	result, err := r.db.ExecContext(ctx, query,
		course.Title, course.Description, course.CategoryID, course.Level,
		course.Duration, course.Rating, course.StudentsCount, course.Thumbnail,
		course.Price, course.Status, course.ID,
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

func (r *CourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM courses WHERE id = $1`

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

func (r *CourseRepository) ListByCategory(ctx context.Context, categoryID uuid.UUID, offset, limit int) ([]*models.Course, error) {
	query := `
		SELECT id, title, description, category_id, level, duration,
			   rating, students_count, thumbnail, price, status, created_by,
			   created_at, updated_at
		FROM courses
		WHERE category_id = $1
		ORDER BY created_at DESC
		OFFSET $2 LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := &models.Course{}
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.CategoryID,
			&course.Level, &course.Duration, &course.Rating, &course.StudentsCount,
			&course.Thumbnail, &course.Price, &course.Status, &course.CreatedBy,
			&course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}
