package morestrings

import (
	"strings"
)

// JoinNonEmpty concatenates the non empty elements of its first argument to create a single string.
// The beg string is placen at the beginning, unless there are no elements.
// The separator string sep is placed between elements in the resulting string.
func JoinNonEmpty(beg, sep string, elems ...string) string {
	var nonEmptyElems = make([]string, 0, len(elems))
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
