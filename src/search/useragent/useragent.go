package useragent

import (
	"math/rand"
	"time"
)

type userAgentWithHeaders struct {
	UserAgent       string
	SecCHUA         string
	SecCHUAMobile   string
	SecCHUAPlatform string
}

// UserAgents used when making requests and their corresponding Sec-Ch-Ua headers.
var userAgentArray = [...]userAgentWithHeaders{
	// Chrome 127.0.0, Mac OS X
	{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.3",
		`"Google Chrome";v="127", "Chromium";v="127", "Not=A?Brand";v="24"`,
		"?0",
		`"Macintosh"`,
	},
	// Chrome 128.0.0, Windows
	{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.3",
		`"Google Chrome";v="128", "Chromium";v="128", "Not=A?Brand";v="24"`,
		"?0",
		`"Windows"`,
	},
	// Chrome 127.0.0, Windows
	{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.3",
		`"Google Chrome";v="127", "Chromium";v="127", "Not=A?Brand";v="24"`,
		"?0",
		`"Windows"`,
	},
	// Edge 128.0.0, Windows
	{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.",
		`"Microsoft Edge";v="128", "Edg";v="128"`,
		"?0",
		`"Windows"`,
	},
}

func randomUserAgentStruct() userAgentWithHeaders {
	// WARNING: Will stop working after year 2262.
	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)
	return userAgentArray[randGen.Intn(len(userAgentArray))]
}

func RandomUserAgent() string {
	randomUA := randomUserAgentStruct()
	return randomUA.UserAgent
}

func RandomUserAgentWithHeaders() userAgentWithHeaders {
	return randomUserAgentStruct()
}
