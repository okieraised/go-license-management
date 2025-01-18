package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRSA2048KeyPair(t *testing.T) {
	privateKey, publicKey, err := NewRSA2048PKCS1KeyPair()
	assert.NoError(t, err)
	fmt.Println(privateKey)
	fmt.Println(publicKey)
}

func TestNewLicenseKeyWithRSA2048(t *testing.T) {
	privateKey, _, err := NewRSA2048PKCS1KeyPair()
	assert.NoError(t, err)

	licenseKey, err := NewLicenseKeyWithRSA2048PKCS1(privateKey, "dadada")
	assert.NoError(t, err)

	fmt.Println(licenseKey)
}

func TestVerifyLicenseKeyWithRSA2048(t *testing.T) {
	privateKey, publicKey, err := NewRSA2048PKCS1KeyPair()
	assert.NoError(t, err)
	fmt.Println(privateKey)
	fmt.Println(publicKey)

	licenseKey, err := NewLicenseKeyWithRSA2048PKCS1(privateKey, "dadada")
	assert.NoError(t, err)

	fmt.Println(licenseKey)

	valid, _, err := VerifyLicenseKeyWithRSA2048PKCS1(publicKey, licenseKey)
	assert.NoError(t, err)
	fmt.Println(valid)
}
