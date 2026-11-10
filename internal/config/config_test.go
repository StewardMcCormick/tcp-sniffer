package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/viper"
)

const validConfigYAML = `capturing:
  interface: "en0"
  max_snapshot_length: 96
  promiscuous_mode: true
  timeout: "25ms"
  engine:
    type: "sequential"
`

const emptyInterfaceYAML = `capturing:
  interface: ""
  max_snapshot_length: 96
  promiscuous_mode: false
  timeout: "25ms"
  engine:
    type: "sequential"
`

const invalidEngineYAML = `capturing:
  interface: "en0"
  max_snapshot_length: 96
  promiscuous_mode: false
  timeout: "25ms"
  engine:
    type: "shared"
`

func validAppConfig() *AppConfig {
	return &AppConfig{
		Capturing: SnifferConfig{
			Interface:      "en0",
			SnapshotLength: 96,
			Promisc:        true,
			Timeout:        25 * time.Millisecond,
			Engine:         EngineConfig{Type: EngineTypeSequential},
		},
	}
}

func TestEngineType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		typ     EngineType
		wantErr error
	}{
		{name: "sequential", typ: EngineTypeSequential, wantErr: nil},
		{name: "shared", typ: EngineType("shared"), wantErr: ErrInvalidEngineType},
		{name: "empty", typ: EngineType(""), wantErr: ErrInvalidEngineType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.typ.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EngineType(%q).Validate() = %v, want %v", tt.typ, err, tt.wantErr)
			}
		})
	}
}

func TestEngineConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     EngineConfig
		wantErr error
	}{
		{name: "valid", cfg: EngineConfig{Type: EngineTypeSequential}, wantErr: nil},
		{name: "invalid type", cfg: EngineConfig{Type: EngineType("shared")}, wantErr: ErrInvalidEngineConfig},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("EngineConfig.Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestSnifferConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     SnifferConfig
		wantErr []error
	}{
		{
			name: "valid",
			cfg: SnifferConfig{
				Interface: "en0",
				Engine:    EngineConfig{Type: EngineTypeSequential},
			},
			wantErr: nil,
		},
		{
			name: "empty interface",
			cfg: SnifferConfig{
				Interface: "",
				Engine:    EngineConfig{Type: EngineTypeSequential},
			},
			wantErr: []error{ErrInvalidInterfaceName},
		},
		{
			name: "invalid engine type",
			cfg: SnifferConfig{
				Interface: "en0",
				Engine:    EngineConfig{Type: EngineType("shared")},
			},
			wantErr: []error{ErrInvalidEngineConfig},
		},
		{
			name: "multiple errors",
			cfg: SnifferConfig{
				Interface: "",
				Engine:    EngineConfig{Type: EngineType("shared")},
			},
			wantErr: []error{ErrInvalidInterfaceName, ErrInvalidEngineConfig},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if len(tt.wantErr) == 0 {
				if err != nil {
					t.Fatalf("SnifferConfig.Validate() = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("SnifferConfig.Validate() = nil, want error")
			}
			if !errors.Is(err, ErrInvalidSnifferConfig) {
				t.Errorf("SnifferConfig.Validate() = %v, want wrapping %v", err, ErrInvalidSnifferConfig)
			}
			for _, want := range tt.wantErr {
				if !errors.Is(err, want) {
					t.Errorf("SnifferConfig.Validate() = %v, want wrapping %v", err, want)
				}
			}
		})
	}
}

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     AppConfig
		wantErr error
	}{
		{
			name: "valid",
			cfg: AppConfig{
				Capturing: SnifferConfig{
					Interface: "en0",
					Engine:    EngineConfig{Type: EngineTypeSequential},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid capturing",
			cfg: AppConfig{
				Capturing: SnifferConfig{
					Interface: "",
					Engine:    EngineConfig{Type: EngineTypeSequential},
				},
			},
			wantErr: ErrInvalidAppConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AppConfig.Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		cli     map[string]any
		want    *AppConfig
		wantErr error
	}{
		{
			name:    "valid",
			content: validConfigYAML,
			want:    validAppConfig(),
		},
		{
			name:    "empty interface",
			content: emptyInterfaceYAML,
			wantErr: ErrInvalidInterfaceName,
		},
		{
			name:    "invalid engine type",
			content: invalidEngineYAML,
			wantErr: ErrInvalidEngineConfig,
		},
		{
			name:    "malformed file",
			content: "capturing: [unclosed",
			wantErr: ErrCannotReadConfig,
		},
		{
			name: "cli only",
			cli: map[string]any{
				"capturing.interface":   "en0",
				"capturing.engine.type": "sequential",
			},
			want: &AppConfig{
				Capturing: SnifferConfig{
					Interface: "en0",
					Engine:    EngineConfig{Type: EngineTypeSequential},
				},
			},
		},
		{
			name:    "missing all configurations",
			wantErr: ErrInvalidInterfaceName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()

			path := filepath.Join(t.TempDir(), "config.yaml")
			if tt.content != "" {
				if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
					t.Fatalf("write config file: %v", err)
				}
			}
			v.SetConfigFile(path)

			for key, val := range tt.cli {
				v.Set(key, val)
			}

			got, err := load(v)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("load() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("load() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		content  string
		want     *AppConfig
		wantErr  error
	}{
		{
			name:     "valid custom path and name",
			fileName: "custom.yml",
			content:  validConfigYAML,
			want:     validAppConfig(),
			wantErr:  nil,
		},
		{
			name:     "malformed file",
			fileName: "custom.yml",
			content:  "capturing: [unclosed",
			want:     nil,
			wantErr:  ErrCannotReadConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			dir := t.TempDir()
			path := filepath.Join(dir, tt.fileName)
			if tt.content != "" {
				if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
					t.Fatalf("write config file: %v", err)
				}
			}

			got, err := Load(path)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Load() error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

