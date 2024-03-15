package router

import "strings"

func isWildcardOrigin(origin string) bool {
	return strings.Contains(origin, "*")
}

func UnderWildcard(origin, wildcardOrigin string) bool {
	// expects that wildcard appears only once
	// returns slice of 2 elements if * is first or last char
	// returns slice of 3 elements where * is the middle (index of 1)
	split := strings.SplitN(wildcardOrigin, "*", 2)

	if split[0] == "*" { // wildcard is the first character
		if strings.HasSuffix(origin, split[1]) {
			return true
		}
	} else if split[1] == "*" { // wildcard is the last character
		if strings.HasPrefix(origin, split[0]) {
			return true
		}
	} else {
		if strings.HasPrefix(origin, split[0]) && strings.HasSuffix(origin, split[2]) {
			return true
		}
	}

	return false
}

func CheckOrigin(frontendUrls []string) func(origin string) bool {
	return func(origin string) bool {
		for _, allowedOrigin := range frontendUrls {
			// TODO: make sure only one * can exist when loading frontend urls config
			dynamic := isWildcardOrigin(allowedOrigin)

			// allowed origin doesn't have wildcard
			if !dynamic && origin == allowedOrigin {
				return true
			}

			// allowedOrigin has wildcard
			if dynamic && UnderWildcard(origin, allowedOrigin) {
				return true
			}
		}

		// if no allowed origins match
		return false
	}
}
