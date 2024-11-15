package duckduckgo

import (
	"context"
	"testing"

	"github.com/hearchco/agent/src/search/engines/_engines_test"
)

func TestSuggest(t *testing.T) {
	se := New()
	se.InitSuggester(context.Background())
	_engines_test.CheckSuggest(t, se, "test")
}
