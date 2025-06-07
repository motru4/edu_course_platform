package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresDB создает новое подключение к базе данных PostgreSQL
func NewPostgresDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки подключения к базе данных: %w", err)
	}

	return db, nil
}
