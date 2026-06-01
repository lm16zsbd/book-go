package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"strings"
)

func AESDecrypt(hexKey string, data []byte) ([]byte, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	iv, err := hex.DecodeString(DefaultIV())
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	return pkcs7Unpad(data), nil
}

func AESEncrypt(hexKey string, data []byte) ([]byte, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	iv, err := hex.DecodeString(DefaultIV())
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padded := pkcs7Pad(data, aes.BlockSize)
	encrypted := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, padded)
	return encrypted, nil
}

func RSAEncrypt(publicKeyPEM, payload string) (string, error) {
	key := strings.ReplaceAll(publicKeyPEM, "\\n", "\n")
	if !strings.Contains(key, "-----BEGIN") {
		decoded, err := base64.StdEncoding.DecodeString(publicKeyPEM)
		if err == nil {
			str := string(decoded)
			if strings.Contains(str, "-----BEGIN") {
				key = str
			}
		}
	}

	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return "", nil
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", nil
	}

	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, []byte(payload), nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DefaultIV() string {
	return "00000000000000000000000000000000"
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	padding := int(data[length-1])
	if padding > length {
		return data
	}
	return data[:length-padding]
}
