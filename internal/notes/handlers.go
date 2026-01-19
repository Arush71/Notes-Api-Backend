package notes

import (
	"net/http"
	"notes-api/internal/helpers"
	"strings"

	"github.com/google/uuid"
)

func (s *Service) NotesCollectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetAllNotesHandler(w, r)
	case http.MethodPost:
		s.CreateNewNoteHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}
func (s *Service) NoteItemHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	uuid_id, err := uuid.Parse(id)
	if id == "" || err != nil {
		http.Error(w, "couldn't extract id", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		s.GetNoteHandler(w, r, uuid_id)
	case http.MethodDelete:
		s.DeleteNoteHandler(w, r, uuid_id)
	case http.MethodPut:
		s.UpdateNoteHandler(w, r, uuid_id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Service) AuthRouteHandler(w http.ResponseWriter, r *http.Request) {
	reqPath := strings.TrimPrefix(r.URL.Path, "/app/auth/")
	switch {
	case reqPath == "login" && http.MethodPost == r.Method:
		s.LoginHandler(w, r)
	case reqPath == "register" && http.MethodPost == r.Method:
		s.RegisterHandler(w, r)
	case reqPath != "register" && reqPath != "login":
		helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{
			Error: "route not found",
		})
	default:
		helpers.WriteError(w, http.StatusMethodNotAllowed, helpers.ErrorResponse{
			Error: "method not allowed",
		})
	}
}
