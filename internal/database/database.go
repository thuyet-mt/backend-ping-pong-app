package database

import (
	"database/sql"
	"fmt"

	log "github.com/jeanphorn/log4go"
	_ "github.com/lib/pq"

	// internal packages
	"backend-ping-pong-app/internal/config"
)

func OpenPostgresDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	log.Info("Connecting to PostgreSQL database...")
	log.Debug("Postgres config: host=%s port=%s db=%s user=%s",
		cfg.Host, cfg.Port, cfg.Name, cfg.User,
	)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Error("Failed to open postgres connection: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping postgres database: %v", err)
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Info("Successfully connected to PostgreSQL database")
	return db, nil
}
