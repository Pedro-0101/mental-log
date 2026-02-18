package config

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {

	slog.Info("Starting database")

	db, err := sql.Open("sqlite3", "./mental-log.db")
	if err != nil {
		return nil, err
	}

	slog.Info("Database connected")

	slog.Info("Creating database tables")

	slog.Info("Creating notes table")
	db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)

	slog.Info("Notes table created")

	slog.Info("Database tables created")

	slog.Info("Database started")

	return db, nil
}
