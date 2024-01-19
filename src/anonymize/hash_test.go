package anonymize_test

import (
	"testing"

	"github.com/hearchco/hearchco/src/anonymize"
)

func TestHashToSHA256B64(t *testing.T) {
	// original string, expected hash (sha256 returns binary and is encoded to base64)
	tests := []testPair{
		{"", "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU="},
		{"a", "ypeBEsobvcr6wjGzmiPcTaeG7_gUfE5yuYB3ha_uSLs="},
		{"ab", "-44g_C5MPySMYMOb1lLzwTRymLuXe4tNWQO4UFViBgM="},
		{"abc", "ungWv48Bz-pBQUDeXa4iI7ADYaOWF3qctBD_YfIAFa0="},
	}

	for _, test := range tests {
		hash := anonymize.HashToSHA256B64(test.orig)
		if hash != test.expected {
			t.Errorf("HashToSHA256B64(\"%v\") = \"%v\", want \"%v\"", test.orig, hash, test.expected)
		}
	}
}
