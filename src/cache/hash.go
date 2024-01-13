package cache

import "crypto/sha256"

func HashString(s string) string {
	hashedString := sha256.Sum256([]byte(s))
	return string(hashedString[:])
}
