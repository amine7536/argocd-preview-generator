package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppsConfig struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name     string `yaml:"name"`
	ImageTag string `yaml:"image_tag"`
}

func Load(path string) (*AppsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var cfg AppsConfig
	if unmarshalErr := yaml.Unmarshal(data, &cfg); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, unmarshalErr)
	}

	return &cfg, nil
}
