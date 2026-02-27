package config

import (
	"fmt"
	"log/slog"

	"github.com/Pedro-0101/mental-dump/internal/domain"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {

	slog.Info("Starting database")

	db, err := gorm.Open(sqlite.Open("./mental-dump.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	slog.Info("Database connected")

	slog.Info("Auto migrating database tables")

	err = db.AutoMigrate(
		&domain.Status{},
		&domain.User{},
		&domain.Folder{},
		&domain.Note{},
		&domain.Entry{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate tables: %w", err)
	}

	slog.Info("Database tables initialized successfully")

	return db, nil
}
