package utils

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// NewEd25519KeyPair generates the private signing key and the public verify key using Ed25519 algorithm
func NewEd25519KeyPair() (string, string, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", err
	}

	// Export the private key in PKCS#8 DER format
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	signingKey := hex.EncodeToString(privateKeyBytes)

	// Export the public key in SPKI DER format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	verifyKey := hex.EncodeToString(publicKeyBytes)

	return signingKey, verifyKey, nil
}

// NewLicenseKeyWithEd25519 generates new license key using Ed25519 algorithm
func NewLicenseKeyWithEd25519(signingKey string, data any) (string, error) {
	privateKeyBytes, err := hex.DecodeString(signingKey)
	if err != nil {
		return "", err
	}

	decodedPrivateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return "", err
	}

	// Assert that it is of type ed25519.PrivateKey
	privateKey, ok := decodedPrivateKey.(ed25519.PrivateKey)
	if !ok {
		return "", errors.New("decoded key is not of type ed25519.PrivateKey")
	}

	bData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Sign the data with the private key
	signature := ed25519.Sign(privateKey, bData)
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	encodedData := base64.StdEncoding.EncodeToString(bData)

	// Combine the encoded data and signature to create the license key
	licenseKey := fmt.Sprintf("%s.%s", encodedSignature, encodedData)

	return licenseKey, nil
}

// VerifyLicenseKeyWithEd25519 verifies a license key against the provided public key using Ed25519 algorithm
func VerifyLicenseKeyWithEd25519(verifyKey string, licenseKey string) (bool, []byte, error) {
	parts := strings.Split(licenseKey, ".")
	if len(parts) != 2 {
		return false, nil, errors.New("invalid license key format")
	}
	encodedData := parts[1]
	encodedSignature := parts[0]

	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return false, nil, err
	}

	signature, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		return false, nil, err
	}

	publicKeyBytes, err := hex.DecodeString(verifyKey)
	if err != nil {
		return false, nil, err
	}

	decodedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return false, nil, err
	}

	publicKey, ok := decodedPublicKey.(ed25519.PublicKey)
	if !ok {
		return false, nil, errors.New("decoded key is not of type ed25519.PublicKey")
	}

	return ed25519.Verify(publicKey, data, signature), data, nil
}
