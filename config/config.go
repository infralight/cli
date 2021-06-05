package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Profile             string
	URL                 string
	AuthorizationHeader string
	AccessKey           string
	SecretKey           string
}

var (
	ErrConfigNotFound = errors.New("profile configuration not found")
)

func LoadConfig(profile string) (c Config, err error) {
	path, err := xdg.SearchConfigFile(
		fmt.Sprintf("infralight/%s.toml", profile),
	)
	if err != nil {
		return c, ErrConfigNotFound
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf(
			"failed reading config file %s: %w",
			path, err,
		)
	}

	err = toml.Unmarshal(data, &c)
	if err != nil {
		return c, fmt.Errorf(
			"failed parsing config file %s: %w",
			path, err,
		)
	}

	return c, nil
}

func (c Config) Save() (path string, err error) {
	path, err = xdg.ConfigFile(
		fmt.Sprintf("infralight/%s.toml", c.Profile),
	)
	if err != nil {
		return path, fmt.Errorf(
			"failed getting location for file: %w",
			err,
		)
	}

	b, err := toml.Marshal(c)
	if err != nil {
		return path, fmt.Errorf(
			"failed encoding configuration: %w",
			err,
		)
	}

	err = os.WriteFile(path, b, 0600)
	if err != nil {
		return path, fmt.Errorf(
			"failed writing %s: %w",
			path, err,
		)
	}

	return path, nil
}
