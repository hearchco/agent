package config

// ReaderConfig is format in which the config is read from the config file and environment variables.
type ReaderConfig struct {
	Server    ReaderServer                  `koanf:"server"`
	REngines  map[string]ReaderEngineConfig `koanf:"engines"`
	RExchange ReaderExchange                `koanf:"exchange"`
}
type Config struct {
	Server   Server
	Engines  EngineConfig
	Exchange Exchange
}
