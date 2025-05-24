package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, userID, courseID uuid.UUID, rating int, text string) error {
	query := `
		INSERT INTO reviews (id, user_id, course_id, rating, text, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err := r.db.ExecContext(ctx, query, uuid.New(), userID, courseID, rating, text)
	if err != nil {
		return err
	}

	// Обновляем средний рейтинг курса
	updateQuery := `
		UPDATE courses
		SET rating = (
			SELECT AVG(rating)::float
			FROM reviews
			WHERE course_id = $1
		)
		WHERE id = $1
	`

	_, err = r.db.ExecContext(ctx, updateQuery, courseID)
	return err
}

func (r *ReviewRepository) GetByCourse(ctx context.Context, courseID uuid.UUID) ([]struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Rating    int
	Text      string
	CreatedAt string
}, error) {
	query := `
		SELECT id, user_id, rating, text, created_at
		FROM reviews
		WHERE course_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []struct {
		ID        uuid.UUID
		UserID    uuid.UUID
		Rating    int
		Text      string
		CreatedAt string
	}

	for rows.Next() {
		var review struct {
			ID        uuid.UUID
			UserID    uuid.UUID
			Rating    int
			Text      string
			CreatedAt string
		}
		if err := rows.Scan(&review.ID, &review.UserID, &review.Rating, &review.Text, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *ReviewRepository) Update(ctx context.Context, id uuid.UUID, rating int, text string) error {
	query := `
		UPDATE reviews
		SET rating = $1, text = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, rating, text, id)
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

func (r *ReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`

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
