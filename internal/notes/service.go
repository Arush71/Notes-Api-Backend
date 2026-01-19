package notes

import "notes-api/internal/db"

type Service struct {
	Q           *db.Queries
	TokenSecret string
}
