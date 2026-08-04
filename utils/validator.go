package utils

import (
	"net/mail"
)

// return bool

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
