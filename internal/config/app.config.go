package config

import (
	"gorm.io/gorm"

	"github.com/Pedro-0101/mental-dump/internal/repo"
	"github.com/Pedro-0101/mental-dump/internal/service"
)

type App struct {
	DB *gorm.DB

	NoteRepo    *repo.NoteRepo
	NoteService *service.NoteService
}

func NewApp() (*App, error) {

	// Database
	db, err := NewDB()
	if err != nil {
		return nil, err
	}

	// Repositories
	noteRepo := repo.NewNoteRepo(db)

	// Services
	noteService := service.NewNoteService(noteRepo)
	entryService := service.NewEntryService()

	// Fyne Config
	name := "Mental dump"
	version := "0.0.1"
	theme := "dark"
	fontSize := float32(12.0)
	width := 800
	height := 600

	fyneConfig := NewFyneConfig(name, version, theme, fontSize, width, height, noteService)

	fyneConfig.Start(entryService)

	return &App{
		DB:          db,
		NoteRepo:    noteRepo,
		NoteService: noteService,
	}, nil
}

func (a *App) Close() error {
	sqlDB, err := a.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
