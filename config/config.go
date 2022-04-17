package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Storage StorageConfig `json:"storage"`
}

type StorageConfig struct {
	// number of days the data is stored
	TTL     int    `json:"ttl"`
	BaseDir string `json:"dir"`
}

func ParseConfig(configPath string) (*Config, error) {
	fileBody, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(fileBody, &cfg)
	if err != nil {
		return nil, err
	}

	err = validateConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Storage.TTL < 1 {
		return errors.New("storage config: ttl must be > 1")
	}
	if len(cfg.Storage.BaseDir) == 0 {
		return errors.New("storage config: base dir name must be not empty")
	}
	return nil
}
