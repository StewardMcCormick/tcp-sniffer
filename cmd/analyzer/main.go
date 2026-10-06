package main

import (
	"fmt"
	"log"

	"github.com/StewardMcCormick/tcp-analyzer/config"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.Capturing.Interface)
}
