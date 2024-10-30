package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "abcd1234"
	hashed, err := HashPassword(password)
	assert.NoError(t, err)

	fmt.Println(hashed)
}
