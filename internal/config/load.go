// Package config loads RADIUS Director configuration documents.
package config

import (
	"os"

	"github.com/gobcn/radius-director/internal/model"
	"gopkg.in/yaml.v3"
)

// Load reads a YAML configuration document.
func Load(path string) (model.Configuration, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Configuration{}, err
	}

	var configuration model.Configuration
	if err := yaml.Unmarshal(data, &configuration); err != nil {
		return model.Configuration{}, err
	}

	return configuration, nil
}
