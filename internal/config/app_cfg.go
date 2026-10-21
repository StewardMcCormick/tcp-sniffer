package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var ErrInvalidAppConfig error = errors.New("invalid application config")

var ErrCannotReadConfig error = errors.New("cannot read configuration file")

type AppConfig struct {
	Capturing SnifferConfig `mapstructure:"capturing"`
}

func (ac *AppConfig) Validate() error {
	if err := ac.Capturing.Validate(); err != nil {
		return fmt.Errorf("%w - %w", ErrInvalidAppConfig, err)
	}

	return nil
}

func Load(configPath string) (*AppConfig, error) {
	viper.SetConfigFile(configPath)

	return load(viper.GetViper())
}

func load(v *viper.Viper) (*AppConfig, error) {
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w - %w", ErrCannotReadConfig, err)
		}
	}

	cfg := &AppConfig{}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config - %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("cannot load config - %w", err)
	}

	return cfg, nil
}
