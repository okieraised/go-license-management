package tokens

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"go-license-management/internal/utils"
)

func GenerateToken(key []byte, content interface{}) (string, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(content)
	if err != nil {
		return "", err
	}

	com, err := utils.Compress(buf.Bytes())
	if err != nil {
		return "", err
	}

	cypherText, err := utils.Encrypt(com, key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cypherText), nil
}

func DecryptToken(key []byte, token string) (interface{}, error) {
	cypherText, err := hex.DecodeString(token)
	if err != nil {

		return nil, err
	}

	com, err := utils.Decrypt(cypherText, key)
	if err != nil {
		return nil, err
	}

	decom, err := utils.Decompress(com)
	if err != nil {
		return nil, err
	}

	var original string
	dec := gob.NewDecoder(bytes.NewReader(decom))
	err = dec.Decode(&original)
	if err != nil {
		return nil, err
	}

	return original, nil
}
