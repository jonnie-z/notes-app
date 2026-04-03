package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
)

const PORT = ":8080"

// const DATA_FILE = "./db.json"
// const TEMP_DATA_FILE = DATA_FILE + ".tmp"

type StoreType int

const (
	StoreJSON StoreType = iota
	StoreInMemory
)

type App struct {
	store NoteRepository
}

func newApp(storeType StoreType) App {
	var store NoteRepository

	switch storeType {
	case StoreJSON:
		store = newNoteStore()
	case StoreInMemory:
		store = newInMemoryStore()
	}

	return App{
		store: store,
	}
}

type NoteRepository interface {
	GetAll() ([]Note, error)
	Create(text string) (Note, error)
	Update(id int, text string) (Note, error)
	Delete(id int) error
}

type NoteStore struct {
	notes         []Note
	nextID        int
	notesFile     string
	tempNotesFile string

	mu sync.Mutex
}

type InMemoryStore struct {
	notes []Note
	nextID int
	
	mu sync.Mutex
}

func newInMemoryStore() *InMemoryStore {
	inMemoryStore := &InMemoryStore{
		notes: []Note{},
		nextID: 0,
	}

	return inMemoryStore
}

func (i *InMemoryStore) Create(text string) (Note, error) {
	note := Note{}
	json.NewDecoder(strings.NewReader(text)).Decode(&note)

	i.appendToStore(&note)

	return note, nil
}

func (i *InMemoryStore) GetAll() ([]Note, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	notesCopy := make([]Note, len(i.notes))
	copy(notesCopy, i.notes)

	return notesCopy, nil
}

func (i *InMemoryStore) Update(id int, text string) (Note, error) {
	idx := getNoteIdx(i.notes, id)

	if idx != -1 {
		newNote := Note{}
		json.NewDecoder(strings.NewReader(text)).Decode(&newNote)

		i.mu.Lock()
		defer i.mu.Unlock()
		i.notes[idx].Text = newNote.Text

		return newNote, nil
	} else {
		return Note{}, fmt.Errorf("Id '%d' not found!", id)
	}
}


func (i *InMemoryStore) Delete(id int) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	idx := getNoteIdx(i.notes, id)

	if idx == -1 {
		return fmt.Errorf("Note with ID '%d' not found!", id)
	}

	i.notes = removeNoteAtIndex(i.notes, idx)

	return nil
}

func (i *InMemoryStore) appendToStore(note *Note) {
	i.mu.Lock()
	defer i.mu.Unlock()

	note.ID = i.nextID
	i.nextID = i.nextID + 1

	i.notes = append(i.notes, *note)
}

// func (i *InMemoryStore) setNextID() {
// 	i.mu.Lock()
// 	defer i.mu.Unlock()
// 	highest := math.MinInt64

// 	for _, note := range i.notes {
// 		if note.ID > highest {
// 			highest = note.ID
// 		}
// 	}

// 	if highest < 0 {
// 		highest = -1
// 	}

// 	i.nextID = highest + 1
// }

func newNoteStore() *NoteStore {
	notesFile := "./db.json"

	noteStore := &NoteStore{
		notes: []Note{},
		nextID: math.MinInt64,
		notesFile: notesFile,
		tempNotesFile: notesFile + ".tmp",
	}

	noteStore.loadNotes()
	noteStore.setNextID()

	return noteStore
}

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func (n *NoteStore) setNextID() {
	n.mu.Lock()
	defer n.mu.Unlock()
	highest := math.MinInt64

	for _, note := range n.notes {
		if note.ID > highest {
			highest = note.ID
		}
	}

	if highest < 0 {
		highest = -1
	}

	n.nextID = highest + 1
}

func (n *NoteStore) loadNotes() {
	file, err := os.Open(n.notesFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Data file not found. Initializing empty slice.")
			n.notes = []Note{}
		} else {
			fmt.Printf("Error opening file: %v\n", err)
		}

		return
	}

	defer file.Close()
	fmt.Println("Data file opened successfully. Loading into slice.")

	n.notes = []Note{}
	if err := json.NewDecoder(file).Decode(&n.notes); err != nil {
		log.Fatal("Error decoding data file!")
	}

	fmt.Println("Data file successully loaded into Notes.")
}

func (n *NoteStore) saveNotes() {
	file, err := os.Create(n.tempNotesFile)
	if err != nil {
		log.Fatalf("Error creating temp data file: %v", err)
	}

	if err := json.NewEncoder(file).Encode(n.notes); err != nil {
		fmt.Printf("Error encoding temp data file: %v; removing temp file\n", err)
		file.Close()
		os.Remove(n.tempNotesFile)
		return
	}

	file.Close()

	if err = os.Rename(n.tempNotesFile, n.notesFile); err != nil {
		fmt.Printf("Error replacing data file: %v\n", err)
	}
	fmt.Println("Sucessfully saved data file")
}

func main() {
	args := os.Args[1:]
	var storeType StoreType

	if len(args) == 0 {
		storeType = StoreJSON
	} else {
		switch args[0] {
		case "json":
			storeType = StoreJSON
		case "mem":
			storeType = StoreInMemory
		}
	}

	app := newApp(storeType)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/notes", app.getNotesHandler)
	mux.HandleFunc("POST /api/notes", app.postNoteHandler)
	mux.HandleFunc("DELETE /api/notes/{id}", app.deleteNoteHandler)
	mux.HandleFunc("PUT /api/notes/{id}", app.putNotesHandler)

	fmt.Println("Starting Server on", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}
}
