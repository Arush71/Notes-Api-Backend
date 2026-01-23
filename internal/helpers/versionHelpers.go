package helpers

import (
	"context"
	"errors"
	"notes-api/internal/db"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("UNAUTHORIZED")

func ValidateNote(noteId string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(noteId))
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil

}

func ValidateVersion(versionNumber string) (int, error) {
	num, err := strconv.Atoi(versionNumber)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func UserOwned(noteId, userId uuid.UUID, Q *db.Queries, context context.Context) error {
	if Q == nil {
		return errors.New("no query provided.")
	}
	checkUser, err := Q.CheckIfUserOwned(context, db.CheckIfUserOwnedParams{ID: noteId, OwnerID: userId})
	if err != nil {
		return err
	}
	if !checkUser {
		return ErrUnauthorized
	}
	return nil
}
