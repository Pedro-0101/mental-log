package repo

import (
	"github.com/Pedro-0101/mental-dump/internal/domain"
	"gorm.io/gorm"
)

type NoteRepo struct {
	DB *gorm.DB
}

func NewNoteRepo(db *gorm.DB) *NoteRepo {
	return &NoteRepo{DB: db}
}

func (r *NoteRepo) Create(title, tags string, folderID *int64, userID int64) (*domain.Note, error) {
	newNote := &domain.Note{
		Title:    title,
		Tags:     tags,
		FolderID: folderID,
		UserID:   userID,
		StatusID: 1,
	}

	if err := r.DB.Create(newNote).Error; err != nil {
		return nil, err
	}

	return newNote, nil
}

func (r *NoteRepo) FindAll() ([]domain.Note, error) {
	var notes []domain.Note

	if err := r.DB.Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *NoteRepo) FindByID(id int64) (*domain.Note, error) {
	var note domain.Note

	if err := r.DB.First(&note, id).Error; err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *NoteRepo) Update(note *domain.Note) error {
	// Let GORM save the provided note fields
	return r.DB.Save(note).Error
}

func (r *NoteRepo) Delete(id int64) error {
	return r.DB.Delete(&domain.Note{}, id).Error
}

func (r *NoteRepo) FindByFolderID(folderID int64) ([]domain.Note, error) {
	var notes []domain.Note
	if err := r.DB.Where("folder_id = ?", folderID).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// FindRootNotes returns notes that are not inside any folder.
func (r *NoteRepo) FindRootNotes() ([]domain.Note, error) {
	var notes []domain.Note
	if err := r.DB.Where("folder_id IS NULL").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}
