package scraper

import (
	"github.com/gocolly/colly/v2"
)

// OnHTML registers a function. Function will be executed on every HTML
// element matched by the GoQuery Selector parameter.
// GoQuery Selector is a selector used by https://github.com/PuerkitoBio/goquery.
func (e *EngineBase) OnHTML(goquerySelector string, f colly.HTMLCallback) {
	e.collector.OnHTML(goquerySelector, f)
}

// OnResponse registers a function. Function will be executed on every response.
func (e *EngineBase) OnResponse(f colly.ResponseCallback) {
	e.collector.OnResponse(f)
}

// OnRequest registers a function. Function will be executed on every
// request made by the Collector.
func (e *EngineBase) OnRequest(f colly.RequestCallback) {
	e.collector.OnRequest(f)
}

// Wait returns when the collector jobs are finished.
func (e EngineBase) Wait() {
	e.collector.Wait()
}
