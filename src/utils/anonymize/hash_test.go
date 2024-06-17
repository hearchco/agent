package anonymize

import (
	"testing"
)

func TestHashToSHA256B64(t *testing.T) {
	// original string, expected hash (sha256 returns binary and is encoded to base64)
	tests := []testPair{
		{"", "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU="},
		{"banana death", "e8kN64XJ4Icr6Tl9VYrBRj50UJCPlyillODm3vVNk2g="},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "LYwvbZeMohcStfbeNsnTH6jpak-l2P-LAYjfuefBcbs="},
		{"Ćao hrčko!! 🐹", "_Y3KWzrx2UkeTp8b--48L6OFgv51JWPlZArjoFOrmbw="},
	}

	for _, test := range tests {
		hash := HashToSHA256B64(test.orig)
		if hash != test.expected {
			t.Errorf("HashToSHA256B64(%q) = %q, want %q", test.orig, hash, test.expected)
		}
	}
}
