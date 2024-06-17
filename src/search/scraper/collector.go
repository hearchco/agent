package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/useragent"
)

func (e *EngineBase) initCollector(ctx context.Context) {
	// Get a random user agent with it's Sec-CH-UA headers.
	ua := useragent.RandomUserAgentWithHeaders()

	// Initialize the collector.
	e.collector = colly.NewCollector(
		colly.StdlibContext(ctx),
		colly.Async(),
		colly.MaxDepth(1),
		colly.IgnoreRobotsTxt(),
		colly.UserAgent(ua.UserAgent),
		colly.Headers(map[string]string{
			"Accept":             "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"Accept-Encoding":    "gzip, deflate, br",
			"Accept-Language":    "en-US,en;q=0.9",
			"Sec-Ch-Ua":          ua.SecCHUA,
			"Sec-Ch-Ua-Mobile":   ua.SecCHUAMobile,
			"Sec-Ch-Ua-Platform": ua.SecCHUAPlatform,
			"Sec-Fetch-Dest":     "document",
			"Sec-Fetch-Mode":     "navigate",
			"Sec-Fetch-Site":     "none",
		}),
	)
}

func (e *EngineBase) initLimitRule(timings config.CategoryTimings) {
	limitRule := colly.LimitRule{
		DomainGlob:  "*",
		Delay:       timings.Delay,
		RandomDelay: timings.RandomDelay,
		Parallelism: timings.Parallelism,
	}
	if err := e.collector.Limit(&limitRule); err != nil {
		log.Panic().
			Caller().
			Err(err).
			Str("limitRule", fmt.Sprintf("%v", limitRule)).
			Msg("Failed adding new limit rule")
		// ^PANIC
	}
}

func (e *EngineBase) initCollectorOnRequest(ctx context.Context) {
	e.collector.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			if IsTimeoutError(err) {
				log.Trace().
					Caller().
					Err(err).
					Str("engine", e.Name.String()).
					Msg("Context timeout error")
			} else {
				log.Error().
					Caller().
					Err(err).
					Str("engine", e.Name.String()).
					Msg("Context error")
			}
			r.Abort()
			return
		}
	})
}

func (e *EngineBase) initCollectorOnResponse() {
	e.collector.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Headers.Get("Content-Encoding"), "br") {
			reader := brotli.NewReader(bytes.NewReader(r.Body))

			body, err := io.ReadAll(reader)
			if err != nil {
				log.Error().
					Caller().
					Err(err).
					Str("engine", e.Name.String()).
					Msg("Failed to decode brotli response")
				return
			}

			r.Body = body
		}
	})
}

func (e *EngineBase) initCollectorOnError() {
	e.collector.OnError(func(r *colly.Response, err error) {
		if IsTimeoutError(err) {
			log.Trace().
				Caller().
				// Err(err). // Timeout error produces Get "url" error with the query.
				Str("engine", e.Name.String()).
				// Str("url", urll). // Can't reliably anonymize it (because it's engine dependent).
				Msg("Request timeout error for url")
		} else {
			log.Error().
				Caller().
				Err(err).
				Str("engine", e.Name.String()).
				// Str("url", urll). // Can't reliably anonymize it (because it's engine dependent).
				Bytes("response", r.Body). // WARN: Query can be present, depending on the response from the engine.
				Msg("Request error for url")
		}
	})
}
