package config

import (
	"time"

	"github.com/joeshaw/envdecode"
)

type (
	// Config stores complete configuration
	Config struct {
		HTTP  HTTPConfig
		App   AppConfig
		Cloud CloudConfig
	}

	// HTTPConfig stores HTTP configuration
	HTTPConfig struct {
		Hostname     string        `env:"HOSTNAME,default=0.0.0.0"`
		Port         uint16        `env:"PORT,default=8080"`
		ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT,default=5s"`
		WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT,default=10s"`
		IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT,default=2m"`
	}

	// AppConfig stores application configuration
	AppConfig struct {
		Timeout time.Duration `env:"APP_TIMEOUT,default=3s"`
	}

	// CloudConfig stores cloud configuration
	CloudConfig struct {
		Project     string `env:"GCP_PROJECT,default=apitest-359415"`
		Credentials string `env:"GOOGLE_APPLICATION_CREDENTIALS"`
		PubSubTopic string `env:"APP_PUBSUB_TOPIC,default=events"`
	}
)

// GetConfig loads and returns configuration
func GetConfig() (Config, error) {
	var cfg Config
	err := envdecode.StrictDecode(&cfg)
	return cfg, err
}
