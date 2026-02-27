package view

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (n *NoteView) CreateNoteButton(window fyne.Window) *widget.Button {
	return widget.NewButton("Add Note", func() {
		n.CreateNote(window)
	})
}

func (n *NoteView) CreateSaveButton() *widget.Button {
	return widget.NewButton("Save", func() {
		n.dbService.SaveDatabase()
	})
}

func (n *NoteView) CreateNote(window fyne.Window) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter note title...")

	dialog.ShowForm(
		"New Note",
		"Create",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Title", entry),
		},
		func(response bool) {
			if response {
				slog.Info("Creating note", "Note", entry.Text)
				newNote, err := n.noteService.CreateNote(entry.Text)
				if err != nil {
					slog.Error("Error creating note", "error", err)
					return
				}
				slog.Info("Note created", "note", newNote.Title)
				n.mainContent.Objects = []fyne.CanvasObject{n.RenderNoteContent(newNote.ID)}
				n.mainContent.Refresh()
				n.RefreshList()
			}
		},
		window,
	)
}
