package license_key

//// NewJWTLicenseKeyWithRSA2048 generates new jwt license key using RSA2048 algorithm
//func NewJWTLicenseKeyWithRSA2048(signingKey string, data any) (string, error) {
//	bData, err := json.Marshal(data)
//	if err != nil {
//		return "", err
//	}
//
//	// Decode the private key string
//	block, _ := pem.Decode([]byte(signingKey))
//
//	if block == nil || block.Type != RSAPrivateKeyStr {
//		return "", errors.New("failed to decode PEM block containing private key")
//	}
//
//	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		return "", err
//	}
//
//	claims := jwt.MapClaims{
//		"sub":        "user_id_1234",
//		"exp":        time.Now().Add(time.Hour * 24).Unix(), // License valid for 24 hours
//		"iat":        time.Now().Unix(),
//		"scope":      "basic_license",
//		"license_id": "license_123456",
//	}
//
//	// Create a new JWT token with claims and sign it with RS256
//	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
//	tokenString, err := token.SignedString(privateKey)
//	if err != nil {
//		return "", err
//	}
//
//	return "", nil
//}
