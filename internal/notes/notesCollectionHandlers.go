package notes

import (
	"database/sql"
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
	tx, err := s.DB.BeginTx(r.Context(), &sql.TxOptions{ReadOnly: false, Isolation: sql.LevelDefault})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_server_error"})
		return
	}
	defer tx.Rollback()
	qtx := s.Q.WithTx(tx)
	notes, err := qtx.CreateNewNote(r.Context(), db.CreateNewNoteParams{Title: "", Content: "", OwnerID: user})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{
			Error: "internal_server_error",
		})
		return
	}
	if _, err = qtx.CreateAVersion(r.Context(), db.CreateAVersionParams{NoteID: notes.ID, Title: "", Content: "", VersionNumber: 1}); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{
			Error: "internal_server_error",
		})
		return
	}
	if err = tx.Commit(); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{
			Error: "internal_server_error",
		})
		return
	}

	helpers.WriteJson(w, http.StatusCreated, notes)
}
