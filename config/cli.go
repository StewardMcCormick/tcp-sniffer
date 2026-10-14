package config

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

type uint16Value uint16

func (v *uint16Value) String() string {
	return strconv.FormatUint(uint64(*v), 10)
}

func (v *uint16Value) Set(s string) error {
	val, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return fmt.Errorf("number should be from 0 to 65535")
	}
	*v = uint16Value(uint16(val))
	return nil
}

func parseCLIArgsInCfg(cfg *AppConfig) {
	flagIfaceName := flag.String("i", "", "capturing interface name")
	flagSnapshotLength := uint16Value(65535)
	flag.Var(&flagSnapshotLength, "snaplen", "max captured snapshot length")
	flagPromiscuousMode := flag.Bool("promisc", false, "enable promiscuous mode")
	flagTimeout := flag.Duration("t", 25*time.Millisecond, "RX_RING timeout, if '0' - immediate mode")

	flag.Parse()

	set := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		set[f.Name] = true
	})

	if set["i"] {
		cfg.Capturing.Interface = *flagIfaceName
	}
	if set["snaplen"] {
		cfg.Capturing.SnapshotLength = uint16(flagSnapshotLength)
	}
	if set["promisc"] {
		cfg.Capturing.Promisc = *flagPromiscuousMode
	}
	if set["t"] {
		cfg.Capturing.Timeout = *flagTimeout
	}
}
