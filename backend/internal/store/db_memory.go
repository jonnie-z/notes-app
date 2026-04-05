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

func (i *InMemoryStore) GetAll() ([]Note, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	notesCopy := make([]Note, len(i.notes))
	copy(notesCopy, i.notes)

	return notesCopy, nil
}

func (i *InMemoryStore) Search(query string) ([]Note, error) {
	if query == "" {
		return i.GetAll()
	}

	notes, err := i.GetAll()
	if err != nil {
		return nil, err
	}

	var result []Note
	for _, note := range notes {
		if strings.Contains(strings.ToLower(note.Body), strings.ToLower(query)) {
			result = append(result, note)
		}
	}

	return result, nil
}

func (i *InMemoryStore) Create(body string) (Note, error) {
	note := Note{}
	json.NewDecoder(strings.NewReader(body)).Decode(&note)

	i.appendToStore(&note)

	return note, nil
}

func (i *InMemoryStore) Update(id int, body string) (Note, error) {
	idx := getNoteIdx(i.notes, id)

	if idx != -1 {
		newNote := Note{}
		json.NewDecoder(strings.NewReader(body)).Decode(&newNote)

		i.mu.Lock()
		defer i.mu.Unlock()
		i.notes[idx].Body = newNote.Body

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