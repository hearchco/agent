package main

import (
	"log"
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

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}
