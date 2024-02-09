package _sedefaults

import (
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func PageFromContext(ctx *colly.Context, seName engines.Name) int {
	var pageStr string = ctx.Get("page")
	page, converr := strconv.Atoi(pageStr)
	if converr != nil {
		log.Panic().
			Err(converr).
			Str("engine", seName.String()).
			Str("page", pageStr).
			Msg("_sedefaults.PageFromContext(): failed to convert page number to int")
		// ^PANIC
	}
	return page
}
