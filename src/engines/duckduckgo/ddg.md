# Duck Duck Go
Send a [post request](https://github.com/gocolly/colly/issues/175#issuecomment-400024313) to https://lite.duckduckgo.com/lite/ with body: `q=<query>&dc=<rank of first result on page>`. It will return 20-22 results.


http://api.jquery.com/index/


```
c.OnHTML("div.filters > table > tbody", func(e *colly.HTMLElement) {
        var linkHref  string
        var linkText  string
        var titleText string
        var descText  string
        counter := 1
        
        e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
            sel := el.DOM
            switch counter % 4 {
            case 1:
                linkHref, _ = sel.Find("a.result-link").Attr("href")
                titleText = strings.TrimSpace(sel.Find("a.result-link").Text())
            case 2:
                descText = strings.TrimSpace(sel.Find("td.result-snippet").Text())
            case 3:
                linkText = strings.TrimSpace(sel.Find("span.link-text").Text())
                if strings.Contains(linkHref, "https") {
                    linkText = "https://" + linkText;
                } else {
                    linkText = "http://" + linkText;
                }
            case 0:
                rank += 1
                if linkText != "" && linkText != "#" && titleText != "" {
                    result := structures.Result{
                        Rank:        filteredRank,
                        URL:         linkText,
                        Title:       titleText,
                        Description: descText,
                    }
                    results = append(results, result)
                    filteredRank += 1
                }

                // check if there is a next button at the end.
                nextPageHref, _ := sel.Find("a #pnnext").Attr("href")
                nextPageLink = strings.TrimSpace(nextPageHref)
            }

            counter += 1
        })

    })

```