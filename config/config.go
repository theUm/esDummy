package config

import "github.com/theUm/esDummy/elastic"

type Config struct {
	Log                 LoggerConfig
	ElasticConfig       elastic.Config
	HealthCheckHTTPPort int `env:"HEALTHCHECK_PORT" envDefault:"8000"`
}

type LoggerConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	Pretty   bool   `env:"LOG_PRETTY" envDefault:"true"`
}
