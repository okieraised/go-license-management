package service

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/permissions"
	"time"
)

func (svc *AuthenticationService) generateSuperadminJWT(ctx *gin.Context, master *entities.Master) (string, int64, error) {

	jwtPermissions := make([]string, 0)
	for k, _ := range permissions.SuperAdminPermissionMapper {
		jwtPermissions = append(jwtPermissions, k)
	}

	now := time.Now()
	exp := now.Add(time.Hour).Unix()
	claims := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub":         master.Username,   // Subject (user identifier)
		"iss":         constants.AppName, // Issuer
		"aud":         master.RoleName,   // Audience (user role)
		"exp":         exp,               // Expiration time
		"iat":         now.Unix(),
		"nbf":         now.Unix(),
		"tenant":      "*",
		"status":      constants.AccountStatusActive,
		"permissions": jwtPermissions,
	})

	privateKey, err := base64.StdEncoding.DecodeString(master.Ed25519PrivateKey)
	if err != nil {
		return "", 0, err
	}

	decodedPrivateKey, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return "", 0, err
	}

	privateKey, ok := decodedPrivateKey.(ed25519.PrivateKey)
	if !ok {
		return "", 0, errors.New("decoded key is not of type ed25519.PrivateKey")
	}

	tokenString, err := claims.SignedString(decodedPrivateKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, exp, nil
}

func (svc *AuthenticationService) generateJWT(ctx *gin.Context, tenant *entities.Tenant, account *entities.Account) (string, int64, error) {

	jwtPermissions := make([]string, 0)
	switch account.RoleName {
	case constants.RoleUser:
		for k, _ := range permissions.UserPermissionMapper {
			jwtPermissions = append(jwtPermissions, k)
		}
	case constants.RoleAdmin:
		for k, _ := range permissions.AdminPermissionMapper {
			jwtPermissions = append(jwtPermissions, k)
		}
	case constants.RoleSuperAdmin:
		for k, _ := range permissions.SuperAdminPermissionMapper {
			jwtPermissions = append(jwtPermissions, k)
		}
	}

	now := time.Now()
	exp := now.Add(time.Hour).Unix()
	claims := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub":         account.Username,  // Subject (user identifier)
		"iss":         constants.AppName, // Issuer
		"aud":         account.RoleName,  // Audience (user role)
		"exp":         exp,               // Expiration time
		"iat":         now.Unix(),
		"nbf":         now.Unix(),
		"tenant":      account.TenantName,
		"status":      account.Status,
		"permissions": jwtPermissions,
	})

	privateKey, err := hex.DecodeString(tenant.Ed25519PrivateKey)
	if err != nil {
		return "", 0, err
	}

	decodedPrivateKey, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return "", 0, err
	}

	privateKey, ok := decodedPrivateKey.(ed25519.PrivateKey)
	if !ok {
		return "", 0, errors.New("decoded key is not of type ed25519.PrivateKey")
	}

	tokenString, err := claims.SignedString(decodedPrivateKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, exp, nil
}
