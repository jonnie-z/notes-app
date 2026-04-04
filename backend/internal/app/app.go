package app

import (
	"log"

	"github.com/jonnie-z/notes-app/internal/store"
)

type App struct {
	Store store.NoteRepository
}

func NewApp(storeType store.StoreType) *App {
	dsn := "./notes.db"
	var appStore store.NoteRepository

	switch storeType {
	case store.StoreJSON:
		appStore = store.NewNoteStore()
	case store.StoreInMemory:
		appStore = store.NewInMemoryStore()
	case store.StoreSQL:
		s, err := store.NewSQLiteStore(dsn)
		if err != nil {
			log.Fatalf("Err creating sql store! %s", err.Error())
		}

		appStore = s
	}

	return &App{
		Store: appStore,
	}
}