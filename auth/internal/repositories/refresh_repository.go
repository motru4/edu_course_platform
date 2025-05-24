package repositories

import (
	"database/sql"
	"fmt"

	"auth-service/internal/models"

	"github.com/google/uuid"
)

type RefreshRepository struct {
	db *sql.DB
}

func NewRefreshRepository(db *sql.DB) *RefreshRepository {
	return &RefreshRepository{db: db}
}

func (r *RefreshRepository) Create(session *models.RefreshSession) error {
	_, err := r.db.Exec(`
		INSERT INTO refresh_sessions (id, user_id, refresh_token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, session.ID, session.UserID, session.RefreshToken, session.ExpiresAt, session.CreatedAt)

	if err != nil {
		return fmt.Errorf("error creating refresh session: %w", err)
	}

	return nil
}

func (r *RefreshRepository) GetByToken(token string) (*models.RefreshSession, error) {
	var session models.RefreshSession
	err := r.db.QueryRow(`
		SELECT id, user_id, refresh_token, expires_at, created_at
		FROM refresh_sessions
		WHERE refresh_token = $1
	`, token).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting refresh session: %w", err)
	}

	return &session, nil
}

func (r *RefreshRepository) Update(sessionID uuid.UUID, newToken string, expiresAt int64) error {
	result, err := r.db.Exec(`
		UPDATE refresh_sessions 
		SET refresh_token = $1, expires_at = $2
		WHERE id = $3
	`, newToken, expiresAt, sessionID)

	if err != nil {
		return fmt.Errorf("error updating refresh session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh session not found")
	}

	return nil
}

func (r *RefreshRepository) Delete(sessionID uuid.UUID) error {
	_, err := r.db.Exec(`
		DELETE FROM refresh_sessions WHERE id = $1
	`, sessionID)

	if err != nil {
		return fmt.Errorf("error deleting refresh session: %w", err)
	}

	return nil
}

func (r *RefreshRepository) DeleteExpired(currentTime int64) error {
	_, err := r.db.Exec(`
		DELETE FROM refresh_sessions WHERE expires_at < $1
	`, currentTime)

	if err != nil {
		return fmt.Errorf("error deleting expired sessions: %w", err)
	}

	return nil
}

func (r *RefreshRepository) DeleteUserSessions(userID uuid.UUID) error {
	_, err := r.db.Exec(`
		DELETE FROM refresh_sessions WHERE user_id = $1
	`, userID)

	if err != nil {
		return fmt.Errorf("error deleting user sessions: %w", err)
	}

	return nil
}

func (r *RefreshRepository) DeleteAllUserSessions(userID uuid.UUID) error {
	query := `DELETE FROM refresh_sessions WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
