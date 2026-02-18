package repo

import (
	"database/sql"
	"time"

	"github.com/Pedro-0101/mental-log/internal/domain"
)

type NoteRepo struct {
	DB *sql.DB
}

func NewNoteRepo(db *sql.DB) *NoteRepo {
	return &NoteRepo{DB: db}
}

func (r *NoteRepo) Create(title string) (*domain.Note, error) {
	newNote := &domain.Note{
		Title: title,
	}

	result, err := r.DB.Exec("INSERT INTO notes (title) VALUES (?)", title)
	if err != nil {
		return nil, err
	}

	newNote.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return newNote, nil
}

func (r *NoteRepo) FindAll() ([]domain.Note, error) {
	notes := []domain.Note{}

	rows, err := r.DB.Query("SELECT id, title, created_at, updated_at FROM notes")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var note domain.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *NoteRepo) FindByID(id int64) (*domain.Note, error) {
	var note domain.Note

	row := r.DB.QueryRow("SELECT id, title, COALESCE(content, ''), created_at, updated_at FROM notes WHERE id = ?", id)
	if err := row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepo) Update(note *domain.Note) error {
	_, err := r.DB.Exec("UPDATE notes SET content = ?, updated_at = ? WHERE id = ?", note.Content, time.Now(), note.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *NoteRepo) Delete(id int64) error {
	_, err := r.DB.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
