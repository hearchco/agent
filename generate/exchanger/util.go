package main

import (
	"log"
	"os"
	"strings"
)

func validConst(v Value) bool {
	lowerName := strings.ToLower(v.name)
	return lowerName != "undefined" && isDirectory(lowerName)
}

// isDirectory reports whether the named file is a directory.
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isDirectoryFatal(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
		// ^FATAL
	}
	return info.IsDir()
}
