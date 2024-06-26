package result

type ResultOutput any

func ConvertToOutput(results []Result, salt string) []ResultOutput {
	var output = make([]ResultOutput, 0, len(results))
	for _, r := range results {
		output = append(output, r.ConvertToOutput(salt))
	}
	return output
}

func ConvertSuggestionsToOutput(suggestions []Suggestion) []string {
	var output = make([]string, 0, len(suggestions))
	for _, s := range suggestions {
		output = append(output, s.Value())
	}
	return output
}
