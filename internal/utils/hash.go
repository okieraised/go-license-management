package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"strconv"
	"strings"
)

const (
	SALT_BYTE_SIZE    = 24
	HASH_BYTE_SIZE    = 24
	PBKDF2_ITERATIONS = 1000
)

func Hash(password string) (string, error) {
	salt := make([]byte, SALT_BYTE_SIZE)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", errors.New("failed to generate random salt")
	}

	hbts := pbkdf2.Key([]byte(password), salt, PBKDF2_ITERATIONS, HASH_BYTE_SIZE, sha1.New)

	return fmt.Sprintf("%v:%v:%v",
		PBKDF2_ITERATIONS,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hbts)), nil
}

func Verify(raw, hash string) (bool, error) {
	hparts := strings.Split(hash, ":")

	itr, err := strconv.Atoi(hparts[0])
	if err != nil {
		return false, err
	}
	salt, err := base64.StdEncoding.DecodeString(hparts[1])
	if err != nil {

		return false, err
	}

	hsh, err := base64.StdEncoding.DecodeString(hparts[2])
	if err != nil {
		return false, err
	}

	rhash := pbkdf2.Key([]byte(raw), salt, itr, len(hsh), sha1.New)
	return equal(rhash, hsh), nil
}

// bytes comparisons
func equal(h1, h2 []byte) bool {
	diff := uint32(len(h1)) ^ uint32(len(h2))
	for i := 0; i < len(h1) && i < len(h2); i++ {
		diff |= uint32(h1[i] ^ h2[i])
	}

	return diff == 0
}

func HashPassword(password string) (string, error) {
	bPassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bPassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil

}

func Compare(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))
	return err == nil
}
