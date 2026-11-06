package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StewardMcCormick/tcp-analyzer/internal/config"
	"github.com/StewardMcCormick/tcp-analyzer/internal/sniffer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgPath string

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

	capturer, err := sniffer.NewPcapCapturer(cfg.Capturing)
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

	if err := viper.BindPFlag("capturing.interface", rootCmd.Flags().Lookup("interface")); err != nil {
		log.Fatalf("[ERROR] Cannot bind flag - %v", err)
	}
	if err := viper.BindPFlag("capturing.max_snapshot_length", rootCmd.Flags().Lookup("snaplen")); err != nil {
		log.Fatalf("[ERROR] Cannot bind flag - %v", err)
	}
	if err := viper.BindPFlag("capturing.promiscuous_mode", rootCmd.Flags().Lookup("promics")); err != nil {
		log.Fatalf("[ERROR] Cannot bind flag - %v", err)
	}
	if err := viper.BindPFlag("capturing.timeout", rootCmd.Flags().Lookup("timeout")); err != nil {
		log.Fatalf("[ERROR] Cannot bind flag - %v", err)
	}
	if err := viper.BindPFlag("capturing.engine.type", rootCmd.Flags().Lookup("engine.type")); err != nil {
		log.Fatalf("[ERROR] Cannot bind flag - %v", err)
	}
}
