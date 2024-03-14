package result

func FirstNchars(str string, n int) string {
	v := []rune(str)
	if n < 0 || n >= len(v) {
		return str
	}
	return string(v[:n])
}

// modifies the passed slice of results,
// changes the description of the results to be at most N characters long
func Shorten(results []Result, n int) {
	suffix := "..."
	if n < 0 {
		return
	} else if n-len(suffix) <= 0 {
		suffix = "" // no room for suffix
	}

	// can't use _, result := range short because we need to modify the elements in slice
	for i := range results {
		result := &results[i]
		if len(result.Description) > n {
			descShort := FirstNchars(result.Description, n-len(suffix))
			result.Description = descShort + suffix
		}
	}
}
