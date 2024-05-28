package cache

import "fmt"

func combineIntoKey(s ...string) string {
	var key string
	for i, v := range s {
		if i == 0 {
			key = v
		} else {
			key = fmt.Sprintf("%v_%v", key, v)
		}
	}
	return key
}
