package helpers

import (
	"errors"
	"net/mail"

	"github.com/lib/pq"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func HandleUserValidationJson(email string, password string) *ErrorResponse {
	emailValid := ValidateEmail(email)
	passwordValid := ValidatePassword(password)
	if emailValid && passwordValid {
		return nil
	}
	errRespnse := &ErrorResponse{
		Error:  "invalid_fields",
		Fields: make(map[string]string),
	}
	if !passwordValid {
		(*errRespnse).Fields["password"] = "password should be 8 or more characters"
	}
	if !emailValid {
		errRespnse.Fields["email"] = "email is invalid."
	}
	return errRespnse
}

type errStruct struct {
	ViolationError bool
	NoRowError     bool
}

func HandleErrors(err error) *errStruct {
	var pqerr *pq.Error
	if errors.As(err, &pqerr) {
		if pqerr.Code == "23505" {
			errS := &errStruct{
				ViolationError: true,
			}
			return errS
		}
	}
	return nil
}
