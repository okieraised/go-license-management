package middlewares

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/database/postgres"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/utils"
	"net/http"
	"strings"
)

func JWTValidationMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader(constants.AuthorizationHeader)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		authHdrPart := strings.Split(authHeader, " ")
		switch len(authHdrPart) {
		case 2:
			if authHdrPart[0] != constants.AuthorizationTypeBearer {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
				return
			}

			tenantName := ctx.Param("tenant_name")
			if tenantName == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
				return
			}

			tenant := &entities.Tenant{
				Name: tenantName,
			}
			err := postgres.GetInstance().NewSelect().Model(tenant).WherePK().Scan(ctx)
			if err != nil {
				logging.GetInstance().GetLogger().Error(err.Error())
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
				return
			}

			publicKey, err := hex.DecodeString(tenant.Ed25519PublicKey)
			if err != nil {
				logging.GetInstance().GetLogger().Error(err.Error())
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
				return
			}

			decodedPublicKey, err := x509.ParsePKIXPublicKey(publicKey)
			if err != nil {
				logging.GetInstance().GetLogger().Error(err.Error())
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
				return
			}

			publicKey, ok := decodedPublicKey.(ed25519.PublicKey)
			if !ok {
				logging.GetInstance().GetLogger().Error("public key not of type ed25519")
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, nil)
				return
			}

			parsedToken, err := jwt.Parse(authHdrPart[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*utils.SigningMethodEdDSA); !ok {
					logging.GetInstance().GetLogger().Error(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return publicKey, nil
			})
			if err != nil {
				logging.GetInstance().GetLogger().Error("public key not of type ed25519")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
				return
			}

			permissionsClaims, ok := parsedToken.Claims.(jwt.MapClaims)["permissions"]
			if !ok {
				logging.GetInstance().GetLogger().Error("missing [permissions] claims")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
				return
			}

			permissions, ok := permissionsClaims.([]interface{})
			if !ok {
				logging.GetInstance().GetLogger().Error("invalid [permissions] claims")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
				return
			}

			ctx.Set(constants.ContextValuePermissions, permissions)
		default:
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		ctx.Next()
	}
}
