package anonymize

import (
	"math/rand"
	"time"
)

func deduplicate(orig string) string {
	dedupStr := ""
	encountered := make(map[rune]bool)

	for _, char := range orig {
		if encountered[char] {
			continue
		} else {
			encountered[char] = true
			dedupStr += string(char)
		}
	}

	return dedupStr
}

func randomize(orig string) string {
	inRune := []rune(orig)

	rand.New(rand.NewSource(time.Now().Unix()))
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune)
}

func String(orig string) string {
	return randomize(deduplicate(orig))
}
