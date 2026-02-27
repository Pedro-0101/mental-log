package service

import (
	"log/slog"

	"github.com/Pedro-0101/mental-dump/internal/domain"
	"github.com/Pedro-0101/mental-dump/internal/repo"
)

type NoteService struct {
	noteRepo      *repo.NoteRepo
	folderService *FolderService
}

func NewNoteService(noteRepo *repo.NoteRepo, folderService *FolderService) *NoteService {
	return &NoteService{noteRepo: noteRepo, folderService: folderService}
}

func (s *NoteService) CreateNote(title, tags string, folderID *int64) (*domain.Note, error) {
	mockedUserID := int64(1)
	newNote, err := s.noteRepo.Create(title, tags, folderID, mockedUserID)
	if err != nil {
		slog.Error("Error creating note", "error", err)
		return nil, err
	}

	// Propagate tags to parent folders
	s.folderService.PropagateTags(folderID, tags)

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

func (s *NoteService) FindByFolderID(folderID int64) ([]domain.Note, error) {
	return s.noteRepo.FindByFolderID(folderID)
}

func (s *NoteService) FindRootNotes() ([]domain.Note, error) {
	return s.noteRepo.FindRootNotes()
}
