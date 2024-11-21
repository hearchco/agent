package yep

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"strings"
// 	"sync/atomic"

// 	"github.com/gocolly/colly/v2"
// 	"github.com/rs/zerolog/log"

// 	"github.com/hearchco/agent/src/search/engines/options"
// 	"github.com/hearchco/agent/src/search/result"
// 	"github.com/hearchco/agent/src/search/scraper"
// 	"github.com/hearchco/agent/src/search/scraper/parse"
// 	"github.com/hearchco/agent/src/utils/anonymize"
// 	"github.com/hearchco/agent/src/utils/morestrings"
// )

// func (se Engine) WebSearch(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
// 	foundResults := atomic.Bool{}
// 	retErrors := make([]error, 0, opts.Pages.Max)
// 	pageRankCounter := scraper.NewPageRankCounter(opts.Pages.Max)

// 	se.OnRequest(func(r *colly.Request) {
// 		r.Headers.Set("Accept", "*/*")
// 	})

// 	se.OnResponse(func(r *colly.Response) {
// 		body := string(r.Body)
// 		prefix := "[\"Ok\","
// 		suffix := ']'
// 		index := strings.Index(body, prefix)

// 		if index != 0 || body[len(body)-1] != byte(suffix) {
// 			log.Error().
// 				Caller().
// 				Str("engine", se.Name.String()).
// 				Str("body", body).
// 				Str("prefix", prefix).
// 				Str("suffix", string(suffix)).
// 				Msg("Failed parsing response, couldn't find start/end of JSON")
// 			return
// 		}

// 		// starts after prefix and ends before suffix
// 		// so after "[\"Ok\"," and before "]"
// 		resultsJson := body[len(prefix) : len(body)-1]
// 		var content jsonResponse
// 		if err := json.Unmarshal([]byte(resultsJson), &content); err != nil {
// 			log.Error().
// 				Caller().
// 				Err(err).
// 				Str("engine", se.Name.String()).
// 				Str("body", body).
// 				Str("content", resultsJson).
// 				Msg("Failed parsing response, couldn't unmarshal JSON")
// 			return
// 		}

// 		pageIndex := se.PageFromContext(r.Request.Ctx)
// 		page := pageIndex + opts.Pages.Start + 1

// 		for _, jsonResult := range content.Results {
// 			if jsonResult.TType != "Organic" {
// 				continue
// 			}

// 			goodURL, goodTitle, goodDesc := parse.SanitizeFields(jsonResult.URL, jsonResult.Title, jsonResult.Snippet)

// 			r, err := result.ConstructResult(se.Name, goodURL, goodTitle, goodDesc, page, pageRankCounter.GetPlusOne(pageIndex))
// 			if err != nil {
// 				log.Error().
// 					Caller().
// 					Err(err).
// 					Msg("Failed to construct result")
// 			} else {
// 				log.Trace().
// 					Caller().
// 					Int("page", page).
// 					Int("rank", pageRankCounter.GetPlusOne(pageIndex)).
// 					Str("result", fmt.Sprintf("%v", r)).
// 					Msg("Sending result to channel")
// 				resChan <- r
// 				pageRankCounter.Increment(pageIndex)
// 				if !foundResults.Load() {
// 					foundResults.Store(true)
// 				}
// 			}
// 		}
// 	})

// 	// Static params.
// 	paramLocale := localeParamString(opts.Locale)
// 	paramSafeSearch := safeSearchParamString(opts.SafeSearch)

// 	for i := range opts.Pages.Max {
// 		pageNum := i + opts.Pages.Start
// 		ctx := colly.NewContext()
// 		ctx.Put("page", strconv.Itoa(i))

// 		// Dynamic params.
// 		paramPage := ""
// 		if pageNum > 0 {
// 			paramPage = fmt.Sprintf("%v=%v", paramKeyPage, (pageNum+2)*10+1)
// 		}

// 		combinedParamsLeft := morestrings.JoinNonEmpty("?", "&", paramClient, paramLocale, paramPage, paramNo_correct)
// 		combinedParamsRight := morestrings.JoinNonEmpty("&", "&", paramSafeSearch, paramType)

// 		// Non standard order of params required
// 		urll := fmt.Sprintf("%v%v&q=%v%v", searchURL, combinedParamsLeft, query, combinedParamsRight)
// 		anonUrll := fmt.Sprintf("%v%v&q=%v%v", searchURL, combinedParamsLeft, anonymize.String(query), combinedParamsRight)

// 		if err := se.Get(ctx, urll, anonUrll); err != nil {
// 			retErrors = append(retErrors, err)
// 		}
// 	}

// 	se.Wait()
// 	close(resChan)
// 	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
// }
