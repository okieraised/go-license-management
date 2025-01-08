package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type test struct {
	xx string
}

func main() {

	hash := sha256.Sum256([]byte(`{"request_id":"4f1684f1-dd10-4eaf-a7aa-2aefc52847c0","code":"00000","message":"OK","server_time":1735866210,"data":{"username":"used1relu22g2dd4","role_name":"user","email":"122eel32d2gd2d3@gmail.com","first_name":"dfec92f3-d60c-48e7-b123-a3e48bbe4829","last_name":"user","status":"active","metadata":null,"created_at":"2025-01-03T08:03:30.897248571+07:00","updated_at":"2025-01-03T08:03:30.897248571+07:00"}}`))

	fmt.Println(hex.EncodeToString(hash[:]))

	var f *test

	if f == nil {
		fmt.Println("f is nil")
	}

}
