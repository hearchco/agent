package useragent

import (
	"math/rand"
	"time"
)

// lowercase private list of user agents
var defaultUserAgentList = [...]string{
	// Chrome 119.0.0, Windows
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	// Chrome 119.0.0, Mac OS X
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	// Edge 119.0.0, Windows
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
}

func RandomUserAgent() string {
	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)
	return defaultUserAgentList[randGen.Intn(len(defaultUserAgentList))]
}
