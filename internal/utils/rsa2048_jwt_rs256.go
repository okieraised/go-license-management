package utils

import (
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// NewLicenseKeyWithJWTRS256 generates new license key in jwt format using RS256 signing method.
func NewLicenseKeyWithJWTRS256(signingKey string, data any) (string, error) {

	privateKeyPem, err := base64.StdEncoding.DecodeString(signingKey)
	if err != nil {
		return "", err
	}

	mapped, ok := data.(map[string]interface{})
	if !ok {
		return "", errors.New("invalid jwt data structure, must be map[string]interface{}")
	}

	var exp = time.Now().Add(time.Duration(2592000) * time.Second)
	expiry, ok := mapped["expiry"].(time.Time)
	if !ok {
		return "", errors.New("invalid jwt expiry field, must be time.Time{}")
	}
	if !expiry.IsZero() {
		exp = expiry
	}

	createdAt, ok := mapped["created_at"].(time.Time)
	if !ok {
		createdAt = time.Now()
	}

	tenantID, ok := mapped["tenant_id"]
	if !ok {
		return "", errors.New("invalid jwt tenant_id field, must be string")
	}

	licenseID, ok := mapped["license_id"]
	if !ok {
		return "", errors.New("invalid jwt license_id field, must be string")
	}

	productID, ok := mapped["product_id"]
	if !ok {
		return "", errors.New("invalid jwt product_id field, must be string")
	}

	metadata, ok := mapped["metadata"]
	if !ok {
		metadata = map[string]interface{}{}
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":      "go-license-management",
		"aud":      tenantID,
		"jit":      licenseID,
		"sub":      productID,
		"exp":      exp.Unix(),
		"iat":      createdAt.Unix(),
		"nbf":      createdAt.Unix(),
		"metadata": metadata,
	})

	pkey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPem)
	if err != nil {
		return "", err
	}

	licenseKey, err := claims.SignedString(pkey)
	if err != nil {
		return "", err
	}

	return licenseKey, nil
}
