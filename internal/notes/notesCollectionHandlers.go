package notes

import (
	"log"
	"net/http"
	"notes-api/internal/db"
)

func (s *Service) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := s.Q.GetAllNotes(r.Context())
	if err != nil {
		log.Println("GetAllNotes failed:", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if notes == nil {
		notes = []db.GetAllNotesRow{}
	}
	WriteJson(w, 200, notes)
}
func (s *Service) CreateNewNoteHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := s.Q.CreateNewNote(r.Context(), db.CreateNewNoteParams{Title: "", Content: ""})
	if err != nil {
		log.Println("CreateNewNote failed:", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJson(w, http.StatusCreated, notes)
}
