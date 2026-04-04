package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (a *API) GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	notes, _ := a.App.Store.GetAll()
	fmt.Printf("notes: %#v\n", notes)
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	// randomNumber := rand.IntN(100)
	// if randomNumber > 65 {
	// 	http.Error(w, "OPE IT FAILED WHOOPS", http.StatusInternalServerError)
	// 	return
	// }

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	if !json.Valid(body) {
		http.Error(w, "Invalid JSON!", http.StatusBadRequest)
		return
	}

	note, _ := a.App.Store.Create(string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		// notes = slices.DeleteFunc(notes, func(n Note) bool {
		// 	return n.ID == id
		// })

		err = a.App.Store.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *API) PutNotesHandler(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		newNote, err := a.App.Store.Update(id, string(bodyBytes))
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