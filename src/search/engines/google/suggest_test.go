package google

import (
	"context"
	"testing"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/engines/_engines_test"
)

func TestSuggest(t *testing.T) {
	se := New()
	se.InitSuggester(context.Background(), config.CategoryTimings{})
	_engines_test.CheckSuggest(t, se, "test")
}
