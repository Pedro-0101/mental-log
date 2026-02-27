package service

import "github.com/Pedro-0101/mental-dump/internal/domain"

type EntryService struct {
	entries []domain.Entry
}

func NewEntryService() *EntryService {
	return &EntryService{}
}

func (e *EntryService) CreateEntry(entry domain.Entry) {
	e.entries = append(e.entries, entry)
}

func (e *EntryService) GetEntries() []domain.Entry {
	return e.entries
}
