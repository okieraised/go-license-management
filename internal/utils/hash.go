package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password using default cost
func HashPassword(password string) (string, error) {
	bPassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bPassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil

}

// CompareHashedPassword compares input plaintext with its possible hash value
func CompareHashedPassword(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))
	return err == nil
}
