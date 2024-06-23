package category

import (
	"fmt"
)

type Name string

const (
	UNDEFINED   Name = "undefined"
	SUGGESTIONS Name = "suggestions"
	GENERAL     Name = "general"
	IMAGES      Name = "images"
	SCIENCE     Name = "science"
	THOROUGH    Name = "thorough"
)

func (cat Name) String() string {
	return string(cat)
}

// Converts a string to a category name if it exists.
// If the string is empty, then GENERAL is returned.
// Otherwise returns UNDEFINED.
func FromString(cat string) (Name, error) {
	switch cat {
	case "", GENERAL.String():
		return GENERAL, nil
	case IMAGES.String():
		return IMAGES, nil
	case SCIENCE.String():
		return SCIENCE, nil
	case THOROUGH.String():
		return THOROUGH, nil
	case SUGGESTIONS.String():
		return UNDEFINED, fmt.Errorf("category %q is not allowed", cat)
	default:
		return UNDEFINED, fmt.Errorf("category %q is not defined", cat)
	}
}
