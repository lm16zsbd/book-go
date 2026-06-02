package crypto

import (
	"os"
	"testing"
)

func TestRSAEncrypt(t *testing.T) {
	pubKey := os.Getenv("RSA_PUB_KEY")
	aesKey := os.Getenv("AES_KEY")
	if pubKey == "" || aesKey == "" {
		t.Skip("RSA_PUB_KEY or AES_KEY not set")
	}

	result1, err := RSAEncrypt(pubKey, aesKey)
	if err != nil {
		t.Fatalf("RSAEncrypt failed: %v", err)
	}
	t.Logf("wrapKey (run 1): %s", result1)

	result2, err := RSAEncrypt(pubKey, aesKey)
	if err != nil {
		t.Fatalf("RSAEncrypt failed: %v", err)
	}

	if result1 != result2 {
		t.Fatal("wrapKey is NOT deterministic — same inputs produced different outputs")
	}
	t.Log("wrapKey IS deterministic (same input = same output)")
}
