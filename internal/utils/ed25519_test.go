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
	token := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhZG1pbiIsImV4cCI6MTczNTI5MDQ2NCwiaWF0IjoxNzM1Mjg2ODY0LCJpc" +
		"3MiOiJnby1saWNlbnNlLW1hbmFnZW1lbnQiLCJuYmYiOjE3MzUyODY4NjQsInBlcm1pc3Npb25zIjpbInByb2R1Y3RfdG9rZW5zLmdlbmVyYXRl" +
		"IiwicHJvZHVjdC51cGRhdGUiLCJsaWNlbnNlLXVzZXJzLmF0dGFjaCIsInRlbmFudC5jcmVhdGUiLCJhZG1pbi51cGRhdGUiLCJ1c2VyLnVwZGF" +
		"0ZSIsIm1hY2hpbmUuY2hlY2stb3V0IiwiYWRtaW4uZGVsZXRlIiwicHJvZHVjdC5kZWxldGUiLCJwb2xpY3kuZGVsZXRlIiwibWFjaGluZS5jcm" +
		"VhdGUiLCJ1c2VyLmRlbGV0ZSIsImxpY2Vuc2UtdXNlcnMuZGV0YWNoIiwibWFjaGluZS5kZWxldGUiLCJwcm9kdWN0LnJlYWQiLCJwb2xpY3kuY" +
		"3JlYXRlIiwibGljZW5zZS10b2tlbnMuZ2VuZXJhdGUiLCJtYWNoaW5lLnJlYWQiLCJ1c2VyLnVuYmFuIiwiZW50aXRsZW1lbnQuY3JlYXRlIiwi" +
		"bGljZW5zZS5jaGVjay1pbiIsImxpY2Vuc2Uuc3VzcGVuZCIsInRlbmFudC5kZWxldGUiLCJ1c2VyLmNyZWF0ZSIsImFkbWluLnJlYWQiLCJwb2x" +
		"pY3kucmVhZCIsImxpY2Vuc2UtdXNhZ2UuaW5jcmVtZW50IiwibGljZW5zZS11c2FnZS5yZXNldCIsInRlbmFudC5yZWFkIiwicG9saWN5LnVwZG" +
		"F0ZSIsImxpY2Vuc2UuY2hlY2stb3V0IiwibGljZW5zZS1lbnRpdGxlbWVudHMuZGV0YWNoIiwibGljZW5zZS1lbnRpdGxlbWVudHMuYXR0YWNoI" +
		"iwibWFjaGluZS1oZWFydGJlYXQucmVzZXQiLCJwb2xpY3lfZW50aXRsZW1lbnRzLmF0dGFjaCIsImxpY2Vuc2UuY3JlYXRlIiwibGljZW5zZS5k" +
		"ZWxldGUiLCJsaWNlbnNlLnJlYWQiLCJ1c2VyX3Bhc3N3b3JkLnJlc2V0IiwicHJvZHVjdC5jcmVhdGUiLCJsaWNlbnNlLnJlaW5zdGF0ZSIsImF" +
		"kbWluLmNyZWF0ZSIsInVzZXIuYmFuIiwidXNlcl9wYXNzd29yZC51cGRhdGUiLCJsaWNlbnNlLnJldm9rZSIsInRlbmFudC51cGRhdGUiLCJwb2" +
		"xpY3lfZW50aXRsZW1lbnRzLmRldGFjaCIsIm1hY2hpbmUtaGVhcnRiZWF0LnBpbmciLCJ1c2VyLnJlYWQiLCJlbnRpdGxlbWVudC5yZWFkIiwiZ" +
		"W50aXRsZW1lbnQuZGVsZXRlIiwiZW50aXRsZW1lbnQudXBkYXRlIiwibGljZW5zZS51cGRhdGUiLCJsaWNlbnNlLXBvbGljeS51cGRhdGUiLCJs" +
		"aWNlbnNlLnJlbmV3IiwibGljZW5zZS52YWxpZGF0ZSIsImxpY2Vuc2UtdXNhZ2UuZGVjcmVtZW50IiwibWFjaGluZS51cGRhdGUiXSwic3ViIjo" +
		"idXNlcjEifQ.XNHjszQoW1QUTthXZVuKuZhu2Pu1qoOamrcCkTsAMT30Ftgfr1pPfKC84HiW_Drg1O3BrdG6AzSog_kviH9LCw"

	publicKey, err := hex.DecodeString("302a300506032b6570032100462c6f6e52b3b35f8546a8447865c09b6ace033101c43bddd0fcc181a0fd776b")
	assert.NoError(t, err)

	decodedPublicKey, err := x509.ParsePKIXPublicKey(publicKey)
	assert.NoError(t, err)

	publicKey, ok := decodedPublicKey.(ed25519.PublicKey)
	assert.True(t, ok)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is ED25519
		if _, ok := token.Method.(*SigningMethodEdDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	fmt.Println(parsedToken.Claims.(jwt.MapClaims)["permissions"].([]interface{}))

}
