package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	return db, nil
}

func InitDB(db *sql.DB) error {
	initSQL := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL UNIQUE,
	isAdmin BOOL NOT NULL DEFAULT false,
	createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS notes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL DEFAULT '',
	content TEXT NOT NULL DEFAULT '',
	userId INTEGER NOT NULL,
	createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (userId) REFERENCES users (id)
);

CREATE TRIGGER IF NOT EXISTS update_updatedAt_users
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users
    SET updatedAt = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
CREATE TRIGGER IF NOT EXISTS update_updatedAt_notes
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
    UPDATE notes
    SET updatedAt = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;`

	if _, err := db.Exec(initSQL); err != nil {
		return fmt.Errorf("failed to execute init SQL: %w", err)
	}
	return nil
}
