package database

import (
	"database/sql"
	"fmt"

	"vibecoded/internal/models"
)

func (r *Repo) FindPasswordCollision(hashedPassword string) (string, error) {
	query := `
SELECT
	id, username, password, isAdmin, createdAt, updatedAt
FROM users
WHERE password = ?;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("r.db.Prepare: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(hashedPassword)

	var user models.UserDBItem
	if err := row.Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("row.Scan: %w", err)
	}

	return user.Username, nil
}

func (r *Repo) NewUser(newUser *models.UserInternal) (*models.UserDBItem, error) {
	query := `
INSERT INTO users (
	username, password, isAdmin
) VALUES (?, ?, ?)
RETURNING
	id, username, password, isAdmin, createdAt, updatedAt;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("r.db.Prepare: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(
		newUser.Username,
		newUser.HashedPassword,
		newUser.IsAdmin)
	var user models.UserDBItem
	if err := row.Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	return &user, nil
}

func (r *Repo) Login(username, hashedPassword string) (*models.UserDBItem, error) {
	query := `
SELECT
	id, username, password, isAdmin, createdAt, updatedAt
FROM users
WHERE username = ? AND password = ?`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("r.db.Prepare: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(username, hashedPassword)

	var user models.UserDBItem
	if err := row.Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	return &user, nil
}

func (r *Repo) GetUserDBItemByUsername(username string) (*models.UserDBItem, error) {
	query := `
SELECT
	id, username, password, isAdmin, createdAt, updatedAt
FROM users
WHERE username = ?;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("r.db.Prepare: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(username)

	var user models.UserDBItem
	if err := row.Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	return &user, nil
}

func (r *Repo) GetUserDBItemByID(id int) (*models.UserDBItem, error) {
	query := `
SELECT
	id, username, password, isAdmin, createdAt, updatedAt
FROM users
WHERE id = ?;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("r.db.Prepare: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	var user models.UserDBItem
	if err := row.Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	return &user, nil
}
