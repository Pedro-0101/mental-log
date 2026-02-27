package view

import (
	"image/color"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-dump/internal/domain"
)

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

func (n *NoteView) RenderNoteList(notes []domain.Note, addNoteButton, saveButton *widget.Button) fyne.CanvasObject {
	n.notes = notes
	n.listContainer = container.NewStack()
	n.mainContent = container.NewStack(widget.NewLabel("Select a note to view its content"))

	title := canvas.NewText("Notes", theme.Color(theme.ColorNameForeground))
	title.TextStyle.Bold = true
	title.TextSize = 20
	header := container.NewPadded(title)

	sidebarSpacer := canvas.NewRectangle(color.Transparent)
	sidebarSpacer.SetMinSize(fyne.NewSize(250, 0))

	n.RefreshList()

	sidebar := container.NewHBox(
		container.NewStack(sidebarSpacer, container.NewBorder(container.NewVBox(header, saveButton, addNoteButton), nil, nil, nil, n.listContainer)),
		widget.NewSeparator(),
	)

	return container.NewBorder(nil, nil, sidebar, nil, n.mainContent)
}
