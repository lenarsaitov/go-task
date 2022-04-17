package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lenarsaitov/go-task/internals/config"
	"github.com/lenarsaitov/go-task/pkg/logging"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) *sqlx.DB {
	logger := logging.GetLogger()

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Name, cfg.Postgres.Password, cfg.Postgres.SSL))
	if err != nil {
		logger.Fatalf("Failed to open postgres connection: %s", err)
		return nil
	}

	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		logger.Fatalf("Failed cheking ping postgres connection: %s", err)
		return nil
	}

	return db
}
