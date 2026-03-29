package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

const PORT = ":8080"

var notes = []Note{}
var nextID = 0

// type Response struct {
// 	Notes []Note `json:"notes"`
// }

type Note struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func setNextID() {
	highest := math.MinInt64

	for _, note := range notes {
		if note.ID > highest {
			highest = note.ID
		}
	}

	if highest < 0 { highest = -1 }

	nextID = highest + 1
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("notes: %#v\n", notes)
	// resp := Response{
	// 	Notes: []Note{},
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func postNotes(w http.ResponseWriter, r *http.Request) {
	note := Note{}
	json.NewDecoder(r.Body).Decode(&note)

	note.ID = nextID
	nextID = nextID + 1

	notes = append(notes, note)

	w.WriteHeader(http.StatusCreated)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		// notes = slices.DeleteFunc(notes, func(n Note) bool {
		// 	return n.ID == id
		// })

		noteIdx := getNoteIdx(id)
		
		if noteIdx != -1 {
			notes = append(notes[:noteIdx], notes[noteIdx+1:]...)
		} else {
			http.Error(w, "Id not found!", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getNoteIdx(id int) int {
	for i, note := range notes {
		if note.ID == id {
			return i
		}
	}

	return -1
}

func main() {
	setNextID()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/notes", getNotes)
	mux.HandleFunc("POST /api/notes", postNotes)
	mux.HandleFunc("DELETE /api/notes/{id}", deleteNote)

	fmt.Println("Starting Server on", PORT)
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}
}
