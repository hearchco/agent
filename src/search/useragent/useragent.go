package useragent

import (
	"math/rand"
	"time"
)

// lowercase private list of user agents
var defaultUserAgentList = [...]string{
	// Chrome: Windows, MacOS, Linux, Android
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
	// Edge: Windows
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.0.0",
	// Firefox: Windows
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0",
}

func RandomUserAgent() string {
	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)
	return defaultUserAgentList[randGen.Intn(len(defaultUserAgentList))]
}
