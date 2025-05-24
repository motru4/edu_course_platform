package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type PurchaseRepository struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) CreatePurchase(ctx context.Context, userID, courseID uuid.UUID) error {
	query := `
		INSERT INTO purchased_courses (user_id, course_id, purchased_at)
		VALUES ($1, $2, NOW())
	`

	_, err := r.db.ExecContext(ctx, query, userID, courseID)
	return err
}

func (r *PurchaseRepository) HasPurchased(ctx context.Context, userID, courseID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM purchased_courses
			WHERE user_id = $1 AND course_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, courseID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PurchaseRepository) GetPurchasedCourses(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT course_id
		FROM purchased_courses
		WHERE user_id = $1
		ORDER BY purchased_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courseIDs []uuid.UUID
	for rows.Next() {
		var courseID uuid.UUID
		if err := rows.Scan(&courseID); err != nil {
			return nil, err
		}
		courseIDs = append(courseIDs, courseID)
	}

	return courseIDs, nil
}

func (r *PurchaseRepository) UpdateCompletedAt(ctx context.Context, userID, courseID uuid.UUID, completedAt time.Time) error {
	query := `
		UPDATE purchased_courses
		SET completed_at = $1
		WHERE user_id = $2 AND course_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, completedAt, userID, courseID)
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
