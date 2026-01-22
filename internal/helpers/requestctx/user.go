package requestctx

import (
	"net/http"

	"github.com/google/uuid"
)

type UserKeyType struct{}

var UserKey = UserKeyType{}

// getter

func GetUserFromRequest(r *http.Request) (uuid.UUID, bool) {
	id, ok := r.Context().Value(UserKey).(uuid.UUID)
	return id, ok
}
