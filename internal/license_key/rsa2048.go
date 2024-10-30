package license_key

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

const RSAPrivateKeyStr = "RSA PRIVATE KEY"
const RSAPublicKeyStr = "RSA PUBLIC KEY"

// NewRSA2048KeyPair generates the private key and the public key pair using RSA2048 algorithm
func NewRSA2048KeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// Encode the private key to PEM format (PKCS1)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  RSAPrivateKeyStr,
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Encode the public key to PEM format (PKCS1)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  RSAPublicKeyStr,
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

// NewLicenseKeyWithRSA2048 generates new license key using RSA2048 algorithm
func NewLicenseKeyWithRSA2048(signingKey string, data any) (string, error) {
	bData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Encode the original data to base64
	encodedData := base64.StdEncoding.EncodeToString(bData)

	// Sign the data using the private key with SHA-256 hashing
	hash := sha256.New()
	hash.Write(bData)
	hashed := hash.Sum(nil)

	// Decode the private key string
	block, _ := pem.Decode([]byte(signingKey))

	if block == nil || block.Type != RSAPrivateKeyStr {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	// Encode the signature in base64
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	// Combine the encoded data and signature to create the license key
	licenseKey := fmt.Sprintf("%s.%s", encodedData, encodedSignature)

	return licenseKey, nil
}

// VerifyLicenseKeyWithRSA2048 verifies a license key against the provided public key using Ed25519 algorithm
func VerifyLicenseKeyWithRSA2048(verifyKey string, licenseKey string) (bool, []byte, error) {
	parts := strings.Split(licenseKey, ".")
	if len(parts) != 2 {
		return false, nil, errors.New("invalid license key format")
	}
	encodedData := parts[0]
	encodedSignature := parts[1]

	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return false, nil, err
	}

	// Decode signature
	signature, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		return false, nil, err
	}

	// Decode the private key string
	block, _ := pem.Decode([]byte(verifyKey))

	if block == nil || block.Type != RSAPublicKeyStr {
		return false, nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return false, nil, err
	}

	// Check sum of data
	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, nil, err
	}

	return true, data, nil
}
