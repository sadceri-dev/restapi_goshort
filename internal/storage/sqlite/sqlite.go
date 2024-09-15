package sqlite

import (
	"LINKSHORTENER/internal/storage"
	"LINKSHORTENER/internal/storage/sqlite"
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)


type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %W", op, err)
	}


	stmt, err := db.Prepare('
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx)alias ON url(alias);
	')
	if err != nil {
		return nil, fmt.Errorf("%s: %W", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %W", op, err)
	}

	return &Storage{db: db}, nil
}


func  (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err :=s.db.Prepare("INSERT INTO url(url,alias) values(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %W", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstrainUnique {
			return 0, fmt.Errorf("%s: %W", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %W", op, err)
	}

	return id, nil
}


func (s *Storage) GetURL(alias string) (string, error) {
    const op = "storage.sqlite.GetURL"

    stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
    if err != nil {
        return "", fmt.Errorf("%s: prepare statement: %w", op, err)
    }

    var resURL string
    
    err = stmt.QueryRow(alias).Scan(&resURL)
    if errors.Is(err, sql.ErrNoRows) {
        return "", storage.ErrURLNotFound
    }
    if err != nil {
        return "", fmt.Errorf("%s: execute statement: %w", op, err)
    }

    return resURL, nil
}
