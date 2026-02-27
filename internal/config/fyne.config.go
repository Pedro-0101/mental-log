package config

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/Pedro-0101/mental-dump/internal/domain"
	"github.com/Pedro-0101/mental-dump/internal/service"
	"github.com/Pedro-0101/mental-dump/internal/view"
)

type FyneConfig struct {
	Title    string
	Icon     string
	Theme    fyne.Theme
	FontSize float32
	Width    int
	Height   int

	noteService *service.NoteService
}

func NewFyneConfig(
	title string,
	icon string,
	themeString string,
	fontSize float32,
	width int,
	height int,

	noteService *service.NoteService,
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
		Width:    width,
		Height:   height,

		noteService: noteService,
	}
}

func (f *FyneConfig) Start(notes []domain.Note) {

	slog.Info("Starting Fyne app")

	a := app.New()
	a.Settings().SetTheme(f.Theme)

	window := a.NewWindow(f.Title)
	window.Resize(fyne.NewSize(float32(f.Width), float32(f.Height)))

	noteView := view.NewNoteView(f.noteService)
	addNoteButton := noteView.CreateNoteButton(window)
	saveButton := noteView.CreateSaveButton()

	layout := noteView.RenderNoteList(notes, addNoteButton, saveButton)

	window.SetContent(layout)
	window.ShowAndRun()
}
