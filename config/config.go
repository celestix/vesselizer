package config

import (
	"os"
	"vessel/errors"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// name of the vessel
	Name string `yaml:"name"`
	// base image used by the vessel
	BaseImage string `yaml:"base"`
	// build file for building the vessel
	BuildFile string `yaml:"buildfile"`
	// entrypoint of the vessel
	Entrypoint string `yaml:"entrypoint"`
}

func ParseFile(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var c Config
	err = yaml.NewDecoder(file).Decode(&c)
	if err != nil {
		return nil, err
	}
	if c.BaseImage == "" {
		return nil, errors.ErrBaseImageNotSpecified
	}
	if c.BuildFile == "" {
		return nil, errors.ErrBuildFileNotSpecified
	}
	if c.Entrypoint == "" {
		return nil, errors.ErrEntrypointNotSpecified
	}
	if _, ok := Images[c.BaseImage]; !ok {
		return nil, errors.ErrSpecifiedBaseImageNotFound
	}
	return &c, nil
}
