package httpapi

import "net/http"

func (api *API) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/notes", api.GetNotesHandler)
	mux.HandleFunc("POST /api/notes", api.PostNoteHandler)
	mux.HandleFunc("DELETE /api/notes/{id}", api.DeleteNoteHandler)
	mux.HandleFunc("PUT /api/notes/{id}", api.PutNotesHandler)

	return mux
}
