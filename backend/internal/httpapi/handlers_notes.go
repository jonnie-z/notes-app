package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jonnie-z/notes-app/internal/store"
)

const (
	PAGE_MIN      int = 1
	PAGE_SIZE_MIN int = 10
)

type NotesPage struct {
	Notes    []store.Note `json:"notes"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
	Total    int          `json:"total"`
}

func (a *API) GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	result := NotesPage{
		Notes:    []store.Note{},
		Page:     0,
		PageSize: 0,
		Total:    0,
	}
	var err error

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		result.Page = PAGE_MIN
	} else {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil {
			result.Page = PAGE_MIN
		} else {
			result.Page = pageInt
		}
	}

	pageSizeStr := r.URL.Query().Get("pageSize")
	if pageSizeStr == "" {
		result.PageSize = PAGE_SIZE_MIN
	} else {
		pageSizeInt, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			result.PageSize = PAGE_SIZE_MIN
		} else {
			result.PageSize = pageSizeInt
		}
	}

	if result.Page < 1 {
		result.Page = PAGE_MIN
	}
	if result.PageSize < 1 {
		result.PageSize = PAGE_SIZE_MIN
	}

	q := r.URL.Query().Get("query")
	result.Notes, result.Total, err = a.App.Store.List(q, result.Page, result.PageSize)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("notes: %#v\n", result.Notes)
	if err := json.NewEncoder(w).Encode(result); err != nil {
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

		if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
		if err := json.NewEncoder(w).Encode(newNote); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
