package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

// NewPostgresDB : подключение к базе данных
func NewPostgresDB(config Config) (*pgxpool.Pool, error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", config.Username, config.Password, config.Host, config.Port, config.Name, config.SSLMode)
	dbPool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}
	return dbPool, nil
}
