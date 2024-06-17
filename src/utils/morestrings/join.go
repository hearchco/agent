package morestrings

import (
	"strings"
)

// JoinNonEmpty concatenates the non empty elements of its first argument to create a single string. The separator
// string sep is placed between elements in the resulting string.
func JoinNonEmpty(elems []string, beg string, sep string) string {
	nonEmptyElems := []string{}
	for _, elem := range elems {
		if elem != "" {
			nonEmptyElems = append(nonEmptyElems, elem)
		}
	}

	if len(nonEmptyElems) == 0 {
		return ""
	} else if len(nonEmptyElems) == 1 {
		return beg + nonEmptyElems[0]
	} else {
		return beg + strings.Join(nonEmptyElems, sep)
	}
}
