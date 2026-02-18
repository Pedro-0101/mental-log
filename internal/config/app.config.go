package config

import (
	"database/sql"

	"github.com/Pedro-0101/mental-log/internal/repo"
	"github.com/Pedro-0101/mental-log/internal/service"
)

type App struct {
	DB *sql.DB

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

	// Fyne
	fyneConfig := NewFyneConfig("Mental Log", "", "dark", 12, 800, 600, noteService)

	notes, err := noteService.FindAll()
	if err != nil {
		return nil, err
	}

	fyneConfig.Start(notes)

	return &App{
		DB:          db,
		NoteRepo:    noteRepo,
		NoteService: noteService,
	}, nil
}

func (a *App) Close() error {
	return a.DB.Close()
}
