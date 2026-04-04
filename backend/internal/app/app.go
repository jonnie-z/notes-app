package app

import "github.com/jonnie-z/notes-app/internal/store"

type App struct {
	Store store.NoteRepository
}

func NewApp(storeType store.StoreType) *App {
	var appStore store.NoteRepository

	switch storeType {
	case store.StoreJSON:
		appStore = store.NewNoteStore()
	case store.StoreInMemory:
		appStore = store.NewInMemoryStore()
	}

	return &App{
		Store: appStore,
	}
}