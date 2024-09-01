package useragent

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
)

var browsers = [...]string{"chrome", "edge"}
var versions = [...]int{127, 128}

type userAgentWithHeaders struct {
	UserAgent       string
	SecCHUA         string
	SecCHUAMobile   string
	SecCHUAPlatform string
}

func userAgentStruct(browser string, version int) userAgentWithHeaders {
	if !slices.Contains(browsers[:], browser) {
		log.Panic().
			Str("browser", browser).
			Msg("Invalid browser")
		// ^PANIC - This should never happen
	}
	if !slices.Contains(versions[:], version) {
		log.Panic().
			Int("version", version).
			Msg("Invalid version")
		// ^PANIC - This should never happen
	}

	const userAgentTemplate = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.0 Safari/537.36"
	userAgent := fmt.Sprintf(userAgentTemplate, version)
	if browser == "edge" {
		userAgent = fmt.Sprintf("%s Edg/%d.0.0.0", userAgent, version)
	}

	const secCHUATemplate = `"Chromium";v="%d", "Not;A=Brand";v="24", "%s";v="%d"`
	secCHUA := fmt.Sprintf(secCHUATemplate, version, "Google Chrome", version)
	if browser == "edge" {
		secCHUA = fmt.Sprintf(secCHUATemplate, version, "Microsoft Edge", version)
	}

	return userAgentWithHeaders{
		userAgent,
		secCHUA,
		"?0",
		`"Windows"`,
	}
}

func randomUserAgentStruct() userAgentWithHeaders {
	// WARNING: Will stop working after year 2262.
	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)
	return userAgentStruct(browsers[randGen.Intn(len(browsers))], versions[randGen.Intn(len(versions))])
}

func RandomUserAgent() string {
	randomUA := randomUserAgentStruct()
	return randomUA.UserAgent
}

func RandomUserAgentWithHeaders() userAgentWithHeaders {
	return randomUserAgentStruct()
}
