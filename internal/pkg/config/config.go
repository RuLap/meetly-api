package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Log struct {
		Level      string            `yaml:"level"`
		File       string            `yaml:"file"`
		LokiURL    string            `yaml:"loki_url"`
		LokiLabels map[string]string `yaml:"loki_labels"`
	} `yaml:"log"`

	HTTPServer struct {
		Address     string        `yaml:"address"`
		Timeout     time.Duration `yaml:"timeout"`
		IdleTimeout time.Duration `yaml:"idle_timeout"`
	} `yaml:"http_server"`

	PostgresConnString string `yaml:"postgres_conn_string"`

	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

func Load(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("cannot parse config: %v", err)
	}

	cfg.Log.LokiURL = os.ExpandEnv(cfg.Log.LokiURL)
	cfg.PostgresConnString = os.ExpandEnv(cfg.PostgresConnString)
	cfg.JWT.Secret = os.ExpandEnv(cfg.JWT.Secret)

	return &cfg
}
