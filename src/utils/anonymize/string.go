package anonymize

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

// Anonymize string
func String(orig string) string {
	return shuffle(deduplicate(orig))
}

// Anonymize substring of a string
func Substring(orig string, ssToAnon string) string {
	return strings.ReplaceAll(orig, ssToAnon, String(ssToAnon))
}

// Remove duplicate characters from string.
func deduplicate(orig string) string {
	dedupStr := ""
	encountered := make(map[rune]bool)

	for _, char := range orig {
		if !encountered[char] {
			encountered[char] = true
			dedupStr += string(char)
		}
	}

	return dedupStr
}

// Shuffle string because deduplicate retains the order of letters.
func shuffle(orig string) string {
	inRune := []rune(orig)

	// WARNING: In year 2262, this will break.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune)
}

// Sort string characters lexicographically.
func sortString(orig string) string {
	// Convert the string to a slice of characters.
	characters := strings.Split(orig, "")
	sort.Strings(characters)
	return strings.Join(characters, "")
}
