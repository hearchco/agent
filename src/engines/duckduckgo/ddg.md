# Duck Duck Go
Send a [POST request](https://github.com/gocolly/colly/issues/175#issuecomment-400024313) to `https://lite.duckduckgo.com/lite/` with body: `q=<query>&dc=<rank of first result on page>`. It will return 20-22 results. GET requests could be used like `https://lite.duckduckgo.com/lite/?q=<query>&dc=<rank of first result on page>`.

First request could be: col.PostRaw(Info.URL, []byte("q="+query+"&dc=1"))

This may be useful: http://api.jquery.com/index/

The href on the title sometimes contains telemetry, and is not a valid URL then. That's why we fetch the scheme from it, and append it to the span text.