package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StewardMcCormick/tcp-analyzer/internal/config"
	"github.com/StewardMcCormick/tcp-analyzer/internal/logger"
	"github.com/StewardMcCormick/tcp-analyzer/internal/sniffer"
	"github.com/StewardMcCormick/tcp-analyzer/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgPath string

var log *slog.Logger

var rootCmd = &cobra.Command{
	Use:           "mccsniff",
	Short:         "TCP sniffer by Stw Mcc.",
	Long:          "TCP sniffer pet-project by Steward McCormick.",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func run() error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("%w - %w", ErrApplicationCannotRun, err)
	}

	log = logger.New(cfg.Logger, os.Stderr)

	flowStorage := storage.NewRwMutexMapStorage()
	capturer, err := sniffer.NewPcapCapturer(cfg.Capturing, log, flowStorage)
	if err != nil {
		return fmt.Errorf("%w - %w", ErrApplicationCannotRun, err)
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- capturer.Run()
	}()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	case <-sig:
	}

	return capturer.Stop()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yaml", "config filepath")

	rootCmd.Flags().StringP("interface", "i", "", "capturing interface name.")
	rootCmd.Flags().Uint16P("snaplen", "s", 65535, "maximum captured snapshot length.")
	rootCmd.Flags().BoolP("promics", "p", false, "is promiscuous_mode enabled.")
	rootCmd.Flags().DurationP("timeout", "t", 25*time.Millisecond, "RX_RING timeout, if '0' - immediate mode.")
	rootCmd.Flags().StringP("engine.type", "e", "", "capturing Engine type")
	rootCmd.Flags().String("log.level", "info", "slog logging level")
	rootCmd.Flags().String("log.format", "text", "slog logging format")

	if err := viper.BindPFlag("capturing.interface", rootCmd.Flags().Lookup("interface")); err != nil {
		log.Error("Cannot bind flag",
			"error", err,
		)
	}
	if err := viper.BindPFlag("capturing.max_snapshot_length", rootCmd.Flags().Lookup("snaplen")); err != nil {
		log.Error("Cannot bind flag",
			"error", err,
		)
	}
	if err := viper.BindPFlag("capturing.promiscuous_mode", rootCmd.Flags().Lookup("promics")); err != nil {
		log.Error("Cannot bind flag",
			"error", err,
		)
	}
	if err := viper.BindPFlag("capturing.timeout", rootCmd.Flags().Lookup("timeout")); err != nil {
		log.Error("Cannot bind flag",
			"error", err,
		)
	}
	if err := viper.BindPFlag("capturing.engine.type", rootCmd.Flags().Lookup("engine.type")); err != nil {
		log.Error("Cannot bind flag",
			"error", err,
		)
	}
	if err := viper.BindPFlag("logger.level", rootCmd.Flags().Lookup("log.level")); err != nil {
		log.Error("Cannot bind flag",
			"error", err,
		)
	}
	if err := viper.BindPFlag("logger.format", rootCmd.Flags().Lookup("log.format")); err != nil {
		log.Error("Cannot bind flag",
			"error", err,
		)
	}
}
