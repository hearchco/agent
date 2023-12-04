package swisscows

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func generateNonce(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	const alphabet string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	var nonce string = ""
	for i := 0; i < length; i++ {
		randInd := r.Intn(length)
		nonce += string(alphabet[randInd])
	}

	return nonce
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
	var res string = ""
	for i := 0; i < len(str); i++ {
		if unicode.IsUpper(rune(str[i])) {
			res += string(unicode.ToLower(rune(str[i])))
		} else {
			res += string(unicode.ToUpper(rune(str[i])))
		}
	}
	return res
}

// performs rot13 and also switches capitalization of each character
func rot13(str string) string {
	var result string = ""
	for i := 0; i < len(str); i++ {
		result += string(rot13Byte(str[i]))
	}
	return result
}

func rot13Switch(str string) string {
	return switchCapitalization(rot13(str))
}

func generateSignature(params string, nonce string) (string, error) {
	var rot13Nonce string = rot13Switch(nonce)
	var data string = "/web/search" + params + rot13Nonce

	var encData string = hashToSHA256B64(data)
	encData = strings.ReplaceAll(encData, "=", "")
	encData = strings.ReplaceAll(encData, "+", "-")
	encData = strings.ReplaceAll(encData, "/", "_")

	//log.Debug().Msgf("Final: %v", encData)

	return string(encData), nil
}

// returns nonce, signature
func generateAuth(params string) (string, string, error) {
	params = strings.ReplaceAll(params, "+", " ")

	nonce := generateNonce(32)
	auth, err := generateSignature(params, nonce)
	if err != nil {
		return "", "", fmt.Errorf("generateAuth(): %w", err)
	}

	return nonce, auth, nil
}
