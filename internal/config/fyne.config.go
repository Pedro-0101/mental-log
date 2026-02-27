package config

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/Pedro-0101/mental-dump/internal/service"
	"github.com/Pedro-0101/mental-dump/internal/view/entry"
	"github.com/Pedro-0101/mental-dump/internal/view/folder"
)

type FyneConfig struct {
	Title    string
	Icon     string
	Theme    fyne.Theme
	FontSize float32
	Size     fyneSize

	noteService   *service.NoteService
	folderService *service.FolderService
}

type fyneSize struct {
	Width  float32
	Height float32
}

func NewFyneConfig(
	title string,
	icon string,
	themeString string,
	fontSize float32,
	width int,
	height int,

	noteService *service.NoteService,
	folderService *service.FolderService,
) *FyneConfig {

	var themeValue fyne.Theme
	if themeString == "dark" {
		themeValue = theme.DarkTheme()
	} else {
		themeValue = theme.LightTheme()
	}

	return &FyneConfig{

		Title:    title,
		Icon:     icon,
		Theme:    themeValue,
		FontSize: fontSize,
		Size:     fyneSize{Width: float32(width), Height: float32(height)},

		noteService:   noteService,
		folderService: folderService,
	}
}

func (f *FyneConfig) Start(entryService *service.EntryService) {

	slog.Info("Starting Fyne app")

	a := app.New()
	a.Settings().SetTheme(f.Theme)

	window := a.NewWindow(f.Title)
	window.Resize(fyne.NewSize(f.Size.Width, f.Size.Height))

	// Configure entry view
	var entryWindowSize fyneSize = fyneSize{
		Width:  f.Size.Width * 0.80,
		Height: f.Size.Height * 0.99,
	}

	// Configure folder view
	var folderWindowSize fyneSize = fyneSize{
		Width:  f.Size.Width * 0.20,
		Height: f.Size.Height * 0.99,
	}

	folderView := folder.NewListFolder(f.folderService, f.noteService, folderWindowSize.Width, folderWindowSize.Height, window)
	entryView := entry.NewEntryView(entryService, entryWindowSize.Width, entryWindowSize.Height)

	content := container.NewHSplit(folderView.RenderList(), entryView.RenderEntry())
	content.SetOffset(0.2)

	window.SetContent(content)
	window.ShowAndRun()
}
