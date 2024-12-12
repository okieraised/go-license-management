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
