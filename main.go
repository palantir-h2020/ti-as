package main

import (
	"PalantirIpAnonymizationService/backend"
	"PalantirIpAnonymizationService/logger"
	"PalantirIpAnonymizationService/utils"
	"os"
	"path/filepath"
)

func main() {
	// Initialize logger
	var log = logger.New(false)

	if len(os.Args) == 1 {
		log.Error("At least one argument must be provided.")
		log.Error("Example Usage: ")
		log.Error("\tPalantirIPAnonymizationService <CONFIG_FILE>")

		os.Exit(-1)
	}

	// Get program arguments. First arg is program name
	configFile := os.Args[1]

	// Get project directory
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	log.Info(exPath)

	// Load service configuration
	var config utils.Config
	config.LoadConf(configFile)

	// Initialize and start service
	a := backend.App{}
	a.Initialize(&config)
	a.Run(config.Server.Host + ":" + config.Server.Port)
}
