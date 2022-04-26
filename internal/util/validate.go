package util

import (
	"fmt"
	"github.com/ianaspivac/news-site-go/internal/httperr"
	"net/http"
	"net/mail"
	"unicode"
)

func ValidateMail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return httperr.New(err.Error(), http.StatusBadRequest)
	}
	return nil
}

func ValidatePassword(pass string) error {
	if len(pass) < 8 {
		return httperr.ValidationError("password", "should be at least 8 characters long")
	}

	var (
		numberPresent  bool
		upperPresent   bool
		specialPresent bool
	)
	for _, c := range pass {
		switch {
		case unicode.IsNumber(c):
			numberPresent = true
		case unicode.IsUpper(c):
			upperPresent = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			specialPresent = true
		case unicode.IsLetter(c):
			continue
		default:
			return httperr.ValidationError("password", fmt.Sprintf("unsupported character: %c", c))
		}
	}

	if !numberPresent {
		return httperr.ValidationError("password", "should contain at least one number")
	} else if !upperPresent {
		return httperr.ValidationError("password", "should contain at least one uppercase character")
	} else if !specialPresent {
		return httperr.ValidationError("password", "should contain at least one special character")
	}

	return nil
}
