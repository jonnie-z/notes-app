package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
)

type NoteStore struct {
	notes         []Note
	nextID        int
	notesFile     string
	tempNotesFile string

	mu sync.Mutex
}

func NewNoteStore() *NoteStore {
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

func (n *NoteStore) GetAll() ([]Note, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	dst := make([]Note, len(n.notes))
	copy(dst, n.notes)

	return dst, nil
}

func (n *NoteStore) Create(text string) (Note, error) {
	note := Note{}
	json.NewDecoder(strings.NewReader(text)).Decode(&note)

	n.appendNote(&note)
	n.saveNotes()

	return note, nil
}

func (n *NoteStore) appendNote(note *Note) {
	n.mu.Lock()
	defer n.mu.Unlock()

	note.ID = n.nextID
	n.nextID = n.nextID + 1

	n.notes = append(n.notes, *note)
}

func (n *NoteStore) Delete(id int) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	noteIdx := getNoteIdx(n.notes, id)

	if noteIdx != -1 {
		n.notes = removeNoteAtIndex(n.notes, noteIdx)
		n.saveNotes()

		return nil
	} else {
		return fmt.Errorf("Id '%d' not found!", id)
	}
}

func (n *NoteStore) Update(id int, text string) (Note, error) {
	idx := getNoteIdx(n.notes, id)

	if idx != -1 {
		newNote := Note{}
		json.NewDecoder(strings.NewReader(text)).Decode(&newNote)

		n.mu.Lock()
		defer n.mu.Unlock()
		n.notes[idx].Text = newNote.Text
		n.saveNotes()

		return newNote, nil
	} else {
		return Note{}, fmt.Errorf("Id '%d' not found!", id)
	}
}