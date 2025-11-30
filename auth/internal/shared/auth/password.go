package auth

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/DangeL187/erax"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", erax.Wrap(err, "failed to hash password")
	}

	return string(bytes), nil
}

func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}
