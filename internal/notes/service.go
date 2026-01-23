package notes

import (
	"database/sql"
	"notes-api/internal/db"
)

type Service struct {
	Q           *db.Queries
	TokenSecret string
	DB          *sql.DB
}
