package middleware

import (
	"context"
	"net/http"
	"notes-api/internal/auth"
	"notes-api/internal/helpers"
	"notes-api/internal/helpers/requestctx"
)

func MakeAuthMiddleWare(tokenSecret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			access_token, err := auth.GetBearerToken(r.Header)
			if err != nil {
				helpers.WriteError(w, http.StatusUnauthorized, helpers.ErrorResponse{
					Error:   "unauthorized",
					Message: "Mission or invalid authentication token.",
				})
				return
			}
			user, err := auth.ValidateJWT(access_token, tokenSecret)
			if err != nil {
				helpers.WriteError(w, http.StatusUnauthorized, helpers.ErrorResponse{
					Error:   "unauthorized",
					Message: "invalid token",
				})
				return
			}
			ctx := context.WithValue(r.Context(), requestctx.UserKey, user)
			r = r.WithContext(ctx)
			next(w, r)
		}
	}
}
