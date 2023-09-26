package main

import (
	"os"
	"strings"
)

func (s Values) Delete(v Value) []Value {
	index := s.IndexOf(v)
	ret := make([]Value, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func (s Values) IndexOf(x Value) int {
	for i, v := range s {
		if x == v {
			return i
		}
	}
	return -1
}

func validConst(v Value) bool {
	lowerName := strings.ToLower(v.name)
	return lowerName != "undefined" && exists(lowerName)
}

// exists returns whether the given file or directory exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
