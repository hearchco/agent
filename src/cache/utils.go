package cache

import "fmt"

func combineIntoKey(s ...any) string {
	var key string
	for i, v := range s {
		if i == 0 {
			key = fmt.Sprintf("%v", v)
		} else {
			key = fmt.Sprintf("%v_%v", key, v)
		}
	}
	return key
}
