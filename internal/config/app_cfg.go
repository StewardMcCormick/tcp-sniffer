package config

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
)

var ErrInvalidAppConfig error = errors.New("invalid application config")

type AppConfig struct {
	Capturing SnifferConfig `mapstructure:"capturing"`
}

func (ac *AppConfig) Validate() error {
	if err := ac.Capturing.Validate(); err != nil {
		return fmt.Errorf("%w - %w", ErrInvalidAppConfig, err)
	}

	return nil
}

func InitConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	if err := cleanenv.ReadConfig("config.yaml", cfg); err != nil {
		if !errors.Is(err, io.EOF) {
			log.Fatalf("[ERROR] Cannot read 'config.yaml' file - %v", err)
		}
	}

	if err := cfg.Capturing.Validate(); err != nil {
		log.Fatalf("[ERROR] Invalid configuration - %v", err)
	}

	return cfg, nil
}

func Load() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[WARNING] Cannot read configuration file - %v", err)
	}

	cfg := &AppConfig{}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("[ERROR] Cannot parse configuration - %v", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("cannot load config - %w", err)
	}

	return cfg, nil
}
