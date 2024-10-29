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
		fmt.Printf("wrong hash %v", hash)
		return false, errors.New("wrong hash, iteration is invalid")
	}
	salt, err := base64.StdEncoding.DecodeString(hparts[1])
	if err != nil {
		fmt.Print("wrong hash, salt error:", err)
		return false, errors.New("wrong hash, salt error:" + err.Error())
	}

	hsh, err := base64.StdEncoding.DecodeString(hparts[2])
	if err != nil {
		fmt.Print("wrong hash, hash error:", err)
		return false, errors.New("wrong hash, hash error:" + err.Error())
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

// Hash password

func hashPassword(password string) (string, error) {

	// Convert password string to byte slice

	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err

}

// Check if two passwords match using Bcrypt's CompareHashAndPassword

// which return nil on success and an error on failure.

func doPasswordsMatch(hashedPassword, currPassword string) bool {

	err := bcrypt.CompareHashAndPassword(

		[]byte(hashedPassword), []byte(currPassword))

	return err == nil

}