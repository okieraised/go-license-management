package service

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/utils"
	"time"
)

func (svc *AuthenticationService) generateJWT(ctx *gin.Context, tenant *entities.Tenant, account *entities.Account) (string, int64, error) {

	jwtPermissions := make([]string, 0)
	switch account.RoleName {
	case constants.RoleUser:
		for k, _ := range constants.UserPermissionMapper {
			jwtPermissions = append(jwtPermissions, k)
		}
	case constants.RoleAdmin:
		for k, _ := range constants.AdminPermissionMapper {
			jwtPermissions = append(jwtPermissions, k)
		}
	case constants.RoleSuperAdmin:
		for k, _ := range constants.SuperAdminPermissionMapper {
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

	signer := &utils.ED25519Signer{}
	signer.SetPrivateKey(privateKey)

	tokenString, err := claims.SignedString(signer)
	if err != nil {
		return "", 0, err
	}

	return tokenString, exp, nil
}
