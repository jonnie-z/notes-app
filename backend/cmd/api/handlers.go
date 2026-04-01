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

	notes := a.store.GetAll()
	fmt.Printf("notes: %#v\n", notes)
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (n *NoteStore) GetAll() []Note {
	n.mu.Lock()
	defer n.mu.Unlock()

	dst := make([]Note, len(n.notes))
	copy(dst, n.notes)

	return dst
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

	note := a.store.Create(string(bodyBytes))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (n *NoteStore) Create(body string) Note {
	note := Note{}
	json.NewDecoder(strings.NewReader(body)).Decode(&note)

	n.appendAndSave(&note)

	return note
}

func (n *NoteStore) appendAndSave(note *Note) {
	n.mu.Lock()
	defer n.mu.Unlock()

	note.ID = n.nextID
	n.nextID = n.nextID + 1

	n.notes = append(n.notes, *note)
	n.saveNotes()
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
	noteIdx := n.getNoteIdx(id)

	if noteIdx != -1 {
		n.removeNoteAtIndex(noteIdx)
		return nil
	} else {
		return fmt.Errorf("Id '%d' not found!", id)
	}
}

func (n *NoteStore) removeNoteAtIndex(idx int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.notes = append(n.notes[:idx], n.notes[idx+1:]...)
	n.saveNotes()
}

func (n *NoteStore) getNoteIdx(id int) int {
	n.mu.Lock()
	defer n.mu.Unlock()

	for i, note := range n.notes {
		if note.ID == id {
			return i
		}
	}

	return -1
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

func (n *NoteStore) Update(id int, body string) (Note, error) {
	noteIdx := n.getNoteIdx(id)

	if noteIdx != -1 {
		newNote := Note{}
		json.NewDecoder(strings.NewReader(body)).Decode(&newNote)

		n.updateNoteAtIndex(noteIdx, newNote)

		return newNote, nil
	} else {
		return Note{}, fmt.Errorf("Id '%d' not found!", id)
	}
}

func (n *NoteStore) updateNoteAtIndex(idx int, newNote Note) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.notes[idx].Text = newNote.Text
	n.saveNotes()
}
