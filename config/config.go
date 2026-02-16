package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppsConfig struct {
	Namespace string    `yaml:"namespace"`
	Services  []Service `yaml:"services"`
	Infra     []Infra   `yaml:"infra"`
}

type Service struct {
	Name     string `yaml:"name"`
	ImageTag string `yaml:"image_tag"`
}

type Infra struct {
	Name           string      `yaml:"name"`
	Chart          string      `yaml:"chart"`
	RepoURL        string      `yaml:"repoURL"`
	TargetRevision string      `yaml:"targetRevision"`
	Values         interface{} `yaml:"values"`
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
