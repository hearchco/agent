package result

func firstNchars(str string, n int) string {
	v := []rune(str)
	if n >= len(v) {
		return str
	}
	return string(v[:n])
}

func Shorten(results []Result) []Result {
	resultsShort := make([]Result, len(results))
	copy(resultsShort, results)
	for _, result := range resultsShort {
		if len(result.Description) >= 400 {
			descShort := firstNchars(result.Description, 397)
			result.Description = descShort + "..."
		}
	}
	return resultsShort
}
