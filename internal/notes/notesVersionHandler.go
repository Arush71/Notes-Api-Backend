package notes

import (
	"database/sql"
	"errors"
	"net/http"
	"notes-api/internal/db"
	"notes-api/internal/helpers"
	"notes-api/internal/helpers/requestctx"
	"time"
)

func (s *Service) GetAllVersions(w http.ResponseWriter, r *http.Request) {
	userId, ok := requestctx.GetUserFromRequest(r)
	if !ok {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "Internal_server_error."})
		return
	}
	noteId, err := helpers.ValidateNote(r.PathValue("id"))
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "Bad_request", Message: "Can't extract id"})
		return
	}
	if err = helpers.UserOwned(noteId, userId, s.Q, r.Context()); err != nil {
		if errors.Is(err, helpers.ErrUnauthorized) {
			helpers.WriteError(w, http.StatusForbidden, helpers.ErrorResponse{Error: "FORBIDDEN"})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "something went wrong please try again."})
		return
	}
	versions, err := s.Q.GetAllVersions(r.Context(), noteId)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "something went wrong please try again."})
		return
	}
	if versions == nil {
		versions = []db.GetAllVersionsRow{}
	}
	helpers.WriteJson(w, http.StatusOK, versions)
}
func (s *Service) GetAVersion(w http.ResponseWriter, r *http.Request) {
	userId, ok := requestctx.GetUserFromRequest(r)
	if !ok {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "Internal_server_error."})
		return
	}
	noteId, err := helpers.ValidateNote(r.PathValue("id"))
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "Bad_request", Message: "Can't extract id"})
		return
	}
	versionNum, err := helpers.ValidateVersion(r.PathValue("version_number"))
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "Bad_request", Message: "Can't extract version number"})
		return
	}
	if err = helpers.UserOwned(noteId, userId, s.Q, r.Context()); err != nil {
		if errors.Is(err, helpers.ErrUnauthorized) {
			helpers.WriteError(w, http.StatusForbidden, helpers.ErrorResponse{Error: "FORBIDDEN"})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "something went wrong please try again."})
		return
	}
	version, err := s.Q.GetNoteVersion(r.Context(), db.GetNoteVersionParams{NoteID: noteId, VersionNumber: int32(versionNum)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{Error: "resource_not_found", Message: "Version Not found"})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "try again latter."})
		return
	}
	helpers.WriteJson(w, http.StatusOK, version)
}

func (s *Service) RollBackVersion(w http.ResponseWriter, r *http.Request) {
	userId, ok := requestctx.GetUserFromRequest(r)
	if !ok {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "Internal_server_error."})
		return
	}
	noteId, err := helpers.ValidateNote(r.PathValue("id"))
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "Bad_request", Message: "Can't extract id"})
		return
	}
	versionNum, err := helpers.ValidateVersion(r.PathValue("version_number"))
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "Bad_request", Message: "Can't extract version number"})
		return
	}
	tx, err := s.DB.BeginTx(r.Context(), nil)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_server_error"})
		return
	}
	defer tx.Rollback()
	qtx := s.Q.WithTx(tx)
	if err = helpers.UserOwned(noteId, userId, qtx, r.Context()); err != nil {
		if errors.Is(err, helpers.ErrUnauthorized) {
			helpers.WriteError(w, http.StatusForbidden, helpers.ErrorResponse{Error: "FORBIDDEN"})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "something went wrong please try again."})
		return
	}

	noteVersion, err := qtx.GetNoteVersion(r.Context(), db.GetNoteVersionParams{NoteID: noteId, VersionNumber: int32(versionNum)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{Error: "resource_not_found", Message: "Version Not found"})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "try again latter."})
		return
	}
	currentTopVersion, err := qtx.GetCurrentHighestVersion(r.Context(), noteId)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "try again latter."})
		return
	}
	newVersion := currentTopVersion + 1
	if _, err = qtx.CreateAVersion(r.Context(), db.CreateAVersionParams{
		NoteID:        noteVersion.NoteID,
		VersionNumber: newVersion,
		Title:         noteVersion.Title,
		Content:       noteVersion.Content,
	}); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "try again latter."})
		return
	}

	// create actually note
	note, err := qtx.UpdateNote(r.Context(), db.UpdateNoteParams{ID: noteVersion.NoteID, OwnerID: userId, Title: noteVersion.Title, Content: noteVersion.Content})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{
				Error: "resource_not_found",
			})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "try again latter."})
		return
	}

	if err = tx.Commit(); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "internal_error", Message: "try again latter."})
		return
	}

	// return struct

	type RData struct {
		Title          string    `json:"title"`
		Content        string    `json:"content"`
		UpdatedAt      time.Time `json:"updated_at"`
		CurrentVersion int32     `json:"current_version"`
	}
	helpers.WriteJson(w, http.StatusOK, RData{
		Title:          note.Title,
		Content:        note.Content,
		UpdatedAt:      note.UpdatedAt,
		CurrentVersion: newVersion,
	})
}
