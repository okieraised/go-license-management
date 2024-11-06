package license_key

//// NewJWTRS256KeyPair generates the private key and the public key pair using RS256 algorithm
//func NewJWTRS256KeyPair() (string, string, error) {
//	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
//	if err != nil {
//		return "", "", err
//	}
//
//	// Encode the private key to PEM format (PKCS1)
//	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
//		Type:  RSAPrivateKeyStr,
//		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
//	})
//
//	// Encode the public key to PEM format (PKCS1)
//	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
//		Type:  RSAPublicKeyStr,
//		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
//	})
//
//	return string(privateKeyPEM), string(publicKeyPEM), nil
//}
