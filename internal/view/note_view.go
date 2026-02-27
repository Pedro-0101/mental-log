package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-dump/internal/domain"
	"github.com/Pedro-0101/mental-dump/internal/service"
)

type NoteView struct {
	noteService   *service.NoteService
	dbService     *service.DBService
	notes         []domain.Note
	list          *widget.List
	listContainer *fyne.Container
	mainContent   *fyne.Container
}

func NewNoteView(noteService *service.NoteService) *NoteView {
	return &NoteView{noteService: noteService, dbService: service.NewDBService()}
}
