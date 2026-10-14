package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/StewardMcCormick/tcp-analyzer/internal/sniffer"
	"github.com/ilyakaznacheev/cleanenv"
)

var ErrInvalidConfig = errors.New("invalid config")

type AppConfig struct {
	Capturing sniffer.Config `yaml:"capturing"`
}

func (a *AppConfig) validateConfig() error {
	if a.Capturing.Interface == "" {
		return fmt.Errorf(
			"%w - interface name cannot be empty - use '-i' flag or set it in 'config.yaml",
			ErrInvalidConfig,
		)
	}

	return nil
}

func InitConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	if err := cleanenv.ReadConfig("config.yaml", cfg); err != nil {
		log.Print("[WARNING] Cannot read 'config.yaml' file")
	}

	parseCLIArgsInCfg(cfg)

	if err := cfg.validateConfig(); err != nil {
		return nil, fmt.Errorf("cannot read config - %w", err)
	}

	return cfg, nil
}
