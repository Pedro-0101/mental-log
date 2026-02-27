package folder

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-dump/internal/domain"
	"github.com/Pedro-0101/mental-dump/internal/service"
)

type CreateFolderButton struct {
	Button *widget.Button
}

func NewCreateFolderButton(
	folderService *service.FolderService,
	window fyne.Window,
	onCreated func(),
) *CreateFolderButton {
	btn := widget.NewButton("+ New Folder", func() {
		showCreateFolderModal(folderService, window, onCreated)
	})
	return &CreateFolderButton{Button: btn}
}

func showCreateFolderModal(
	folderService *service.FolderService,
	window fyne.Window,
	onCreated func(),
) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Folder name")

	// Build the parent folder dropdown
	allFolders, err := folderService.FindAll()
	if err != nil {
		slog.Error("Error loading folders for modal", "error", err)
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
		widget.NewFormItem("Parent Folder", parentSelect),
	}

	d := dialog.NewForm(
		"Create Folder",
		"Create",
		"Cancel",
		formItems,
		func(ok bool) {
			if !ok {
				return
			}
			title := titleEntry.Text
			if title == "" {
				slog.Warn("Folder title is empty, skipping creation")
				return
			}

			var parentID *int64
			if parentSelect.Selected != "(None - Root)" {
				parentID = parentMap[parentSelect.Selected]
			}

			_, err := folderService.CreateFolder(title, parentID)
			if err != nil {
				slog.Error("Error creating folder", "error", err)
				return
			}
			slog.Info("Folder created", "title", title)

			if onCreated != nil {
				onCreated()
			}
		},
		window,
	)

	d.Resize(fyne.NewSize(400, 200))
	d.Show()
}
