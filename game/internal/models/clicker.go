package models

import (
	"time"

	"github.com/google/uuid"
)

// ClickerStats представляет статистику пользователя в игре-кликере
type ClickerStats struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	TotalClicks   int64      `json:"total_clicks" db:"total_clicks"`
	ClicksPerSec  *float64   `json:"clicks_per_second,omitempty" db:"clicks_per_second"`
	LastClickTime *time.Time `json:"last_click_time,omitempty" db:"last_click_time"`
	LastSaveTime  *time.Time `json:"last_save_time,omitempty" db:"last_save_time"`
	LastSaveCount *int       `json:"last_save_count,omitempty" db:"last_save_count"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// ClickerSession представляет одну игровую сессию кликера
type ClickerSession struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	StartTime  time.Time  `json:"start_time" db:"start_time"`
	EndTime    *time.Time `json:"end_time,omitempty" db:"end_time"`
	ClickCount int        `json:"click_count" db:"click_count"`
	AverageCPS *float64   `json:"average_cps,omitempty" db:"average_cps"`
	MaxCPS     *float64   `json:"max_cps,omitempty" db:"max_cps"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// LeaderboardEntry представляет запись в таблице лидеров
type LeaderboardEntry struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Score     int64     `json:"score" db:"score"`
	Rank      *int      `json:"rank" db:"rank"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ClickRequest представляет запрос на сохранение кликов
type ClickRequest struct {
	ClickCount      int     `json:"click_count" binding:"required,min=1"`
	SessionTime     float64 `json:"session_time" binding:"required,min=0"`
	ClientTimestamp int64   `json:"client_timestamp" binding:"required"`
}

// ClickResponse представляет ответ на запрос сохранения кликов
type ClickResponse struct {
	TotalClicks int64  `json:"total_clicks"`
	Status      string `json:"status"`
}

// LeaderboardResponse представляет ответ с таблицей лидеров
type LeaderboardResponse struct {
	Entries  []LeaderboardEntry `json:"entries"`
	UserRank *LeaderboardEntry  `json:"user_rank,omitempty"`
}

// StatsResponse представляет ответ со статистикой пользователя
type StatsResponse struct {
	Stats          ClickerStats     `json:"stats"`
	RecentSessions []ClickerSession `json:"recent_sessions,omitempty"`
}
