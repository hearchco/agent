package anonymize_test

import (
	"testing"

	"github.com/hearchco/hearchco/src/anonymize"
)

func TestDeduplicate(t *testing.T) {
	// original string, expected deduplicated string
	tests := []testPair{
		{"", ""},
		{"gmail", "gmail"},
		{"banana death", "ban deth"},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "Lorem ipsudlta,cngbq.UvxDhfE"},
	}

	for _, test := range tests {
		deduplicated := anonymize.Deduplicate(test.orig)
		if deduplicated != test.expected {
			t.Errorf("deduplicate(\"%v\") = \"%v\", want \"%v\"", test.orig, deduplicated, test.expected)
		}
	}
}

func TestSortString(t *testing.T) {
	// original string, sorted string
	tests := []testPair{
		{"", ""},
		{"gmail", "agilm"},
		{"banana death", " aaaabdehnnt"},
		{
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"                  ,,.Laaaaaaabccccddddddddeeeeeeeeeeeggiiiiiiiiiiilllllmmmmmmnnnnnoooooooooopppqrrrrrrsssssstttttttttuuuuuu",
		},
	}

	for _, test := range tests {
		sorted := anonymize.SortString(test.orig)

		if sorted != test.expected {
			t.Errorf("SortString(\"%v\") = \"%v\", want \"%v\"", test.orig, sorted, test.expected)
		}
	}
}

func TestShuffle(t *testing.T) {
	// original string, sorted string
	tests := []testPair{
		{"", ""},
		{"gmail", "agilm"},
		{"banana death", " aaaabdehnnt"},
		{
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"                  ,,.Laaaaaaabccccddddddddeeeeeeeeeeeggiiiiiiiiiiilllllmmmmmmnnnnnoooooooooopppqrrrrrrsssssstttttttttuuuuuu",
		},
	}

	for _, test := range tests {
		shuffled := anonymize.Shuffle(test.orig)
		shuffledSorted := anonymize.SortString(shuffled)

		if shuffledSorted != test.expected {
			t.Errorf("SortString(Shuffle(\"%v\")) = \"%v\", want \"%v\"", test.orig, shuffledSorted, test.expected)
		}
	}
}
