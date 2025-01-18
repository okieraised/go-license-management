package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Encrypt(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

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
