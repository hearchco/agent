package config

import (
	"github.com/hearchco/agent/src/search/category"
)

// ReaderConfig is format in which the config is read from the config file and environment variables.
type ReaderConfig struct {
	Server      ReaderServer                     `koanf:"server"`
	RCategories map[category.Name]ReaderCategory `koanf:"categories"`
	RExchange   ReaderExchange                   `koanf:"exchange"`
}
type Config struct {
	Server     Server
	Categories map[category.Name]Category
	Exchange   Exchange
}
