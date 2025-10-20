package main

import (
	"log"

	"github.com/Woland-prj/dilemator/config"
	"github.com/Woland-prj/dilemator/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
