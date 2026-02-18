package view

import (
	"image/color"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-log/internal/domain"
	"github.com/Pedro-0101/mental-log/internal/service"
)

type autoSaveEntry struct {
	widget.Entry
	onTrigger func()
}

func newAutoSaveEntry(multiline bool, onTrigger func()) *autoSaveEntry {
	e := &autoSaveEntry{onTrigger: onTrigger}
	e.MultiLine = multiline
	e.ExtendBaseWidget(e)
	return e
}

func (e *autoSaveEntry) TypedKey(key *fyne.KeyEvent) {
	e.Entry.TypedKey(key)
	if key.Name == fyne.KeySpace || key.Name == fyne.KeyEnter || key.Name == fyne.KeyReturn {
		if e.onTrigger != nil {
			e.onTrigger()
		}
	}
}

type NoteView struct {
	noteService   *service.NoteService
	notes         []domain.Note
	list          *widget.List
	listContainer *fyne.Container
	mainContent   *fyne.Container
}

func NewNoteView(noteService *service.NoteService) *NoteView {
	return &NoteView{noteService: noteService}
}

func (n *NoteView) CreateNoteButton(window fyne.Window) *widget.Button {
	return widget.NewButton("Add Note", func() {
		n.CreateNote(window)
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
				n.RefreshList()
			}
		},
		window,
	)
}

func (n *NoteView) RefreshList() {
	notes, err := n.noteService.FindAll()
	if err != nil {
		slog.Error("Error fetching notes", "error", err)
		return
	}
	n.notes = notes

	if len(n.notes) == 0 {
		placeholder := widget.NewLabel("No notes found")
		placeholder.Alignment = fyne.TextAlignLeading
		n.listContainer.Objects = []fyne.CanvasObject{placeholder}
	} else {
		if n.list == nil {
			n.list = widget.NewList(
				func() int {
					return len(n.notes)
				},
				func() fyne.CanvasObject {
					label := widget.NewLabel("")
					label.Alignment = fyne.TextAlignLeading

					btn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
					btn.Importance = widget.DangerImportance

					return container.NewBorder(nil, nil, nil, btn, label)
				},
				func(i widget.ListItemID, o fyne.CanvasObject) {
					container := o.(*fyne.Container)
					var btn *widget.Button
					var label *widget.Label

					for _, obj := range container.Objects {
						if b, ok := obj.(*widget.Button); ok {
							btn = b
						} else if l, ok := obj.(*widget.Label); ok {
							label = l
						}
					}

					if label != nil {
						label.SetText(n.notes[i].Title)
					}
					if btn != nil {
						id := n.notes[i].ID
						btn.OnTapped = func() {
							n.DeleteNote(id)
						}
					}
				},
			)
			n.list.OnSelected = func(id widget.ListItemID) {
				content := n.RenderNoteContent(n.notes[id].ID)
				if content != nil {
					n.mainContent.Objects = []fyne.CanvasObject{content}
					n.mainContent.Refresh()
				}
			}
		}
		n.listContainer.Objects = []fyne.CanvasObject{n.list}
		n.list.Refresh()
	}
	n.listContainer.Refresh()
}

func (n *NoteView) UpdateNote(id int64, title, content string) {
	// Find the old note to see if title changed
	var oldTitle string
	for _, note := range n.notes {
		if note.ID == id {
			oldTitle = note.Title
			break
		}
	}

	note := &domain.Note{
		ID:      id,
		Title:   title,
		Content: content,
	}

	err := n.noteService.UpdateNote(note)
	if err != nil {
		slog.Error("Error updating note", "error", err)
		return
	}

	slog.Info("Note updated (sync)", "id", id)

	// Only refresh the list if the title changed
	if oldTitle != title {
		n.RefreshList()
	}
}

func (n *NoteView) DeleteNote(id int64) {
	err := n.noteService.DeleteNote(id)
	if err != nil {
		slog.Error("Error deleting note", "error", err)
		return
	}
	slog.Info("Note deleted", "note", id)
	n.RefreshList()
}

func (n *NoteView) RenderNoteList(notes []domain.Note, addNoteButton *widget.Button) fyne.CanvasObject {
	n.notes = notes
	n.listContainer = container.NewStack()
	n.mainContent = container.NewStack(widget.NewLabel("Select a note to view its content"))

	title := canvas.NewText("Notes", theme.ForegroundColor())
	title.TextStyle.Bold = true
	title.TextSize = 20
	header := container.NewPadded(title)

	sidebarSpacer := canvas.NewRectangle(color.Transparent)
	sidebarSpacer.SetMinSize(fyne.NewSize(250, 0))

	n.RefreshList()

	sidebar := container.NewHBox(
		container.NewStack(sidebarSpacer, container.NewBorder(container.NewVBox(header, addNoteButton), nil, nil, nil, n.listContainer)),
		widget.NewSeparator(),
	)

	return container.NewBorder(nil, nil, sidebar, nil, n.mainContent)
}

func (n *NoteView) RenderNoteContent(noteId int64) fyne.CanvasObject {
	slog.Info("Rendering note content", "noteId", noteId)

	note, err := n.noteService.FindByID(noteId)
	if err != nil {
		slog.Error("Error fetching note", "error", err)
		return nil
	}

	titleEntry := newAutoSaveEntry(false, nil)
	titleEntry.SetText(note.Title)
	titleEntry.TextStyle.Bold = true

	contentEntry := newAutoSaveEntry(true, nil)
	contentEntry.SetText(note.Content)
	contentEntry.SetPlaceHolder("Start typing...")

	saveFunc := func() {
		n.UpdateNote(note.ID, titleEntry.Text, contentEntry.Text)
	}

	titleEntry.onTrigger = saveFunc
	contentEntry.onTrigger = saveFunc

	editArea := container.NewBorder(titleEntry, nil, nil, nil, contentEntry)

	return container.NewPadded(editArea)
}
