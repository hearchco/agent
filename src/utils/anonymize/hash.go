package anonymize

import (
	"crypto/hmac"
	"encoding/base64"
	"crypto/sha256"
)

func CalculateHashBase64(message string) string {
	hasher := sha256.New()
	hasher.Write([]byte(message))
	hashedBinary := hasher.Sum(nil)
	hashedString := base64.URLEncoding.EncodeToString(hashedBinary)
	return hashedString
}

func CalculateMACBase64(message string, key string) string {
	hasher := hmac.New(sha256.New, []byte(key)) 
	hasher.Write([]byte(message))
	hashedBinary := hasher.Sum(nil)
	hashedString := base64.URLEncoding.EncodeToString(hashedBinary)
	return hashedString
}

func VerifyMACBase64(tag string, orig string, key string) bool {
	return tag == CalculateMACBase64(orig, key)
}
