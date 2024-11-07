package utils

import (
	"fmt"
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

	valid, _, err := VerifyLicenseKeyWithEd25519(verifyKey, licenseKey)
	assert.NoError(t, err)

	fmt.Println(valid)
}
