package config

import (
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
)

type Settings struct {
	RequestedResultsPerPage int      `koanf:"requestedresults"`
	Shortcut                string   `koanf:"shortcut"`
	Proxies                 []string `koanf:"proxies"`
}

// ReaderConfig is format in which the config is read from the config file
type ReaderConfig struct {
	Server      ReaderServer                     `koanf:"server"`
	RCategories map[category.Name]ReaderCategory `koanf:"categories"`
	Settings    map[string]Settings              `koanf:"settings"`
}
type Config struct {
	Server     Server
	Categories map[category.Name]Category
	Settings   map[engines.Name]Settings
}
