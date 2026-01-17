package notes

import (
	"net/http"
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
