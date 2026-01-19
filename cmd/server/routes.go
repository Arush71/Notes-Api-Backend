package main

import (
	"net/http"
	"notes-api/internal/notes"
)

func registerRouter(mux *http.ServeMux, notesService *notes.Service) {
	mux.HandleFunc("/app/notes", notesService.NotesCollectionHandler)
	mux.HandleFunc("/app/notes/{id}", notesService.NoteItemHandler)
	mux.HandleFunc("/app/auth/", notesService.AuthRouteHandler)
}
