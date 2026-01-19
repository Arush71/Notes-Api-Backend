package notes

import (
	"database/sql"
	"errors"
	"net/http"
	"notes-api/internal/auth"
	"notes-api/internal/db"
	"notes-api/internal/helpers"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := helpers.ReadJson(r, &req); err != nil || req.Email == nil || req.Password == nil {
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
				Error:   "invalid_json",
				Message: "must be a valid json.",
			})
			return
		}
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Error:   "missing_fields",
			Message: "email and password are required",
		})
		return
	}
	email := strings.TrimSpace(*req.Email)
	password := strings.TrimSpace(*req.Password)
	if err := helpers.HandleUserValidationJson(email, password); err != nil {
		helpers.WriteError(w, http.StatusUnprocessableEntity, *err)
		return
	}
	// validation done.

	user, err := s.Q.FindUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusUnauthorized, helpers.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email or password.",
			})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{})
		return
	}
	// email checkup sorted.

	isMatched, err := auth.CheckPassHash(password, user.HashedPassword)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{})
		return
	}
	if !isMatched {
		helpers.WriteError(w, http.StatusUnauthorized, helpers.ErrorResponse{
			Error:   "invalid_credentials",
			Message: "Invalid email or password.",
		})
		return
	}
	// user verified.

	jwt, err := auth.MakeJWT(user.ID, s.TokenSecret, time.Hour)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{})
		return
	}
	// jwt created
	type responseStruct struct {
		AccessToken string    `json:"access_token"`
		UserId      uuid.UUID `json:"user_id"`
		ExpiresIn   int       `json:"expires_in"`
	}
	helpers.WriteJson(w, http.StatusOK, responseStruct{
		AccessToken: jwt,
		UserId:      user.ID,
		ExpiresIn:   3600,
	})
}
func (s *Service) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := helpers.ReadJson(r, &req); err != nil || req.Email == nil || req.Password == nil {
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
				Error:   "invalid_json",
				Message: "must be a valid json.",
			})
			return
		}
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Error:   "missing_fields",
			Message: "email and password are required",
		})
		return
	}
	email := strings.TrimSpace(*req.Email)
	password := strings.TrimSpace(*req.Password)
	if err := helpers.HandleUserValidationJson(email, password); err != nil {
		helpers.WriteError(w, http.StatusUnprocessableEntity, *err)
		return
	}
	// validation done.
	hashed_password, err := auth.HashPassword(password)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{
			Error:   "internal_error",
			Message: "Something went wrong. Please try again later.",
		})
	}
	user, err := s.Q.CreateUser(r.Context(), db.CreateUserParams{
		Email:          email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		handleErr := helpers.HandleErrors(err)
		if handleErr != nil && handleErr.ViolationError {
			helpers.WriteError(w, http.StatusConflict, helpers.ErrorResponse{
				Error:   "user_exists",
				Message: "An account with this email, already exists.",
			})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{})
		return
	}

	helpers.WriteJson(w, http.StatusCreated, user)
}
