package utils

import (
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/MartinMurithi/storeforge/pkg/errors"
)

var e164Regex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return apperrors.ErrInvalidEmailFormat
	}
	return nil
}

func ValidatePhone(phone string) (string, error) {
	trimmed := strings.TrimSpace(phone)
	if trimmed == "" {
		return "", apperrors.ErrPhoneRequired
	}
	if !e164Regex.MatchString(trimmed) {
		return "", apperrors.ErrInvalidPhoneNumber
	}

	// remove leading '+' for storage
	if trimmed[0] == '+' {
		trimmed = trimmed[1:]
	}
	return trimmed, nil
}


func IsValidDate(format, dateString string) bool {
	t, err := time.Parse(format, dateString)
	if err != nil {
		return false
	}
	return t.Format(format) == dateString
}
