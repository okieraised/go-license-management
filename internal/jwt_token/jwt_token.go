package jwt_token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(userEmail, username, userRole, tenantName, privateKey string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   username,
		"iss":   tenantName,
		"aud":   userRole,
		"email": userEmail,
		"role":  userRole,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := claims.SignedString([]byte(privateKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//func ValidateJWT(tokenString string) (*jwt.Token, error) {
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		return secretKey, nil
//	})
//
//	// Check for verification errors
//	if err != nil {
//		return nil, err
//	}
//
//	// Check if the token is valid
//	if !token.Valid {
//		return nil, fmt.Errorf("invalid token")
//	}
//
//	// Return the verified token
//	return token, nil
//
//}
