package repositories

import (
	"context"
	"course2/internal/models"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, first_name, last_name, avatar, role, confirmed, settings,
			   created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Avatar, &user.Role, &user.Confirmed, &settingsJSON,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if settingsJSON != nil {
		if err := json.Unmarshal(settingsJSON, &user.Settings); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, profile *models.UserProfile) error {
	settingsJSON, err := json.Marshal(profile.Settings)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, avatar = $3, settings = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		profile.FirstName, profile.LastName, profile.Avatar,
		settingsJSON, profile.ID,
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

func (r *UserRepository) GetTotalXP(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM xp_entries
		WHERE user_id = $1
	`

	var totalXP int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&totalXP)
	if err != nil {
		return 0, err
	}

	return totalXP, nil
}

func (r *UserRepository) GetPurchasedCourses(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
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
