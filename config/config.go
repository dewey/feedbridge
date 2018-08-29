package config

// Config holds our app configuration
type Config struct {
	RefreshInterval   int    `env:"REFRESH_INTERVAL" envDefault:"15"`
	CacheExpiration   int    `env:"CACHE_EXPIRATION" envDefault:"30"`
	CacheExpiredPurge int    `env:"CACHE_EXPIRED_PURGE" envDefault:"60"`
	StorageBackend    string `env:"STORAGE_BACKEND" envDefault:"memory"`
	StoragePath       string `env:"STORAGE_PATH" envDefault:"/feedbridge-data"`
	Environment       string `env:"ENVIRONMENT" envDefault:"develop"`
	Port              int    `env:"PORT" envDefault:"8080"`
}
