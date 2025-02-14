package middlewares

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-license-management/internal/cerrors"
	"go-license-management/internal/config"
	"go-license-management/internal/constants"
	"go-license-management/internal/infrastructure/database/entities"
	"go-license-management/internal/infrastructure/database/postgres"
	"go-license-management/internal/infrastructure/logging"
	"go-license-management/internal/response"
	"net/http"
	"strings"
	"time"
)

func JWTValidationMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logging.GetInstance().GetLogger().Info("validating jwt token")

		authHeader := ctx.GetHeader(constants.AuthorizationHeader)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.NewResponse(ctx).ToResponse(
					cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
					"missing authorization header",
					nil,
					nil,
					nil,
				),
			)
			return
		}

		authHdrPart := strings.Split(authHeader, " ")
		switch len(authHdrPart) {
		case 2:
			if authHdrPart[0] != constants.AuthorizationTypeBearer {
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"invalid bearer authorization header",
						nil,
						nil,
						nil,
					),
				)
				return
			}

			tenantName := ctx.Param("tenant_name")
			if tenantName == "" {
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"invalid request url",
						nil,
						nil,
						nil,
					),
				)
				return
			}

			tenant := &entities.Tenant{
				Name: tenantName,
			}
			err := postgres.GetInstance().NewSelect().Model(tenant).WherePK().Scan(ctx)
			if err != nil {
				logging.GetInstance().GetLogger().Error(err.Error())
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer],
						cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer],
						nil,
						nil,
						nil,
					),
				)
				return
			}

			publicKey, err := base64.StdEncoding.DecodeString(tenant.Ed25519PublicKey)
			if err != nil {
				logging.GetInstance().GetLogger().Error(err.Error())
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer],
						cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer],
						nil,
						nil,
						nil,
					),
				)
				return
			}

			decodedPublicKey, err := x509.ParsePKIXPublicKey(publicKey)
			if err != nil {
				logging.GetInstance().GetLogger().Error(err.Error())
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer],
						cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer],
						nil,
						nil,
						nil,
					),
				)
				return
			}

			var parsedToken *jwt.Token
			parsedToken, err = jwt.Parse(authHdrPart[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
					logging.GetInstance().GetLogger().Error(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return decodedPublicKey, nil
			})
			if err != nil {
				master := &entities.Master{
					Username: config.SuperAdminUsername,
				}
				err = postgres.GetInstance().NewSelect().Model(master).WherePK().Scan(ctx)
				if err != nil {
					logging.GetInstance().GetLogger().Error(err.Error())
					ctx.AbortWithStatusJSON(
						http.StatusInternalServerError,
						response.NewResponse(ctx).ToResponse(
							cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer],
							cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer],
							nil,
							nil,
							nil,
						),
					)
					return
				}

				publicKey, err = base64.StdEncoding.DecodeString(master.Ed25519PublicKey)
				if err != nil {
					logging.GetInstance().GetLogger().Error(err.Error())
					ctx.AbortWithStatusJSON(
						http.StatusInternalServerError,
						response.NewResponse(ctx).ToResponse(
							cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer],
							cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer],
							nil,
							nil,
							nil,
						),
					)
					return
				}

				decodedPublicKey, err = x509.ParsePKIXPublicKey(publicKey)
				if err != nil {
					logging.GetInstance().GetLogger().Error(err.Error())
					ctx.AbortWithStatusJSON(
						http.StatusInternalServerError,
						response.NewResponse(ctx).ToResponse(
							cerrors.ErrCodeMapper[cerrors.ErrGenericInternalServer],
							cerrors.ErrMessageMapper[cerrors.ErrGenericInternalServer],
							nil,
							nil,
							nil,
						),
					)
					return
				}
				parsedToken, err = jwt.Parse(authHdrPart[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
						logging.GetInstance().GetLogger().Error(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return decodedPublicKey, nil
				})
				if err != nil {
					logging.GetInstance().GetLogger().Error(err.Error())
					ctx.AbortWithStatusJSON(
						http.StatusUnauthorized,
						response.NewResponse(ctx).ToResponse(
							cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
							cerrors.ErrMessageMapper[cerrors.ErrGenericUnauthorized],
							nil,
							nil,
							nil,
						),
					)
					return
				}
			}

			exp, err := parsedToken.Claims.GetExpirationTime()
			if err != nil {
				logging.GetInstance().GetLogger().Error(err.Error())
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"invalid [exp] claim",
						nil,
						nil,
						nil,
					),
				)
				return
			}

			if exp.Before(time.Now()) {
				logging.GetInstance().GetLogger().Error("token has expired")
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"token has expired",
						nil,
						nil,
						nil,
					),
				)
				return
			}

			audience, err := parsedToken.Claims.GetAudience()
			if err != nil {
				logging.GetInstance().GetLogger().Error(err.Error())
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"invalid [aud] claims",
						nil,
						nil,
						nil,
					),
				)
				return
			}
			ctx.Set(constants.ContextValueAudience, audience)

			subject, err := parsedToken.Claims.GetSubject()
			if err != nil {
				logging.GetInstance().GetLogger().Error(err.Error())
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"invalid [sub] claims",
						nil,
						nil,
						nil,
					),
				)
				return
			}
			ctx.Set(constants.ContextValueSubject, subject)

			tenantClaims, ok := parsedToken.Claims.(jwt.MapClaims)["tenant"]
			if !ok {
				logging.GetInstance().GetLogger().Error("missing [tenant] claims")
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"missing [tenant] claims",
						nil,
						nil,
						nil,
					),
				)
				return
			}

			tenantCtx, ok := tenantClaims.(interface{})
			if !ok {
				logging.GetInstance().GetLogger().Error("invalid [tenant] claims")
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"invalid [tenant] claims",
						nil,
						nil,
						nil,
					),
				)
				return
			}
			ctx.Set(constants.ContextValueTenant, tenantCtx)

			statusClaims, ok := parsedToken.Claims.(jwt.MapClaims)["status"]
			if !ok {
				logging.GetInstance().GetLogger().Error("missing [status] claims")
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"missing [status] claims",
						nil,
						nil,
						nil,
					),
				)
				return
			}
			statusCtx, ok := statusClaims.(interface{})
			if !ok {
				logging.GetInstance().GetLogger().Error("invalid [status] claims")
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericUnauthorized],
						"invalid [status] claims",
						nil,
						nil,
						nil,
					),
				)
				return
			}

			status := statusCtx.(string)
			if status == constants.AccountStatusBanned {
				ctx.AbortWithStatusJSON(
					http.StatusForbidden,
					response.NewResponse(ctx).ToResponse(
						cerrors.ErrCodeMapper[cerrors.ErrGenericPermission],
						"account has been banned",
						nil,
						nil,
						nil,
					),
				)
				return
			}

		default:
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{})
			return
		}

		ctx.Next()
	}
}
