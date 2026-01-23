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
	// versioning
	mux.HandleFunc("GET /app/notes/{id}/versions", authMiddleware(notesService.GetAllVersions))
	mux.HandleFunc("GET /app/notes/{id}/versions/{version_number}", authMiddleware(notesService.GetAVersion))
	mux.HandleFunc("POST /app/notes/{id}/versions/{version_number}/rollback", authMiddleware(notesService.RollBackVersion))
}
