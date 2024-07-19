package config

import (
	"time"
)

// ReaderServer is format in which the config is read from the config file and environment variables.
type ReaderServer struct {
	// Environment in which the server is running (normal or lambda).
	Environment string `koanf:"environment"`
	// Port on which the API server listens.
	Port int `koanf:"port"`
	// URLs used for CORS (wildcards allowed).
	// Comma separated.
	FrontendUrls string `koanf:"frontendurls"`
	// Cache settings.
	Cache ReaderCache `koanf:"cache"`
	// Image proxy settings.
	ImageProxy ReaderImageProxy `koanf:"imageproxy"`
}
type Server struct {
	// Environment in which the server is running (normal or lambda).
	Environment string
	// Port on which the API server listens.
	Port int
	// URLs used for CORS (wildcards allowed).
	FrontendUrls []string
	// Cache settings.
	Cache Cache
	// Image proxy settings.
	ImageProxy ImageProxy
}

// ReaderCache is format in which the config is read from the config file and environment variables.
type ReaderCache struct {
	// Can be "none" or "redis".
	Type string `koanf:"type"`
	// Prefix to use for cache keys.
	KeyPrefix string `koanf:"keyprefix"`
	// Has no effect if type is "none".
	TTL ReaderTTL `koanf:"ttl"`
	// Redis specific settings.
	Redis Redis `koanf:"redis"`
}
type Cache struct {
	// Can be "none" or "redis".
	Type string
	// Prefix to use for cache keys.
	KeyPrefix string
	// Has no effect if type is "none".
	TTL TTL
	// Redis specific settings.
	Redis Redis
	// DynamoDB specific settings.
	DynamoDB DynamoDB
}

// ReaderTTL is format in which the config is read from the config file and environment variables.
// In <number><unit> format.
// Example: 1s, 1m, 1h, 1d, 1w, 1M, 1y.
// If unit is not specified, it is assumed to be milliseconds.
type ReaderTTL struct {
	// How long to store the results in cache.
	// Setting this to 0 caches the results forever.
	Results string `koanf:"time"`
	// How long to store the currencies in cache.
	// Setting this to 0 caches the currencies forever.
	Currencies string `koanf:"currencies"`
}
type TTL struct {
	// How long to store the results in cache.
	// Setting this to 0 caches the results forever.
	Results time.Duration
	// How long to store the currencies in cache.
	// Setting this to 0 caches the currencies forever.
	Currencies time.Duration
}

type Redis struct {
	Host     string `koanf:"host"`
	Port     uint16 `koanf:"port"`
	Password string `koanf:"password"`
	Database uint8  `koanf:"database"`
}

type DynamoDB struct {
	Region string `koanf:"region"`
	Table  string `koanf:"table"`
	// Endpoint is used for local testing.
	Endpoint string `koanf:"endpoint"`
}

// ReaderProxy is format in which the config is read from the config file and environment variables.
// In <number><unit> format.
// Example: 1s, 1m, 1h, 1d, 1w, 1M, 1y.
// If unit is not specified, it is assumed to be milliseconds.
type ReaderImageProxy struct {
	SecretKey string `koanf:"secretkey"`
	Timeout   string `koanf:"timeout"`
}
type ImageProxy struct {
	Salt    string
	Timeout time.Duration
}
