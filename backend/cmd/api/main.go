package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
)

const PORT = ":8080"

// const DATA_FILE = "./db.json"
// const TEMP_DATA_FILE = DATA_FILE + ".tmp"

type App struct {
	store *NoteStore
}

func newApp() App {
	return App{
		store: newNoteStore(),
	}
}

type NoteStore struct {
	notes         []Note
	nextID        int
	notesFile     string
	tempNotesFile string

	mu sync.Mutex
}

func newNoteStore() *NoteStore {
	notesFile := "./db.json"
	return &NoteStore{
		notes: []Note{},
		nextID: math.MinInt64,
		notesFile: notesFile,
		tempNotesFile: notesFile + ".tmp",
	}
}

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func (n *NoteStore) setNextID() {
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
	app := newApp()
	app.store.loadNotes()
	app.store.setNextID()

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
