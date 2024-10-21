package cache

func combineIntoKey(elem ...string) string {
	key := ""
	for _, e := range elem {
		key += e
	}
	return key
}
