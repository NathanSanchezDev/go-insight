// internal/config/config.go
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		User           string `yaml:"user"`
		Name           string `yaml:"name"`
		Port           int    `yaml:"port"`
		MaxConnections int    `yaml:"max_connections"`
	} `yaml:"database"`

	RateLimit struct {
		RequestsPerMinute int `yaml:"requests_per_minute"`
		WindowMinutes     int `yaml:"window_minutes"`
	} `yaml:"rate_limiting"`

	Features struct {
		DebugLogging      bool `yaml:"debug_logging"`
		PrometheusEnabled bool `yaml:"prometheus_enabled"`
		JWTAuthEnabled    bool `yaml:"jwt_auth_enabled"`
	} `yaml:"features"`

	Monitoring struct {
		Prometheus struct {
			Path           string `yaml:"path"`
			ScrapeInterval string `yaml:"scrape_interval"`
		} `yaml:"prometheus"`
	} `yaml:"monitoring"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile("config/app.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
