package result

import (
	"testing"
)

type testPair struct {
	orig     string
	expected string
}

func TestFirstNcharsNegative(t *testing.T) {
	tests := []testPair{
		{"", ""},
		{"banana death", "banana death"},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."},
		{"Ćao 🐹 hrčko!!", "Ćao 🐹 hrčko!!"},
	}

	for _, test := range tests {
		v := firstNchars(test.orig, -1)
		if v != test.expected {
			t.Errorf("FirstNChars(%q) = %q, want %q", test.orig, v, test.expected)
		}
	}
}

func TestFirstNcharsZero(t *testing.T) {
	tests := []testPair{
		{"", ""},
		{"banana death", ""},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", ""},
		{"Ćao 🐹 hrčko!!", ""},
	}

	for _, test := range tests {
		v := firstNchars(test.orig, 0)
		if v != test.expected {
			t.Errorf("FirstNChars(%q) = %q, want %q", test.orig, v, test.expected)
		}
	}
}

func TestFirstNchars1(t *testing.T) {
	tests := []testPair{
		{"", ""},
		{"banana death", "b"},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "L"},
		{"Ćao 🐹 hrčko!!", "Ć"},
	}

	for _, test := range tests {
		v := firstNchars(test.orig, 1)
		if v != test.expected {
			t.Errorf("FirstNChars(%q) = %q, want %q", test.orig, v, test.expected)
		}
	}
}

func TestFirstNchars10(t *testing.T) {
	tests := []testPair{
		{"", ""},
		{"banana death", "banana dea"},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "Lorem ipsu"},
		{"Ćao 🐹 hrčko!!", "Ćao 🐹 hrčk"},
	}

	for _, test := range tests {
		v := firstNchars(test.orig, 10)
		if v != test.expected {
			t.Errorf("FirstNChars(%q) = %q, want %q", test.orig, v, test.expected)
		}
	}
}

func TestShortenNegative(t *testing.T) {
	tests := []testPair{
		// 0 characters -> 0 characters (nothing changes)
		{"", ""},
		// 304 characters -> 304 characters (nothing changes)
		{"Knowing the word count of a text can be important. For example, if an author has to write a minimum or maximum amount of words for an article, essay, report, story, book, paper, you name it. WordCounter will help to make sure its word count reaches a specific requirement or stays within a certain limit.", "Knowing the word count of a text can be important. For example, if an author has to write a minimum or maximum amount of words for an article, essay, report, story, book, paper, you name it. WordCounter will help to make sure its word count reaches a specific requirement or stays within a certain limit."},
		// 400 characters -> 400 characters (nothing changes)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit. You can also copy and paste it.", "Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit. You can also copy and paste it."},
		// 402 characters -> 402 characters (nothing changes)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit them. You can also copy and paste.", "Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit them. You can also copy and paste."},
		// 445 characters -> 445 characters (nothing changes)
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."},
	}

	// Create test results.
	var results = make([]General, 0, len(tests))
	for _, test := range tests {
		v := General{
			generalJSON: generalJSON{
				Title:       test.orig,
				Description: test.orig,
			},
		}
		results = append(results, v)
	}

	// Shorten the results.
	for i := range results {
		results[i] = *results[i].Shorten(-1, -1).(*General)
	}

	// Check if the results are shortened as expected.
	for i, test := range tests {
		v := results[i]
		if v.Title() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Title(), len(v.Title()), test.expected, len(test.expected))
		}
		if v.Description() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Description(), len(v.Description()), test.expected, len(test.expected))
		}
	}
}

func TestShortenZero(t *testing.T) {
	tests := []testPair{
		// 0 characters -> 0 characters
		{"", ""},
		// 304 characters -> 0 characters (no room for suffix)
		{"Knowing the word count of a text can be important. For example, if an author has to write a minimum or maximum amount of words for an article, essay, report, story, book, paper, you name it. WordCounter will help to make sure its word count reaches a specific requirement or stays within a certain limit.", ""},
		// 400 characters -> 0 characters (no room for suffix)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit. You can also copy and paste it.", ""},
		// 402 characters -> 0 characters (no room for suffix)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit them. You can also copy and paste.", ""},
		// 445 characters -> 0 characters (no room for suffix)
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", ""},
	}

	// Create test results.
	var results = make([]General, 0, len(tests))
	for _, test := range tests {
		v := General{
			generalJSON: generalJSON{
				Title:       test.orig,
				Description: test.orig,
			},
		}
		results = append(results, v)
	}

	// Shorten the results.
	for i := range results {
		results[i] = *results[i].Shorten(0, 0).(*General)
	}

	// Check if the results are shortened as expected.
	for i, test := range tests {
		v := results[i]
		if v.Title() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Title(), len(v.Title()), test.expected, len(test.expected))
		}
		if v.Description() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Description(), len(v.Description()), test.expected, len(test.expected))
		}
	}
}

func TestShorten1(t *testing.T) {
	tests := []testPair{
		// 0 characters -> 0 characters
		{"", ""},
		// 304 characters -> 1 character (no room for suffix)
		{"Knowing the word count of a text can be important. For example, if an author has to write a minimum or maximum amount of words for an article, essay, report, story, book, paper, you name it. WordCounter will help to make sure its word count reaches a specific requirement or stays within a certain limit.", "K"},
		// 400 characters -> 1 character (no room for suffix)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit. You can also copy and paste it.", "A"},
		// 402 characters -> 1 character (no room for suffix)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit them. You can also copy and paste.", "A"},
		// 445 characters -> 1 character (no room for suffix)
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "L"},
	}

	// Create test results.
	var results = make([]General, 0, len(tests))
	for _, test := range tests {
		v := General{
			generalJSON: generalJSON{
				Title:       test.orig,
				Description: test.orig,
			},
		}
		results = append(results, v)
	}

	// Shorten the results.
	for i := range results {
		results[i] = *results[i].Shorten(1, 1).(*General)
	}

	// Check if the results are shortened as expected.
	for i, test := range tests {
		v := results[i]
		if v.Title() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Title(), len(v.Title()), test.expected, len(test.expected))
		}
		if v.Description() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Description(), len(v.Description()), test.expected, len(test.expected))
		}
	}
}

