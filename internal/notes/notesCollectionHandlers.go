package notes

import (
	"net/http"
	"notes-api/internal/db"
	"notes-api/internal/helpers"
	"notes-api/internal/helpers/requestctx"
)

func (s *Service) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := requestctx.GetUserFromRequest(r)
	if !ok {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "Internal_server_error."})
		return
	}
	notes, err := s.Q.GetAllNotes(r.Context(), user)
	if err != nil {
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
	user, ok := requestctx.GetUserFromRequest(r)
	if !ok {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "Internal_server_error."})
		return
	}
	notes, err := s.Q.CreateNewNote(r.Context(), db.CreateNewNoteParams{Title: "", Content: "", OwnerID: user})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{
			Error: "internal_server_error",
		})
		return
	}
	helpers.WriteJson(w, http.StatusCreated, notes)
}
