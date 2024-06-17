package scraper

import (
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
)

func (e EngineBase) PageFromContext(ctx *colly.Context) int {
	var pageStr string = ctx.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Panic().
			Caller().
			Err(err).
			Str("engine", e.Name.String()).
			Str("page", pageStr).
			Msg("Failed to convert page number to int")
		// ^PANIC
	}
	return page
}
