package rand

import (
	"crypto/rand"
	"encoding/hex"
)

func Token(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
