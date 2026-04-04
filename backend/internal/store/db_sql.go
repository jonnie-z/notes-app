package store

import "database/sql"

type SQLiteStore struct {
	db *sql.DB
}

