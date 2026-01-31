package database

import (
	"fmt"
	"strings"

	"vibecoded/internal/models"
)

func (r *Repo) NewNote(newNote *models.NoteInternal) (*models.NoteDBItem, error) {
	query := `
INSERT INTO notes (
	title, content, userId
) VALUES (?, ?, ?)
RETURNING
	id, title, content, userId, createdAt, updatedAt;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("r.db.Prepare: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(
		newNote.Title,
		newNote.Content,
		newNote.UserID)
	var note models.NoteDBItem
	if err := row.Scan(&note.ID, &note.Title, &note.Content, &note.UserID, &note.CreatedAt, &note.UpdatedAt); err != nil {
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	return &note, nil
}

func (r *Repo) UpdateNote(note *models.NoteUpdate) error {
	var b strings.Builder
	b.WriteString(`
UPDATE notes
SET`)

	var columnsToSet []string
	var values []any
	if note.Title != "" {
		columnsToSet = append(columnsToSet, "title = ?")
		values = append(values, note.Title)
	}
	if note.Content != "" {
		columnsToSet = append(columnsToSet, "content = ?")
		values = append(values, note.Content)
	}

	if len(columnsToSet) == 0 {
		return fmt.Errorf("no columns to update for note %d", note.ID)
	}
	b.WriteString(" " + strings.Join(columnsToSet, ", "))

	b.WriteString(" WHERE id = ?;")
	values = append(values, note.ID)

	stmt, err := r.db.Prepare(b.String())
	if err != nil {
		return fmt.Errorf("r.db.Prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("stmt.Exec: %w", err)
	}

	return nil
}

func (r *Repo) GetNoteDBItemsByUserID(userID int) ([]*models.NoteDBItem, error) {
	query := `
SELECT
	id, title, content, userId, createdAt, updatedAt
FROM notes
WHERE userId = ?;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("r.db.Prepare: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, fmt.Errorf("stmt.Query: %w", err)
	}

	var notes []*models.NoteDBItem
	for rows.Next() {
		var note models.NoteDBItem
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.UserID, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		notes = append(notes, &note)
	}

	return notes, nil
}
