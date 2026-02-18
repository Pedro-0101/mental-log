package service

import (
	"log/slog"

	"github.com/Pedro-0101/mental-log/internal/domain"
	"github.com/Pedro-0101/mental-log/internal/repo"
)

type NoteService struct {
	noteRepo *repo.NoteRepo
}

func NewNoteService(noteRepo *repo.NoteRepo) *NoteService {
	return &NoteService{noteRepo: noteRepo}
}

func (s *NoteService) CreateNote(title string) (*domain.Note, error) {
	newNote, err := s.noteRepo.Create(title)
	if err != nil {
		slog.Error("Error creating note", "error", err)
		return nil, err
	}

	return newNote, nil
}

func (s *NoteService) FindAll() ([]domain.Note, error) {
	return s.noteRepo.FindAll()
}

func (s *NoteService) FindByID(id int64) (*domain.Note, error) {

	slog.Info("Finding note by id", "id", id)

	return s.noteRepo.FindByID(id)
}

func (s *NoteService) UpdateNote(note *domain.Note) error {
	return s.noteRepo.Update(note)
}

func (s *NoteService) DeleteNote(id int64) error {
	return s.noteRepo.Delete(id)
}
