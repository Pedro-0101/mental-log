package entry

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Pedro-0101/mental-dump/internal/service"
)

type EntryView struct {
	service     *service.EntryService
	dumpContent string
	width       float32
	height      float32
}

type submitEntry struct {
	widget.Entry
}

func NewSubmitEntry() *submitEntry {
	e := &submitEntry{}
	e.MultiLine = true
	e.ExtendBaseWidget(e)
	return e
}

func (e *submitEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
		if e.OnSubmitted != nil {
			slog.Info("Submitting entry")
			e.OnSubmitted(e.Text)
		}
		return
	}
	e.Entry.TypedKey(key)
}

func NewEntryView(service *service.EntryService, width float32, height float32) *EntryView {
	return &EntryView{service: service, width: width, height: height}
}

func (e *EntryView) RenderEntry() fyne.CanvasObject {

	container := container.NewWithoutLayout()

	// Text entry for dumping thoughts
	entry := NewSubmitEntry()

	entry.Resize(fyne.NewSize(e.width, e.height))

	entry.SetPlaceHolder("Enter your entry here...")

	// On submitted save entry on sorting stack and clear entry
	entry.OnSubmitted = func(s string) {
		e.saveEntry(entry)
	}

	container.Add(entry)

	return container

}

func (e *EntryView) saveEntry(entry *submitEntry) {

	// Add entry on save entry stack and clear entry
	slog.Info("Saving dump", "dump", entry.Text)
	e.dumpContent = ""
	entry.SetText("")

}
