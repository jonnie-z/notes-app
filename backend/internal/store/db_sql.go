package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

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

func (s *SQLiteStore) GetAll() ([]Note, error) { 
	// TODO: sqlc?
	result := []Note{}

	q := "SELECT id, body FROM notes;"
	rows, err := s.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note Note

		if err := rows.Scan(&note.ID, &note.Body); err != nil {
			return nil, err
		}

		result = append(result, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SQLiteStore) Search(query string) ([]Note, error) {
	result := []Note{}
	q := "SELECT id, body FROM notes WHERE body LIKE ?"
	param := "%" + query + "%"

	rows, err := s.DB.Query(q, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Body); err != nil {
			return nil, err
		}

		result = append(result, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}


	return result, nil
}

func (s *SQLiteStore) GetByID(id int) (Note, error) {
	var result Note
	q := "SELECT id, body FROM notes;"
	rows, err := s.DB.Query(q)
	if err != nil {
		return Note{}, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&result.ID, &result.Body); err != nil {
			return Note{}, err
		}
	}

	if err := rows.Err(); err != nil {
		return Note{}, err
	}

	return result, nil
}

func (s *SQLiteStore) Create(body string) (Note, error) {
	note := Note{}
	json.NewDecoder(strings.NewReader(body)).Decode(&note)

	res, err := s.DB.Exec("INSERT INTO notes (body) VALUES (?);", note.Body)
	if err != nil {
		return Note{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Note{}, err
	}

	note.ID = int(id)

	fmt.Printf("createdNote: %#v\n", note)
	
	return note, nil 
}

func (s *SQLiteStore) Update(id int, body string) (Note, error) { 
	input := Note{}
	json.NewDecoder(strings.NewReader(body)).Decode(&input)

	res, err := s.DB.Exec("UPDATE notes SET body = ? WHERE id = ?;", input.Body, id)
	if err != nil {
		return Note{}, err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return Note{}, err
	}

	if rowsUpdated != 1 {
		return Note{}, errors.New("rowsUpdated != 1!")
	}


	output, err := s.GetByID(input.ID)
	if err != nil {
		return Note{}, err
	}
	
	return output, nil
}

func (s *SQLiteStore) Delete(id int) error { 
	res, err := s.DB.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsDeleted != 1 {
		return errors.New("rowsDeleted not 1!")
	}

	return nil
}