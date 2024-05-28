package category

import "fmt"

var catMap = map[string]Name{
	"general":  GENERAL,
	"images":   IMAGES,
	"science":  SCIENCE,
	"sci":      SCIENCE,
	"thorough": THOROUGH,
	"slow":     THOROUGH,
}

// converts a string to a category name if it exists
// if the string is empty, then GENERAL is returned
// otherwise returns UNDEFINED
func FromString(cat string) (Name, error) {
	if cat == "" {
		return GENERAL, nil
	}

	catName, ok := catMap[cat]
	if !ok {
		return UNDEFINED, fmt.Errorf("category %q is not defined", cat)
	}

	return catName, nil
}
