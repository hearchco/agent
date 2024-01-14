package anonymize

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashToSHA256B64(orig string) string {
	// hash string with sha256
	hasher := sha256.New()
	hasher.Write([]byte(orig))
	hashedBinary := hasher.Sum(nil)

	// convert binary to string
	hashedString := base64.URLEncoding.EncodeToString(hashedBinary)

	return hashedString
}
