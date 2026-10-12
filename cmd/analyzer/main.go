package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/StewardMcCormick/tcp-analyzer/config"
	"github.com/StewardMcCormick/tcp-analyzer/internal/sniffer"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	engine := sniffer.NewSequentialEngine()
	capturer, err := sniffer.NewCapturer(cfg.Capturing.Interface, engine)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := capturer.Run(); err != nil {
			panic(err)
		}
	}()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	if err = capturer.Stop(); err != nil {
		panic(err)
	}
}
