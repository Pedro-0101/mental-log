package folder

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-dump/internal/domain"
	"github.com/Pedro-0101/mental-dump/internal/service"
)

type CreateNoteButton struct {
	Button *widget.Button
}

func NewCreateNoteButton(
	noteService *service.NoteService,
	folderService *service.FolderService,
	window fyne.Window,
	onCreated func(),
) *CreateNoteButton {
	btn := widget.NewButton("+ New Note", func() {
		showCreateNoteModal(noteService, folderService, window, onCreated)
	})
	return &CreateNoteButton{Button: btn}
}

func showCreateNoteModal(
	noteService *service.NoteService,
	folderService *service.FolderService,
	window fyne.Window,
	onCreated func(),
) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Note title")

	tagsEntry := widget.NewEntry()
	tagsEntry.SetPlaceHolder("Tags (comma-separated, e.g. go, fyne)")

	// Build the parent folder dropdown
	allFolders, err := folderService.FindAll()
	if err != nil {
		slog.Error("Error loading folders for note modal", "error", err)
		allFolders = []domain.Folder{}
	}

	parentOptions := []string{"(None - Root)"}
	parentMap := map[string]*int64{}
	for _, f := range allFolders {
		parentOptions = append(parentOptions, f.Title)
		id := f.ID
		parentMap[f.Title] = &id
	}

	parentSelect := widget.NewSelect(parentOptions, nil)
	parentSelect.SetSelected("(None - Root)")

	formItems := []*widget.FormItem{
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Tags", tagsEntry),
		widget.NewFormItem("Parent Folder", parentSelect),
	}

	d := dialog.NewForm(
		"Create Note",
		"Create",
		"Cancel",
		formItems,
		func(ok bool) {
			if !ok {
				return
			}
			title := titleEntry.Text
			if title == "" {
				slog.Warn("Note title is empty, skipping creation")
				return
			}

			tags := tagsEntry.Text

			var folderID *int64
			if parentSelect.Selected != "(None - Root)" {
				folderID = parentMap[parentSelect.Selected]
			}

			_, err := noteService.CreateNote(title, tags, folderID)
			if err != nil {
				slog.Error("Error creating note", "error", err)
				return
			}
			slog.Info("Note created", "title", title, "tags", tags)

			if onCreated != nil {
				onCreated()
			}
		},
		window,
	)

	d.Resize(fyne.NewSize(400, 250))
	d.Show()
}
