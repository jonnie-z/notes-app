package app

import (
	"github.com/jonnie-z/notes-app/internal/store"
)

type NoteRepository interface {
	GetAll() ([]store.Note, error)
	Search(query string) ([]store.Note, error)
	Create(body string) (store.Note, error)
	Update(id int, body string) (store.Note, error)
	Delete(id int) error
}

type App struct {
	Store NoteRepository

	Port      string
	NotesFile string
}
