package main

import (
	"errors"
	"log"
)

var ErrApplicationCannotRun = errors.New("cannot run application")

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("[ERROR] - %v", err)
	}
}
