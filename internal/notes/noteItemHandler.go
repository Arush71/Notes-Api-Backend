package notes

import (
	"database/sql"
	"errors"
	"net/http"
	"notes-api/internal/db"

	"github.com/google/uuid"
)

func (s *Service) GetNoteHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	note, err := s.Q.GetANote(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "couldn't find a note of that id.")
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	WriteJson(w, http.StatusOK, note)
}
func (s *Service) DeleteNoteHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	note, err := s.Q.DeleteNote(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "couldn't find a note of that id.")
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	WriteJson(w, http.StatusOK, note)
}
func (s *Service) UpdateNoteHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	type reqT struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}
	var req reqT
	if er := ReadJson(r, &req); er != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}
	if req.Content == nil || req.Title == nil {
		writeError(w, http.StatusBadRequest, "both content and id should be provided.")
		return
	}
	note, err := s.Q.UpdateNote(r.Context(), db.UpdateNoteParams{
		ID:      id,
		Title:   *req.Title,
		Content: *req.Content,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "couldn't find a note of that id.")
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	WriteJson(w, http.StatusOK, note)
}
