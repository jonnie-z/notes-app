package app

import (
	"github.com/jonnie-z/notes-app/internal/store"
)

type NoteRepository interface {
	List(query string, page int, pageSize int) ([]store.Note, int, error)
	Create(body string) (store.Note, error)
	Update(id int, body string) (store.Note, error)
	Delete(id int) error
}

type App struct {
	Store NoteRepository

	Port      string
	NotesFile string
}
