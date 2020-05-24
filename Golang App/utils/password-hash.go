package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// CreatePassword will create password using bcrypt
func CreatePassword(passwordString string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordString), 8)
	if err != nil {
		return "", errors.New("Error occurred while creating a Hash")
	}

	return string(hashedPassword), nil
}

// ComparePasswords will create password using bcrypt
func ComparePasswords(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("The '" + password + "' and '" + hashedPassword + "' strings don't match")
	}
	return nil
}
