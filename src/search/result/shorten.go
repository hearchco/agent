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
	for i := range resultsShort {
		if len(resultsShort[i].Description) >= 400 {
			descShort := firstNchars(resultsShort[i].Description, 397)
			resultsShort[i].Description = descShort + "..."
		}
	}
	return resultsShort
}
