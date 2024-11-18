package utils

import (
	"bytes"
	"math/rand"
	"time"
)

func GenerateToken() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const segmentLength = 8
	const totalLength = segmentLength * 5

	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var buffer bytes.Buffer

	for i := 0; i < totalLength; i++ {
		buffer.WriteByte(charset[rd.Intn(len(charset))])
		if (i+1)%segmentLength == 0 && i != totalLength-1 {
			buffer.WriteByte('-')
		}
	}

	return buffer.String()
}
