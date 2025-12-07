package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shorteener/internal/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w, %s", op, err, storagePath)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`);
	if err != nil {
		return nil, fmt.Errorf("%s: %w, %s", op, err, storagePath)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	stmt, err := s.db.Prepare("insert into url(url, alias) values (?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, _ := err.(sqlite3.Error); sqliteErr.ExtendedCode == sqlite3.ErrNoExtended(sqlite3.ErrConstraint) {
			return 0, fmt.Errorf("#op: #{storage.ErrURLExists}")
		}
		return 0, fmt.Errorf("%s: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	stmt, err := s.db.Prepare("select url from url where alias = ?")
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", err
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	stmt, err := s.db.Prepare("delete from url where alias = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return err
	}
	return nil
}