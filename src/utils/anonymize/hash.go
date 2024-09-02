package anonymize

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hearchco/agent/src/utils/moretime"
)

// Format used for the timestamps.
const timestampFormat = time.RFC3339

// Returns the hash of the message.
func CalculateHashBase64(message string) string {
	hasher := sha256.New()
	hasher.Write([]byte(message))
	hashedBinary := hasher.Sum(nil)
	hashedString := base64.URLEncoding.EncodeToString(hashedBinary)
	return hashedString
}

// Returns the hash of the message and the timestamp used to generate it.
func CalculateHMACBase64(message, key string, t time.Time) (string, string) {
	hasher := hmac.New(sha256.New, []byte(key))
	timestamp := base64.URLEncoding.EncodeToString([]byte(t.Format(timestampFormat)))

	hasher.Write([]byte(timestamp))
	hasher.Write([]byte(message))
	hashedBinary := hasher.Sum(nil)

	hashedString := base64.URLEncoding.EncodeToString(hashedBinary)
	return hashedString, timestamp
}

// Returns whether the tag is valid for the given message, timestamp and key.
func VerifyHMACBase64(tag, orig, key, timestampB64 string) (bool, error) {
	timestamp, err := base64.URLEncoding.DecodeString(timestampB64)
	if err != nil {
		return false, fmt.Errorf("error decoding timestamp: %v", err)
	}

	t, err := time.Parse(timestampFormat, string(timestamp))
	if err != nil {
		return false, fmt.Errorf("error parsing timestamp: %v", err)
	}

	// TODO: Make duration of the timestamp configurable.
	if time.Since(t) > moretime.Day {
		return false, nil
	}

	verificator, _ := CalculateHMACBase64(orig, key, t)
	return tag == verificator, nil
}
