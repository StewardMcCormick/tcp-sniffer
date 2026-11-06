package config

import (
	"errors"
	"fmt"
	"time"
)

type EngineType string

const (
	EngineTypeSequential EngineType = "sequential"
)

var ErrInvalidEngineType = errors.New("invalid engine type - should be one of 'sequential' or 'shared'")

func (e *EngineType) Validate() error {
	switch *e {
	case EngineTypeSequential:
		return nil
	default:
		return ErrInvalidEngineType
	}
}

var ErrInvalidEngineConfig = errors.New("invalid capturing engine config")

type EngineConfig struct {
	Type EngineType `mapstructure:"type"`
}

func (ec *EngineConfig) Validate() error {
	if err := ec.Type.Validate(); err != nil {
		return fmt.Errorf("%w - %w", ErrInvalidEngineConfig, err)
	}

	return nil
}

var (
	ErrInvalidSnifferConfig = errors.New("invalid sniffer config")
	ErrInvalidInterfaceName = errors.New("invalid interface name")
)

type SnifferConfig struct {
	Interface      string        `mapstructure:"interface"`
	SnapshotLength uint16        `mapstructure:"max_snapshot_length"`
	Promisc        bool          `mapstructure:"promiscuous_mode"`
	Timeout        time.Duration `mapstructure:"timeout"`
	Engine         EngineConfig  `mapstructure:"engine"`
}

func (sc *SnifferConfig) Validate() error {
	errs := make([]error, 0)

	if sc.Interface == "" {
		errs = append(errs, fmt.Errorf("%w - should not be empty", ErrInvalidInterfaceName))
	}

	if err := sc.Engine.Validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return nil
	}

	if sc.Timeout <= 0 {
		sc.Timeout = -10 * time.Millisecond
	}

	return fmt.Errorf("%w - %w", ErrInvalidSnifferConfig, errors.Join(errs...))
}
