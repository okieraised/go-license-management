package utils

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEd25519KeyPair(t *testing.T) {
	signingKey, verifyKey, err := NewEd25519KeyPair()
	assert.NoError(t, err)

	fmt.Println(signingKey)
	fmt.Println(verifyKey)
}

func TestLicenseKeyWithEd25519(t *testing.T) {
	licenseKey, err := NewLicenseKeyWithEd25519("302e020100300506032b6570042204201411064cece60c82fe80dd7ca82c6239fd7f2094fcfb5f27405e23b38c7cae9a", "sart")
	assert.NoError(t, err)

	fmt.Println(licenseKey)
}

func TestVerifyLicenseKeyWithEd25519(t *testing.T) {
	signingKey, verifyKey, err := NewEd25519KeyPair()
	assert.NoError(t, err)

	fmt.Println(signingKey)
	fmt.Println(verifyKey)

	licenseKey, err := NewLicenseKeyWithEd25519(signingKey, "sart")
	assert.NoError(t, err)

	fmt.Println(licenseKey)

	valid, data, err := VerifyLicenseKeyWithEd25519(verifyKey, licenseKey)
	assert.NoError(t, err)

	fmt.Println(valid)
	fmt.Println(data, string(data))
}

func TestDecodeJWT(t *testing.T) {
	token := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhZG1pbiIsImV4cCI6MTczNTYyMzk5MSwiaWF0IjoxNzM1NjIwMzkxLCJpc3MiOiJnby1saWNlbnNlLW1hbmFnZW1lbnQiLCJuYmYiOjE3MzU2MjAzOTEsInBlcm1pc3Npb25zIjpbImxpY2Vuc2UtZW50aXRsZW1lbnRzLmRldGFjaCIsImxpY2Vuc2UtdXNlcnMuYXR0YWNoIiwicG9saWN5X2VudGl0bGVtZW50cy5hdHRhY2giLCJwb2xpY3lfZW50aXRsZW1lbnRzLmRldGFjaCIsImxpY2Vuc2UtZW50aXRsZW1lbnRzLmF0dGFjaCIsIm1hY2hpbmUtaGVhcnRiZWF0LnJlc2V0IiwidGVuYW50LnJlYWQiLCJtYWNoaW5lLmNoZWNrLW91dCIsIm1hY2hpbmUuZGVsZXRlIiwibWFjaGluZS5yZWFkIiwiYWRtaW4uZGVsZXRlIiwidXNlcl9wYXNzd29yZC5yZXNldCIsImVudGl0bGVtZW50LnJlYWQiLCJsaWNlbnNlLmNoZWNrLW91dCIsInVzZXIuZGVsZXRlIiwicG9saWN5LmRlbGV0ZSIsImxpY2Vuc2UudXBkYXRlIiwibWFjaGluZS51cGRhdGUiLCJwcm9kdWN0X3Rva2Vucy5nZW5lcmF0ZSIsInBvbGljeS5jcmVhdGUiLCJsaWNlbnNlLXVzYWdlLnJlc2V0IiwibGljZW5zZS1wb2xpY3kudXBkYXRlIiwibWFjaGluZS5jcmVhdGUiLCJ0ZW5hbnQudXBkYXRlIiwiYWRtaW4uY3JlYXRlIiwidXNlci5iYW4iLCJ1c2VyX3Bhc3N3b3JkLnVwZGF0ZSIsImxpY2Vuc2UucmV2b2tlIiwibGljZW5zZS11c2FnZS5kZWNyZW1lbnQiLCJsaWNlbnNlLXRva2Vucy5nZW5lcmF0ZSIsInByb2R1Y3QuY3JlYXRlIiwicHJvZHVjdC5kZWxldGUiLCJwcm9kdWN0LnVwZGF0ZSIsImxpY2Vuc2UuY3JlYXRlIiwibGljZW5zZS11c2FnZS5pbmNyZW1lbnQiLCJtYWNoaW5lLWhlYXJ0YmVhdC5waW5nIiwicG9saWN5LnJlYWQiLCJsaWNlbnNlLnJlYWQiLCJ1c2VyLnJlYWQiLCJ1c2VyLnVuYmFuIiwidXNlci51cGRhdGUiLCJlbnRpdGxlbWVudC5jcmVhdGUiLCJsaWNlbnNlLmNoZWNrLWluIiwibGljZW5zZS11c2Vycy5kZXRhY2giLCJhZG1pbi51cGRhdGUiLCJwb2xpY3kudXBkYXRlIiwibGljZW5zZS5yZW5ldyIsInRlbmFudC5jcmVhdGUiLCJwcm9kdWN0LnJlYWQiLCJsaWNlbnNlLnZhbGlkYXRlIiwidXNlci5jcmVhdGUiLCJlbnRpdGxlbWVudC51cGRhdGUiLCJsaWNlbnNlLnJlaW5zdGF0ZSIsImxpY2Vuc2Uuc3VzcGVuZCIsImFkbWluLnJlYWQiLCJlbnRpdGxlbWVudC5kZWxldGUiLCJ0ZW5hbnQuZGVsZXRlIiwibGljZW5zZS5kZWxldGUiXSwic3ViIjoidXNlcjEiLCJ0ZW5hbnQiOiJ0ZXN0In0.vVSbLQgPTyezexIS5mfE0FRqcGCNOxY6ZxoZJW60tXAGNWdr2gjjjD0B9W-zfoCWv_N5J6VhV-AgfCG9rs4xBQ"

	publicKey, err := hex.DecodeString("302a300506032b6570032100462c6f6e52b3b35f8546a8447865c09b6ace033101c43bddd0fcc181a0fd776b")
	assert.NoError(t, err)

	decodedPublicKey, err := x509.ParsePKIXPublicKey(publicKey)
	assert.NoError(t, err)

	publicKey, ok := decodedPublicKey.(ed25519.PublicKey)
	assert.True(t, ok)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return decodedPublicKey, nil
	})
	assert.NoError(t, err)

	fmt.Println(parsedToken.Claims.(jwt.MapClaims))

}
