package repositories

import (
	"database/sql"
	"fmt"

	"auth-service/internal/models"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (id, email, password_hash, role, created_at, password_changed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, user.ID, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.CreatedAt)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(`
		SELECT id, email, password_hash, role, confirmed, google_id, created_at, password_changed_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Confirmed,
		&user.GoogleID,
		&user.CreatedAt,
		&user.PasswordChangedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(`
		SELECT id, email, password_hash, role, confirmed, google_id, created_at, password_changed_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Confirmed,
		&user.GoogleID,
		&user.CreatedAt,
		&user.PasswordChangedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) CheckEmailExists(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, email).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}

	return exists, nil
}

func (r *UserRepository) UpdateConfirmation(userID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE users SET confirmed = true WHERE id = $1
	`, userID)

	if err != nil {
		return fmt.Errorf("error updating user confirmation: %w", err)
	}

	return nil
}

func (r *UserRepository) UpdateGoogleID(userID uuid.UUID, googleID string) error {
	_, err := r.db.Exec(`
		UPDATE users SET google_id = $1 WHERE id = $2
	`, googleID, userID)

	if err != nil {
		return fmt.Errorf("error updating google id: %w", err)
	}

	return nil
}

func (r *UserRepository) UpdatePassword(userID uuid.UUID, hashedPassword string) error {
	query := `UPDATE users SET password_hash = $1, password_changed_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, hashedPassword, userID)
	return err
}

func (r *UserRepository) Update(user *models.User) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET password_hash = $1, 
			password_changed_at = $2
		WHERE id = $3
	`, user.PasswordHash, user.PasswordChangedAt, user.ID)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}
