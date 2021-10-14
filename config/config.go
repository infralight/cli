package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/infralight/cli/version"
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
	c, err = loadProductConfig(version.Product, profile)
	if err != nil {
		// try to load configuration from the old product name
		var oldErr error
		c, oldErr = loadProductConfig(version.OldProduct, profile)
		if oldErr != nil {
			return c, err
		}
	}

	return c, nil
}

func loadProductConfig(product, profile string) (c Config, err error) {
	path, err := xdg.SearchConfigFile(
		fmt.Sprintf("%s/%s.toml", strings.ToLower(product), profile),
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
		fmt.Sprintf("%s/%s.toml", strings.ToLower(version.Product), c.Profile),
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
