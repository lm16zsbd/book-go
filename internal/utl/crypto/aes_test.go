package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"hash"
	"os"
	"strings"
	"testing"
)

var errPEMDecode = errors.New("failed to decode PEM block")
var errNotRSA = errors.New("not an RSA private key")

// rsaPrivateKeyPEM converts a raw base64 private key (with literal \n) into PEM.
func rsaPrivateKeyPEM(raw string) string {
	key := raw
	if !strings.Contains(key, "-----BEGIN") {
		clean := strings.NewReplacer("\\n", "", "\\r", "", " ", "", "\n", "", "\r", "").Replace(raw)
		decoded, _ := base64.StdEncoding.DecodeString(clean)
		key = string(decoded)
	}
	return strings.ReplaceAll(key, "\\n", "\n")
}

func rsaDecryptOAEP(privateKeyPEM, ciphertextB64 string, hash func() hash.Hash) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errPEMDecode
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
	}
	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return "", errNotRSA
	}
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}
	plaintext, err := rsa.DecryptOAEP(hash(), rand.Reader, rsaPriv, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func TestGoWrapKeyMatchesNodeJS(t *testing.T) {
	pubKey := os.Getenv("RSA_PUB_KEY")
	priKey := os.Getenv("RSA_PRI_KEY")
	aesKey := os.Getenv("AES_KEY")
	if pubKey == "" || priKey == "" || aesKey == "" {
		t.Skip("env not set")
	}

	goWrap, err := RSAEncrypt(pubKey, aesKey)
	if err != nil {
		t.Fatalf("RSAEncrypt failed: %v", err)
	}
	t.Logf("Go wrapKey: %s", goWrap)

	priPEM := rsaPrivateKeyPEM(priKey)

	// Must decrypt with SHA-1 (matching Node.js + frontend node-forge default)
	dec, sha1Err := rsaDecryptOAEP(priPEM, goWrap, sha1.New)
	if sha1Err != nil {
		t.Fatalf("SHA-1 decrypt FAILED: %v", sha1Err)
	}
	if dec != aesKey {
		t.Fatalf("SHA-1 decrypted='%s', expected AES_KEY='%s'", dec, aesKey)
	}
	t.Logf("SHA-1 OK: decrypted matches AES_KEY ✅")

	// SHA-256 must fail (since Go now uses SHA-1)
	_, sha256Err := rsaDecryptOAEP(priPEM, goWrap, sha256.New)
	if sha256Err != nil {
		t.Logf("SHA-256 correctly FAILED (Go uses SHA-1) ✅")
	} else {
		t.Fatal("SHA-256 should have failed but didn't — Go might still be using SHA-256!")
	}
}
