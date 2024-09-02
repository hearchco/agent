package anonymize

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func CalculateHashBase64(message string) string {
	hasher := sha256.New()
	hasher.Write([]byte(message))
	hashedBinary := hasher.Sum(nil)
	hashedString := base64.URLEncoding.EncodeToString(hashedBinary)
	return hashedString
}

func CalculateHMACBase64(message string, key string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(message))
	hashedBinary := hasher.Sum(nil)
	hashedString := base64.URLEncoding.EncodeToString(hashedBinary)
	return hashedString
}

func VerifyHMACBase64(tag string, orig string, key string) bool {
	return tag == CalculateHMACBase64(orig, key)
}
