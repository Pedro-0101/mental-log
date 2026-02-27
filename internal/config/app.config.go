package config

import (
	"gorm.io/gorm"

	"github.com/Pedro-0101/mental-dump/internal/repo"
	"github.com/Pedro-0101/mental-dump/internal/service"
)

type App struct {
	DB *gorm.DB

	NoteRepo      *repo.NoteRepo
	NoteService   *service.NoteService
	FolderRepo    *repo.FolderRepo
	FolderService *service.FolderService
}

func NewApp() (*App, error) {

	// Database
	db, err := NewDB()
	if err != nil {
		return nil, err
	}

	// Repositories
	noteRepo := repo.NewNoteRepo(db)
	folderRepo := repo.NewFolderRepo(db)

	// Services
	folderService := service.NewFolderService(folderRepo)
	noteService := service.NewNoteService(noteRepo, folderService)
	entryService := service.NewEntryService()

	// Fyne Config
	name := "Mental dump"
	version := "0.0.1"
	theme := "dark"
	fontSize := float32(12.0)
	width := 800
	height := 600

	fyneConfig := NewFyneConfig(name, version, theme, fontSize, width, height, noteService, folderService)

	fyneConfig.Start(entryService)

	return &App{
		DB:            db,
		NoteRepo:      noteRepo,
		NoteService:   noteService,
		FolderRepo:    folderRepo,
		FolderService: folderService,
	}, nil
}

func (a *App) Close() error {
	sqlDB, err := a.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
