package category

import (
	"slices"

	"github.com/hearchco/agent/src/search/engines"
)

// Returns true if the category contains any disabled engines.
// Otherwise, returns false.
func (c Category) ContainsDisabledEngines(disabledEngines []engines.Name) bool {
	for _, eng := range disabledEngines {
		if slices.Contains(c.Engines, eng) {
			return true
		}
	}

	return false
}
