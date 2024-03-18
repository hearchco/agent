package router

// import (
// 	"strings"
// )

// func isWildcardOrigin(origin string) bool {
// 	return strings.Contains(origin, "*")
// }

// func UnderWildcard(origin, wildcardOrigin string) bool {
// 	if wildcardOrigin == "*" {
// 		return true
// 	}

// 	if strings.Count(wildcardOrigin, "*") != 1 {
// 		return false
// 	}

// 	if origin == "" || wildcardOrigin == "" {
// 		return false
// 	}

// 	// expects that wildcard appears only once
// 	// returns slice of 2 elements
// 	// first element is empty string if wildcard is first char
// 	// second element is empty string if wildcard is last char
// 	// first element is substring before wildcard and second is substring after
// 	split := strings.SplitN(wildcardOrigin, "*", 2)

// 	if split[0] == "" { // wildcard is the first character
// 		if strings.HasSuffix(origin, split[1]) {
// 			return true
// 		}
// 	} else if split[1] == "" { // wildcard is the last character
// 		if strings.HasPrefix(origin, split[0]) {
// 			return true
// 		}
// 	} else {
// 		if strings.HasPrefix(origin, split[0]) && strings.HasSuffix(origin, split[1]) {
// 			return true
// 		}
// 	}

// 	return false
// }

// func CheckOrigin(frontendUrls []string) func(origin string) bool {
// 	return func(origin string) bool {
// 		for _, allowedOrigin := range frontendUrls {
// 			// TODO: make sure only one * can exist when loading frontend urls config
// 			dynamic := isWildcardOrigin(allowedOrigin)

// 			// allowed origin doesn't have wildcard
// 			if !dynamic && origin == allowedOrigin {
// 				return true
// 			}

// 			// allowedOrigin has wildcard
// 			if dynamic && UnderWildcard(origin, allowedOrigin) {
// 				return true
// 			}
// 		}

// 		// if no allowed origins match
// 		return false
// 	}
// }
