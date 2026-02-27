package view

import (
	"image/color"
	"log/slog"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-dump/internal/domain"
	"github.com/Pedro-0101/mental-dump/internal/helper"
	"github.com/Pedro-0101/mental-dump/internal/service"
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
	if key.Name == fyne.KeyEnter || key.Name == fyne.KeyReturn {
		if e.onTrigger != nil {
			e.onTrigger()
		}
		return
	}
	e.Entry.TypedKey(key)
}

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

func (n *NoteView) UpdateNote(note *domain.Note) {
	err := n.noteService.UpdateNote(note)
	if err != nil {
		slog.Error("Error updating note", "error", err)
		return
	}

	slog.Info("Note updated (sync)", "id", note.ID)
}

func (n *NoteView) DeleteNote(id int64) {
	err := n.noteService.DeleteNote(id)
	if err != nil {
		slog.Error("Error deleting note", "error", err)
		return
	}
	slog.Info("Note deleted", "note", id)
	n.RefreshList()
	n.mainContent.Objects = []fyne.CanvasObject{widget.NewLabel("Select a note to view its content")}
	n.mainContent.Refresh()
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

type paragraphBlock struct {
	Timestamp string
	Text      string
}

func parseContentBlocks(content string) []paragraphBlock {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	re := regexp.MustCompile(`\[(\d{2}/\d{2} \d{2}:\d{2})\]`)
	indices := re.FindAllStringIndex(content, -1)
	matches := re.FindAllStringSubmatch(content, -1)

	if len(indices) == 0 {
		return []paragraphBlock{{Timestamp: "", Text: strings.TrimSpace(content)}}
	}

	var blocks []paragraphBlock
	for i, idx := range indices {
		ts := matches[i][1]
		textStart := idx[1]
		if textStart < len(content) && content[textStart] == '\n' {
			textStart++
		}

		var textEnd int
		if i+1 < len(indices) {
			textEnd = indices[i+1][0]
		} else {
			textEnd = len(content)
		}

		text := strings.TrimRight(content[textStart:textEnd], "\n")
		if text != "" {
			blocks = append(blocks, paragraphBlock{Timestamp: ts, Text: text})
		}
	}

	return blocks
}

func buildParagraphWidget(block paragraphBlock) fyne.CanvasObject {
	var dateLabel *canvas.Text
	if block.Timestamp != "" {
		dateLabel = canvas.NewText("["+block.Timestamp+"] | Pedro Paulino", color.NRGBA{R: 130, G: 130, B: 130, A: 255})
		dateLabel.TextSize = 11
		dateLabel.TextStyle.Italic = true
	}

	textLabel := widget.NewLabel(block.Text)
	textLabel.Wrapping = fyne.TextWrapWord

	if dateLabel != nil {
		return container.NewVBox(dateLabel, textLabel)
	}
	return container.NewVBox(textLabel)
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
	paragraphs := container.NewVBox()

	blocks := parseContentBlocks(note.Content)
	for _, block := range blocks {
		paragraphs.Add(buildParagraphWidget(block))
	}

	contentEntry := newAutoSaveEntry(true, nil)
	contentEntry.SetPlaceHolder("Digite e pressione Enter para salvar...")

	scrollContent := container.NewVBox(
		paragraphs,
		layout.NewSpacer(),
		contentEntry,
	)
	scrollArea := container.NewVScroll(scrollContent)

	commitParagraph := func() {
		text := strings.TrimSpace(contentEntry.Text)
		if text == "" {
			return
		}

		text = strings.ReplaceAll(text, "\n", " ")

		timestamp := helper.GetCurrentTimeStr()
		newBlock := "[" + timestamp + "]\n" + text + "\n"

		if note.Content != "" && !strings.HasSuffix(note.Content, "\n") {
			note.Content += "\n"
		}
		note.Content += newBlock

		n.UpdateNote(note)

		paragraphs.Add(buildParagraphWidget(paragraphBlock{Timestamp: timestamp, Text: text}))
		paragraphs.Refresh()
		contentEntry.SetText("")
		contentEntry.Refresh()

		scrollArea.ScrollToBottom()

		slog.Info("Paragraph committed", "noteId", noteId, "timestamp", timestamp)
	}

	contentEntry.onTrigger = commitParagraph

	titleEntry.onTrigger = func() {
		note.Title = titleEntry.Text
		err := n.noteService.UpdateNote(note)
		if err != nil {
			slog.Error("Error updating note title", "error", err)
		}
		n.RefreshList()
	}

	editArea := container.NewBorder(titleEntry, nil, nil, nil, scrollArea)

	return container.NewPadded(editArea)
}
