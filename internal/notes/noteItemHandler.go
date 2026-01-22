package notes

import (
	"database/sql"
	"errors"
	"net/http"
	"notes-api/internal/db"
	"notes-api/internal/helpers"
	"notes-api/internal/helpers/requestctx"

	"github.com/google/uuid"
)

func (s *Service) GetNoteHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	user, ok := requestctx.GetUserFromRequest(r)
	if !ok {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "Internal_server_error."})
		return
	}
	note, err := s.Q.GetANote(r.Context(), db.GetANoteParams{
		ID:      id,
		OwnerID: user,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{
				Error:   "not_found",
				Message: "resource not found.",
			})
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	helpers.WriteJson(w, http.StatusOK, note)
}
func (s *Service) DeleteNoteHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	user, ok := requestctx.GetUserFromRequest(r)
	if !ok {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "Internal_server_error."})
		return
	}
	_, err := s.Q.DeleteNote(r.Context(), db.DeleteNoteParams{
		ID:      id,
		OwnerID: user,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{
				Error:   "not_found",
				Message: "resource not found.",
			})
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (s *Service) UpdateNoteHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	user, ok := requestctx.GetUserFromRequest(r)
	if !ok {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "Internal_server_error."})
		return
	}
	type reqT struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}
	var req reqT
	if er := helpers.ReadJson(r, &req); er != nil {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Error: "bad request",
		})
		return
	}
	if req.Content == nil || req.Title == nil {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Error: "both content and title should be provided.",
		})
		return
	}
	note, err := s.Q.UpdateNote(r.Context(), db.UpdateNoteParams{
		ID:      id,
		Title:   *req.Title,
		Content: *req.Content,
		OwnerID: user,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{
				Error: "resource_not_found",
			})
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	helpers.WriteJson(w, http.StatusOK, note)
}
