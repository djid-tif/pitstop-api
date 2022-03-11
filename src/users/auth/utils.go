package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cryptCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
