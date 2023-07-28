package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/gocolly/colly/v2/queue"
)

// Search returns a list of search results.
func Search(ctx context.Context, seUrl string, query string, dom Dom, opts Options) ([]Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := RateLimit.Wait(ctx); err != nil {
		return nil, err
	}

	c := colly.NewCollector(colly.MaxDepth(1))
	if len(opts) == 0 {
		opts = append(opts, Options{})
	}

	if opts[0].UserAgent == "" {
		c.UserAgent = DefaultUserAgent
	} else {
		c.UserAgent = opts[0].UserAgent
	}

	q, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})

	limit := opts[0].Limit

	results := []Result{}
	nextPageLink := ""
	var rErr error
	filteredRank := 1
	rank := 1

	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			rErr = err
			return
		}
		if opts[0].FollowNextPage && nextPageLink != "" {
			req, err := r.New("GET", nextPageLink, nil)
			if err == nil {
				q.AddRequest(req)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		rErr = err
	})

	// https://www.w3schools.com/cssref/css_selectors.asp
	c.OnHTML(dom.Result, func(e *colly.HTMLElement) {

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

		// check if there is a next button at the end.
		nextPageHref, _ := sel.Find(dom.NextPage).Attr("href")
		nextPageLink = strings.TrimSpace(nextPageHref)

	})

	/*
		c.OnHTML(dom.NextPage, func(e *colly.HTMLElement) {

			sel := e.DOM

			// check if there is a next button at the end.
			if nextPageHref, exists := sel.Attr("href"); exists {
				start := getStart(strings.TrimSpace(nextPageHref))
				nextPageLink = buildUrl(seUrl, query, limit, start)
				q.AddURL(nextPageLink)
			} else {
				nextPageLink = ""
			}
		})
	*/

	url := buildUrl(seUrl, query, limit, 0)

	if opts[0].ProxyAddr != "" {
		rp, err := proxy.RoundRobinProxySwitcher(opts[0].ProxyAddr)
		if err != nil {
			return nil, err
		}
		c.SetProxyFunc(rp)
	}

	q.AddURL(url)
	q.Run(c)

	if rErr != nil {
		if strings.Contains(rErr.Error(), "Too Many Requests") {
			return nil, ErrBlocked
		}
		return nil, rErr
	}

	// Reduce results to max limit
	if opts[0].Limit != 0 && len(results) > opts[0].Limit {
		return results[:opts[0].Limit], nil
	}

	return results, nil
}

func buildUrl(seUrl string, query string, limit int, start int) string {
	query = strings.Trim(query, " ")
	query = strings.Replace(query, " ", "+", -1)

	var url string

	if start == 0 {
		url = fmt.Sprintf("%s%s", seUrl, query)
	} else {
		url = fmt.Sprintf("%s%s&start=%d", seUrl, query, start)
	}

	if limit != 0 {
		url = fmt.Sprintf("%s&num=%d", url, limit)
	}

	return url
}
