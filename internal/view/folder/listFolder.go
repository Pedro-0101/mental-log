package folder

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-dump/internal/service"
)

type ListFolder struct {
	folderService *service.FolderService
	noteService   *service.NoteService
	width         float32
	height        float32
	tree          *widget.Tree
	window        fyne.Window
}

func NewListFolder(
	folderService *service.FolderService,
	noteService *service.NoteService,
	width float32,
	height float32,
	window fyne.Window,
) *ListFolder {
	return &ListFolder{
		folderService: folderService,
		noteService:   noteService,
		width:         width,
		height:        height,
		window:        window,
	}
}

func (l *ListFolder) RenderList() fyne.CanvasObject {
	// Build the tree first so the button callback can refresh it
	l.tree = widget.NewTree(
		// childUIDs: returns children of a given tree node
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			if id == "" {
				// Root: show root folders + root notes
				var ids []widget.TreeNodeID

				folders, _ := l.folderService.FindRootFolders()
				for _, f := range folders {
					ids = append(ids, fmt.Sprintf("folder:%d", f.ID))
				}

				notes, _ := l.noteService.FindRootNotes()
				for _, n := range notes {
					ids = append(ids, fmt.Sprintf("note:%d", n.ID))
				}
				return ids
			}

			parts := strings.SplitN(id, ":", 2)
			if parts[0] == "folder" {
				folderID, _ := strconv.ParseInt(parts[1], 10, 64)

				var ids []widget.TreeNodeID

				// Subfolders
				subfolders, _ := l.folderService.FindByParentID(folderID)
				for _, sf := range subfolders {
					ids = append(ids, fmt.Sprintf("folder:%d", sf.ID))
				}

				// Notes inside this folder
				notes, _ := l.noteService.FindByFolderID(folderID)
				for _, n := range notes {
					ids = append(ids, fmt.Sprintf("note:%d", n.ID))
				}
				return ids
			}

			return nil
		},
		// isBranch
		func(id widget.TreeNodeID) bool {
			return id == "" || strings.HasPrefix(id, "folder:")
		},
		// create
		func(branch bool) fyne.CanvasObject {
			if branch {
				return container.NewHBox(widget.NewIcon(theme.FolderIcon()), widget.NewLabel("Folder"))
			}
			return container.NewHBox(widget.NewIcon(theme.FileIcon()), widget.NewLabel("Note"))
		},
		// update
		func(id widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			box := obj.(*fyne.Container)
			label := box.Objects[1].(*widget.Label)
			icon := box.Objects[0].(*widget.Icon)

			parts := strings.SplitN(id, ":", 2)
			if parts[0] == "folder" {
				folderID, _ := strconv.ParseInt(parts[1], 10, 64)
				folders, _ := l.folderService.FindAll()
				for _, f := range folders {
					if f.ID == folderID {
						label.SetText(f.Title)
						break
					}
				}
				icon.SetResource(theme.FolderIcon())
			} else if parts[0] == "note" {
				noteID, _ := strconv.ParseInt(parts[1], 10, 64)
				note, err := l.noteService.FindByID(noteID)
				if err == nil {
					label.SetText(note.Title)
				}
				icon.SetResource(theme.FileIcon())
			}
		},
	)

	l.tree.OnSelected = func(id widget.TreeNodeID) {
		slog.Info("Tree node selected", "id", id)
	}

	// Create folder button with modal
	createFolderBtn := NewCreateFolderButton(l.folderService, l.window, func() {
		l.tree.Refresh()
	})
	createFolderBtn.Button.Resize(fyne.NewSize(l.width, l.height))

	// Create note button with modal
	createNoteBtn := NewCreateNoteButton(l.noteService, l.folderService, l.window, func() {
		l.tree.Refresh()
	})
	createNoteBtn.Button.Resize(fyne.NewSize(l.width, l.height))

	// Render actions buttons
	actionButtons := container.NewVBox(
		widget.NewLabelWithStyle("Actions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		createFolderBtn.Button,
		createNoteBtn.Button,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Explorer", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	return container.NewBorder(actionButtons, nil, nil, nil, l.tree)
}
