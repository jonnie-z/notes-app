package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (a *App) getNotesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	notes, _ := a.store.GetAll()
	fmt.Printf("notes: %#v\n", notes)
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (n *NoteStore) GetAll() ([]Note, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	dst := make([]Note, len(n.notes))
	copy(dst, n.notes)

	return dst, nil
}

func (a *App) postNoteHandler(w http.ResponseWriter, r *http.Request) {
	// randomNumber := rand.IntN(100)
	// if randomNumber > 65 {
	// 	http.Error(w, "OPE IT FAILED WHOOPS", http.StatusInternalServerError)
	// 	return
	// }

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	note, _ := a.store.Create(string(bodyBytes))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func (a *App) deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		// notes = slices.DeleteFunc(notes, func(n Note) bool {
		// 	return n.ID == id
		// })

		err = a.store.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func (a *App) putNotesHandler(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		newNote, err := a.store.Update(id, string(bodyBytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newNote); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
