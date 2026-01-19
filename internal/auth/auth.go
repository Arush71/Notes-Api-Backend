package auth

import (
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(pass string) (string, error) {
	hashed, err := argon2id.CreateHash(pass, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashed, nil
}

func CheckPassHash(pass, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(pass, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "notesy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userId.String(),
	})
	token_str, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return token_str, nil
}
