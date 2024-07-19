package swisscows

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/hearchco/agent/src/utils/anonymize"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"

// Returns nonce and signature.
func generateAuth(params string) (string, string, error) {
	paramsWOP := strings.ReplaceAll(params, "+", " ")
	nonce := generateNonce(32)

	auth, err := generateSignature(paramsWOP, nonce)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate auth (nonce and signature): %w", err)
	}

	return nonce, auth, nil
}

func generateNonce(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	nonce := ""
	for range length {
		randInd := r.Intn(length)
		nonce += string(alphabet[randInd])
	}

	return nonce
}

func generateSignature(params string, nonce string) (string, error) {
	rot13Nonce := rot13Switch(nonce)
	data := "/web/search" + params + rot13Nonce
	encData := anonymize.CalculateHashBase64(data)
	encData = strings.ReplaceAll(encData, "=", "")
	encData = strings.ReplaceAll(encData, "+", "-")
	encData = strings.ReplaceAll(encData, "/", "_")

	return encData, nil
}

func rot13Switch(str string) string {
	return switchCapitalization(rot13(str))
}

// Performs rot13 and switches capitalization of each character.
func rot13(str string) string {
	result := ""

	for i := range len(str) {
		result += string(rot13Byte(str[i]))
	}

	return result
}

func rot13Byte(b byte) byte {
	var a, z byte

	switch {
	case 'a' <= b && b <= 'z':
		a, z = 'a', 'z'
	case 'A' <= b && b <= 'Z':
		a, z = 'A', 'Z'
	default:
		return b
	}

	return (b-a+13)%(z-a+1) + a
}

func switchCapitalization(str string) string {
	res := ""

	for i := range len(str) {
		if unicode.IsUpper(rune(str[i])) {
			res += string(unicode.ToLower(rune(str[i])))
		} else {
			res += string(unicode.ToUpper(rune(str[i])))
		}
	}

	return res
}
