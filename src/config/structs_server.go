package config

import "time"

// ReaderServer is format in which the config is read from the config file
type ReaderServer struct {
	// Environment in which the server is running (normal or lambda)
	Environment string `koanf:"environment"`
	// Port on which the API server listens
	Port int `koanf:"port"`
	// URLs used for CORS (wildcards allowed)
	// comma separated
	FrontendUrls string `koanf:"frontendurls"`
	// Cache settings
	Cache ReaderCache `koanf:"cache"`
	// Image proxy settings
	Proxy ReaderImageProxy `koanf:"proxy"`
}
type Server struct {
	// Environment in which the server is running (normal or lambda)
	Environment string
	// Port on which the API server listens
	Port int
	// URLs used for CORS (wildcards allowed)
	FrontendUrls []string
	// Cache settings
	Cache Cache
	// Image proxy settings
	Proxy ImageProxy
}

// ReaderCache is format in which the config is read from the config file
type ReaderCache struct {
	// Can be "none", "badger" or "redis"
	Type string `koanf:"type"`
	// Has no effect if type is "none"
	TTL ReaderTTL `koanf:"ttl"`
	// Badger specific settings
	Badger Badger `koanf:"badger"`
	// Redis specific settings
	Redis Redis `koanf:"redis"`
}
type Cache struct {
	// Can be "none", "badger" or "redis"
	Type string
	// Has no effect if type is "none"
	TTL TTL
	// Badger specific settings
	Badger Badger
	// Redis specific settings
	Redis Redis
}

// ReaderTTL is format in which the config is read from the config file
// In <number><unit> format
// Example: 1s, 1m, 1h, 1d, 1w, 1M, 1y
// If unit is not specified, it is assumed to be milliseconds
type ReaderTTL struct {
	// how long to store the results in cache
	// setting this to 0 caches the results forever
	Time string `koanf:"time"`
	// if the remaining TTL when retrieving from cache is less than this, update the cache entry and reset the TTL
	// setting this to 0 disables this feature
	// setting this to the same value (or higher) as Results will update the cache entry every time
	RefreshTime string `koanf:"refreshtime"`
}
type TTL struct {
	// How long to store the results in cache
	// Setting this to 0 caches the results forever
	Time time.Duration
	// If the remaining TTL when retrieving from cache is less than this, update the cache entry and reset the TTL
	// Setting this to 0 disables this feature
	// Setting this to the same value (or higher) as Results will update the cache entry every time
	RefreshTime time.Duration
}

type Badger struct {
	// Setting this to false will result in badger not persisting the cache to disk making it run "in-memory"
	Persist bool `koanf:"persist"`
}

type Redis struct {
	Host     string `koanf:"host"`
	Port     uint16 `koanf:"port"`
	Password string `koanf:"password"`
	Database uint8  `koanf:"database"`
}

// ReaderProxy is format in which the config is read from the config file
type ReaderImageProxy struct {
	Salt     string                   `koanf:"salt"`
	Timeouts ReaderImageProxyTimeouts `koanf:"timeouts"`
}
type ImageProxy struct {
	Salt     string
	Timeouts ImageProxyTimeouts
}

// ReaderProxyTimeouts is format in which the config is read from the config file
// In <number><unit> format
// Example: 1s, 1m, 1h, 1d, 1w, 1M, 1y
// If unit is not specified, it is assumed to be milliseconds
type ReaderImageProxyTimeouts struct {
	Dial         string `koanf:"dial"`
	KeepAlive    string `koanf:"keepalive"`
	TLSHandshake string `koanf:"tlshandshake"`
}
type ImageProxyTimeouts struct {
	Dial         time.Duration
	KeepAlive    time.Duration
	TLSHandshake time.Duration
}
