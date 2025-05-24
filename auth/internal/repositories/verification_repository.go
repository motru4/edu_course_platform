package repositories

import (
	"auth-service/internal/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type VerificationRepository struct {
	db *sql.DB
}

func NewVerificationRepository(db *sql.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) Create(verification *models.VerificationCode) error {
	_, err := r.db.Exec(`
		INSERT INTO verification_codes (id, user_id, email, code, type, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, verification.ID, verification.UserID, verification.Email, verification.Code,
		verification.Type, verification.ExpiresAt, verification.CreatedAt)

	if err != nil {
		return fmt.Errorf("error creating verification code: %w", err)
	}

	return nil
}

func (r *VerificationRepository) GetActiveCode(email string, verificationType models.VerificationType) (*models.VerificationCode, error) {
	var code models.VerificationCode
	query := `
		SELECT id, user_id, email, code, type, used, expires_at, created_at 
		FROM verification_codes 
		WHERE email = $1 
		AND type = $2 
		AND used = FALSE 
		AND expires_at > NOW()
		ORDER BY created_at DESC 
		LIMIT 1
	`
	err := r.db.QueryRow(query, email, verificationType).Scan(
		&code.ID,
		&code.UserID,
		&code.Email,
		&code.Code,
		&code.Type,
		&code.Used,
		&code.ExpiresAt,
		&code.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &code, nil
}

func (r *VerificationRepository) DeleteExpiredCodes() error {
	_, err := r.db.Exec(`
		DELETE FROM verification_codes
		WHERE expires_at <= $1
	`, time.Now())

	if err != nil {
		return fmt.Errorf("error deleting expired codes: %w", err)
	}

	return nil
}

func (r *VerificationRepository) MarkAsUsed(id uuid.UUID) error {
	query := `UPDATE verification_codes SET used = TRUE WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
