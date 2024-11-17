package searchtype

import (
	"fmt"
)

type Name string

const (
	WEB         Name = "web"
	IMAGES      Name = "images"
	SUGGESTIONS Name = "suggestions"
)

func (st Name) String() string {
	return string(st)
}

// Converts a string to a search type name if it exists.
// Otherwise returns an error.
func FromString(st string) (Name, error) {
	switch st {
	case WEB.String():
		return WEB, nil
	case IMAGES.String():
		return IMAGES, nil
	case SUGGESTIONS.String():
		return SUGGESTIONS, nil
	default:
		return "", fmt.Errorf("search type %q is not defined", st)
	}
}
