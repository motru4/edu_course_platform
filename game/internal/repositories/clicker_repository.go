package repositories

import (
	"context"
	"database/sql"
	"time"

	"game/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ClickerRepository представляет репозиторий для работы с данными кликера
type ClickerRepository struct {
	db *sqlx.DB
}

// NewClickerRepository создает новый экземпляр ClickerRepository
func NewClickerRepository(db *sqlx.DB) *ClickerRepository {
	return &ClickerRepository{
		db: db,
	}
}

// GetOrCreateStats получает или создает статистику кликера для пользователя
func (r *ClickerRepository) GetOrCreateStats(ctx context.Context, userID uuid.UUID) (*models.ClickerStats, error) {
	stats := &models.ClickerStats{}

	query := `
		SELECT id, user_id, total_clicks, clicks_per_second, last_click_time, 
		       last_save_time, last_save_count, created_at, updated_at
		FROM clicker_stats
		WHERE user_id = $1
	`

	err := r.db.GetContext(ctx, stats, query, userID)
	if err == sql.ErrNoRows {
		// Создаем новую запись статистики
		newStats := &models.ClickerStats{
			ID:          uuid.New(),
			UserID:      userID,
			TotalClicks: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		insertQuery := `
			INSERT INTO clicker_stats (id, user_id, total_clicks, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, user_id, total_clicks, clicks_per_second, last_click_time, 
			          last_save_time, last_save_count, created_at, updated_at
		`

		err = r.db.GetContext(ctx, stats, insertQuery,
			newStats.ID, newStats.UserID, newStats.TotalClicks,
			newStats.CreatedAt, newStats.UpdatedAt)

		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return stats, nil
}

// UpdateStats обновляет статистику кликера для пользователя
func (r *ClickerRepository) UpdateStats(ctx context.Context, stats *models.ClickerStats) error {
	stats.UpdatedAt = time.Now()

	var clicksPerSec *float64
	if stats.ClicksPerSec != nil {
		clicksPerSec = stats.ClicksPerSec
	}

	var lastSaveCount *int
	if stats.LastSaveCount != nil {
		lastSaveCount = stats.LastSaveCount
	}

	query := `
		UPDATE clicker_stats
		SET total_clicks = $1, clicks_per_second = $2, last_click_time = $3, 
		    last_save_time = $4, last_save_count = $5, updated_at = $6
		WHERE user_id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		stats.TotalClicks, clicksPerSec, stats.LastClickTime,
		stats.LastSaveTime, lastSaveCount, stats.UpdatedAt, stats.UserID)

	return err
}

// SaveSession сохраняет сессию кликера
func (r *ClickerRepository) SaveSession(ctx context.Context, session *models.ClickerSession) error {
	query := `
		INSERT INTO clicker_sessions (id, user_id, start_time, end_time, click_count, average_cps, max_cps, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.StartTime, session.EndTime,
		session.ClickCount, session.AverageCPS, session.MaxCPS, session.CreatedAt)

	return err
}

// GetRecentSessions получает последние сессии кликера для пользователя
func (r *ClickerRepository) GetRecentSessions(ctx context.Context, userID uuid.UUID, limit int) ([]models.ClickerSession, error) {
	sessions := []models.ClickerSession{}

	query := `
		SELECT id, user_id, start_time, end_time, click_count, average_cps, max_cps, created_at
		FROM clicker_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	err := r.db.SelectContext(ctx, &sessions, query, userID, limit)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// UpdateLeaderboard обновляет или создает запись в таблице лидеров
func (r *ClickerRepository) UpdateLeaderboard(ctx context.Context, entry *models.LeaderboardEntry) error {
	query := `
		INSERT INTO clicker_leaderboard (id, user_id, username, score, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id)
		DO UPDATE SET
			score = $4,
			updated_at = $5
		RETURNING id, user_id, username, score, rank, updated_at
	`

	now := time.Now()
	entry.UpdatedAt = now

	err := r.db.GetContext(ctx, entry, query,
		entry.ID, entry.UserID, entry.Username, entry.Score, now)

	return err
}

// GetLeaderboard получает таблицу лидеров
func (r *ClickerRepository) GetLeaderboard(ctx context.Context, limit int, offset int) ([]models.LeaderboardEntry, error) {
	entries := []models.LeaderboardEntry{}

	query := `
		SELECT id, user_id, username, score, rank, updated_at
		FROM clicker_leaderboard
		ORDER BY score DESC
		LIMIT $1 OFFSET $2
	`

	err := r.db.SelectContext(ctx, &entries, query, limit, offset)
	if err != nil {
		return nil, err
	}

	// Обновляем ранги
	for i := range entries {
		rank := offset + i + 1
		entries[i].Rank = &rank
	}

	return entries, nil
}

// GetUserRank получает ранг пользователя в таблице лидеров
func (r *ClickerRepository) GetUserRank(ctx context.Context, userID uuid.UUID) (*models.LeaderboardEntry, error) {
	entry := &models.LeaderboardEntry{}

	query := `
		SELECT l.id, l.user_id, l.username, l.score, l.updated_at,
		       (SELECT COUNT(*) FROM clicker_leaderboard WHERE score > l.score) + 1 AS rank
		FROM clicker_leaderboard l
		WHERE l.user_id = $1
	`

	err := r.db.GetContext(ctx, entry, query, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return entry, nil
}

// RefreshLeaderboardRanks обновляет ранги в таблице лидеров
func (r *ClickerRepository) RefreshLeaderboardRanks(ctx context.Context) error {
	query := `
		WITH ranked AS (
			SELECT id, ROW_NUMBER() OVER (ORDER BY score DESC) AS new_rank
			FROM clicker_leaderboard
		)
		UPDATE clicker_leaderboard cl
		SET rank = r.new_rank
		FROM ranked r
		WHERE cl.id = r.id
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}
