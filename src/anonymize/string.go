package anonymize

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

// remove duplicate characters from string
func Deduplicate(orig string) string {
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

// sort string characters lexicographically
func SortString(orig string) string {
	// Convert the string to a slice of characters
	characters := strings.Split(orig, "")

	// Sort the slice
	sort.Strings(characters)

	// Join the sorted slice back into a string
	return strings.Join(characters, "")
}

// shuffle string because deduplicate retains the order of letters
func Shuffle(orig string) string {
	inRune := []rune(orig)

	// WARNING: in year 2262, this will break
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune)
}

// anonymize string
func String(orig string) string {
	return Shuffle(Deduplicate(orig))
}

// anonymize substring of string
func Substring(orig string, ssToAnon string) string {
	anonSubstring := Shuffle(Deduplicate(ssToAnon))
	return strings.ReplaceAll(orig, ssToAnon, anonSubstring)
}
