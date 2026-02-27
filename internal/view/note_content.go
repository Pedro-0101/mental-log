package view

import (
	"log/slog"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-dump/internal/domain"
	"github.com/Pedro-0101/mental-dump/internal/helper"
)

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
