package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/gocolly/colly/v2/queue"
	"github.com/tminaorg/brzaguza/src/search/useragent"
)

// Search returns a list of search results.
func Search(ctx context.Context, searchEngineURL string, query string, dom DOMPaths, options Options) ([]Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := RateLimit.Wait(ctx); err != nil {
		return nil, err
	}

	collector := colly.NewCollector(colly.MaxDepth(1))

	if options.UserAgent == "" {
		collector.UserAgent = useragent.DefaultUserAgent()
	} else {
		collector.UserAgent = options.UserAgent
	}

	requestQueue, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})

	limit := options.Limit

	results := []Result{}
	//nextPageLink := ""
	var rErr error
	filteredRank := 1
	rank := 1

	collector.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			rErr = err
			return
		}

		/*
			if options.FollowNextPage && nextPageLink != "" {
				req, err := r.New("GET", nextPageLink, nil)
				if err == nil {
					requestQueue.AddRequest(req)
				}
			}
		*/
	})

	collector.OnError(func(r *colly.Response, err error) {
		rErr = err
	})

	// https://www.w3schools.com/cssref/css_selectors.asp
	collector.OnHTML(dom.Result, func(e *colly.HTMLElement) {

		sel := e.DOM

		linkHref, _ := sel.Find(dom.Link).Attr("href")
		linkText := strings.TrimSpace(linkHref)
		titleText := strings.TrimSpace(sel.Find(dom.Title).Text())
		descText := strings.TrimSpace(sel.Find(dom.Description).Text())

		rank += 1
		if linkText != "" && linkText != "#" && titleText != "" {
			result := Result{
				Rank:        filteredRank,
				URL:         linkText,
				Title:       titleText,
				Description: descText,
			}
			results = append(results, result)
			filteredRank += 1
		}

		/*
			// check if there is a next button at the end.
			nextPageHref, _ := sel.Find(dom.NextPage).Attr("href")
			nextPageLink = strings.TrimSpace(nextPageHref)
		*/

	})

	/*
		collector.OnHTML(dom.NextPage, func(e *colly.HTMLElement) {

			sel := e.DOM

			// check if there is a next button at the end.
			if nextPageHref, exists := sel.Attr("href"); exists {
				start := getStart(strings.TrimSpace(nextPageHref))
				nextPageLink = buildUrl(searchEngineURL, query, limit, start)
				requestQueue.AddURL(nextPageLink)
			} else {
				nextPageLink = ""
			}
		})
	*/

	url := buildUrl(searchEngineURL, query, limit, 0)

	if options.ProxyAddr != "" {
		rp, err := proxy.RoundRobinProxySwitcher(options.ProxyAddr)
		if err != nil {
			return nil, err
		}
		collector.SetProxyFunc(rp)
	}

	requestQueue.AddURL(url)
	requestQueue.Run(collector)

	if rErr != nil {
		if strings.Contains(rErr.Error(), "Too Many Requests") {
			return nil, ErrRateLimited
		}
		return nil, rErr
	}

	// Reduce results to max limit
	if options.Limit != 0 && len(results) > options.Limit {
		return results[:options.Limit], nil
	}

	return results, nil
}

func buildUrl(searchEngineUrl string, query string, limit int, start int) string {
	query = strings.Trim(query, " ")
	query = strings.Replace(query, " ", "+", -1)

	var url string

	if start == 0 {
		url = fmt.Sprintf("%s%s", searchEngineUrl, query)
	} else {
		url = fmt.Sprintf("%s%s&start=%d", searchEngineUrl, query, start)
	}

	if limit != 0 {
		url = fmt.Sprintf("%s&num=%d", url, limit)
	}

	return url
}
