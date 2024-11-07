package utils

import (
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go-license-management/internal/infrastructure/models/license_key"
	"math"
	"time"
)

// NewLicenseKeyWithJWTRS256 generates new license key in jwt format using RS256 signing method.
func NewLicenseKeyWithJWTRS256(signingKey string, data *license_key.LicenseKeyContent) (string, error) {

	block, _ := pem.Decode([]byte(signingKey))

	if block == nil || block.Type != RSAPrivateKeyStr {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	var exp int64 = math.MaxInt64
	if !data.Expiry.IsZero() {
		exp = data.Expiry.Unix()
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "go-license-management",
		"aud": DerefPointer(data.TenantID),
		"jit": DerefPointer(data.LicenseID),
		"sub": DerefPointer(data.ProductID),
		"exp": exp,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
	})

	licenseKey, err := claims.SignedString(block.Bytes)
	if err != nil {
		return "", err
	}

	return licenseKey, nil
}
