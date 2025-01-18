package utils

import (
	"fmt"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	for _ = range 10 {
		fmt.Println(GenerateToken())
	}
}
