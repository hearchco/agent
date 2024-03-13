package result_test

import (
	"testing"

	"github.com/hearchco/hearchco/src/search/result"
)

type testPair struct {
	orig     string
	expected string
}

func TestFirstNchars(t *testing.T) {
	// original string, expected string
	tests := []testPair{
		{"", ""},
		{"banana death", "banana dea"},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "Lorem ipsu"},
		{"Ćao 🐹 hrčko!!", "Ćao 🐹 hrčk"},
	}

	for _, test := range tests {
		v := result.FirstNchars(test.orig, 10)
		if v != test.expected {
			t.Errorf("FirstNChars(%q) = %q, want %q", test.orig, v, test.expected)
		}
	}
}

func TestShorten(t *testing.T) {
	// original string, expected string
	tests := []testPair{
		// 0 characters -> 0 characters
		{"", ""},
		// 304 characters -> 304 characters
		{"Knowing the word count of a text can be important. For example, if an author has to write a minimum or maximum amount of words for an article, essay, report, story, book, paper, you name it. WordCounter will help to make sure its word count reaches a specific requirement or stays within a certain limit.", "Knowing the word count of a text can be important. For example, if an author has to write a minimum or maximum amount of words for an article, essay, report, story, book, paper, you name it. WordCounter will help to make sure its word count reaches a specific requirement or stays within a certain limit."},
		// 400 characters -> 400 characters
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit. You can also copy and paste it.", "Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit. You can also copy and paste it."},
		// 402 characters -> 400 characters with ... as the last 3 characters
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit them. You can also copy and paste.", "Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit them. You can also copy and p..."},
		// 445 characters -> 400 characters with ... as the last 3 characters
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa ..."},
	}

	// create test results
	var results = make([]result.Result, 0, len(tests))
	for _, test := range tests {
		v := result.Result{
			Description: test.orig,
		}
		results = append(results, v)
	}

	// shorten the descriptions
	result.Shorten(results, 400)

	// check if the descriptions are shortened as expected
	for i, test := range tests {
		v := results[i].Description
		if v != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v, len(v), test.expected, len(test.expected))
		}
	}
}
