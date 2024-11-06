package tokens

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestGenerateNewToken(t *testing.T) {
	for _ = range 10 {
		//content := map[string]interface{}{
		//	"product_id": "asdasdadsa",
		//	"iat":        time.Now().Unix(),
		//	"perm":       []string{"*"},
		//}

		content := "1e4fb4f4-4123-4400-b605-cc186878ecbe" + strconv.FormatInt(time.Now().Unix(), 10)

		token, err := GenerateToken([]byte("1e4fb4f4-4123-4400-b605-cc186878ecbe"), content)
		assert.NoError(t, err)
		fmt.Println("token", token)

		original, err := DecryptToken([]byte("1e4fb4f4-4123-4400-b605-cc186878ecbe"), token)
		assert.NoError(t, err)
		fmt.Println("original", original)
	}
}
