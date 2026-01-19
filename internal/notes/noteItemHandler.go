package notes

import (
	"database/sql"
	"errors"
	"net/http"
	"notes-api/internal/db"
	"notes-api/internal/helpers"

	"github.com/google/uuid"
)

func (s *Service) GetNoteHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	note, err := s.Q.GetANote(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{
				Error: "couldn't find a note of that id.",
			})
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	helpers.WriteJson(w, http.StatusOK, note)
}
func (s *Service) DeleteNoteHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	note, err := s.Q.DeleteNote(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{
				Error: "couldn't find a note of that id.",
			})
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	helpers.WriteJson(w, http.StatusOK, note)
}
func (s *Service) UpdateNoteHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
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
			Error: "both content and id should be provided.",
		})
		return
	}
	note, err := s.Q.UpdateNote(r.Context(), db.UpdateNoteParams{
		ID:      id,
		Title:   *req.Title,
		Content: *req.Content,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{
				Error: "couldn't find a note of that id.",
			})
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	helpers.WriteJson(w, http.StatusOK, note)
}
