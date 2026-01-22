package main

import (
	"net/http"
	"notes-api/internal/middleware"
	"notes-api/internal/notes"
)

func registerRouter(mux *http.ServeMux, notesService *notes.Service) {
	authMiddleware := middleware.MakeAuthMiddleWare(notesService.TokenSecret)
	mux.HandleFunc("/app/notes", authMiddleware(notesService.NotesCollectionHandler))
	mux.HandleFunc("/app/notes/{id}", authMiddleware(notesService.NoteItemHandler))
	mux.HandleFunc("/app/auth/", notesService.AuthRouteHandler)
}

// todo: add auth middleware
