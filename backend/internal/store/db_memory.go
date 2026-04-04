package store

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type InMemoryStore struct {
	notes  []Note
	nextID int

	mu sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	inMemoryStore := &InMemoryStore{
		notes:  []Note{},
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