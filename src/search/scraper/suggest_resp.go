package scraper

import (
	"encoding/json"
	"fmt"
)

// Converts a opensearch.xml compatible suggestions API JSON to a slice of suggestions.
func SuggestRespToSuggestions(data []byte) ([]string, error) {
	// Define a structure that matches the JSON structure.
	var resp []any

	// Unmarshal the JSON data.
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Check the structure and extract the slice of strings.
	if len(resp) < 2 {
		return nil, fmt.Errorf("unexpected JSON structure")
	}

	// Assert the second element is a slice.
	strSlice, ok := resp[1].([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected type for second element")
	}

	// Error if no suggestions returned.
	if len(strSlice) == 0 {
		return nil, fmt.Errorf("empty suggestions")
	}

	// Convert to slice of strings.
	suggs := make([]string, 0, len(strSlice))
	for _, item := range strSlice {
		if sug, ok := item.(string); !ok {
			return nil, fmt.Errorf("unexpected type in string slice")
		} else {
			suggs = append(suggs, sug)
		}

	}

	return suggs, nil
}
