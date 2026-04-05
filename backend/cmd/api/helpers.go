package main

import (
	"log"
	"os"

	"github.com/jonnie-z/notes-app/internal/app"
	"github.com/jonnie-z/notes-app/internal/store"
)

func env(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	return v
}

func newApp(storeType store.StoreType) *app.App {
	var appStore app.NoteRepository
	port := env("PORT", ":8080")
	dsn := env("DSN", "./notes.db")

	app := &app.App{}
	app.Port = port
	app.NotesFile = env("NOTES_FILE", "db.json")

	switch storeType {
	case store.StoreJSON:
		appStore = store.NewNoteStore(app.NotesFile)
	case store.StoreInMemory:
		appStore = store.NewInMemoryStore()
	case store.StoreSQL:
		s, err := store.NewSQLiteStore(dsn)
		if err != nil {
			log.Fatalf("Err creating sql store! %s", err.Error())
		}

		appStore = s
	}

	app.Store = appStore

	return app
}