func TestShorten2(t *testing.T) {
	tests := []testPair{
		// 0 characters -> 0 characters
		{"", ""},
		// 304 characters -> 2 characters (no room for suffix)
		{"Knowing the word count of a text can be important. For example, if an author has to write a minimum or maximum amount of words for an article, essay, report, story, book, paper, you name it. WordCounter will help to make sure its word count reaches a specific requirement or stays within a certain limit.", "Kn"},
		// 400 characters -> 2 characters (no room for suffix)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit. You can also copy and paste it.", "Ap"},
		// 402 characters -> 2 characters (no room for suffix)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit them. You can also copy and paste.", "Ap"},
		// 445 characters -> 2 characters (no room for suffix)
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "Lo"},
	}

	// Create test results.
	var results = make([]General, 0, len(tests))
	for _, test := range tests {
		v := General{
			generalJSON: generalJSON{
				Title:       test.orig,
				Description: test.orig,
			},
		}
		results = append(results, v)
	}

	// Shorten the results.
	for i := range results {
		results[i] = *results[i].Shorten(2, 2).(*General)
	}

	// Check if the results are shortened as expected.
	for i, test := range tests {
		v := results[i]
		if v.Title() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Title(), len(v.Title()), test.expected, len(test.expected))
		}
		if v.Description() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Description(), len(v.Description()), test.expected, len(test.expected))
		}
	}
}

func TestShorten3(t *testing.T) {
	tests := []testPair{
		// 0 characters -> 0 characters
		{"", ""},
		// 304 characters -> 3 characters (no room for suffix)
		{"Knowing the word count of a text can be important. For example, if an author has to write a minimum or maximum amount of words for an article, essay, report, story, book, paper, you name it. WordCounter will help to make sure its word count reaches a specific requirement or stays within a certain limit.", "Kno"},
		// 400 characters -> 3 characters (no room for suffix)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit. You can also copy and paste it.", "Apa"},
		// 402 characters -> 3 characters (no room for suffix)
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit them. You can also copy and paste.", "Apa"},
		// 445 characters -> 3 characters (no room for suffix)
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "Lor"},
	}

	// Create test results.
	var results = make([]General, 0, len(tests))
	for _, test := range tests {
		v := General{
			generalJSON: generalJSON{
				Title:       test.orig,
				Description: test.orig,
			},
		}
		results = append(results, v)
	}

	// Shorten the results.
	for i := range results {
		results[i] = *results[i].Shorten(3, 3).(*General)
	}

	// Check if the results are shortened as expected.
	for i, test := range tests {
		v := results[i]
		if v.Title() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Title(), len(v.Title()), test.expected, len(test.expected))
		}
		if v.Description() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Description(), len(v.Description()), test.expected, len(test.expected))
		}
	}
}

func TestShorten4(t *testing.T) {
	tests := []testPair{
		// 0 characters -> 0 characters
		{"", ""},
		// 304 characters -> 4 characters with ... as the last 3 characters
		{"Knowing the word count of a text can be important. For example, if an author has to write a minimum or maximum amount of words for an article, essay, report, story, book, paper, you name it. WordCounter will help to make sure its word count reaches a specific requirement or stays within a certain limit.", "K..."},
		// 400 characters -> 4 characters with ... as the last 3 characters
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit. You can also copy and paste it.", "A..."},
		// 402 characters -> 4 characters with ... as the last 3 characters
		{"Apart from counting words and characters, our online editor can help you to improve word choice and writing style, and, optionally, help you to detect grammar mistakes and plagiarism. To check word count simply place your cursor into the text box above and start typing. You'll see the number of characters and words increase or decrease as you type, delete, and edit them. You can also copy and paste.", "A..."},
		// 445 characters -> 4 characters with ... as the last 3 characters
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", "L..."},
	}

	// Create test results.
	var results = make([]General, 0, len(tests))
	for _, test := range tests {
		v := General{
			generalJSON: generalJSON{
				Title:       test.orig,
				Description: test.orig,
			},
		}
		results = append(results, v)
	}

	// Shorten the results.
	for i := range results {
		results[i] = *results[i].Shorten(4, 4).(*General)
	}

	// Check if the results are shortened as expected.
	for i, test := range tests {
		v := results[i]
		if v.Title() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Title(), len(v.Title()), test.expected, len(test.expected))
		}
		if v.Description() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Description(), len(v.Description()), test.expected, len(test.expected))
		}
	}
}

func TestShorten400(t *testing.T) {
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

	// Create test results.
	var results = make([]General, 0, len(tests))
	for _, test := range tests {
		v := General{
			generalJSON: generalJSON{
				Title:       test.orig,
				Description: test.orig,
			},
		}
		results = append(results, v)
	}

	// Shorten the results.
	for i := range results {
		results[i] = *results[i].Shorten(400, 400).(*General)
	}

	// Check if the results are shortened as expected.
	for i, test := range tests {
		v := results[i]
		if v.Title() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Title(), len(v.Title()), test.expected, len(test.expected))
		}
		if v.Description() != test.expected {
			t.Errorf("\n\tShorten(%q)\n\tlen = %v\n\n\tGot: %q\n\tlen = %v\n\n\tWant: %q\n\tlen = %v", test.orig, len(test.orig), v.Description(), len(v.Description()), test.expected, len(test.expected))
		}
	}
}
