package search

import (
	"fmt"

	"github.com/hearchco/agent/src/search/engines/options"
)

func validateParams(query string, opts options.Options) error {
	if query == "" {
		return fmt.Errorf("query can't be empty")
	}
	if opts.Locale == "" {
		return fmt.Errorf("locale can't be empty")
	}
	if opts.Pages.Start < 0 {
		return fmt.Errorf("pages start can't be negative")
	}
	if opts.Pages.Max < 1 {
		return fmt.Errorf("pages max can't be less than 1")
	}

	return nil
}
