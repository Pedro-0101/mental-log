package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
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
