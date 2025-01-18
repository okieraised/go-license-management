package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	cypher, err := Encrypt([]byte("hehe"), []byte("secret"))
	assert.NoError(t, err)

	fmt.Println(string(cypher))

}

func TestHashPassword(t *testing.T) {
	password := "abcd1234"
	hashed, err := HashPassword(password)
	assert.NoError(t, err)

	fmt.Println(hashed)
}
