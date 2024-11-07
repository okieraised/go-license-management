package utils

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"go-license-management/internal/infrastructure/models/license_key"
	"math"
	"time"
)

// NewLicenseKeyWithJWTRS256 generates new license key in jwt format using RS256 signing method.
func NewLicenseKeyWithJWTRS256(signingKey string, data *license_key.LicenseKeyContent) (string, error) {

	privateKeyPem, err := base64.URLEncoding.DecodeString(signingKey)
	if err != nil {
		return "", err
	}

	var exp int64 = math.MaxInt32
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
