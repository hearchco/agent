package router_test

import (
	"testing"

	"github.com/hearchco/hearchco/src/router"
)

type Tests struct {
	Origin         string
	WildcardOrigin string
	Expected       bool
}

type E2ETests struct {
	Origin         string
	AllowedOrigins []string
	Expected       bool
}

func TestUnderWildcard(t *testing.T) {
	tests := []Tests{
		// good inputs
		{"https://hearch.co", "*", true},
		{"https://preview.hearch.co", "*hearch.co", true},
		{"https://preview.hearch.co", "https://preview.*", true},
		{"https://preview.hearch.co", "https://*.hearch.co", true},
		{"https://example.org", "*", true},
		{"https://preview.example.org", "*hearch.co", false},
		{"https://preview.example.org", "https://preview.*", true},
		{"https://preview.example.org", "https://*.hearch.co", false},
		{"https://hearch.co", "*", true},
		{"https://preview.hearch.co", "*example.org", false},
		{"https://preview.hearch.co", "https://preview.*", true},
		{"https://preview.hearch.co", "https://*.example.org", false},
		{"https://staging.hearch.co", "https://preview.*", false},
		// bad inputs
		{"https://hearch.co", "**", false},
		{"https://hearch.co", "***", false},
		{"https://hearch.co", "**.example.org", false},
		{"https://hearch.co", "https://example.**", false},
		{"https://hearch.co", "*.example.*", false},
		// very bad inputs
		{"", "", false},
		{"https://hearch.co", "", false},
		{"", "https://hearch.co", false},
		// real use cases
		{"https://hearch.co", "https://hearch.co", false}, // since we don't accept OriginWildcard without wildcard
		{"https://hearch.co", "https://*hearch.co", true},
		{"https://feat-image-search.hearch.co", "https://*hearch.co", true},
		{"https://feat-image-search.hearch.co", "https://*.hearch.co", true},
		{"http://localhost:5173", "http://localhost*", true},
	}

	for _, test := range tests {
		ok := router.UnderWildcard(test.Origin, test.WildcardOrigin)
		if ok != test.Expected {
			t.Errorf("UnderWildcard(%q, %q) = \"%v\", want \"%v\"", test.Origin, test.WildcardOrigin, ok, test.Expected)
		}
	}
}

func TestCheckOrigin(t *testing.T) {
	allowedOrigins := []string{
		"http://localhost:5173",
		"https://*hearch.co",
		"https://*hearchco-frontend.pages.dev",
	}

	tests := []E2ETests{
		{"http://localhost:5173", allowedOrigins, true},
		{"http://localhost:8080", allowedOrigins, false},
		{"https://hearch.co", allowedOrigins, true},
		{"https://feat-image-search.hearch.co", allowedOrigins, true},
		{"https://example.org", allowedOrigins, false},
	}

	for _, test := range tests {
		fco := router.CheckOrigin(test.AllowedOrigins)
		ok := fco(test.Origin)
		if ok != test.Expected {
			t.Errorf("CheckOrigin(%q) for %q = \"%v\", want \"%v\"", test.AllowedOrigins, test.Origin, ok, test.Expected)
		}
	}
}
