package main

import (
	"errors"
)

var ErrApplicationCannotRun = errors.New("cannot run application")

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("Cannot start application",
			"error", err,
		)
	}
}
