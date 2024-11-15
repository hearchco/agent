package category

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/utils/moretime"
)

func Base64ToCategoryType(b64 string) (Category, error) {
	cj, err := Base64ToCategoryJSON(b64)
	if err != nil {
		return Category{}, fmt.Errorf("failed to convert base64 to category JSON: %w", err)
	}

	return cj.ToCategoryType()
}

func Base64ToCategoryJSON(b64 string) (CategoryJSON, error) {
	s, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		return CategoryJSON{}, fmt.Errorf("failed to decode base64: %w (%v)", err, b64)
	}

	var cj CategoryJSON
	if err := json.Unmarshal(s, &cj); err != nil {
		return CategoryJSON{}, fmt.Errorf("failed to unmarshal category JSON: %w (%v)", err, string(s))
	}

	return cj, nil
}

// Converts the category JSON into a more program friendly category type.
// Returns an error if any issues occur during the conversion.
func (cj CategoryJSON) ToCategoryType() (Category, error) {
	// Initialize the engines slices.
	engEnabled := make([]engines.Name, 0)
	engRequired := make([]engines.Name, 0)
	engRequiredByOrigin := make([]engines.Name, 0)
	engPreferred := make([]engines.Name, 0)
	engPreferredByOrigin := make([]engines.Name, 0)

	// Set the engines slices according to the provided JSON.
	for nameS, conf := range cj.Engines {
		name, err := engines.NameString(nameS)
		if err != nil {
			return Category{}, fmt.Errorf("failed converting string to engine name: %w", err)
		}

		if conf.Enabled {
			engEnabled = append(engEnabled, name)

			if conf.Required {
				engRequired = append(engRequired, name)
			} else if conf.RequiredByOrigin {
				engRequiredByOrigin = append(engRequiredByOrigin, name)
			} else if conf.Preferred {
				engPreferred = append(engPreferred, name)
			} else if conf.PreferredByOrigin {
				engPreferredByOrigin = append(engPreferredByOrigin, name)
			}
		}
	}

	// Timings config.
	timings := Timings{
		PreferredTimeout: moretime.ConvertFromFancyTime(cj.Timings.PreferredTimeout),
		HardTimeout:      moretime.ConvertFromFancyTime(cj.Timings.HardTimeout),
	}

	// Set the category config.
	return Category{
		Engines:                  engEnabled,
		RequiredEngines:          engRequired,
		RequiredByOriginEngines:  engRequiredByOrigin,
		PreferredEngines:         engPreferred,
		PreferredByOriginEngines: engPreferredByOrigin,
		Ranking:                  cj.Ranking, // Stays the same.
		Timings:                  timings,
	}, nil
}
