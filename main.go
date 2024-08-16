package main

import (
	"log"
	"path/filepath"

	"github.com/denniswon/tcex/app"
)

func main() {
	configFile, err := filepath.Abs(".env")
	if err != nil {
		log.Fatalf("[!] Failed to find `.env` : %s\n", err.Error())
	}
	app.Run(configFile)
}
