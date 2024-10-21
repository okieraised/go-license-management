package middlewares

//func JWTValidationMW() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//
//		authHeader := ctx.GetHeader(comconstants.AuthorizationHeader)
//		if authHeader == "" {
//			ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
//			return
//		}
//
//		authHdrPart := strings.Split(authHeader, " ")
//		switch len(authHdrPart) {
//		case 2:
//			if authHdrPart[0] != comconstants.AuthorizationTypeBearer {
//				ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
//				return
//			}
//
//			publicKey := ""
//
//			_, err := verifyAccessToken(authHdrPart[1], publicKey)
//			if err != nil {
//				ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
//				return
//			}
//
//		default:
//			ctx.AbortWithStatusJSON(http.StatusUnauthorized, nil)
//			return
//		}
//
//		ctx.Next()
//	}
//}
//
//func verifyAccessToken(accessToken, publicKey string) (jwt.MapClaims, error) {
//	tokenParts := strings.Split(accessToken, ".")
//	if len(tokenParts) != 3 {
//		return nil, errors.New(comerrors.ErrInvalidTokenFormat.Error())
//	}
//
//	// Parse public key
//	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
//	if err != nil {
//		return nil, err
//	}
//
//	seg, err := jwt.NewParser().DecodeSegment(tokenParts[2])
//	if err != nil {
//		return nil, err
//	}
//
//	// Verify token with public key
//	err = jwt.SigningMethodRS256.Verify(strings.Join(tokenParts[0:2], "."), seg, key)
//	if err != nil {
//		return nil, err
//	}
//
//	// Parse token
//	tok, err := jwt.Parse(accessToken, func(jwtToken *jwt.Token) (interface{}, error) {
//		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
//			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
//		}
//		return key, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	claims, ok := tok.Claims.(jwt.MapClaims)
//	if !ok || !tok.Valid {
//		return nil, comerrors.ErrCannotMapClaimsFromToken
//	}
//
//	exp, err := claims.GetExpirationTime()
//	if err != nil {
//		return nil, err
//	}
//	if exp.Unix() < time.Now().Unix() {
//		return nil, comerrors.ErrTokenHasExpired
//	}
//
//	return claims, nil
//}
