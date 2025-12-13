package lib

import (
	"net/mail"
		"errors"
	"regexp"
	"strings"
)

var (
	// E.164 format: optional +, country code + subscriber number
	// Total digits: 8–15
	e164Regex = regexp.MustCompile(`^\+?[1-9]\d{7,14}$`)
)

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidatePhone(phone string) (string, error) {
	normalized := strings.TrimSpace(phone)

	if normalized == "" {
		return "", errors.New("phone number is required")
	}

	if !e164Regex.MatchString(normalized) {
		return "", errors.New("invalid phone number format")
	}

	//remove leading '+' for storage consistency
	if normalized[0] == '+' {
		normalized = normalized[1:]
	}

	return normalized, nil
}
