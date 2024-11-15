package category

import (
	"slices"

	"github.com/hearchco/agent/src/search/engines"
)

// Remove the specified engines from the Category.
// Passed as pointer to modify the original.
func (c *Category) DisableEngines(disabledEngines []engines.Name) {
	c.Engines = slices.DeleteFunc(c.Engines, func(e engines.Name) bool {
		return slices.Contains(disabledEngines, e)
	})
	c.RequiredEngines = slices.DeleteFunc(c.RequiredEngines, func(e engines.Name) bool {
		return slices.Contains(disabledEngines, e)
	})
	c.RequiredByOriginEngines = slices.DeleteFunc(c.RequiredByOriginEngines, func(e engines.Name) bool {
		return slices.Contains(disabledEngines, e)
	})
	c.PreferredEngines = slices.DeleteFunc(c.PreferredEngines, func(e engines.Name) bool {
		return slices.Contains(disabledEngines, e)
	})
	c.PreferredByOriginEngines = slices.DeleteFunc(c.PreferredByOriginEngines, func(e engines.Name) bool {
		return slices.Contains(disabledEngines, e)
	})
}
