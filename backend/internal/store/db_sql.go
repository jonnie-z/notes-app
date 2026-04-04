package store

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	DB *sql.DB
}

func NewSQLiteStore(dsn string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &SQLiteStore{
		DB: db,
	}

	store.DB.Exec(`CREATE TABLE IF NOT EXISTS notes (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	body TEXT NOT NULL
	);`)

	return store, nil
}

func (s *SQLiteStore) GetAll() ([]Note, error) { return []Note{}, nil }
func (s *SQLiteStore) Create(text string) (Note, error) { return Note{}, nil }
func (s *SQLiteStore) Update(id int, text string) (Note, error) { return Note{}, nil }
func (s *SQLiteStore) Delete(id int) error { return nil }