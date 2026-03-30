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
	notes         []Note
	nextID        int
	notesFile     string
	tempNotesFile string

	mu sync.Mutex
}

func newApp() App {
	notesFile := "./db.json"
	return App{
		notes:         []Note{},
		nextID:        0,
		notesFile:     notesFile,
		tempNotesFile: notesFile + ".tmp",
	}
}

type Note struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (a *App) setNextID() {
	highest := math.MinInt64

	for _, note := range a.notes {
		if note.ID > highest {
			highest = note.ID
		}
	}

	if highest < 0 {
		highest = -1
	}

	a.nextID = highest + 1
}

func (a *App) loadNotes() {
	file, err := os.Open(a.notesFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Data file not found. Initializing empty slice.")
			a.notes = []Note{}
		} else {
			fmt.Printf("Error opening file: %v\n", err)
		}

		return
	}

	defer file.Close()
	fmt.Println("Data file opened successfully. Loading into slice.")

	a.notes = []Note{}
	if err := json.NewDecoder(file).Decode(&a.notes); err != nil {
		log.Fatal("Error decoding data file!")
	}

	fmt.Println("Data file successully loaded into Notes.")
}

func (a *App) saveNotes() {
	file, err := os.Create(a.tempNotesFile)
	if err != nil {
		log.Fatalf("Error creating temp data file: %v", err)
	}

	if err := json.NewEncoder(file).Encode(a.notes); err != nil {
		fmt.Printf("Error encoding temp data file: %v; removing temp file\n", err)
		file.Close()
		os.Remove(a.tempNotesFile)
		return
	}

	file.Close()

	if err = os.Rename(a.tempNotesFile, a.notesFile); err != nil {
		fmt.Printf("Error replacing data file: %v\n", err)
	}
	fmt.Println("Sucessfully saved data file")
}

func main() {
	app := newApp()
	app.loadNotes()
	app.setNextID()

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
