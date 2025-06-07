package services

import (
	"context"
	"errors"
	"math"
	"time"

	"game/internal/models"
	"game/internal/repositories"

	"github.com/google/uuid"
)

// Константы для проверки на читерство
const (
	MaxClicksPerSecond     = 20.0  // Максимальное количество кликов в секунду
	SuspiciousClicksPerSec = 15.0  // Подозрительное количество кликов в секунду
	MaxSessionTime         = 300.0 // Максимальное время сессии в секундах (5 минут)
	MaxTimeDrift           = 10.0  // Максимальное расхождение во времени клиента и сервера в секундах
)

var (
	ErrInvalidClickRate = errors.New("недопустимая скорость кликов")
	ErrInvalidSession   = errors.New("недопустимая сессия")
	ErrInvalidTime      = errors.New("недопустимое время")
)

// ClickerService представляет сервис для работы с кликером
type ClickerService struct {
	repo *repositories.ClickerRepository
}

// NewClickerService создает новый экземпляр ClickerService
func NewClickerService(repo *repositories.ClickerRepository) *ClickerService {
	return &ClickerService{
		repo: repo,
	}
}

// SaveClicks сохраняет клики пользователя
func (s *ClickerService) SaveClicks(ctx context.Context, userID uuid.UUID, req *models.ClickRequest) (*models.ClickResponse, error) {
	// Проверка на читерство
	if err := s.validateClickRequest(req); err != nil {
		return &models.ClickResponse{
			Status: "rejected",
		}, err
	}

	// Получаем или создаем статистику пользователя
	stats, err := s.repo.GetOrCreateStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Вычисляем скорость кликов
	clicksPerSecond := float64(req.ClickCount) / req.SessionTime
	now := time.Now()
	saveCount := req.ClickCount

	// Обновляем статистику
	stats.TotalClicks += int64(req.ClickCount)
	stats.ClicksPerSec = &clicksPerSecond
	stats.LastClickTime = &now
	stats.LastSaveTime = &now
	stats.LastSaveCount = &saveCount

	// Сохраняем обновленную статистику
	err = s.repo.UpdateStats(ctx, stats)
	if err != nil {
		return nil, err
	}

	// Создаем новую сессию
	sessionStart := time.Now().Add(-time.Duration(req.SessionTime * float64(time.Second)))
	avgCPS := clicksPerSecond
	session := &models.ClickerSession{
		ID:         uuid.New(),
		UserID:     userID,
		StartTime:  sessionStart,
		EndTime:    &now,
		ClickCount: req.ClickCount,
		AverageCPS: &avgCPS,
		CreatedAt:  now,
	}

	// Сохраняем сессию
	err = s.repo.SaveSession(ctx, session)
	if err != nil {
		return nil, err
	}

	// Обновляем таблицу лидеров
	leaderboardEntry := &models.LeaderboardEntry{
		ID:     uuid.New(),
		UserID: userID,
		Score:  stats.TotalClicks,
	}

	err = s.repo.UpdateLeaderboard(ctx, leaderboardEntry)
	if err != nil {
		return nil, err
	}

	return &models.ClickResponse{
		TotalClicks: stats.TotalClicks,
		Status:      "accepted",
	}, nil
}

// GetStats получает статистику пользователя
func (s *ClickerService) GetStats(ctx context.Context, userID uuid.UUID) (*models.StatsResponse, error) {
	// Получаем статистику пользователя
	stats, err := s.repo.GetOrCreateStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем последние сессии
	sessions, err := s.repo.GetRecentSessions(ctx, userID, 10)
	if err != nil {
		return nil, err
	}

	return &models.StatsResponse{
		Stats:          *stats,
		RecentSessions: sessions,
	}, nil
}

// GetLeaderboard получает таблицу лидеров
func (s *ClickerService) GetLeaderboard(ctx context.Context, userID uuid.UUID, limit, offset int) (*models.LeaderboardResponse, error) {
	// Получаем таблицу лидеров
	entries, err := s.repo.GetLeaderboard(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Получаем ранг пользователя
	userRank, err := s.repo.GetUserRank(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.LeaderboardResponse{
		Entries:  entries,
		UserRank: userRank,
	}, nil
}

// RefreshLeaderboard обновляет ранги в таблице лидеров
func (s *ClickerService) RefreshLeaderboard(ctx context.Context) error {
	return s.repo.RefreshLeaderboardRanks(ctx)
}

// validateClickRequest проверяет запрос на сохранение кликов на наличие признаков читерства
func (s *ClickerService) validateClickRequest(req *models.ClickRequest) error {
	// Проверка скорости кликов
	clicksPerSecond := float64(req.ClickCount) / req.SessionTime
	if clicksPerSecond > MaxClicksPerSecond {
		return ErrInvalidClickRate
	}

	// Проверка времени сессии
	if req.SessionTime <= 0 || req.SessionTime > MaxSessionTime {
		return ErrInvalidSession
	}

	// Проверка расхождения времени клиента и сервера
	clientTime := time.Unix(req.ClientTimestamp, 0)
	serverTime := time.Now()
	timeDiff := math.Abs(serverTime.Sub(clientTime).Seconds())

	if timeDiff > MaxTimeDrift {
		return ErrInvalidTime
	}

	return nil
}
