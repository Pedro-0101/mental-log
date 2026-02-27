package config

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "modernc.org/sqlite"
)

func NewDB() (*sql.DB, error) {

	slog.Info("Starting database")

	db, err := sql.Open("sqlite", "./mental-log.db")
	if err != nil {
		return nil, err
	}

	slog.Info("Database connected")

	slog.Info("Creating database tables")

	slog.Info("Creating notes table")
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create notes table: %w", err)
	}

	slog.Info("Database tables initialized successfully")

	return db, nil
}
