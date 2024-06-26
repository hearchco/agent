package googlescholar

import (
	"github.com/hearchco/agent/src/search/scraper"
)

type Engine struct {
	scraper.EngineBase
}

func New() *Engine {
	return &Engine{scraper.EngineBase{
		Name:    seName,
		Origins: origins[:],
	}}
}
