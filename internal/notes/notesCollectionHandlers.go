package notes

import (
	"log"
	"net/http"
	"notes-api/internal/db"
	"notes-api/internal/helpers"
)

func (s *Service) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := s.Q.GetAllNotes(r.Context())
	if err != nil {
		log.Println("GetAllNotes failed:", err)
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{
			Error: "internal server error",
		})
		return
	}
	if notes == nil {
		notes = []db.GetAllNotesRow{}
	}
	helpers.WriteJson(w, 200, notes)
}
func (s *Service) CreateNewNoteHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := s.Q.CreateNewNote(r.Context(), db.CreateNewNoteParams{Title: "", Content: ""})
	if err != nil {
		log.Println("CreateNewNote failed:", err)
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{
			Error: "internal server error",
		})
		return
	}
	helpers.WriteJson(w, http.StatusCreated, notes)
}
