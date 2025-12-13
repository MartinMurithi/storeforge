package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hashpassword(password string) (string, error) {
	const COST int = 12

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), COST)

	if err != nil {
		return "", fmt.Errorf("error hashing password %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
