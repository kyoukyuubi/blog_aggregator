package config

import (
	"encoding/json"
	"os"
)

func Read() (Config, error) {
	var cfg Config

	path, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}