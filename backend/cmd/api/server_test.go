package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jonnie-z/notes-app/internal/app"
	"github.com/jonnie-z/notes-app/internal/httpapi"
	"github.com/jonnie-z/notes-app/internal/store"
)

type inputPostNote struct{
	body string
}

func TestPostNotes(t *testing.T) {
	app := app.NewApp(store.StoreInMemory)
	api := httpapi.API{App: app}

	tests := []struct {
		input inputPostNote
		expected store.Note
	}{
		{inputPostNote{"{\"text\":\"hello\"}"}, store.Note{ID:0,Text:"hello"}},
	}

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/notes", strings.NewReader(test.input.body))
		w := httptest.NewRecorder()

		api.PostNoteHandler(w, r)

		var actual store.Note
		json.NewDecoder(w.Body).Decode(&actual)

		if actual != test.expected {
			t.Errorf("Expected '%v' but got '%v'", test.expected, actual)
		}

		if actual.ID != test.expected.ID {
			t.Errorf("Expected note ID to be '%d' but got '%v'", test.expected.ID, actual.ID)
		}

		if actual.Text != test.expected.Text {
			t.Errorf("Expected note Text to be '%s' but got '%s'", test.expected.Text, actual.Text)
		}

		if w.Code != http.StatusCreated {
			t.Fatalf("expected '%d', got '%d'", http.StatusCreated, w.Code)
		}

		notes, _ := app.Store.GetAll()
		if len(notes) != 1 {
			t.Fatalf("expected '1' note, got '%d'", len(notes))
		}
	}
}

func TestPostNotesError(t *testing.T) {
	app := app.NewApp(store.StoreInMemory)
	api := httpapi.API{App: app}

	r := httptest.NewRequest(http.MethodPost, "/api/notes", strings.NewReader("lol"))
	w := httptest.NewRecorder()

	api.PostNoteHandler(w, r)

	var actual store.Note
	json.NewDecoder(w.Body).Decode(&actual)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected '%d', got '%d'", http.StatusBadRequest, w.Code)
	}

	notes, _ := app.Store.GetAll()
	if len(notes) != 0 {
		t.Fatalf("expected '0' note, got '%d'", len(notes))
	}
}

func TestDeleteNote(t *testing.T) {
	app := app.NewApp(store.StoreInMemory)
	api := httpapi.API{App: app}
	app.Store.Create("hello")

	r := httptest.NewRequest(http.MethodDelete, "/api/notes/0", nil)
	r.SetPathValue("id", "0")
	w := httptest.NewRecorder()

	api.DeleteNoteHandler(w, r)

	if w.Code != http.StatusOK {
		t.Logf("%v", w.Body.String())
		t.Fatalf("expected '%d', got '%d'", http.StatusOK, w.Code)
	}

	notes, _ := app.Store.GetAll()
	if len(notes) != 0 {
		t.Fatalf("expected '0' note, got '%d'", len(notes))
	}
}