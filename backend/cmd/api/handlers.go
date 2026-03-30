package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (a *App) getNotesHandler(w http.ResponseWriter, r *http.Request) {
	// resp := Response{
	// 	Notes: []Note{},
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	notes := getNotesCopy(a)
	fmt.Printf("notes: %#v\n", notes)
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getNotesCopy(a *App) []Note {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	dst := make([]Note, len(a.notes))
	copy(dst, a.notes)

	return dst
}

func (a *App) postNoteHandler(w http.ResponseWriter, r *http.Request) {
	note := Note{}
	json.NewDecoder(r.Body).Decode(&note)

	appendAndSave(a, note)

	w.WriteHeader(http.StatusCreated)
}

func appendAndSave(a *App, note Note) {
	a.mu.Lock()
	defer a.mu.Unlock()

	note.ID = a.nextID
	a.nextID = a.nextID + 1

	a.notes = append(a.notes, note)
	a.saveNotes()
}

func (a *App) deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		// notes = slices.DeleteFunc(notes, func(n Note) bool {
		// 	return n.ID == id
		// })

		noteIdx := getNoteIdx(a, id)

		if noteIdx != -1 {
			removeNoteAtIndex(a, noteIdx)
		} else {
			http.Error(w, "Id not found!", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) putNotesHandler(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		noteIdx := getNoteIdx(a, id)

		newNote := Note{}
		json.NewDecoder(r.Body).Decode(&newNote)

		if noteIdx != -1 {
			updateNoteAtIndex(a, noteIdx, newNote)
		} else {
			http.Error(w, "Id not found!", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateNoteAtIndex(a *App, idx int, newNote Note) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.notes[idx].Body = newNote.Body
	a.saveNotes()
}

func removeNoteAtIndex(a *App, idx int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.notes = append(a.notes[:idx], a.notes[idx+1:]...)
	a.saveNotes()
}

func getNoteIdx(a *App, id int) int {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i, note := range a.notes {
		if note.ID == id {
			return i
		}
	}

	return -1
}
