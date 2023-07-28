package google

import (
	"context"

	"github.com/tminaorg/brzaguza/src/search"
)

var url string = "https://www.google.com/search?q="

var options search.Options = search.Options{
	UserAgent:      "",
	Limit:          0,
	ProxyAddr:      "",
	FollowNextPage: false,
}

var domQuery search.DomQuery = search.DomQuery{
	Result: "",
}

func Search(ctx context.Context, query string) {
	if ctx == nil {
		ctx = context.Background()
	}

	search.Search(ctx, url, query, domQuery, options)
}
